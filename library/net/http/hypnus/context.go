package hypnus

import (
	"context"
	"encoding/json"
	"math"
	"net/http"

	ecode "github.com/ihornet/go-commom/library/ecode"
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

	Error error // for response

	// store require parameters. ParseForm() just support application/x-www-form-urlencoded,
	// so add the field to store parameters
	Req struct {
		// can't understand go use map[string][]string ???
		Body  map[string]interface{}
		Query map[string]interface{}
	}
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

// serializes the data to json, and reponse to client
func (c *Context) JSON(data interface{}, err error) {

	bcode := ecode.Cause(err)
	obj := &respObj{
		Code:    bcode.Code(),
		Status:  bcode.Status(),
		Message: bcode.Message(),
		Data:    data,
	}

	ret, _ := json.Marshal(obj)
	c.Writer.Write(ret)
}

// response struct
type respObj struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
