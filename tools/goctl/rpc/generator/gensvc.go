package generator

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/collection"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed svc.tpl
var svcTemplate string

// GenSvc generates the servicecontext.go file, which is the resource dependency of a service,
// such as rpc dependency, model dependency, etc.
func (g *Generator) GenSvc(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetSvc()
	svcFilename, err := format.FileNamingFormat(cfg.NamingFormat, "service_context")
	if err != nil {
		return err
	}

	fileName := filepath.Join(dir.Filename, svcFilename+".go")
	text, err := pathx.LoadTemplate(category, svcTemplateFile, svcTemplate)
	if err != nil {
		return err
	}

	imports := collection.NewSet[string]()
	imports.Add(
		fmt.Sprintf(`"%v"`, ctx.GetConfig().Package),
		fmt.Sprintf(`"%v"`, ctx.GetEnt().Package),
		fmt.Sprintf(`"%v"`, pathx.JoinPackages(ctx.GetEnt().Package, "intercept")),
		fmt.Sprintf(`"%v"`, pathx.JoinPackages(ctx.GetEnt().Package, "migrate")),
		fmt.Sprintf(`_ "%v"`, pathx.JoinPackages(ctx.GetEnt().Package, "runtime")),
	)

	return util.With("svc").GoFmt(true).Parse(text).SaveTo(map[string]any{
		"imports": strings.Join(imports.Keys(), pathx.NL),
	}, fileName, false)
}
