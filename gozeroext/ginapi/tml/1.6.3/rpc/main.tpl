package main

import (
	"flag"
	"fmt"

	{{.imports}}

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/heshiyingx/gotool/netext"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/heshiyingx/gotool/gozeroext/zaplog"
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

func main() {
	flag.Parse()
	zapWriter, err := zaplog.NewZapWriter()
	if err != nil {
		logx.Must(err)

	}
	logx.SetWriter(zapWriter)

	var c config.Config
	conf.MustLoad(*configFile, &c)
	port, err := netext.GetAvailablePort()

	if err != nil {
		panic(err)
		return
	}
	c.ListenOn = "0.0.0.0:" + strconv.Itoa(port)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
{{range .serviceNames}}       {{.Pkg}}.Register{{.Service}}Server(grpcServer, {{.ServerPkg}}.New{{.Service}}Server(ctx))
{{end}}
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
