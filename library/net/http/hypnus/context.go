package hypnus

import (
	"context"
	"math"
	"net/http"
)

const (
	_abortIndex int8 = math.MaxInt8 / 2
)

type Context struct {
	context.Context

	Request *http.Request
	Writer  http.ResponseWriter

	handlers []HandlerFunc
	engine   *Engine

	index int8 // control flow

	method string // http method

}

// iterate the handlers
func (c *Context) Next() {

	if c.Request.Method != c.method {
		// TODO handle err
		return
	}

	c.index++
	len := int8(len(c.handlers))
	for ; c.index < len; c.index++ {
		c.handlers[c.index](c)
	}
}

// cancel the handler iteration. Note that this will not stop the current handler
func (c *Context) Abort() {
	c.index = _abortIndex
}
