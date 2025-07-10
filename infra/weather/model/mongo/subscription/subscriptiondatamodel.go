package model

import "github.com/zeromicro/go-zero/core/stores/mon"

var _ SubscriptionDataModel = (*customSubscriptionDataModel)(nil)

type (
	// SubscriptionDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSubscriptionDataModel.
	SubscriptionDataModel interface {
		subscriptionDataModel
	}

	customSubscriptionDataModel struct {
		*defaultSubscriptionDataModel
	}
)

// NewSubscriptionDataModel returns a model for the mongo.
func NewSubscriptionDataModel(url, db, collection string) SubscriptionDataModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customSubscriptionDataModel{
		defaultSubscriptionDataModel: newDefaultSubscriptionDataModel(conn),
	}
}
