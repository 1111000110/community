package svc

import (
	"community.com/conf/databases/kafka"
	"github.com/zeromicro/go-queue/kq"
)

type ModelClient struct {
	KafkaMessageClient *kq.Pusher
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		KafkaMessageClient: kafka.GetKafkaClient("message"), // 获取topic为消息的客户端
	}
}
