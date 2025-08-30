package logic

import (
	"community.com/pkg/snowflakes"
	"community.com/service/message/model/scylla/message"
	"context"

	"community.com/service/message/rpc/internal/svc"
	"community.com/service/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	snowflakes *snowflakes.Snowflake
}

func NewCreateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMessageLogic {
	snowflake, err := snowflakes.NewSnowflake(1, 1)
	if err != nil {
		panic(err)
	}
	return &CreateMessageLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		snowflakes: snowflake,
	}
}

func (l *CreateMessageLogic) CreateMessage(in *__.CreateMessageReq) (*__.CreateMessageResp, error) {
	var err error
	messageId := in.GetMessage().GetMessageId()
	if messageId == 0 {
		messageId, err = l.snowflakes.NextID()
		if err != nil {
			return nil, err
		}
	}
	if err := l.svcCtx.ScyllaClient.CreateMessage(l.ctx, &message.Message{
		MessageId:   messageId,
		SessionId:   in.GetMessage().GetSessionId(),
		SendId:      in.GetMessage().GetSendId(),
		RecipientId: in.GetMessage().GetRecipientId(),
		ReplyId:     in.GetMessage().GetReplyId(),
		CreateTime:  in.GetMessage().GetCreateTime(),
		UpdateTime:  in.GetMessage().GetUpdateTime(),
		Status:      in.GetMessage().GetStatus(),
		Text:        in.GetMessage().GetContent().GetText(),
		MessageType: in.GetMessage().GetContent().GetMessageType(),
		Addition:    in.GetMessage().GetContent().GetAddition(),
	}); err != nil {
		return nil, err
	}
	return &__.CreateMessageResp{}, nil
}
