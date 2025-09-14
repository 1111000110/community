package message

import (
	"community/service/message/rpc/message"
	"community/service/websocketpush/rpc/websocketpush"
	"context"

	"community/service/message/api/internal/svc"
	"community/service/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageCreateLogic {
	return &MessageCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageCreateLogic) MessageCreate(req *types.MessageCreateReq) (resp *types.MessageCreateResp, err error) {
	info, err := l.svcCtx.RpcClient.Message.CreateMessage(l.ctx, &message.CreateMessageReq{
		Message: &message.MessageDetail{
			MessageId:  0,
			SessionId:  req.SessionId,
			SendId:     req.SendId,
			ReplyId:    req.ReplyId,
			CreateTime: 0,
			UpdateTime: 0,
			Status:     req.Status,
			Content: &message.MessageContent{
				Text:        req.Content.Text,
				MessageType: req.Content.MessageType,
				Addition:    req.Content.Addition,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.RpcClient.WebSocketPush.WebSocketPush(l.ctx, &websocketpush.WebSocketPushReq{
		ConnId: []int64{111},
		PushData: &websocketpush.WebSocketPushData{
			NotifyType: "message",
			NotifyVal:  info.GetMessage().String(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &types.MessageCreateResp{}, nil
}
