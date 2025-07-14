package svc

import (
	"community.com/user/model/mysql/user"
	"community.com/user/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type MysqlClient struct {
	UserMysqlClient user.UserModel
}

func NewMysqlClient(c config.Config) *MysqlClient {
	return &MysqlClient{
		UserMysqlClient: user.NewUserModel(sqlx.NewMysql(c.MysqlConf.User)),
	}
}
