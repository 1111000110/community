package logic

import (
	aiclient "community/service/ai/model/mongo/agent"
	"context"

	"community/service/ai/rpc/internal/svc"
	__ "community/service/ai/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentLogic {
	return &CreateAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAgentLogic) CreateAgent(in *__.CreateAgentReq) (*__.CreateAgentResp, error) {
	// 不支持并发调用
	resp := &__.CreateAgentResp{}
	agentId, err := l.svcCtx.ModelClient.MongoAgent.GetLastAgentId(l.ctx) // 获取最后一个id
	if err != nil {
		return resp, err
	}
	in.Agent.AgentId = agentId + 1
	err = l.svcCtx.ModelClient.MongoAgent.Insert(l.ctx, aiclient.RpcAgentToMongoAgent(in.GetAgent()))
	return resp, err
}
