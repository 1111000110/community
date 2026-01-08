package svc

import "community/service/ai/rpc/internal/config"

type ServiceContext struct {
	Config      config.Config
	ModelClient *ModelClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		ModelClient: DefaultModelClient(),
	}
}
