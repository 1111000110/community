package xkafka

import (
	"context"
	"time"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	KafkaConf KafkaConf
}
type KafkaConf struct {
	Community []string
}

var configFile = "/Users/zhangxuan/Data/work/xuan/community/conf/conf.yaml"
var c Config

func init() {
	conf.MustLoad(configFile, &c)
}

// 生产方

func GetKafkaClient(topic string) *kq.Pusher {
	return kq.NewPusher(c.KafkaConf.Community, topic, kq.WithFlushInterval(10*time.Millisecond)) // 10毫秒上报，提高实时性
}

// 消费方

type Consumer interface { // 消费方实现此处理接口
	Consume(ctx context.Context, key, value string) error
}

func DefaultConsumer(group string, topic string, consumer Consumer) []service.Service {
	return GetKafkaConsumer(kq.KqConf{
		Brokers:       c.KafkaConf.Community,
		Group:         group,
		Topic:         topic,
		Offset:        "last",   // 从最新位置开始消费
		Conns:         1,        // 连接数
		Consumers:     8,        // 消费者数量
		Processors:    8,        // 处理器数量
		MinBytes:      1,        // 最小字节数，立即处理消息
		MaxBytes:      10485760, // 最大字节数
		ForceCommit:   true,     // 强制提交
		CommitInOrder: false,    // 不按顺序提交
	}, consumer)
}

func GetKafkaConsumer(kqConf kq.KqConf, consumer Consumer) []service.Service {
	return []service.Service{
		kq.MustNewQueue(kqConf, consumer),
	}
}
