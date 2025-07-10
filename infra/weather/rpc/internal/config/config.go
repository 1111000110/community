package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MongoConf MongoConfig
}
type MongoConfig struct {
	SubscriptionData struct {
		Url        string
		Db         string
		Collection string
	}
}
