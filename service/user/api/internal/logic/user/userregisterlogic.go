package user

import (
	"community.com/service/user/rpc/client/userservice"
	"context"

	"community.com/service/user/api/internal/svc"
	"community.com/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq) (resp *types.UserRegisterResp, err error) {
	respData, err := l.svcCtx.RpcClient.UserServiceClient.UserRegister(l.ctx, &userservice.UserRegisterReq{
		Phone:    req.Phone,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		return nil, err
	}
	return &types.UserRegisterResp{
		UserId: respData.GetUserId(),
		Token:  respData.GetToken(),
	}, nil
}
