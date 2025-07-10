package userservicelogic

import (
	"context"

	"community.com/account/rpc/internal/svc"
	"community.com/account/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteLogic {
	return &UserDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDeleteLogic) UserDelete(in *__.UserDeleteReq) (*__.UserDeleteResp, error) {
	// todo: add your logic here and delete this line

	return &__.UserDeleteResp{}, nil
}
