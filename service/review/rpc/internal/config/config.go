package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		Review []struct {
			Url        string
			DB         string
			Collection string
		}
	}
	RedisConf redis.RedisConf
}
