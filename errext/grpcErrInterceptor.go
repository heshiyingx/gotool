package errs

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

func CodeErrorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp1, err1 := handler(ctx, req)
	if err1 != nil {
		logx.WithContext(ctx).Errorf("err: %v,GRPC_FUNC:%v,req:%#v\n", err, info.FullMethod, err, req)
		return nil, WarpGrpcErr(err1)
	}
	logx.WithContext(ctx).Debugf("GRPC_FUNC:%v,req:%#v,resp:%#v\n", info.FullMethod, req, resp1)
	return resp1, nil
}
func CodeErrorSteamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, ss)
	if err != nil {
		return WarpGrpcErr(err)
	}
	return nil
}
