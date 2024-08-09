package config

import (
    "github.com/toby1991/go-zero-utils/db"
    "github.com/toby1991/go-zero-utils/bizredis"
	"github.com/toby1991/go-zero-utils/bizmemory"
    "github.com/toby1991/go-zero-utils/faktory"
    "github.com/toby1991/go-zero-utils/nsq"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

    Db       db.DbConf             `json:",optional"`
	BizRedis bizredis.BizRedisConf `json:",optional"`
	BizMemory bizmemory.BizMemoryConf `json:",optional"`
	Faktory  faktory.FaktoryConf   `json:",optional"`
	Nsq      nsq.NsqConf           `json:",optional"`
}
