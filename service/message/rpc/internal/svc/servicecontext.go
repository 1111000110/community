package svc

import (
	"community/service/message/rpc/internal/config"
)

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
