package relationsservicelogic

import (
	"context"

	"community.com/user/rpc/internal/svc"
	"community.com/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRelationsGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRelationsGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRelationsGetLogic {
	return &UserRelationsGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRelationsGetLogic) UserRelationsGet(in *__.UserRelationsGetReq) (*__.UserRelationsGetResp, error) {
	// todo: add your logic here and delete this line

	return &__.UserRelationsGetResp{}, nil
}
