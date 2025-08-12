{{.head}}

package {{.filePackage}}

import (
	"context"

	{{.pbPackage}}
	{{if ne .pbPackage .protoGoPackage}}{{.protoGoPackage}}{{end}}
	{{.internalLogicPackage}}
	{{.internalSvcPackage}}

	"github.com/toby1991/go-zero/zrpc"
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

func NewDirect{{.serviceName}}(svcCtx *svc.ServiceContext) {{.serviceName}} {
	return &direct{{.serviceName}}{
		svcCtx: svcCtx,
	}
}

{{.functions}}

{{.directFunctions}}
