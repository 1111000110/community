package svc

import (
	"community/conf/databases/xscylla"
	"community/service/message/model/scylla/message"
)

type ModelClient struct {
	Scylla scyllamessage.MessageModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		Scylla: scyllamessage.NewMessageModel(xscylla.GetScyllaCommunitySession()),
	}
}
