package logic

import (
	"context"

	"community/service/ai/rpc/internal/svc"
	"community/service/ai/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRunAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunAgentLogic {
	return &RunAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RunAgentLogic) RunAgent(in *__.RunAgentReq) (*__.RunAgentResp, error) {
	// todo: add your logic here and delete this line

	return &__.RunAgentResp{}, nil
}
