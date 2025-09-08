package userservicelogic

import (
	"community.com/pkg/tool"
	"community.com/service/user/rpc/internal/svc"
	"community.com/service/user/rpc/pb"
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserLoginLogic) UserLogin(in *__.UserLoginReq) (*__.UserLoginResp, error) {
	userId := in.GetUserId()
	phone := in.GetPhone()
	password := in.GetPassword()
	if password == "" {
		return nil, errors.Errorf("密码不能为空")
	}
	if userId == 0 && phone == "" {
		return &__.UserLoginResp{}, errors.Errorf("手机号和id不能同时为空")
	}
	if userId != 0 {
		userData, err := l.svcCtx.Model.MysqlClient.FindOne(l.ctx, userId)
		if err != nil {
			return &__.UserLoginResp{}, err
		}
		if !tool.ComparePassword(userData.Password, in.GetPassword()) {
			return nil, errors.Errorf("密码错误，请重新输入")
		}
	}
	if phone != "" {
		aesPhone, err := tool.Encrypt(phone)
		if err != nil {
			return &__.UserLoginResp{}, err
		}
		userData, err := l.svcCtx.Model.MysqlClient.FindOneByPhone(l.ctx, aesPhone)
		if err != nil {
			return &__.UserLoginResp{}, err
		}
		userId = userData.UserId
		if !tool.ComparePassword(userData.Password, in.GetPassword()) {
			return nil, errors.Errorf("密码错误，请重新输入")
		}
	}
	token, err := tool.CreateTokenByUserID(in.GetUserId())
	if err != nil {
		return nil, err
	}
	return &__.UserLoginResp{
		Token: token,
	}, nil
}
