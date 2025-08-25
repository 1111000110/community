package kafka

import (
	"github.com/zeromicro/go-zero/core/conf"
)

type Config struct {
	KafkaConf KafkaConfig
}
type KafkaConfig struct {
	Community string
}

type KafkaClient struct {
	Brokers []string
	Config  KafkaConfig
}

var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"
var c Config

func init() {
	conf.MustLoad(configFile, &c)
}

func GetKafkaCommunityClient() *KafkaClient {
	return &KafkaClient{
		Brokers: []string{c.KafkaConf.Community},
		Config:  c.KafkaConf,
	}
}