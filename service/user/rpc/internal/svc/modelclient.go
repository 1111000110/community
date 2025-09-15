package svc

import (
	"community/conf/databases/xmysql"
	"community/service/user/model/mysql/user"
)

type ModelClient struct {
	Mysql user.UserModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		Mysql: user.NewUserModel(xmysql.GetMysqlCommunityClient()),
	}
}
