package breaker

import (
	"errors"
	xtime "go-open/library/time"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

const defaultTimeout = time.Duration(30) * time.Second

var defaultMaxConsecutiveFailures uint32 = 5 // close时， 最大连续错误 ---> 变成open

var (
	ErrTooManyRequests = errors.New("too many requests")
	ErrOpenState       = errors.New("open state")
)

type Counts struct {
	Requests             uint32 // 总请求次数
	TotalSuccesses       uint32 // 总成功次数
	TotalFailures        uint32 // 总失败次数
	ConsecutiveSuccesses uint32 // 连续成功次数
	ConsecutiveFailures  uint32 // 连续失败次数
}

type BreakerConfig struct {
	Name                   string                                  // breaker name
	MaxRequests            uint32                                  // helfOpen时，最大请求次数
	MaxConsecutiveFailures uint32                                  // close时， 最大连续错误 ---> 变成open
	Interval               xtime.Duration                          // close时，定期清理counts
	Timeout                xtime.Duration                          // open时， timeout后进去halfOpen
	ReadyToTrip            func(counts Counts) bool                // close时， 当遇到错误时调用
	OnStateChange          func(name string, from State, to State) // 状态变化时调用
}

type CircuitBreaker struct {
	conf       *BreakerConfig
	mu         sync.Mutex
	state      State
	generation uint64
	counts     Counts
	expiry     time.Time
}

func (counts *Counts) clear() {
	counts.Requests = 0
	counts.TotalSuccesses = 0
	counts.TotalFailures = 0
	counts.ConsecutiveSuccesses = 0
	counts.ConsecutiveFailures = 0
}

func (counts *Counts) onSuccess() {
	counts.TotalSuccesses++
	counts.ConsecutiveSuccesses++
	counts.ConsecutiveFailures = 0
}

func (counts *Counts) onFailure() {
	counts.TotalFailures++
	counts.ConsecutiveFailures++
	counts.ConsecutiveSuccesses++

}

func (counts *Counts) onRequest() {
	counts.Requests++
}

func (breaker *CircuitBreaker) onSuccess() {

	breaker.counts.onSuccess()
	if breaker.state == StateHalfOpen {
		// TODO
		if breaker.counts.ConsecutiveSuccesses >= breaker.conf.MaxRequests {
			breaker.setState(StateClosed, time.Now())
		}
	}
}

func (breaker *CircuitBreaker) onFailure() {
	breaker.counts.onFailure()
	now := time.Now()
	switch breaker.state {
	case StateClosed:
		if breaker.conf.ReadyToTrip(breaker.counts) {
			breaker.setState(StateOpen, now)
		}
	case StateHalfOpen:
		breaker.setState(StateOpen, now)
	}
}

func NewCircuitBreaker(c *BreakerConfig) *CircuitBreaker {

	breaker := new(CircuitBreaker)

	if c.Timeout == 0 {
		c.Timeout = xtime.Duration(defaultTimeout)
	}

	if c.MaxRequests == 0 {
		c.MaxRequests = 1
	}

	if c.MaxConsecutiveFailures > 0 {
		defaultMaxConsecutiveFailures = c.MaxConsecutiveFailures
	}

	if c.ReadyToTrip == nil {
		c.ReadyToTrip = defaultReadyToTrip
	}

	breaker.conf = c
	return breaker
}

func (breaker *CircuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	generation, err := breaker.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := recover(); err != nil {
			breaker.onFailure()
			panic(err)
		}
	}()

	result, err := req()
	breaker.afterRequest(generation, err == nil)
	return result, err
}

func defaultReadyToTrip(counts Counts) bool {
	return counts.ConsecutiveFailures >= defaultMaxConsecutiveFailures
}

func (breaker *CircuitBreaker) beforeRequest() (generation uint64, err error) {

	breaker.mu.Lock()
	defer breaker.mu.Unlock()

	now := time.Now()
	state, generation := breaker.currentState(now)

	if state == StateOpen {
		return generation, ErrOpenState
	} else if state == StateHalfOpen && breaker.counts.Requests >= breaker.conf.MaxRequests {
		return generation, ErrTooManyRequests
	}

	breaker.counts.onRequest()
	return generation, nil
}

func (breaker *CircuitBreaker) afterRequest(before uint64, success bool) {

	breaker.mu.Lock()
	defer breaker.mu.Unlock()

	now := time.Now()
	_, generation := breaker.currentState(now)
	if before != generation {
		return
	}

	if success {
		breaker.onSuccess()
	} else {
		breaker.onFailure()
	}
}

func (breaker *CircuitBreaker) currentState(now time.Time) (State, uint64) {

	switch breaker.state {
	case StateClosed:
		if !breaker.expiry.IsZero() && breaker.expiry.Before(now) {
			breaker.toNewGeneration(now)
		}
	case StateOpen:
		if breaker.expiry.Before(now) {
			breaker.setState(StateHalfOpen, now)
		}
	}

	return breaker.state, breaker.generation
}

func (breaker *CircuitBreaker) setState(state State, now time.Time) {

	if breaker.state == state {
		return
	}

	prev := breaker.state
	breaker.state = state
	breaker.toNewGeneration(now)
	if breaker.conf.OnStateChange != nil {
		breaker.conf.OnStateChange(breaker.conf.Name, prev, breaker.state)
	}
}

func (breaker *CircuitBreaker) toNewGeneration(now time.Time) {
	breaker.generation++
	breaker.counts.clear()
	var zore time.Time
	switch breaker.state {
	case StateClosed:
		if breaker.conf.Interval == 0 {
			breaker.expiry = zore
		} else {
			breaker.expiry = now.Add(time.Duration(breaker.conf.Interval))
		}
	case StateOpen:
		breaker.expiry = now.Add(time.Duration(breaker.conf.Timeout))
	default:
		breaker.expiry = zore
	}
}
