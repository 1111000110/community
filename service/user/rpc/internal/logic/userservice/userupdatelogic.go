package userservicelogic

import (
	"community/pkg/tool"
	"context"
	"time"

	"community/service/user/rpc/internal/svc"
	"community/service/user/rpc/pb"

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
	userData := in.GetUser()
	data, err := l.svcCtx.Model.MysqlClient.FindOne(l.ctx, userData.UserId)
	if err != nil {
		return nil, err
	}
	if in.GetUser().GetPhone() != "" {
		encryptPhone, err := tool.Encrypt(in.GetUser().GetPhone())
		if err != nil {
			return nil, err
		}
		data.Phone = encryptPhone
	}
	if in.GetUser().GetEmail() != "" {
		data.Email = in.GetUser().GetEmail()
	}
	if in.GetUser().GetNickname() != "" {
		data.Nickname = in.GetUser().GetNickname()
	}
	if in.GetUser().GetAvatar() != "" {
		data.Avatar = in.GetUser().GetAvatar()
	}
	if in.GetUser().GetGender() != "" {
		data.Gender = in.GetUser().GetGender()
	}
	if in.GetUser().GetPassword() != "" {
		hashPassword, err := tool.HashPassword(in.GetUser().GetPassword())
		if err != nil {
			return nil, err
		}
		data.Password = hashPassword
	}
	if in.GetUser().GetBirthDate() != 0 {
		data.BirthDate = in.GetUser().GetBirthDate()
	}
	if in.GetUser().GetRole() != "" {
		data.Role = in.GetUser().GetRole()
	}
	if in.GetUser().GetStatus() != 0 {
		data.Status = in.GetUser().GetStatus()
	}
	data.Ut = time.Now().Unix()
	err = l.svcCtx.Model.MysqlClient.Update(l.ctx, data)
	if err != nil {
		return nil, err
	}
	return &__.UserUpdateResp{
		User: userData,
	}, nil
}
