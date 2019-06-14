package breaker

var defaultBreakerGroup *BreakerGroup

type BreakerGroup struct {
	breakers map[string]*CircuitBreaker
	Conf     *BreakerConfig
}

func NewBreakerGroup(conf *BreakerConfig) *BreakerGroup {

	if defaultBreakerGroup != nil {
		return defaultBreakerGroup
	}

	group := new(BreakerGroup)
	group.Conf = conf
	group.breakers = make(map[string]*CircuitBreaker)
	defaultBreakerGroup = group
	return group
}

func (group *BreakerGroup) GetBreaker(uri string) *CircuitBreaker {

	breaker := group.breakers[uri]
	if breaker == nil {
		breaker = NewCircuitBreaker(group.Conf)
		group.breakers[uri] = breaker
	}
	return breaker
}
