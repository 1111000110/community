package svc

import (
	"community/conf/databases/xmysql"
	"community/service/user/model/mysql/user"
)

type ModelClient struct {
	MysqlClient user.UserModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		MysqlClient: user.NewUserModel(xmysql.GetMysqlCommunityClient()),
	}
}
