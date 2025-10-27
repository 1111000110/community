package logic

import (
	mysqlmessage "community/service/message/model/mysql/message"
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
	if err := l.svcCtx.ModelClient.MysqlMessage.Update(l.ctx, &mysqlmessage.Message{
		SessionId:   in.GetMessage().GetSessionId(),
		MessageId:   in.GetMessage().GetMessageId(),
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
