package logic

import (
	mysqlmessage "community/service/message/model/mysql/message"
	"context"

	"community/service/message/rpc/internal/svc"
	"community/service/message/rpc/pb"

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
	data, err := l.svcCtx.ModelClient.MysqlMessage.GetMessageList(l.ctx, in.GetSessionId(), in.GetReq(), in.GetLimit())
	if err != nil {
		return nil, err
	}
	return &__.GetMessageListResp{
		Message: mysqlmessage.ModelsToRpcModels(data),
	}, nil
}
