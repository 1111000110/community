package logic

import (
	"context"

	"community.com/service/message/rpc/internal/svc"
	"community.com/service/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageLogic {
	return &DeleteMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMessageLogic) DeleteMessage(in *__.DeleteMessageReq) (*__.DeleteMessageResp, error) {
	if err := l.svcCtx.ScyllaClient.DeleteMessage(l.ctx, in.GetSessionId(), in.GetMessageId()); err != nil {
		return nil, err
	}
	return &__.DeleteMessageResp{}, nil
}
