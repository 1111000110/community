package logic

import (
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
	// todo: add your logic here and delete this line

	return &__.GetMessageListResp{}, nil
}
