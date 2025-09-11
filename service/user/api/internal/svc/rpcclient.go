package svc

import (
	"community/service/user/api/internal/config"
	"community/service/user/rpc/client/userservice"
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
