package svc

import (
	"community/conf/databases/xscylla"
	"community/service/message/model/scylla/message"
)

type ModelClient struct {
	Scylla message.MessageModel
}

func DefaultModelClient() *ModelClient {
	return &ModelClient{
		Scylla: message.NewMessageModel(xscylla.GetScyllaCommunitySession()),
	}
}
