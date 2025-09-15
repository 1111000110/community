package message

import (
	"community/service/message/rpc/message"
	"context"

	"community/service/message/api/internal/svc"
	"community/service/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageListLogic {
	return &MessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageListLogic) MessageList(req *types.MessageListReq) (resp *types.MessageListResp, err error) {
	resp = &types.MessageListResp{}
	info, err := l.svcCtx.RpcClient.Message.GetMessageList(l.ctx, &message.GetMessageListReq{
		SessionId: req.SessionId,
		Req:       req.Req,
		Limit:     req.Limit,
	})
	if err != nil {
		return resp, err
	}
	for _, d := range info.GetMessage() {
		resp.MessageDetails = append(resp.MessageDetails, types.MessageDetail{
			MessageId:  d.MessageId,
			SessionId:  d.SessionId,
			SendId:     d.SendId,
			ReplyId:    d.ReplyId,
			CreateTime: d.CreateTime,
			UpdateTime: d.UpdateTime,
			Status:     d.Status,
			Content: types.MessageContent{
				Text:        d.Content.Text,
				MessageType: d.Content.MessageType,
				Addition:    d.Content.Addition,
			},
		})
	}
	return
}
