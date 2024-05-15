package main

import (
	"flag"
	"fmt"
	"github.com/toby1991/go-zero-utils/pprof"
	"github.com/zeromicro/go-zero/core/service"

	{{.importPackages}}
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

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

        // api server
	server := rest.MustNewServer(c.RestConf)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	svcGroup.Add(server)

        // start server
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	svcGroup.Start()
}
