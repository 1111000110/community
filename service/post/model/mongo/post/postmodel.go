package model

import (
	"community.com/conf/databases/mongo"
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
)

var _ PostModel = (*customPostModel)(nil)

type (
	// PostModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPostModel.
	PostModel interface {
		postModel
		FindOneByPostId(context.Context, int64) (*Post, error)
	}

	customPostModel struct {
		*defaultPostModel
	}
)

func (m *customPostModel) FindOneByPostId(ctx context.Context, postId int64) (*Post, error) {
	var data Post
	err := m.conn.FindOne(ctx, &data, bson.M{"postId": postId})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewPostModel returns a model for the mongo.
func NewPostModel(url, db, collection string) PostModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customPostModel{
		defaultPostModel: newDefaultPostModel(conn),
	}
}

// NewCommunityModel returns a community model for the mongo. It's Zhang Xuan's local model.
func NewCommunityModel(collection string) PostModel {
	conn := mongo.GetMongoCommunityClient(collection)
	return &customPostModel{
		defaultPostModel: newDefaultPostModel(conn),
	}
}
