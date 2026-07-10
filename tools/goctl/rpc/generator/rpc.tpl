syntax = "proto3";

import "thirdparty/validate/validate.proto";

package {{.package}};
option go_package="./{{.package}}";

// go install github.com/toby1991/go-zero/tools/goctl@latest
// goctl rpc new demo && go generate ./demo/internal/ent && go mod tidy
// goctl rpc protoc demo.proto --go_out=. --go-grpc_out=. --zrpc_out=.
// protoc --validate_out=lang=go:./ *.proto
// go mod tidy
// go generate ./internal/ent
// go mod tidy
// go run main.go

// pagination
message PaginationReq {
  uint64 page = 1;
  uint64 pageSize = 2 [(validate.rules).uint64 = {lt: 50, gt: 0}];
}
message PaginationResp {
  uint64 page = 1;
  uint64 pageSize = 2;
  uint64 total = 3;
}

message OffsetReq {
  uint64 offset = 1;
  uint32 limit = 2 [(validate.rules).uint32.lt = 30];
}

// dto
message UserDto {
    uint64 id = 1;
    string username = 2;
    string email = 3;
}

// get user list
message GetUserListReq {
  PaginationReq paginationQuery = 1;
}
message GetUserListResp {
  repeated UserDto userList = 1;
  PaginationResp pagination = 2;
}


service {{.serviceName}} {
  rpc GetUserList(GetUserListReq) returns (GetUserListResp);
}
