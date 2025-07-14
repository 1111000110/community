package user

import (
	"community.com/pkg/tool"
	"community.com/user/rpc/client/userservice"
	"context"
	"github.com/pkg/errors"

	"community.com/user/api/internal/svc"
	"community.com/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserQueryLogic {
	return &UserQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserQueryLogic) UserQuery(req *types.UserQueryReq) (resp *types.UserQueryResp, err error) {
	userId, err := tool.GetUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	queryUserId := req.QueryUserId
	queryType := req.Type
	logx.Infof("userId:%d,queryType:%s,queryUserId:%d", userId, queryType, queryUserId)
	switch queryType {
	case userservice.GetPrivateInfo:
		if userId != queryUserId {
			return nil, errors.Errorf("权限不足")
		}
	}
	if queryUserId == 0 {
		queryUserId = userId
	}
	data, err := l.svcCtx.RpcClient.UserServiceClient.UserQuery(l.ctx, &userservice.UserQueryReq{
		UserId: queryUserId,
		Type:   queryType,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.UserQueryResp{
		User: types.UserInfo{
			UserBase: types.UserBase{
				UserId:    data.User.GetUserId(),
				NickName:  data.User.GetNickname(),
				Avatar:    data.User.GetAvatar(),
				Gender:    data.User.GetGender(),
				BirthDate: data.User.GetBirthDate(),
			},
			UserPrivate: types.UserPrivate{
				UserId: data.User.GetUserId(),
				Phone:  data.User.GetPhone(),
				Email:  data.User.GetEmail(),
				Role:   data.User.GetRole(),
				Status: data.User.GetStatus(),
			},
		},
	}
	return resp, nil
}
