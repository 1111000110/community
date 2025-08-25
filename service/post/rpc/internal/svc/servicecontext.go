package svc

import (
	model "community.com/service/post/model/mongo/post"
	"community.com/service/post/rpc/internal/config"
)

type ServiceContext struct {
	Config          config.Config
	PostMongoClient model.PostModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		PostMongoClient: model.NewCommunityModel("Post"),
	}
}
