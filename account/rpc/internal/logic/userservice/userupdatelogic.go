package userservicelogic

import (
	"context"

	"community.com/account/rpc/internal/svc"
	"community.com/account/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserUpdateLogic) UserUpdate(in *__.UserUpdateReq) (*__.UserUpdateResp, error) {
	// todo: add your logic here and delete this line

	return &__.UserUpdateResp{}, nil
}
