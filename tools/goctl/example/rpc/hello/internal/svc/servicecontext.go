package svc

import "github.com/toby1991/go-zero/tools/goctl/example/rpc/hello/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
