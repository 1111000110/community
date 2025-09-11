package svc

import (
	"community/conf/databases/scylla"
	"community/service/message/model/scylla/message"
	"community/service/message/rpc/internal/config"
)

type ServiceContext struct {
	Config       config.Config
	ScyllaClient message.MessageModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ScyllaClient: message.NewMessageModel(scylla.GetScyllaCommunitySession()),
	}
}
