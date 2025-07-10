package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Token struct {
		AccessSecret string
		AccessExpire int64
	}
	MysqlConf MysqlConfig
}

type MysqlConfig struct {
	Account string
}
