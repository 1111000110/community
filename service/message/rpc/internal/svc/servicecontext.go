package svc

import (
	"community.com/conf/databases/scylla"
	"community.com/service/message/model/scylla/message"
	"community.com/service/message/rpc/internal/config"
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
