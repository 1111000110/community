package redis

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	RedisConf RedisConfig
}
type RedisConfig struct {
	Community string
}

var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"
var c Config

func init() {
	conf.MustLoad(configFile, &c)
}

func GetRedisCommunityClient() *redis.Redis {
	return redis.New(c.RedisConf.Community)
}