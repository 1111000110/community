package svc

import "community/service/websocketpush/rpc/internal/config"

type ServiceContext struct {
	Config config.Config
	Model  *ModelClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	model := DefaultModelClient()
	return &ServiceContext{
		Config: c,
		Model:  model,
	}
}
