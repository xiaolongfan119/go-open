package rpc

import (
	"context"

	"google.golang.org/grpc"
)

func (s *Server) logging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// TODO logging
		resp, err = handler(ctx, req)
		return
	}
}

func (c *Client) logging() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// TODO logging
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
