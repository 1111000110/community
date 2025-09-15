package message

import (
	"community/service/message/rpc/message"
	"context"

	"community/service/message/api/internal/svc"
	"community/service/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageDeleteByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageDeleteByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageDeleteByIdLogic {
	return &MessageDeleteByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageDeleteByIdLogic) MessageDeleteById(req *types.MessageDeleteByIdReq) (resp *types.MessageDeleteByIdResp, err error) {
	_, err = l.svcCtx.RpcClient.Message.DeleteMessage(l.ctx, &message.DeleteMessageReq{
		SessionId: req.SessionId,
		SendId:    req.SendId,
		MessageId: req.MessageId,
	})
	return
}
