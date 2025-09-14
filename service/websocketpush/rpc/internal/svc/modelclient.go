package svc

import (
	"community/conf/databases/xkafka"
	"community/conf/databases/xredis"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ModelClient struct {
	KafkaMessageClient *kq.Pusher
	RedisClient        *redis.Redis
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		KafkaMessageClient: xkafka.GetKafkaClient("messageNode1"), // 获取topic为消息的客户端
		RedisClient:        xredis.GetRedisCommunityClient(),
	}
}
