package userservicelogic

import (
	"community.com/pkg/tool"
	"community.com/user/model/mysql/user"
	"community.com/user/rpc/internal/svc"
	"community.com/user/rpc/pb"
	"context"
	"github.com/pkg/errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRegisterLogic) UserRegister(in *__.UserRegisterReq) (*__.UserRegisterResp, error) {
	if in.GetPhone() == "" || in.GetPassword() == "" {
		return nil, errors.Errorf("手机号和密码不能为空")
	}
	encryptPhone, err := tool.Encrypt(in.GetPhone())
	if err != nil {
		return nil, err
	}
	hashPassword, err := tool.HashPassword(in.GetPassword())
	if err != nil {
		return nil, err
	}
	insertData := &user.User{
		Password: hashPassword,
		Phone:    encryptPhone,
		Role:     in.GetRole(),
		Ct:       time.Now().Unix(),
		Ut:       time.Now().Unix(),
	}
	if in.GetRole() == "" {
		insertData.Role = "user"
	}
	resultData, err := l.svcCtx.MysqlClient.UserMysqlClient.Insert(l.ctx, insertData)
	if err != nil {
		return nil, err
	}
	userId, err := resultData.LastInsertId()
	if err != nil {
		return nil, err
	}
	token, err := tool.CreateToken(userId, time.Duration(l.svcCtx.Config.Token.AccessExpire)*time.Second, l.svcCtx.Config.Token.AccessSecret)
	if err != nil {
		return nil, err
	}
	return &__.UserRegisterResp{
		Token:  token,
		UserId: userId,
	}, nil
}
