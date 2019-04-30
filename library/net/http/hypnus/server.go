package hypnus

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	xtime "github.com/ihornet/go-commom/library/time"
)

type HandlerFunc func(*Context)

type Engine struct {
	RouterGroup
	conf      *ServerConf
	mux       *http.ServeMux
	metastore map[string]map[string]interface{}
}

type ServerConf struct {
	Network      string
	Address      string
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
}

func DefaultServer(conf *ServerConf) *Engine {
	engine := NewServer(conf)
	return engine
}

func NewServer(conf *ServerConf) *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
		},
		metastore: make(map[string]map[string]interface{}),
		mux:       http.NewServeMux(),
		conf:      conf,
	}

	engine.RouterGroup.engine = engine
	return engine
}

func (engine *Engine) Start() {
	conf := engine.conf
	if conf.Network == "" {
		conf.Network = "tcp"
	}
	l, err := net.Listen(conf.Network, conf.Address)
	if err != nil {
		fmt.Println("[ ERR ] hypnus: listen tcp:", err)
		panic(err)
	}
	server := &http.Server{
		ReadTimeout:  time.Duration(conf.ReadTimeout),
		WriteTimeout: time.Duration(conf.WriteTimeout),
	}

	if err = engine.RunServer(server, l); err != nil {
		fmt.Println("[ ERR ] hypnus: launch sever:", err)
		panic(err)
	}
}

func (engine *Engine) RunServer(server *http.Server, ln net.Listener) (err error) {
	server.Handler = engine.mux
	if err = server.Serve(ln); err != nil {
		return
	}
	return
}

func (engine *Engine) addRoute(method, path string, handlers ...HandlerFunc) {
	if path[0] != '/' {
		panic("hypnus: path must begin with '/'")
	}
	if method == "" {
		panic("hypnus: http method can't be empty")
	}
	if len(handlers) == 0 {
		panic("hypnus: handlers can't be empty")
	}
	if _, ok := engine.metastore[path]; !ok {
		engine.metastore[path] = make(map[string]interface{})
	}

	engine.metastore[path]["method"] = method

	engine.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		c := &Context{
			engine:   engine,
			Request:  req,
			Writer:   w,
			handlers: handlers,
			method:   method,
		}
		engine.handleContext(c)
	})
}

func (engine *Engine) handleContext(c *Context) {
	engine.parseReqParams(c)
	// iterate handlers
	c.Next()
}

func (engine *Engine) parseReqParams(c *Context) {

	req := c.Request
	cType := req.Header.Get("Content-Type")
	c.Req.Body = make(map[string]interface{})

	switch {
	case strings.Contains(cType, "application/json"):
		json.NewDecoder(req.Body).Decode(&c.Req.Body)
	case strings.Contains(cType, "application/x-www-form-urlencoded"):
		req.ParseForm()
		for k, v := range req.PostForm {
			c.Req.Body[k] = v[0]
		}
	default:
		// TODO panic error
	}

	c.Req.Query = make(map[string]interface{})
	if vs, err := url.ParseQuery(req.URL.RawQuery); err == nil {
		for k, v := range vs {
			c.Req.Query[k] = v[0]
		}
	}
}
