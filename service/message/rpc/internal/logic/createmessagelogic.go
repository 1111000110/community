package logic

import (
	"community/pkg/snowflakes"
	"community/service/message/client"
	"community/service/message/model/scylla/message"
	"context"
	"time"

	"community/service/message/rpc/internal/svc"
	"community/service/message/rpc/pb"

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
	messageId, err := l.snowflakes.NextID()
	if err != nil {
		return nil, err
	}
	messageInfo := &scyllamessage.Message{
		MessageId:   messageId,
		SessionId:   in.GetSessionId(),
		SendId:      in.GetSendId(),
		ReplyId:     in.GetReplyId(),
		CreateTime:  time.Now().Unix(),
		UpdateTime:  0,
		Status:      client.MsgStatusNormal,
		Text:        in.GetContent().GetText(),
		MessageType: in.GetContent().GetMessageType(),
		Addition:    in.GetContent().GetAddition(),
	}
	if err = l.svcCtx.ModelClient.Scylla.CreateMessage(l.ctx, messageInfo); err != nil {
		return nil, err
	}
	return &__.CreateMessageResp{
		Message: scyllamessage.ModelToRpcModel(messageInfo),
	}, nil
}
