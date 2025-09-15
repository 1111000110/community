package svc

import (
	"community/service/message/api/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	RpcClient *RpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		RpcClient: NewRpcClient(c),
	}
}
