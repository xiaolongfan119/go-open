package breaker

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	xtime "github.com/ihornet/go-open/library/time"

	"github.com/stretchr/testify/assert"
)

func fail(cb *CircuitBreaker) error {
	msg := "fail"
	_, err := cb.Execute(func() (interface{}, error) { return nil, fmt.Errorf(msg) })
	if err.Error() == msg {
		return nil
	}
	return err
}

func success(cb *CircuitBreaker) error {
	_, err := cb.Execute(func() (interface{}, error) { return nil, nil })
	return err
}

func newCustom() *CircuitBreaker {
	config := new(BreakerConfig)
	config.Name = "cb"
	config.Interval = xtime.Duration(time.Duration(20) * time.Second)
	config.Timeout = xtime.Duration(time.Duration(10) * time.Second)
	config.ReadyToTrip = func(counts Counts) bool {
		return counts.ConsecutiveFailures > 5
	}
	config.OnStateChange = func(name string, from State, to State) {
		fmt.Printf("OnStateChange: name:%s,  from:%v,  to:%v \n", name, from, to)
	}
	return NewCircuitBreaker(config)
}

func TestNewCircuitBreaker(t *testing.T) {
	defaultCB := NewCircuitBreaker(&BreakerConfig{})
	assert.Equal(t, "", defaultCB.conf.Name)
	assert.Equal(t, uint32(1), defaultCB.conf.MaxRequests)
	assert.Equal(t, xtime.Duration(0), defaultCB.conf.Interval)
	//	assert.Equal(t, xtime.Duration(time.Duration(30)*time.Second), defaultCB.conf.Timeout)
	assert.NotNil(t, defaultCB.conf.ReadyToTrip)
	assert.Nil(t, defaultCB.conf.OnStateChange)
	assert.Equal(t, StateClosed, defaultCB.state)
	assert.True(t, defaultCB.expiry.IsZero())

	customCB := newCustom()
	assert.Equal(t, "cb", customCB.conf.Name)
	customCB.setState(StateOpen, time.Now())
	assert.Equal(t, StateOpen, customCB.state)

}

func TestDefaultCircuitBreaker(t *testing.T) {

	config := &BreakerConfig{}
	config.MaxRequests = 2
	config.OnStateChange = func(name string, from State, to State) {
		fmt.Printf("OnStateChange: name:%s,  from:%v,  to:%v \n", name, from, to)
	}

	defaultCB := NewCircuitBreaker(config)
	assert.Equal(t, "", defaultCB.conf.Name)

	for i := 0; i < 5; i++ {
		fail(defaultCB)
	}

	state, _ := defaultCB.currentState(time.Now())
	assert.Equal(t, StateOpen, defaultCB.state)
	assert.Equal(t, uint32(0), defaultCB.counts.ConsecutiveFailures)

	time.Sleep(time.Duration(30) * time.Second)
	state, _ = defaultCB.currentState(time.Now())
	assert.Equal(t, StateHalfOpen, state)

	success(defaultCB)
	state, _ = defaultCB.currentState(time.Now())
	assert.Equal(t, StateHalfOpen, state)

	fail(defaultCB)
	state, _ = defaultCB.currentState(time.Now())
	assert.Equal(t, StateOpen, state)

	time.Sleep(time.Duration(30) * time.Second)
	fmt.Println("#### 1 ###")
	success(defaultCB)
	success(defaultCB)

	fmt.Println("#### 2 ###")

	success(defaultCB)
	success(defaultCB)
	success(defaultCB)
	success(defaultCB)

	state, _ = defaultCB.currentState(time.Now())
	assert.Equal(t, StateClosed, state)

	// fmt.Printf("%v \n", defaultCB.counts)
	// assert.Equal(t, uint32(1), defaultCB.counts.ConsecutiveFailures)
}

func TestCircuitBreakerInParallel(t *testing.T) {

	defaultCB := NewCircuitBreaker(&BreakerConfig{})

	runtime.GOMAXPROCS(runtime.NumCPU())

	ch := make(chan error)

	const numReqs = 10000
	routine := func() {
		for i := 0; i < numReqs; i++ {
			ch <- success(defaultCB)
		}
	}

	const numRoutines = 10
	for i := 0; i < numRoutines; i++ {
		go routine()
	}

	total := uint32(numReqs * numRoutines)
	for i := uint32(0); i < total; i++ {
		err := <-ch
		assert.Nil(t, err)
	}

	assert.Equal(t, Counts{total, total, 0, total, 0}, defaultCB.counts)
}
