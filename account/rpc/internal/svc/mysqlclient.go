package svc

import (
	"community.com/account/model/mysql/account"
	"community.com/account/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type MysqlClient struct {
	AccountMysqlClient account.AccountModel
}

func NewMysqlClient(c config.Config) *MysqlClient {
	return &MysqlClient{
		AccountMysqlClient: account.NewAccountModel(sqlx.NewMysql(c.MysqlConf.Account)),
	}
}
