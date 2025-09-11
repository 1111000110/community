package svc

import (
	"community/service/user/api/internal/config"
	"community/service/user/api/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config     config.Config
	Middleware rest.Middleware
	RpcClient  *RpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	rpcClient := NewRpcClient(c)
	return &ServiceContext{
		Config:     c,
		Middleware: middleware.NewMiddleware().Handle,
		RpcClient:  rpcClient,
	}
}
