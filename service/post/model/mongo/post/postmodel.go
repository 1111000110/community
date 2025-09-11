package model

import (
	"community/conf/databases/mongo"
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
		DeleteOneByPostId(context.Context, int64) error
		FindAllByPostIds(context.Context, []int64) ([]*Post, error)
	}

	customPostModel struct {
		*defaultPostModel
	}
)

func (m *customPostModel) FindAllByPostIds(ctx context.Context, postIds []int64) ([]*Post, error) {
	var data []*Post
	err := m.conn.Find(ctx, &data, bson.M{"postId": bson.M{"$in": postIds}})
	if err != nil {
		return nil, err
	}
	return data, nil
}

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

func (m *customPostModel) DeleteOneByPostId(ctx context.Context, postId int64) error {
	_, err := m.conn.DeleteOne(ctx, &bson.M{"postId": postId})
	return err
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
