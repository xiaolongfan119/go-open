package rpc

import (
	"context"
	"math"
	"sync"
	"time"

	xtime "github.com/xiaolongfan119/go-open/v2/library/time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

var (
	_abortIndex  int8 = math.MaxInt8 / 2
	_defaultConf      = &ClientConfig{
		Dial:    xtime.Duration(time.Second * 20),
		TimeOut: xtime.Duration(time.Millisecond * 250),
	}
)

type ClientConfig struct {
	Dial    xtime.Duration
	TimeOut xtime.Duration
}

type Client struct {
	opt      []grpc.DialOption
	conf     *ClientConfig
	mutex    sync.RWMutex
	handlers []grpc.UnaryClientInterceptor
}

func NewClient(conf *ClientConfig, opt ...grpc.DialOption) *Client {
	c := new(Client)
	c.SetConfig(conf)

	c.UserOpt(grpc.WithBalancerName(roundrobin.Name))
	c.UserOpt(opt...)

	c.Use(c.recovery(), c.logging())
	return c
}

func (c *Client) SetConfig(conf *ClientConfig) (err error) {
	if conf == nil {
		c.conf = _defaultConf
	}
	if conf.Dial <= 0 {
		conf.Dial = xtime.Duration(time.Second * 10)
	}
	if conf.TimeOut <= 0 {
		conf.TimeOut = xtime.Duration(time.Millisecond * 250)
	}
	c.mutex.Lock()
	c.conf = conf
	c.mutex.Unlock()

	return nil
}

func (c *Client) Dial(ctx context.Context, target string) (conn *grpc.ClientConn, err error) {
	c.opt = append(c.opt, grpc.WithBlock())
	c.opt = append(c.opt, grpc.WithInsecure())
	c.opt = append(c.opt, grpc.WithUnaryInterceptor(c.chainUnaryClient()))
	c.mutex.RLock()
	if c.conf.Dial > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(c.conf.Dial))
		defer cancel()
	}

	conn, err = grpc.DialContext(ctx, target, c.opt...)
	err = errors.WithStack(err)
	return
}

func (c *Client) Use(handlers ...grpc.UnaryClientInterceptor) *Client {
	size := len(c.handlers) + len(handlers)
	if size >= int(_abortIndex) {
		// TODO panic error
	}
	mergedHandlers := make([]grpc.UnaryClientInterceptor, size)
	copy(mergedHandlers, c.handlers)
	copy(mergedHandlers[len(c.handlers):], handlers)
	c.handlers = mergedHandlers
	return c
}

func (c *Client) UserOpt(opt ...grpc.DialOption) *Client {
	c.opt = append(c.opt, opt...)
	return c
}

func (c *Client) chainUnaryClient() grpc.UnaryClientInterceptor {
	n := len(c.handlers)
	if n == 0 {
		return func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}

	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		var (
			i            int
			chainHandler grpc.UnaryInvoker
		)

		chainHandler = func(ictx context.Context, imethod string, ireq, ireply interface{}, icc *grpc.ClientConn, iopts ...grpc.CallOption) error {
			if i == n-1 {
				return invoker(ictx, imethod, ireq, ireply, icc, iopts...)
			}

			i++
			return c.handlers[i](ictx, imethod, ireq, ireply, icc, chainHandler, iopts...)
		}
		return c.handlers[0](ctx, method, req, reply, cc, chainHandler, opts...)
	}
}
