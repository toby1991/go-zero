{{.head}}

package {{.filePackage}}

import (
	"context"

	{{.pbPackage}}
	{{if ne .pbPackage .protoGoPackage}}{{.protoGoPackage}}{{end}}
	{{.extraImports}}
	{{.internalLogicPackage}}
	{{.internalSvcPackage}}
	{{.internalConfigPackage}}

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	{{.alias}}

	{{.serviceName}} interface {
		{{.interface}}
	}

	default{{.serviceName}} struct {
		cli zrpc.Client
	}

	direct{{.serviceName}} struct {
		svcCtx *svc.ServiceContext
	}
)

func New{{.serviceName}}(cli zrpc.Client) {{.serviceName}} {
	return &default{{.serviceName}}{
		cli: cli,
	}
}

func NewDirect{{.serviceName}}(configPath string) {{.serviceName}} {
	var c config.Config
	conf.MustLoad(configPath, &c)
	return &direct{{.serviceName}}{
		svcCtx: svc.NewServiceContext(c),
	}
}

{{.functions}}

{{.directFunctions}}
