package hypnus

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	xtime "github.com/ihornet/go-open/library/time"

	log "github.com/ihornet/go-open/library/log"
)

type HandlerFunc func(*Context)

var ServerName = ""

type Engine struct {
	RouterGroup
	conf   *ServerConf
	router *Router
}

type ServerConf struct {
	Port         string
	Host         string
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	Name         string
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
		router: NewRouter(),
		conf:   conf,
	}

	engine.RouterGroup.engine = engine

	ServerName = conf.Name

	return engine
}

func (engine *Engine) Start() {
	conf := engine.conf
	log.Info(fmt.Sprintf("server will launch, listening port %s", conf.Port))
	if err := http.ListenAndServe(conf.Port, engine); err != nil {
		log.Error("", fmt.Sprintf("server launch failed with port %s #### err: %v", conf.Port, err))
		return
	}
}

func (engine *Engine) RunServer(server *http.Server, ln net.Listener) (err error) {
	if err = server.Serve(ln); err != nil {
		return err
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

	engine.router.Handle(method, path, handlers...)
}

func (engine *Engine) handleContext(c *Context) {
	engine.parseReqParams(c)

	if c.Request.URL.Path != "/ping" {
		log.Info(">>>>>>>>>>>>【 ", c.Request.Method, " 】", c.Request.URL.Path,
			" ### param: ", c.Req.Param, " ### query: ", c.Req.Query, " ### body: ", c.Req.Body)
	}

	// iterate handlers
	c.Next()
}

func (engine *Engine) parseReqParams(c *Context) {

	req := c.Request
	cType := req.Header.Get("Content-Type")
	c.Req.Body = make(map[string]string)

	switch {
	case strings.Contains(cType, "application/json"):
		temp := make(map[string]interface{})
		json.NewDecoder(req.Body).Decode(&temp)
		c.Req.Body = convertMap2StrMap(temp)
	//	json.NewDecoder(req.Body).Decode(&c.Req.Body)
	case strings.Contains(cType, "application/x-www-form-urlencoded"):
		req.ParseForm()
		for k, v := range req.PostForm {
			c.Req.Body[k] = v[0]
		}
	default:
		// TODO panic error
	}

	c.Req.Query = make(map[string]string)
	if vs, err := url.ParseQuery(req.URL.RawQuery); err == nil {
		for k, v := range vs {
			c.Req.Query[k] = v[0]
		}
	}
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	defer Recovery()

	path := req.URL.Path
	if root := engine.router.trees[req.Method]; root != nil {
		if handles, ps, _ := root.getValue(path); handles != nil {
			c := &Context{
				engine:   engine,
				Request:  req,
				Writer:   w,
				handlers: handles,
			}
			c.Req.Param = ps.getAll()
			engine.handleContext(c)
			return
		}
	}
	if req.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                    //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,x-requested-with,token") //header的类型
		w.Header().Set("content-type", "application/json")                                    //返回数据格式是json
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,DELETE,POST,OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "1728000")
		w.Write([]byte(""))
	} else {
		if strings.Contains(req.URL.Path, "/debug/pprof/") {
			ProcessInput(req.URL.Path, w)
		} else {
			if req.URL.Path == "/ping" {
				w.Write([]byte("pong"))
			} else {
				log.Warn(fmt.Sprintf(">>>>>>>>> 404【 %s 】 %s", req.Method, req.URL.Path))
				http.NotFound(w, req)
			}
		}
	}
}

func (engine *Engine) ServeFiles(path string, root http.FileSystem) {
	engine.router.ServeFiles(path, root)
}
