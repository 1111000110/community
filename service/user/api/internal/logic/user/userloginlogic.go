package user

import (
	"community/service/user/rpc/client/userservice"
	"context"

	"community/service/user/api/internal/svc"
	"community/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginReq) (resp *types.UserLoginResp, err error) {
	userLoginResp, err := l.svcCtx.RpcClient.UserServiceClient.UserLogin(l.ctx, &userservice.UserLoginReq{
		Phone:    req.Phone,
		UserId:   req.UserId,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.UserLoginResp{
		Token: userLoginResp.Token,
	}
	return
}
