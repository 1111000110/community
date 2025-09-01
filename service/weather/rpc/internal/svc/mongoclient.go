package svc

import (
	model "community.com/service/weather/model/mongo/subscription"
	"community.com/service/weather/rpc/internal/config"
)

type MongoClient struct {
	SubscriptionDataModel model.SubscriptionDataModel
}

func NewMongoClient(c config.Config) *MongoClient {
	return &MongoClient{
		SubscriptionDataModel: model.NewSubscriptionDataModel(c.MongoConf.SubscriptionData.Url, c.MongoConf.SubscriptionData.Db, c.MongoConf.SubscriptionData.Collection),
	}
}
