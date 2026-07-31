package generator

import (
	"strings"
	"testing"

	"github.com/emicklei/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toby1991/go-zero/tools/goctl/rpc/parser"
)

func TestGenLogicFunctionStreamingShapes(t *testing.T) {
	tests := []struct {
		name         string
		rpc          *parser.RPC
		wantRequest  bool
		wantValidate bool
		wantReply    bool
		wantStream   bool
		streamType   string
	}{
		{
			name:         "unary",
			rpc:          &parser.RPC{RPC: &proto.RPC{Name: "Unary", RequestType: "Req", ReturnsType: "Resp"}},
			wantRequest:  true,
			wantValidate: true,
			wantReply:    true,
		},
		{
			name: "server-stream",
			rpc: &parser.RPC{RPC: &proto.RPC{
				Name: "Server", RequestType: "Req", ReturnsType: "Resp", StreamsReturns: true,
			}},
			wantRequest:  true,
			wantValidate: true,
			wantStream:   true,
			streamType:   "pb.Demo_ServerServer",
		},
		{
			name: "client-stream",
			rpc: &parser.RPC{RPC: &proto.RPC{
				Name: "Client", RequestType: "Req", ReturnsType: "Resp", StreamsRequest: true,
			}},
			wantStream: true,
			streamType: "pb.Demo_ClientServer",
		},
		{
			name: "bidi-stream",
			rpc: &parser.RPC{RPC: &proto.RPC{
				Name: "Bidi", RequestType: "Req", ReturnsType: "Resp", StreamsRequest: true, StreamsReturns: true,
			}},
			wantStream: true,
			streamType: "pb.Demo_BidiServer",
		},
	}

	g := NewGenerator("gozero", false)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := g.genLogicFunction("Demo", "pb", "", test.name+"Logic", test.rpc,
				map[string]parser.ImportedProto{})
			require.NoError(t, err)

			assert.Equal(t, test.wantRequest, strings.Contains(got, "in *pb.Req"))
			assert.Equal(t, test.wantValidate, strings.Contains(got, "in.Validate()"))
			assert.Equal(t, test.wantReply, strings.Contains(got, "*pb.Resp"))
			if test.wantStream {
				assert.Contains(t, got, test.streamType)
				if test.wantValidate {
					assert.Contains(t, got, "return err")
					assert.NotContains(t, got, "return nil, err")
				} else {
					assert.NotContains(t, got, "in.Validate()")
					assert.NotContains(t, got, "return err")
					assert.Contains(t, strings.ReplaceAll(got, " ", ""), "returnnil")
				}
			} else {
				assert.Contains(t, got, "return nil, err")
			}
			if !test.wantRequest {
				assert.NotContains(t, got, "in ")
			}
		})
	}
}
