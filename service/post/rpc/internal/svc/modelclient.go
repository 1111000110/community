package svc

import model "community/service/post/model/mongo/post"

type ModelClient struct {
	PostMongoClient model.PostModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		PostMongoClient: model.NewCommunityModel("Post"),
	}
}
