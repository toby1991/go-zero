package generator

import (
	_ "embed"
	"fmt"
	conf "github.com/toby1991/go-zero/tools/goctl/config"
	"github.com/toby1991/go-zero/tools/goctl/rpc/parser"
	"github.com/toby1991/go-zero/tools/goctl/util"
	"github.com/toby1991/go-zero/tools/goctl/util/format"
	"github.com/toby1991/go-zero/tools/goctl/util/pathx"
	"path/filepath"
)

//go:embed ent_generate.tpl
var entGenerateTemplate string

//go:embed ent_schema_user.tpl
var entSchemaUserTemplate string

// GenEnt generates the yaml configuration file of the rpc service,
// including host, port monitoring configuration items and etcd configuration
func (g *Generator) GenEnt(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	if err := g.genEntGenerate(ctx, proto, cfg); err != nil {
		return err
	}
	if err := g.genEntSchemaUser(ctx, proto, cfg); err != nil {
		return err
	}

	return nil
}

func (g *Generator) genEntGenerate(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetEnt()
	entFilename, err := format.FileNamingFormat(cfg.NamingFormat, "generate")
	if err != nil {
		return err
	}
	fileName := filepath.Join(dir.Filename, fmt.Sprintf("%s.go", entFilename))

	text, err := pathx.LoadTemplate(category, entGenerateTemplateFile, entGenerateTemplate)
	if err != nil {
		return err
	}

	return util.With("ent_generate").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{}, fileName, false)
}
func (g *Generator) genEntSchemaUser(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetEnt()

	fname := dir.Filename + "/schema/"
	if err := pathx.MkdirIfNotExist(fname); err != nil {
		return err
	}

	entFilename, err := format.FileNamingFormat(cfg.NamingFormat, "user")
	if err != nil {
		return err
	}
	fileName := filepath.Join(fname, fmt.Sprintf("%s.go", entFilename))

	text, err := pathx.LoadTemplate(category, entSchemaUserTemplateFile, entSchemaUserTemplate)
	if err != nil {
		return err
	}

	return util.With("ent_schema_user").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{}, fileName, false)
}
