package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	model "community/infra/review/model"
	"community/infra/review/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	ReviewModel []model.ReviewModel
	ReviewRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	data := &ServiceContext{
		Config:      c,
		ReviewModel: make([]model.ReviewModel, 0),
		ReviewRedis: redis.MustNewRedis(c.RedisConf),
	}
	for i := 0; i < len(c.DB.Review); i++ {
		data.ReviewModel = append(data.ReviewModel, model.NewReviewModel(c.DB.Review[i].Url, c.DB.Review[i].DB, c.DB.Review[i].Collection))
	}
	return data
}
