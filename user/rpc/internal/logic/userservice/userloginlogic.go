package userservicelogic

import (
	"community.com/pkg/tool"
	"community.com/user/rpc/internal/svc"
	"community.com/user/rpc/pb"
	"context"
	"github.com/pkg/errors"
	"time"

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
		userData, err := l.svcCtx.MysqlClient.UserMysqlClient.FindOne(l.ctx, userId)
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
		userData, err := l.svcCtx.MysqlClient.UserMysqlClient.FindOneByPhone(l.ctx, aesPhone)
		if err != nil {
			return &__.UserLoginResp{}, err
		}
		userId = userData.UserId
		if !tool.ComparePassword(userData.Password, in.GetPassword()) {
			return nil, errors.Errorf("密码错误，请重新输入")
		}
	}
	token, err := tool.CreateToken(in.GetUserId(), time.Duration(l.svcCtx.Config.Token.AccessExpire)*time.Second, l.svcCtx.Config.Token.AccessSecret)
	if err != nil {
		return nil, err
	}
	return &__.UserLoginResp{
		Token: token,
	}, nil
}
