package generator

import (
	_ "embed"
	"fmt"
	"github.com/zeromicro/go-zero/core/collection"
	"path/filepath"
	"strings"

	conf "github.com/toby1991/go-zero/tools/goctl/config"
	"github.com/toby1991/go-zero/tools/goctl/rpc/parser"
	"github.com/toby1991/go-zero/tools/goctl/util"
	"github.com/toby1991/go-zero/tools/goctl/util/format"
	"github.com/toby1991/go-zero/tools/goctl/util/pathx"
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

	configImport := fmt.Sprintf(`"%v"`, ctx.GetConfig().Package)
	entImport := fmt.Sprintf(`"%v"`, ctx.GetEnt().Package)
	entMigrateImport := fmt.Sprintf(`"%v"`, pathx.JoinPackages(ctx.GetEnt().Package, "migrate"))
	entInterceptImport := fmt.Sprintf(`"_ %v"`, pathx.JoinPackages(ctx.GetEnt().Package, "intercept"))
	entRuntimeImport := fmt.Sprintf(`"_ %v"`, pathx.JoinPackages(ctx.GetEnt().Package, "runtime"))

	imports := collection.NewSet()
	imports.AddStr(configImport, entImport, entMigrateImport, entInterceptImport, entRuntimeImport)

	return util.With("svc").GoFmt(true).Parse(text).SaveTo(map[string]any{
		//"imports": fmt.Sprintf(`"%v"`, ctx.GetConfig().Package),
		"imports": strings.Join(imports.KeysStr(), pathx.NL),
	}, fileName, false)
}
