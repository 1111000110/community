package relationsservicelogic

import (
	"context"

	"community.com/account/rpc/internal/svc"
	"community.com/account/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRelationsUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRelationsUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRelationsUpdateLogic {
	return &UserRelationsUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRelationsUpdateLogic) UserRelationsUpdate(in *__.UserRelationsUpdateReq) (*__.UserRelationsUpdateResp, error) {
	// todo: add your logic here and delete this line

	return &__.UserRelationsUpdateResp{}, nil
}
