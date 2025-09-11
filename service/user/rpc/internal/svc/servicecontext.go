package svc

import (
	"community/service/user/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Model  *ModelClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Model:  DefaultModelClient(),
	}
}
