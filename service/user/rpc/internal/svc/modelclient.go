package svc

import (
	"community/conf/databases/mysql"
	"community/service/user/model/mysql/user"
)

type ModelClient struct {
	MysqlClient user.UserModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		MysqlClient: user.NewUserModel(mysql.GetMysqlCommunityClient()),
	}
}
