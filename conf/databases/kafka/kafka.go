package kafka

import (
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	kafkaConfig KafkaConfig
}
type KafkaConfig struct {
	Community []string
}

var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"
var c Config

func init() {
	conf.MustLoad(configFile, &c)
}

func GetKafkaClient(topic string) *kq.Pusher {
	return kq.NewPusher(c.kafkaConfig.Community, topic)
}

type Consumer interface { // 消费方实现此处理接口
	Consume(ctx context.Context, key, value string) error
}

func DefaultConsumer(group string, topic string, consumer Consumer) []service.Service {
	return GetKafkaConsumer(kq.KqConf{
		Brokers: c.kafkaConfig.Community,
		Group:   group,
		Topic:   topic,
	}, consumer)
}

func GetKafkaConsumer(kqConf kq.KqConf, consumer Consumer) []service.Service {
	return []service.Service{
		kq.MustNewQueue(kqConf, consumer),
	}
}
