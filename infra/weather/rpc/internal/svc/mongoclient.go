package svc

import (
	model "community.com/infra/weather/model/mongo/subscription"
	"community.com/infra/weather/rpc/internal/config"
)

type MongoClient struct {
	SubscriptionDataModel model.SubscriptionDataModel
}

func NewMongoClient(c config.Config) *MongoClient {
	return &MongoClient{
		SubscriptionDataModel: model.NewSubscriptionDataModel(c.MongoConf.SubscriptionData.Url, c.MongoConf.SubscriptionData.Db, c.MongoConf.SubscriptionData.Collection),
	}
}
