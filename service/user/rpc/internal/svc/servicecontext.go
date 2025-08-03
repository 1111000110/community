package svc

import (
	"community.com/conf/databases/mysql"
	"community.com/service/user/model/mysql/user"
	"community.com/service/user/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	MysqlClient user.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		MysqlClient: user.NewUserModel(mysql.GetMysqlCommunityClient()),
	}
}
