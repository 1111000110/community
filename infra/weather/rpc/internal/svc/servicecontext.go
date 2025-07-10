package svc

import (
	"community.com/infra/weather/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	MongoClient *MongoClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mongoClient := NewMongoClient(c)
	return &ServiceContext{
		Config:      c,
		MongoClient: mongoClient,
	}
}
