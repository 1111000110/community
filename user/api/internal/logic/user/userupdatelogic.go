package user

import (
	"context"
	"errors"

	"community.com/pkg/tool"
	"community.com/user/api/internal/svc"
	"community.com/user/api/internal/types"
	"community.com/user/rpc/client/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UserUpdateReq) (resp *types.UserUpdateResp, err error) {
	userId, err := tool.GetUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	if userId != req.UserInfo.UserBase.UserId {
		return nil, errors.New("用户ID不匹配")
	}
	_, err = l.svcCtx.RpcClient.UserServiceClient.UserUpdate(l.ctx, &userservice.UserUpdateReq{
		User: &userservice.User{
			UserId:    userId,
			Phone:     req.UserInfo.UserPrivate.Phone,
			Email:     req.UserInfo.UserPrivate.Email,
			Password:  req.UserInfo.UserPrivate.Password,
			Nickname:  req.UserInfo.UserBase.NickName,
			Avatar:    req.UserInfo.UserBase.Avatar,
			Gender:    req.UserInfo.UserBase.Gender,
			BirthDate: req.UserInfo.UserBase.BirthDate,
			Role:      req.UserInfo.UserPrivate.Role,
			Status:    req.UserInfo.UserPrivate.Status,
		},
	})
	if err != nil {
		return nil, err
	}
	return
}
