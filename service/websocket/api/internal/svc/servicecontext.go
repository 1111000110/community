package svc

import (
	"community/conf/databases/xkafka"
	"community/service/websocket/api/internal/config"
	"community/service/websocket/api/internal/hub"

	"github.com/zeromicro/go-zero/core/service"
)

type ServiceContext struct {
	Config        config.Config
	MessageHub    *hub.Hub
	Model         *ModelClient
	KafkaConsumer []service.Service
}

func NewServiceContext(c config.Config) *ServiceContext {
	messageHub := hub.NewHub()
	return &ServiceContext{
		Config:        c,
		MessageHub:    messageHub,
		Model:         DefaultModelClient(),
		KafkaConsumer: xkafka.DefaultConsumer("message", "message", messageHub),
	}
}
