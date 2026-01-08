package logic

import (
	mongoAgent "community/service/ai/model/mongo/agent"
	"context"

	"community/service/ai/rpc/internal/svc"
	"community/service/ai/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAgentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentListLogic {
	return &GetAgentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAgentListLogic) GetAgentList(_ *__.GetAgentListReq) (*__.GetAgentListResp, error) {
	info, err := l.svcCtx.ModelClient.MongoAgent.GetAgentList(l.ctx)
	if err != nil {
		return &__.GetAgentListResp{}, err
	}
	return &__.GetAgentListResp{
		AgentList: mongoAgent.MongoAgentsToRpcAgents(info),
	}, nil
}
