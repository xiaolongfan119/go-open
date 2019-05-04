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
