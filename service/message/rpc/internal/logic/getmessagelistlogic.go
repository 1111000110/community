package logic

import (
	"community.com/service/message/model/scylla/message"
	"context"

	"community.com/service/message/rpc/internal/svc"
	"community.com/service/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessageListLogic) GetMessageList(in *__.GetMessageListReq) (*__.GetMessageListResp, error) {
	data, err := l.svcCtx.ScyllaClient.GetMessageList(l.ctx, in.GetSessionId(), int(in.GetReq()), int(in.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &__.GetMessageListResp{
		Message: message.ModelsToRpcModels(data),
	}, nil
}
