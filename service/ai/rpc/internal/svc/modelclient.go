package svc

import (
	mongoAgent "community/service/ai/model/mongo/agent"
)

type ModelClient struct {
	MongoAgent mongoAgent.AgentModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		MongoAgent: mongoAgent.NewCommunityModel("agent"),
	}
}
