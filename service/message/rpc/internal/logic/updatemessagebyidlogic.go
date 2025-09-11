package logic

import (
	"community/service/message/model/scylla/message"
	"context"

	"community/service/message/rpc/internal/svc"
	"community/service/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMessageByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMessageByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMessageByIdLogic {
	return &UpdateMessageByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMessageByIdLogic) UpdateMessageById(in *__.UpdateMessageByIdReq) (*__.UpdateMessageByIdResp, error) {
	if err := l.svcCtx.ScyllaClient.UpdateMessageById(l.ctx, in.GetMessage().GetSessionId(), in.GetMessage().GetMessageId(), &message.Message{
		SendId:      in.GetMessage().GetSendId(),
		ReplyId:     in.GetMessage().GetReplyId(),
		UpdateTime:  in.GetMessage().GetUpdateTime(),
		Status:      in.GetMessage().GetStatus(),
		Text:        in.GetMessage().GetContent().GetText(),
		MessageType: in.GetMessage().GetContent().GetMessageType(),
		Addition:    in.GetMessage().GetContent().GetAddition(),
	}); err != nil {
		return nil, err
	}
	return &__.UpdateMessageByIdResp{}, nil
}
