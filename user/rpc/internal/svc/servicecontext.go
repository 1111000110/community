package svc

import "community.com/user/rpc/internal/config"

type ServiceContext struct {
	Config      config.Config
	MysqlClient *MysqlClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlClient := NewMysqlClient(c)
	return &ServiceContext{
		Config:      c,
		MysqlClient: mysqlClient,
	}
}
