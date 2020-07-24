package hypnus

import (
	"context"
	"encoding/json"
	"math"
	"net/http"

	ecode "github.com/ihornet/go-open/v2/library/ecode"
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

	Error error // for response

	// store require parameters. ParseForm() just support application/x-www-form-urlencoded,
	// so add the field to store parameters
	Req struct {
		// can't understand go use map[string][]string ???
		Header map[string]string
		Body   map[string]string
		Body2  map[string]interface{}
		Query  map[string]string
		Param  map[string]string
	}
}

// iterate the handlers
func (c *Context) Next() {

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
		Status:  bcode.Code(),
		Success: bcode.Status(),
		Message: bcode.Message(),
		Data:    data,
	}
	c.preHandleJson(obj)
	ret, _ := json.Marshal(obj)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Write(ret)
}

func (c *Context) preHandleJson(obj *respObj) {
	if obj.Success == ecode.SUCCESS && obj.Data == nil {
		obj.Data = struct {
			Success bool `json:"success"`
		}{Success: true}
	} else if obj.Success == ecode.FAILED {
		obj.Data = nil
	}
}

// just for body
func (c *Context) Bind(obj interface{}) error {
	if err := Bind(c.Req.Body, obj); err != nil {
		c.Abort()
		c.JSON(nil, err)
		return err
	}
	return nil
}

// response struct
type respObj struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
