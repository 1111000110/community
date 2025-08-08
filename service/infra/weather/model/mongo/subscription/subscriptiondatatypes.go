package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionData struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID         string             `bson:"user_id" json:"user_id"`
	OpenId         string             `bson:"open_id" json:"open_id"`
	City           string             `bson:"city" json:"city"`
	Time           string             `bson:"time" json:"time"`
	MaxTemperature string             `bson:"max_temperature" json:"max_temperature"`
	MinTemperature string             `bson:"min_temperature" json:"min_temperature"`
	Weather        []string           `bson:"weather" json:"weather"`
	Status         int64              `bson:"status" json:"status"`
	UpdateAt       int64              `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt       int64              `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
