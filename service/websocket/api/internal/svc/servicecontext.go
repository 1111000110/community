package svc

import (
	"community.com/service/websocket/api/internal/config"
	"community.com/service/websocket/api/internal/hub"
	"community.com/service/websocket/api/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config     config.Config
	Middleware rest.Middleware
	MessageHub *hub.Hub
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		Middleware: middleware.NewMiddleware().Handle,
		MessageHub: hub.NewHub(),
	}
}
