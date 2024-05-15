# goctl

[English](readme.md) | 简体中文

goctl 使用见文档 https://go-zero.dev/docs/tutorials/cli/overview

```shell
# 安装工具
go install github.com/toby1991/go-zero/tools/goctl@v1.6.5-utils.8

# 创建项目
mkdir user
go mod init user

# 进入项目，生成 userrpc 服务，生成数据库schema
cd user
goctl rpc new userrpc 
go generate ./userrpc/internal/ent
go mod tidy

# 进入 userrpc 服务，生成 logic 代码
cd userrpc
goctl rpc protoc demo.proto --go_out=. --go-grpc_out=. --zrpc_out=.
protoc --validate_out=lang=go:./ *.proto
go mod tidy

# 编译运行
go run userrpc.go
```