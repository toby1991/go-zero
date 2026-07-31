package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/emicklei/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toby1991/go-zero/core/collection"
	conf "github.com/toby1991/go-zero/tools/goctl/config"
	"github.com/toby1991/go-zero/tools/goctl/rpc/parser"
	"github.com/toby1991/go-zero/tools/goctl/util/stringx"
)

// mockDirContext is a minimal DirContext for unit-testing genCallGroup.
type mockDirContext struct {
	callDir   Dir
	configDir Dir
	logicDir  Dir
	pbDir     Dir
	protoGo   Dir
	svcDir    Dir
}

func (m *mockDirContext) GetCall() Dir                   { return m.callDir }
func (m *mockDirContext) GetEtc() Dir                    { return Dir{} }
func (m *mockDirContext) GetEnt() Dir                    { return Dir{} }
func (m *mockDirContext) GetThirdPartyPb() Dir           { return Dir{} }
func (m *mockDirContext) GetInternal() Dir               { return Dir{} }
func (m *mockDirContext) GetConfig() Dir                 { return m.configDir }
func (m *mockDirContext) GetLogic() Dir                  { return m.logicDir }
func (m *mockDirContext) GetServer() Dir                 { return Dir{} }
func (m *mockDirContext) GetSvc() Dir                    { return m.svcDir }
func (m *mockDirContext) GetPb() Dir                     { return m.pbDir }
func (m *mockDirContext) GetProtoGo() Dir                { return m.protoGo }
func (m *mockDirContext) GetMain() Dir                   { return Dir{} }
func (m *mockDirContext) GetServiceName() stringx.String { return stringx.From("test") }
func (m *mockDirContext) SetPbDir(pbDir, grpcDir string) {}

// TestGenCallGroup_OnlyUsedTypesAliased verifies that in multi-service mode each
// generated client file contains type aliases only for the message types actually
// used by that service's RPCs (fix for issue #5481).
func TestGenCallGroup_OnlyUsedTypesAliased(t *testing.T) {
	tmpDir := t.TempDir()
	callBase := filepath.Join(tmpDir, "call")
	pbBase := filepath.Join(tmpDir, "pb")

	// Pre-create subdirs that genCallGroup will write into.
	require.NoError(t, os.MkdirAll(filepath.Join(callBase, "servicea"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(callBase, "serviceb"), 0755))
	require.NoError(t, os.MkdirAll(pbBase, 0755))

	mctx := &mockDirContext{
		callDir: Dir{
			Filename: callBase,
			Package:  "example.com/multitest/call",
			Base:     "call",
			GetChildPackage: func(childPath string) (string, error) {
				// Return a package path whose Base() is the lowercase service name.
				return filepath.Join(callBase, strings.ToLower(childPath)), nil
			},
		},
		pbDir: Dir{
			Filename: pbBase,
			Package:  "example.com/multitest/pb",
			Base:     "pb",
		},
		protoGo: Dir{
			// Must differ from "servicea"/"serviceb" so isCallPkgSameToPbPkg stays false
			// and alias generation is triggered.
			Filename: pbBase,
			Package:  "example.com/multitest/pb",
			Base:     "pb",
		},
		logicDir: Dir{
			Package: "example.com/multitest/internal/logic",
			GetChildPackage: func(childPath string) (string, error) {
				return "example.com/multitest/internal/logic/" + strings.ToLower(childPath), nil
			},
		},
	}

	// Proto with two services that use completely disjoint message types.
	protoData := parser.Proto{
		Name:      "multi.proto",
		PbPackage: "pb",
		Message: []parser.Message{
			{Message: &proto.Message{Name: "AReq"}},
			{Message: &proto.Message{Name: "AResp"}},
			{Message: &proto.Message{Name: "BReq"}},
			{Message: &proto.Message{Name: "BResp"}},
		},
		Service: parser.Services{
			{
				Service: &proto.Service{Name: "ServiceA"},
				RPC: []*parser.RPC{
					{RPC: &proto.RPC{Name: "DoA", RequestType: "AReq", ReturnsType: "AResp"}},
				},
			},
			{
				Service: &proto.Service{Name: "ServiceB"},
				RPC: []*parser.RPC{
					{RPC: &proto.RPC{Name: "DoB", RequestType: "BReq", ReturnsType: "BResp"}},
				},
			},
		},
	}

	cfg, err := conf.NewConfig("")
	require.NoError(t, err)

	g := NewGenerator("gozero", false)
	require.NoError(t, g.genCallGroup(mctx, protoData, cfg))

	// servicea/servicea.go — aliases for AReq/AResp only
	aContent, err := os.ReadFile(filepath.Join(callBase, "servicea", "servicea.go"))
	require.NoError(t, err)
	aFile := normalizeWS(string(aContent))

	assert.Contains(t, aFile, "AReq = pb.AReq", "ServiceA file should alias AReq")
	assert.Contains(t, aFile, "AResp = pb.AResp", "ServiceA file should alias AResp")
	assert.NotContains(t, aFile, "BReq = pb.BReq", "ServiceA file must not alias BReq")
	assert.NotContains(t, aFile, "BResp = pb.BResp", "ServiceA file must not alias BResp")

	// serviceb/serviceb.go — aliases for BReq/BResp only
	bContent, err := os.ReadFile(filepath.Join(callBase, "serviceb", "serviceb.go"))
	require.NoError(t, err)
	bFile := normalizeWS(string(bContent))

	assert.Contains(t, bFile, "BReq = pb.BReq", "ServiceB file should alias BReq")
	assert.Contains(t, bFile, "BResp = pb.BResp", "ServiceB file should alias BResp")
	assert.NotContains(t, bFile, "AReq = pb.AReq", "ServiceB file must not alias AReq")
	assert.NotContains(t, bFile, "AResp = pb.AResp", "ServiceB file must not alias AResp")
	assert.Contains(t, aFile, `logic "example.com/multitest/internal/logic/servicea"`,
		"ServiceA file should use its child Logic package")
}

// normalizeWS replaces runs of whitespace with a single space.
func normalizeWS(s string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(s, "\n", " \n ")), " ")
}

func TestGetServiceRPCModes(t *testing.T) {
	tests := []struct {
		name          string
		service       parser.Service
		wantUnary     bool
		wantStreaming bool
	}{
		{
			name: "unary-only",
			service: parser.Service{RPC: []*parser.RPC{
				{RPC: &proto.RPC{Name: "Unary", RequestType: "Req", ReturnsType: "Resp"}},
			}},
			wantUnary: true,
		},
		{
			name: "stream-only",
			service: parser.Service{RPC: []*parser.RPC{
				{RPC: &proto.RPC{Name: "Server", RequestType: "Req", ReturnsType: "Resp", StreamsReturns: true}},
				{RPC: &proto.RPC{Name: "Client", RequestType: "Req", ReturnsType: "Resp", StreamsRequest: true}},
			}},
			wantStreaming: true,
		},
		{
			name: "mixed",
			service: parser.Service{RPC: []*parser.RPC{
				{RPC: &proto.RPC{Name: "Unary", RequestType: "Req", ReturnsType: "Resp"}},
				{RPC: &proto.RPC{Name: "Bidi", RequestType: "Req", ReturnsType: "Resp", StreamsRequest: true, StreamsReturns: true}},
			}},
			wantUnary:     true,
			wantStreaming: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hasUnary, hasStreaming := getServiceRPCModes(test.service)
			assert.Equal(t, test.wantUnary, hasUnary)
			assert.Equal(t, test.wantStreaming, hasStreaming)
		})
	}
}

func TestGenDirectFunction_StreamingStubsUseQualifiedClientTypes(t *testing.T) {
	service := parser.Service{
		Service: &proto.Service{Name: "Mixed"},
		RPC: []*parser.RPC{
			{RPC: &proto.RPC{Name: "Unary", RequestType: "Req", ReturnsType: "Resp"}},
			{RPC: &proto.RPC{Name: "Server", RequestType: "Req", ReturnsType: "Resp", StreamsReturns: true}},
			{RPC: &proto.RPC{Name: "Client", RequestType: "Req", ReturnsType: "Resp", StreamsRequest: true}},
			{RPC: &proto.RPC{Name: "Bidi", RequestType: "Req", ReturnsType: "Resp", StreamsRequest: true, StreamsReturns: true}},
		},
	}
	g := NewGenerator("gozero", false)

	directFunctions, err := g.genDirectFunction("pb", "Mixed", service, false)
	require.NoError(t, err)
	direct := strings.Join(directFunctions, "\n")
	statusReturn := "return nil, status.Error(codes.Unimplemented, \"direct mode does not support streaming RPCs\")"
	assert.Equal(t, 3, strings.Count(direct, statusReturn))
	assert.Contains(t, direct, "pb.Mixed_ServerClient")
	assert.Contains(t, direct, "pb.Mixed_ClientClient")
	assert.Contains(t, direct, "pb.Mixed_BidiClient")
	assert.Contains(t, direct, "logic.NewUnaryLogic(ctx, l.svcCtx)")
	assert.NotContains(t, direct, "logic.NewServerLogic")
	assert.NotContains(t, direct, "logic.NewClientLogic")
	assert.NotContains(t, direct, "logic.NewBidiLogic")

	remoteFunctions, err := g.genFunction("pb", "example.com/fixture/pb", "Mixed", service, false,
		map[string]parser.ImportedProto{}, collection.NewSet[string](), collection.NewSet[string]())
	require.NoError(t, err)
	remote := strings.Join(remoteFunctions, "\n")
	for _, clientType := range []string{
		"pb.Mixed_ServerClient",
		"pb.Mixed_ClientClient",
		"pb.Mixed_BidiClient",
	} {
		assert.Contains(t, remote, clientType)
	}
}

func TestGenCallCompatibility_ConditionalImports(t *testing.T) {
	tests := []struct {
		name             string
		rpcs             []*parser.RPC
		wantLogicImport  bool
		wantStreamImport bool
	}{
		{
			name:            "unary-only",
			rpcs:            []*parser.RPC{{RPC: &proto.RPC{Name: "Unary", RequestType: "Req", ReturnsType: "Resp"}}},
			wantLogicImport: true,
		},
		{
			name:             "stream-only",
			rpcs:             []*parser.RPC{{RPC: &proto.RPC{Name: "Server", RequestType: "Req", ReturnsType: "Resp", StreamsReturns: true}}},
			wantStreamImport: true,
		},
		{
			name: "mixed",
			rpcs: []*parser.RPC{
				{RPC: &proto.RPC{Name: "Unary", RequestType: "Req", ReturnsType: "Resp"}},
				{RPC: &proto.RPC{Name: "Server", RequestType: "Req", ReturnsType: "Resp", StreamsReturns: true}},
			},
			wantLogicImport:  true,
			wantStreamImport: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content := generateCompatibilityCall(t, test.rpcs)
			logicImport := `"example.com/fixture/internal/logic"`
			assert.Equal(t, test.wantLogicImport, strings.Contains(content, logicImport))
			assert.Equal(t, test.wantStreamImport, strings.Contains(content, `"google.golang.org/grpc/codes"`))
			assert.Equal(t, test.wantStreamImport, strings.Contains(content, `"google.golang.org/grpc/status"`))
		})
	}
}

func generateCompatibilityCall(t *testing.T, rpcs []*parser.RPC) string {
	t.Helper()
	tmpDir := t.TempDir()
	callDir := filepath.Join(tmpDir, "client")
	require.NoError(t, os.MkdirAll(callDir, 0755))

	ctx := &mockDirContext{
		callDir: Dir{
			Filename: callDir,
			Package:  "example.com/fixture/client",
			Base:     "client",
		},
		configDir: Dir{Package: "example.com/fixture/internal/config"},
		logicDir:  Dir{Package: "example.com/fixture/internal/logic"},
		pbDir:     Dir{Filename: filepath.Join(tmpDir, "pb"), Package: "example.com/fixture/pb", Base: "pb"},
		protoGo:   Dir{Filename: filepath.Join(tmpDir, "grpc"), Package: "example.com/fixture/pb", Base: "pb"},
		svcDir:    Dir{Package: "example.com/fixture/internal/svc"},
	}

	protoData := parser.Proto{
		Name:      "fixture.proto",
		PbPackage: "pb",
		GoPackage: "example.com/fixture/pb",
		Message: []parser.Message{
			{Message: &proto.Message{Name: "Req"}},
			{Message: &proto.Message{Name: "Resp"}},
		},
		Service: parser.Services{{
			Service: &proto.Service{Name: "Fixture"},
			RPC:     rpcs,
		}},
	}

	g := NewGenerator("gozero", false)
	cfg, err := conf.NewConfig("")
	require.NoError(t, err)
	require.NoError(t, g.genCallInCompatibility(ctx, protoData, cfg))

	entries, err := os.ReadDir(callDir)
	require.NoError(t, err)
	var content strings.Builder
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".go" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(callDir, entry.Name()))
		require.NoError(t, err)
		content.Write(data)
	}
	return content.String()
}
