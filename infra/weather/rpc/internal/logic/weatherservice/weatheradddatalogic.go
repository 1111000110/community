package weatherservicelogic

import (
	model "community.com/infra/weather/model/mongo/subscription"
	"community.com/infra/weather/rpc/internal/svc"
	"community.com/infra/weather/rpc/pb"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type WeatherAddDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWeatherAddDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WeatherAddDataLogic {
	return &WeatherAddDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WeatherAddDataLogic) WeatherAddData(in *__.WeatherAddDataReq) (*__.WeatherAddDataResp, error) {
	err := l.svcCtx.MongoClient.SubscriptionDataModel.Insert(l.ctx, &model.SubscriptionData{
		UserID:         in.UserId,
		OpenId:         in.OpenId,
		City:           in.City,
		Time:           in.Time,
		MaxTemperature: in.MaxTemperature,
		MinTemperature: in.MinTemperature,
		Weather:        in.Weather,
		Status:         in.Status,
	})
	if err != nil {
		return nil, err
	}
	return &__.WeatherAddDataResp{}, nil
}
