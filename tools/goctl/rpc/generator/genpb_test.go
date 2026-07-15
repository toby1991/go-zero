package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toby1991/go-zero/tools/goctl/util/pathx"
)

func TestGeneratorBuildProtocCmdFiltersBundledRuntimeProtos(t *testing.T) {
	const baseCmd = "protoc service.proto --go_out=."

	tests := []struct {
		name       string
		importPath string
		wantCmd    string
	}{
		{name: "errors", importPath: "thirdparty/errors/errors.proto", wantCmd: baseCmd},
		{name: "google annotations", importPath: "thirdparty/google/api/annotations.proto", wantCmd: baseCmd},
		{name: "google http", importPath: "thirdparty/google/api/http.proto", wantCmd: baseCmd},
		{name: "google http body", importPath: "thirdparty/google/api/httpbody.proto", wantCmd: baseCmd},
		{name: "google descriptor", importPath: "thirdparty/google/protobuf/descriptor.proto", wantCmd: baseCmd},
		{name: "google duration", importPath: "thirdparty/google/protobuf/duration.proto", wantCmd: baseCmd},
		{name: "google timestamp", importPath: "thirdparty/google/protobuf/timestamp.proto", wantCmd: baseCmd},
		{name: "validate", importPath: "thirdparty/validate/validate.proto", wantCmd: baseCmd},
		{
			name:       "custom thirdparty proto",
			importPath: "thirdparty/acme/types.proto",
			wantCmd:    baseCmd + " thirdparty/acme/types.proto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			src := writeProtoFixture(t, root, "service.proto", []string{tt.importPath})
			writeProtoFixture(t, root, tt.importPath, nil)

			g := new(Generator)
			got, err := g.buildProtocCmd(&ZRpcContext{
				Src:       src,
				ProtocCmd: baseCmd,
			}, root)

			require.NoError(t, err)
			assert.Equal(t, tt.wantCmd, got)
		})
	}
}

func TestGeneratorBuildProtocCmdPreservesOrdinaryAndTransitiveImports(t *testing.T) {
	const baseCmd = "protoc service.proto --go_out=."

	tests := []struct {
		name         string
		rootImports  []string
		dependencies map[string][]string
		wantCmd      string
	}{
		{
			name:    "no imports",
			wantCmd: baseCmd,
		},
		{
			name:        "ordinary import",
			rootImports: []string{"common/types.proto"},
			dependencies: map[string][]string{
				"common/types.proto": nil,
			},
			wantCmd: baseCmd + " common/types.proto",
		},
		{
			name:        "transitive import",
			rootImports: []string{"common/types.proto"},
			dependencies: map[string][]string{
				"common/types.proto": {"common/base.proto"},
				"common/base.proto":  nil,
			},
			wantCmd: baseCmd + " common/types.proto common/base.proto",
		},
		{
			name: "mixed ordinary and bundled transitive imports",
			rootImports: []string{
				"common/types.proto",
				"thirdparty/validate/validate.proto",
			},
			dependencies: map[string][]string{
				"common/types.proto": nil,
				"thirdparty/validate/validate.proto": {
					"thirdparty/google/protobuf/descriptor.proto",
					"thirdparty/google/protobuf/duration.proto",
					"thirdparty/google/protobuf/timestamp.proto",
				},
				"thirdparty/google/protobuf/descriptor.proto": nil,
				"thirdparty/google/protobuf/duration.proto":   nil,
				"thirdparty/google/protobuf/timestamp.proto":  nil,
			},
			wantCmd: baseCmd + " common/types.proto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			src := writeProtoFixture(t, root, "service.proto", tt.rootImports)
			for name, imports := range tt.dependencies {
				writeProtoFixture(t, root, name, imports)
			}

			g := new(Generator)
			got, err := g.buildProtocCmd(&ZRpcContext{
				Src:       src,
				ProtocCmd: baseCmd,
			}, root)

			require.NoError(t, err)
			assert.Equal(t, tt.wantCmd, got)
		})
	}
}

func writeProtoFixture(t *testing.T, root, name string, imports []string) string {
	t.Helper()

	filename := filepath.Join(root, filepath.FromSlash(name))
	require.NoError(t, os.MkdirAll(filepath.Dir(filename), 0o755))

	var source strings.Builder
	source.WriteString("syntax = \"proto3\";\npackage fixture;\n")
	for _, imported := range imports {
		source.WriteString("import \"")
		source.WriteString(imported)
		source.WriteString("\";\n")
	}
	require.NoError(t, os.WriteFile(filename, []byte(source.String()), 0o600))

	return filename
}

func Test_findPbFile(t *testing.T) {
	dir := t.TempDir()
	protoFile := filepath.Join(dir, "greet.proto")
	err := os.WriteFile(protoFile, []byte(`
syntax = "proto3";

package greet;
option go_package="./greet";

message Req{}
message Resp{}
service Greeter {
  rpc greet(Req) returns (Resp);
}
`), 0o666)
	if err != nil {
		t.Log(err)
		return
	}
	t.Run("", func(t *testing.T) {
		output := t.TempDir()
		grpc := filepath.Join(output, "grpc")
		err := pathx.MkdirIfNotExist(grpc)
		if err != nil {
			t.Log(err)
			return
		}
		cmd := exec.Command("protoc", "-I="+filepath.Dir(protoFile), "--go_out="+output, "--go-grpc_out="+grpc, filepath.Base(protoFile))
		cmd.Dir = output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Log(err)
			return
		}
		pbDir, err := findPbFile(output, protoFile, false)
		assert.Nil(t, err)
		pbGo := filepath.Join(pbDir, "greet.pb.go")
		assert.True(t, pathx.FileExists(pbGo))

		grpcDir, err := findPbFile(output, protoFile, true)
		assert.Nil(t, err)
		grpcGo := filepath.Join(grpcDir, "greet_grpc.pb.go")
		assert.True(t, pathx.FileExists(grpcGo))
	})

	t.Run("", func(t *testing.T) {
		output := t.TempDir()
		redirect := filepath.Join(output, "pb")
		grpc := filepath.Join(output, "grpc")
		err := pathx.MkdirIfNotExist(grpc)
		if err != nil {
			t.Log(err)
			return
		}
		cmd := exec.Command("protoc", "-I="+filepath.Dir(protoFile), "--go_out="+output,
			"--go-grpc_out="+grpc, filepath.Base(protoFile), "--go_opt=M"+filepath.Base(protoFile)+"="+redirect)
		cmd.Dir = output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Log(err)
			return
		}
		pbDir, err := findPbFile(output, protoFile, false)
		assert.Nil(t, err)
		pbGo := filepath.Join(pbDir, "greet.pb.go")
		assert.True(t, pathx.FileExists(pbGo))

		grpcDir, err := findPbFile(output, protoFile, true)
		assert.Nil(t, err)
		grpcGo := filepath.Join(grpcDir, "greet_grpc.pb.go")
		assert.True(t, pathx.FileExists(grpcGo))
	})

	t.Run("", func(t *testing.T) {
		output := t.TempDir()
		pbeRedirect := filepath.Join(output, "redirect")
		grpc := filepath.Join(output, "grpc")
		grpcRedirect := filepath.Join(grpc, "redirect")
		err := pathx.MkdirIfNotExist(grpc)
		if err != nil {
			t.Log(err)
			return
		}
		cmd := exec.Command("protoc", "-I="+filepath.Dir(protoFile), "--go_out="+output,
			"--go-grpc_out="+grpc, filepath.Base(protoFile), "--go_opt=M"+filepath.Base(protoFile)+"="+pbeRedirect,
			"--go-grpc_opt=M"+filepath.Base(protoFile)+"="+grpcRedirect)
		cmd.Dir = output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Log(err)
			return
		}
		pbDir, err := findPbFile(output, protoFile, false)
		assert.Nil(t, err)
		pbGo := filepath.Join(pbDir, "greet.pb.go")
		assert.True(t, pathx.FileExists(pbGo))

		grpcDir, err := findPbFile(output, protoFile, true)
		assert.Nil(t, err)
		grpcGo := filepath.Join(grpcDir, "greet_grpc.pb.go")
		assert.True(t, pathx.FileExists(grpcGo))
	})

	t.Run("", func(t *testing.T) {
		output := t.TempDir()
		pbeRedirect := filepath.Join(output, "redirect")
		grpc := filepath.Join(output, "grpc")
		grpcRedirect := filepath.Join(grpc, "redirect")
		err := pathx.MkdirIfNotExist(grpc)
		if err != nil {
			t.Log(err)
			return
		}
		cmd := exec.Command("protoc", "-I="+filepath.Dir(protoFile), "--go_out="+output,
			"--go-grpc_out="+grpc, filepath.Base(protoFile), "--go_opt=M"+filepath.Base(protoFile)+"="+pbeRedirect,
			"--go-grpc_opt=M"+filepath.Base(protoFile)+"="+grpcRedirect, "--go_opt=paths=import", "--go-grpc_opt=paths=source_relative")
		cmd.Dir = output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Log(err)
			return
		}
		pbDir, err := findPbFile(output, protoFile, false)
		assert.Nil(t, err)
		pbGo := filepath.Join(pbDir, "greet.pb.go")
		assert.True(t, pathx.FileExists(pbGo))

		grpcDir, err := findPbFile(output, protoFile, true)
		assert.Nil(t, err)
		grpcGo := filepath.Join(grpcDir, "greet_grpc.pb.go")
		assert.True(t, pathx.FileExists(grpcGo))
	})

	t.Run("", func(t *testing.T) {
		output := t.TempDir()
		pbeRedirect := filepath.Join(output, "redirect")
		grpc := filepath.Join(output, "grpc")
		grpcRedirect := filepath.Join(grpc, "redirect")
		err := pathx.MkdirIfNotExist(grpc)
		if err != nil {
			t.Log(err)
			return
		}
		err = pathx.MkdirIfNotExist(pbeRedirect)
		if err != nil {
			t.Log(err)
			return
		}
		err = pathx.MkdirIfNotExist(grpcRedirect)
		if err != nil {
			t.Log(err)
			return
		}
		cmd := exec.Command("protoc", "-I="+filepath.Dir(protoFile), "--go_out="+output,
			"--go-grpc_out="+grpc, filepath.Base(protoFile), "--go_opt=M"+filepath.Base(protoFile)+"="+pbeRedirect,
			"--go-grpc_opt=M"+filepath.Base(protoFile)+"="+grpcRedirect, "--go_opt=paths=import", "--go-grpc_opt=paths=source_relative",
			"--go_out="+pbeRedirect, "--go-grpc_out="+grpcRedirect)
		cmd.Dir = output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Log(err)
			return
		}
		pbDir, err := findPbFile(output, protoFile, false)
		assert.Nil(t, err)
		pbGo := filepath.Join(pbDir, "greet.pb.go")
		assert.True(t, pathx.FileExists(pbGo))

		grpcDir, err := findPbFile(output, protoFile, true)
		assert.Nil(t, err)
		grpcGo := filepath.Join(grpcDir, "greet_grpc.pb.go")
		assert.True(t, pathx.FileExists(grpcGo))
	})
}
