package logic

import (
	"context"

	"community/service/ai/rpc/internal/svc"
	"community/service/ai/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpsertAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpsertAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpsertAgentLogic {
	return &UpsertAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpsertAgentLogic) UpsertAgent(in *__.UpsertAgentReq) (*__.UpsertAgentResp, error) {
	// todo: add your logic here and delete this line
	return &__.UpsertAgentResp{}, nil
}
