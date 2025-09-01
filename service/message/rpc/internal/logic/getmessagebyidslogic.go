package logic

import (
	"community.com/service/message/model/scylla/message"
	"context"

	"community.com/service/message/rpc/internal/svc"
	"community.com/service/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageByIdsLogic {
	return &GetMessageByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessageByIdsLogic) GetMessageByIds(in *__.GetMessageByIdsReq) (*__.GetMessageByIdsResp, error) {
	data, err := l.svcCtx.ScyllaClient.GetMessageByIds(l.ctx, in.GetSessionId(), in.GetMessageId())
	if err != nil {
		return nil, err
	}
	return &__.GetMessageByIdsResp{
		Message: message.ModelsToRpcModels(data),
	}, nil
}
