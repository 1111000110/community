package logic

import (
	"community/pkg/snowflakes"
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
	messageId := in.GetMessage().GetMessageId()
	if messageId == 0 {
		messageId, err = l.snowflakes.NextID()
		if err != nil {
			return nil, err
		}
	}
	in.Message.MessageId = messageId
	in.Message.CreateTime = time.Now().Unix()
	if err := l.svcCtx.ModelClient.Scylla.CreateMessage(l.ctx, scyllamessage.RpcModelToModel(in.GetMessage())); err != nil {
		return nil, err
	}
	return &__.CreateMessageResp{
		Message: in.GetMessage(),
	}, nil
}
