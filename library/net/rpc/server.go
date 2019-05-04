package rpc

import (
	"context"
	xtime "go-open/library/time"
	"math"
	"sync"
	"time"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

// import (
// 	"bufio"
// 	"encoding/gob"
// 	"fmt"
// 	"io"
// 	"net"
// 	"sync"
// )

// type ServerConfig struct {
// 	Proto string
// 	Addr  string
// }

// type service struct{}

// type Server struct {
// 	lis        net.Listener
// 	serviceMap map[string]*service
// }

// func NewServer(c *ServerConfig) *Server {
// 	if c == nil {
// 		// TODO panic error
// 	}

// 	s := newServer()
// 	go rpcListen(c, s)
// 	return s
// }

// func rpcListen(c *ServerConfig, s *Server) {
// 	l, err := net.Listen(c.Proto, c.Addr)
// 	if err != nil {
// 		// TODO panic error
// 	}
// 	defer func() {
// 		if err := l.Close(); err != nil {
// 			// TODO panic err
// 		}
// 		fmt.Println("rpc listener closed")
// 	}()
// 	fmt.Println("start rpc listen addr:" + c.Addr)
// 	go s.Accept(l)
// }

// func (server *Server) Accept(lis net.Listener) {
// 	for {
// 		conn, err := lis.Accept()
// 		if err != nil {
// 			fmt.Printf("rpc.Server: accept: %s", err.Error())
// 			return
// 		}
// 		go server.ServeConn(conn)
// 	}
// }

// func (server *Server) ServeConn(conn net.Conn) {

// }

// func newServer() *Server {
// 	return &Server{serviceMap: make(map[string]*service)}
// }

// type serverCodec struct {
// 	sending sync.Mutex
// 	resp    Response
// 	req     Request
// 	auth    Auth
// 	rwc     io.ReadWriteCloser
// 	dec     *gob.Decoder
// 	enc     *gob.Encoder
// 	encBuf  *bufio.Writer
// 	addr    net.Addr
// 	closed  bool
// }

// type Response struct{}
// type Request struct{}
// type Auth struct{}

/////////////////////////////////

var (
	_abortIndex int8 = math.MaxInt8 / 2
)

type ServerConfig struct {
	Network string
	Addr    string
	// pings the client to see if the transport is still alive.
	Time xtime.Duration
	// After having pinged for keepalive check, the server waits for a duration
	// of Timeout and if no activity is seen even after that the connection is
	// closed.
	Timeout xtime.Duration
	// timout of indle
	IdleTimeout xtime.Duration
	// maxLifeTime of exist
	MaxLifeTime xtime.Duration
	// ForceCloseWait is an additive period after MaxLifeTime after which the connection will be forcibly closed.
	ForceCloseWait xtime.Duration
}

type Server struct {
	conf     *ServerConfig
	mutex    sync.RWMutex
	server   *grpc.Server
	handlers []grpc.UnaryServerInterceptor
}

func NewServer(conf *ServerConfig, opt ...grpc.ServerOption) (s *Server) {
	if conf == nil {
		//TODO panic error
	}

	s = new(Server)
	if err := s.setConfig(conf); err != nil {
		// TODO panic error
	}

	keepParam := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(conf.IdleTimeout),
		MaxConnectionAgeGrace: time.Duration(conf.ForceCloseWait),
		Time:                  time.Duration(conf.Time),
		Timeout:               time.Duration(conf.Timeout),
	})

	opt = append(opt, keepParam, grpc.UnaryInterceptor(s.interceptor))
	s.server = grpc.NewServer(opt...)
	s.Use(s.recovery(), s.logging())
	return s
}

func (s *Server) setConfig(conf *ServerConfig) (err error) {
	if conf.Addr == "" {
		conf.Addr = "0.0.0.0:90000"
	}
	if conf.Network == "" {
		conf.Network = "tcp"
	}

	if conf.IdleTimeout <= 0 {
		conf.IdleTimeout = xtime.Duration(time.Second * 60)
	}
	if conf.MaxLifeTime <= 0 {
		conf.MaxLifeTime = xtime.Duration(time.Hour * 2)
	}
	if conf.ForceCloseWait <= 0 {
		conf.ForceCloseWait = xtime.Duration(time.Second * 20)
	}

	if conf.Time <= 0 {
		conf.Time = xtime.Duration(time.Second * 60)
	}

	if conf.Timeout <= 0 {
		conf.Time = xtime.Duration(time.Second * 20)
	}

	s.mutex.Lock()
	s.conf = conf
	s.mutex.Unlock()
	return
}

func (s *Server) Use(handlers ...grpc.UnaryServerInterceptor) *Server {
	size := len(s.handlers) + len(handlers)
	if size > int(_abortIndex) {
		panic("rpc: server use too many handler")
	}

	mergedHandlers := make([]grpc.UnaryServerInterceptor, size)
	copy(mergedHandlers, s.handlers)
	copy(mergedHandlers[len(s.handlers):], handlers)
	s.handlers = mergedHandlers
	return s
}

func (s *Server) interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var (
		i     int
		chain grpc.UnaryHandler
	)

	n := len(s.handlers)
	if n == 0 {
		return handler(ctx, req)
	}

	chain = func(ct context.Context, r interface{}) (interface{}, error) {
		if i == n-1 {
			return handler(ct, r)
		}
		i++
		return s.handlers[i](ctx, req, info, chain)
	}

	return s.handlers[0](ctx, req, info, chain)
}

func (s *Server) Register(receiver interface{}) (err error) {

}
