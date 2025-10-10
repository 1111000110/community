package logic

import (
	"community/pkg/xstring"
	"community/service/websocket/client"
	"community/service/websocketpush/rpc/internal/svc"
	"community/service/websocketpush/rpc/pb"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebSocketPushLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWebSocketPushLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebSocketPushLogic {
	return &WebSocketPushLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WebSocketPushLogic) WebSocketPush(in *__.WebSocketPushReq) (*__.WebSocketPushResp, error) {
	redisKey := make([]string, 0)
	for _, connId := range in.GetConnId() {
		client.GetRedisKeyName(xstring.IntToString(connId))
	}
	info, err := l.svcCtx.Model.RedisClient.MgetCtx(l.ctx, redisKey...)
	if err != nil {
		return nil, err
	}
	sendMap := make(map[string][]int64)
	for id, connId := range in.GetConnId() {
		if info[id] != "" {
			sendMap[info[id]] = append(sendMap[info[id]], connId)
		}
	}
	for _, v := range sendMap {
		err = l.svcCtx.Model.KafkaMessageClient.PushWithKey(l.ctx, client.GetKeyByIds(v), in.GetPushData().String())
		if err != nil {
			return nil, err
		}
	}
	return &__.WebSocketPushResp{}, nil
}
