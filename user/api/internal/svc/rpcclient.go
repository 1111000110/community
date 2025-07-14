package svc

import (
	"community.com/user/api/internal/config"
	"community.com/user/rpc/client/userservice"
	"github.com/zeromicro/go-zero/zrpc"
)

type RpcClient struct {
	UserServiceClient userservice.UserService
}

func NewRpcClient(c config.Config) *RpcClient {
	return &RpcClient{
		UserServiceClient: userservice.NewUserService(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
