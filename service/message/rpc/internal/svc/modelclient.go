package svc

import (
	"community/conf/databases/xmysql"
	mysqlmessage "community/service/message/model/mysql/message"
)

type ModelClient struct {
	MysqlMessage mysqlmessage.MessageModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		MysqlMessage: mysqlmessage.NewMessageModel(xmysql.GetMysqlCommunityClient()),
	}
}
