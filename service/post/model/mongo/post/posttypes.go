package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PostId   int64              `bson:"postId,omitempty" json:"postId,omitempty"`   // 帖子ID
	UserId   int64              `bson:"userId,omitempty" json:"userId,omitempty"`   // 用户ID
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`     // 标题
	Content  string             `bson:"content,omitempty" json:"content,omitempty"` // 内容
	Images   string             `bson:"images,omitempty" json:"images,omitempty"`   // 图片
	Theme    string             `bson:"theme,omitempty" json:"theme,omitempty"`     // 主题
	Tags     string             `bson:"tags,omitempty" json:"tags,omitempty"`       // 标签
	Status   int32              `bson:"status,omitempty" json:"status,omitempty"`   // 状态
	UpdateAt int64              `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt int64              `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
