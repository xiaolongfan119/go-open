package rpc

import (
	"context"
	"fmt"
	"go-open/library/ecode"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) recovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:size]
				fmt.Printf("grpc server panic: %v\n%v\n%s\n", req, err, buf)
				err = status.Errorf(codes.Unknown, ecode.ServerErr.Error())
			}
		}()

		resp, err = handler(ctx, req)
		return
	}
}
