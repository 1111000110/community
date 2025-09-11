package svc

import (
	"community/conf/databases/kafka"
	"community/service/websocket/api/internal/config"
	"community/service/websocket/api/internal/hub"
	"community/service/websocket/api/internal/middleware"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config        config.Config
	Middleware    rest.Middleware
	MessageHub    *hub.Hub
	Model         *ModelClient
	KafkaConsumer []service.Service
}

func NewServiceContext(c config.Config) *ServiceContext {
	messageHub := hub.NewHub()
	return &ServiceContext{
		Config:        c,
		Middleware:    middleware.NewMiddleware().Handle,
		MessageHub:    messageHub,
		Model:         DefaultModelClient(),
		KafkaConsumer: kafka.DefaultConsumer("message", "message", messageHub),
	}
}
