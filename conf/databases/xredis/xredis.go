package xredis

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	RedisConf redis.RedisConf
}

var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"
var c Config

func init() {
	conf.MustLoad(configFile, &c)
}

func GetRedisCommunityClient() *redis.Redis {
	return redis.MustNewRedis(c.RedisConf)
}
