package svc

import model "community/service/post/model/mongo/post"

type ModelClient struct {
	MongoPost model.PostModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		MongoPost: model.NewCommunityModel("post"),
	}
}
