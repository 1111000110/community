package svc

import (
	"community/service/message/api/internal/config"
	"community/service/message/rpc/message"
	"community/service/websocketpush/rpc/websocketpush"
	"github.com/zeromicro/go-zero/zrpc"
)

type RpcClient struct {
	Message       message.Message
	WebSocketPush websocketpush.WebsocketPush
}

func NewRpcClient(c config.Config) *RpcClient {
	return &RpcClient{
		Message:       message.NewMessage(zrpc.MustNewClient(c.MessageConf)),
		WebSocketPush: websocketpush.NewWebsocketPush(zrpc.MustNewClient(c.WebSocketPushConf)),
	}
}
