package counter

// import "sync"

// type Counter interface {
// 	Add(int64)
// 	Reset()
// 	Value() int64
// }

// type Group struct {
// 	mu       sync.RWMutex
// 	counters map[string]Counter
// 	New      func() Counter
// }

// func (g *Group) Add(key string, value int64) Counter {
// 	_, ok := g.counters[key]
// 	if !ok {
// 		g.counters[key] = g.New()
// 	}
// 	g.mu.RLock()
// 	g.counters[key].Add(value)
// 	defer g.mu.Unlock()
// 	return g.counters[key]
// }
