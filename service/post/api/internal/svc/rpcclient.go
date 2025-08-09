package svc

import (
	"community.com/service/post/api/internal/config"
	"community.com/service/post/rpc/postservice"
	"community.com/service/user/rpc/client/userservice"
	"github.com/zeromicro/go-zero/zrpc"
)

type RpcClient struct {
	PostClient postservice.PostService
	UserClient userservice.UserService
}

func NewRpcClient(c config.Config) *RpcClient {
	return &RpcClient{
		PostClient: postservice.NewPostService(zrpc.MustNewClient(c.PostRpcConf)),
		UserClient: userservice.NewUserService(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
