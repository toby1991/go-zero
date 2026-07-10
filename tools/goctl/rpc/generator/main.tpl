package main

import (
	"flag"
	"fmt"
	"github.com/toby1991/go-zero-utils/pprof"

	{{.imports}}

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/main.yaml", "the config file")

func main() {
	flag.Parse()

    // config
	var c config.Config
	conf.MustLoad(*configFile, &c)

    // service group
	svcGroup := service.NewServiceGroup()
	defer svcGroup.Stop()

	// pprof
	svcGroup.Add(pprof.PprofServer(18888))

    // rpc server
	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
{{range .serviceNames}}       {{.Pkg}}.Register{{.GRPCService}}Server(grpcServer, {{.ServerPkg}}.New{{.Service}}Server(ctx))
{{end}}
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	svcGroup.Add(s)

    // start server
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	svcGroup.Start()
}
