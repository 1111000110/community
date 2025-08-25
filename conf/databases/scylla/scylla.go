package scylla

import (
	"github.com/zeromicro/go-zero/core/conf"
)

type Config struct {
	ScyllaConf ScyllaConfig
}
type ScyllaConfig struct {
	Community string
	Keyspace string
}

type ScyllaClient struct {
	Hosts    []string
	Keyspace string
	Config   ScyllaConfig
}

var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"
var c Config

func init() {
	conf.MustLoad(configFile, &c)
}

func GetScyllaCommunityClient() *ScyllaClient {
	return &ScyllaClient{
		Hosts:    []string{c.ScyllaConf.Community},
		Keyspace: c.ScyllaConf.Keyspace,
		Config:   c.ScyllaConf,
	}
}