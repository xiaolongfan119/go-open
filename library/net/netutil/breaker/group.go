package breaker

var defaultBreakerGroup *BreakerGroup

type BreakerGroup struct {
	breakers map[string]*CircuitBreaker
}

func NewBreakerGroup() *BreakerGroup {

	if defaultBreakerGroup != nil {
		return defaultBreakerGroup
	}

	group := new(BreakerGroup)
	group.breakers = make(map[string]*CircuitBreaker)
	defaultBreakerGroup = group
	return group
}

func SetBreaker() *BreakerGroup {

}
