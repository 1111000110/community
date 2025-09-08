package svc

import (
	"community.com/conf/databases/mysql"
	"community.com/service/user/model/mysql/user"
)

type ModelClient struct {
	MysqlClient user.UserModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		MysqlClient: user.NewUserModel(mysql.GetMysqlCommunityClient()),
	}
}
