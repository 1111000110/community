package userservicelogic

import (
	"community.com/pkg/tool"
	"community.com/user/rpc/client/userservice"
	"context"

	"community.com/user/rpc/internal/svc"
	"community.com/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserQueryLogic {
	return &UserQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserQueryLogic) UserQuery(in *__.UserQueryReq) (*__.UserQueryResp, error) {
	userData, err := l.svcCtx.MysqlClient.UserMysqlClient.FindOne(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	resp := &__.UserQueryResp{
		User: &__.User{
			UserId:    in.GetUserId(),
			Phone:     "",
			Email:     "",
			Password:  "",
			Nickname:  userData.Nickname,
			Avatar:    userData.Avatar,
			Gender:    userData.Gender,
			BirthDate: userData.BirthDate,
			Role:      "",
			Status:    0,
			CreateAt:  userData.Ct,
			UpdateAt:  userData.Ut,
		},
	}
	switch in.GetType() {
	case userservice.GetPrivateInfo:
		resp.User.Role = userData.Role
		resp.User.Email = userData.Email
		resp.User.Status = userData.Status
		phone, _ := tool.Decrypt(userData.Phone)
		resp.User.Phone = phone
	}
	return resp, nil
}
