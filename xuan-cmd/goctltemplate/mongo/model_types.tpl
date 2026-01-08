package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoBaseModel struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`                 // MongoDB 内置主键
	UpdateTime int64         `bson:"update_time,omitempty" json:"update_time,omitempty"` // 最近更新时间（Unix 时间戳）
	CreateTime int64         `bson:"create_time,omitempty" json:"create_time,omitempty"` // 创建时间（Unix 时间戳）
}

type {{.Type}} struct {
	MongoBaseModel
	// TODO: Fill your own fields
}
