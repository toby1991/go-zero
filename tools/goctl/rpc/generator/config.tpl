package config

import (
    "github.com/zeromicro/go-zero-utils/db"
    "github.com/zeromicro/go-zero-utils/bizredis"
	"github.com/zeromicro/go-zero-utils/bizmemory"
    "github.com/zeromicro/go-zero-utils/queue/faktory"
    "github.com/zeromicro/go-zero-utils/queue/nsq"
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
