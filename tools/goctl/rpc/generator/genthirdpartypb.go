package generator

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"
	conf "github.com/toby1991/go-zero/tools/goctl/config"
	"github.com/toby1991/go-zero/tools/goctl/rpc/parser"
	"github.com/toby1991/go-zero/tools/goctl/util/pathx"
	"io"
	"os"
	"strings"
)

//go:embed thirdparty.tar.gz
var thirdPartyPbTemplate string

// GenEnt generates the yaml configuration file of the rpc service,
// including host, port monitoring configuration items and etcd configuration
func (g *Generator) GenThirdPartyPb(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetThirdPartyPb()

	binary, err := pathx.LoadTemplate(category, thirdPartyPbTemplateFile, thirdPartyPbTemplate)
	if err != nil {
		return err
	}

	return g.extractTarGz(bytes.NewReader([]byte(binary)), strings.TrimSuffix(dir.Filename, "thirdparty")+"/")

}
func (g *Generator) extractTarGz(gzipStream io.Reader, destDir string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)
	var header *tar.Header
	for header, err = tarReader.Next(); err == nil; header, err = tarReader.Next() {
		switch header.Typeflag {
		case tar.TypeDir:
			if err := pathx.MkdirIfNotExist(destDir + header.Name); err != nil {
				return fmt.Errorf("ExtractTarGz: Mkdir() failed: %w", err)
			}
		case tar.TypeReg:
			filename := destDir + header.Name
			_, err := os.Stat(filename)
			if !os.IsNotExist(err) {
				break
			}
			outFile, err := pathx.CreateIfNotExist(filename)
			if err != nil {
				return fmt.Errorf("ExtractTarGz: Create() failed: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				// outFile.Close error omitted as Copy error is more interesting at this point
				outFile.Close()
				return fmt.Errorf("ExtractTarGz: Copy() failed: %w", err)
			}
			if err := outFile.Close(); err != nil {
				return fmt.Errorf("ExtractTarGz: Close() failed: %w", err)
			}
		default:
			return fmt.Errorf("ExtractTarGz: uknown type: %b in %s", header.Typeflag, header.Name)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("ExtractTarGz: Next() failed: %w", err)
	}
	return nil
}
