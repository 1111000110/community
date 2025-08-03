package svc

import (
	"community.com/service/post/api/internal/config"
	"community.com/service/post/rpc/postservice"
	"github.com/zeromicro/go-zero/zrpc"
)

type RpcClient struct {
	UserServiceClient postservice.PostService
}

func NewRpcClient(c config.Config) *RpcClient {
	return &RpcClient{
		UserServiceClient: postservice.NewPostService(zrpc.MustNewClient(c.PostRpcConf)),
	}
}
