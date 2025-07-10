package userservicelogic

import (
	"community.com/account/rpc/internal/svc"
	"community.com/account/rpc/pb"
	"community.com/pkg/tool"
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
		accountData, err := l.svcCtx.MysqlClient.AccountMysqlClient.FindOne(l.ctx, userId)
		if err != nil {
			return &__.UserLoginResp{}, err
		}
		if !tool.ComparePassword(accountData.Password, in.GetPassword()) {
			return nil, errors.Errorf("密码错误，请重新输入")
		}
	}
	if phone != "" {
		aesPhone, err := tool.Encrypt(phone)
		if err != nil {
			return &__.UserLoginResp{}, err
		}
		accountData, err := l.svcCtx.MysqlClient.AccountMysqlClient.FindOneByPhone(l.ctx, aesPhone)
		if err != nil {
			return &__.UserLoginResp{}, err
		}
		userId = accountData.UserId
		if !tool.ComparePassword(accountData.Password, in.GetPassword()) {
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
