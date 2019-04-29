package hypnus

type RouterGroup struct {
	Handlers []HandlerFunc
	basePath string
	engine   *Engine
}

func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers...),
		basePath: JoinPaths(group.basePath, relativePath),
		engine:   group.engine,
	}
}

func (group *RouterGroup) combineHandlers(handlers ...HandlerFunc) []HandlerFunc {
	size := len(group.Handlers) + len(handlers)
	if size > int(_abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make([]HandlerFunc, size)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) {
	group.handle("GET", relativePath, handlers...)
}

func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) {
	group.handle("PUT", relativePath, handlers...)
}

func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) {
	group.handle("POST", relativePath, handlers...)
}

func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) {
	group.handle("DELETE", relativePath, handlers...)
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers ...HandlerFunc) {
	path := JoinPaths(group.basePath, relativePath)
	handlers = group.combineHandlers(handlers...)
	// set path and handler as normal.  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request)
	group.engine.addRoute(httpMethod, path, handlers...)
}
