package model

import (
	"community/conf/databases/xmongo"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var _ AgentModel = (*customAgentModel)(nil)

type (
	// AgentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAgentModel.
	AgentModel interface {
		agentModel
		GetAgentList(ctx context.Context) ([]*Agent, error)
		GetLastAgentId(ctx context.Context) (int64, error)
	}

	customAgentModel struct {
		*defaultAgentModel
	}
)

func (m *customAgentModel) GetAgentList(ctx context.Context) ([]*Agent, error) {
	var data []*Agent
	err := m.conn.Find(ctx, &data, bson.M{})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *customAgentModel) GetLastAgentId(ctx context.Context) (int64, error) {
	var result struct {
		AgentId int64 `bson:"agent_id"`
	}

	opts := options.FindOne().SetSort(bson.M{"agent_id": -1})

	err := m.conn.FindOne(ctx, &result, bson.M{}, opts)
	if err != nil {
		if errors.Is(err, mon.ErrNotFound) || errors.Is(err, mongo.ErrNoDocuments) {
			return 0, nil
		}
		return 0, err
	}

	return result.AgentId, nil
}

// NewAgentModel returns a model for the mongo.
func NewAgentModel(url, db, collection string) AgentModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customAgentModel{
		defaultAgentModel: newDefaultAgentModel(conn),
	}
}

// NewCommunityModel returns a community model for the mongo. It's Zhang Xuan's local model.
func NewCommunityModel(collection string) AgentModel {
	conn := xmongo.GetMongoCommunityClient(collection)
	return &customAgentModel{
		defaultAgentModel: newDefaultAgentModel(conn),
	}
}
