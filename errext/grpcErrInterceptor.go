package errs

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

func CodeErrorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// TODO: fill your logic here
	resp1, err1 := handler(ctx, req)
	if err1 != nil {
		logx.WithContext(ctx).Errorf("GRPC_FUNC:%v err: %v,req:%#v,chickenObj:%#v\n", err, info.FullMethod, err, req)
		return nil, WarpGrpcErr(err1)
	}
	return resp1, nil
}
func CodeErrorSteamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, ss)
	if err != nil {
		return WarpGrpcErr(err)
	}
	return nil
}
