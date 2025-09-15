package svc

import (
	"community/service/user/api/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	RpcClient *RpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	rpcClient := NewRpcClient(c)
	return &ServiceContext{
		Config:    c,
		RpcClient: rpcClient,
	}
}
