package logic

import (
	model "community/review/model"
	"community/review/model/redisKey"
	"community/review/rpc/generateFileV1/internal/svc"
	"community/review/rpc/pb"
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"math"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReviewGetByBaseIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReviewGetByBaseIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewGetByBaseIdLogic {
	return &ReviewGetByBaseIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReviewGetByBaseIdLogic) ReviewGetByBaseId(in *__.ReviewGetByBaseIdReq) (*__.ReviewGetResp, error) {
	// todo: add your logic here and delete this line
	if in.GetBaseId() == 0 || in.GetHeadId() == 0 {
		return &__.ReviewGetResp{}, errors.New("baseId or headId is null")
	}
	redisData := []redis.Pair{}
	redisData, err := l.svcCtx.ReviewRedis.ZrangebyscoreWithScoresAndLimitCtx(l.ctx, rediskey.ReviewKey(in.GetBaseId()), -1, math.MaxInt64, int(in.GetOffset()), int(in.GetLimit()))
	if err != nil {
		return nil, err
	}
	if len(redisData) == 1 && redisData[0].Key == "" && redisData[0].Score == 0 {
		return &__.ReviewGetResp{}, nil
	}
	if len(redisData) == 0 {
		ans, err := l.svcCtx.ReviewModel[int(in.GetHeadId()%(int64(len(l.svcCtx.ReviewModel))))].FindByBaseId(l.ctx, in.GetBaseId())
		if err != nil || ans == nil { //空值解决缓存穿透
			if errors.Is(err, model.ErrNotFound) {
				ok, err := l.svcCtx.ReviewRedis.ZaddnxCtx(l.ctx, rediskey.ReviewKey(in.GetBaseId()), 0, "")
				if err != nil {
					return nil, err
				}
				if ok {
					err := l.svcCtx.ReviewRedis.Expire(rediskey.ReviewKey(in.GetBaseId()), 3)
					if err != nil {
						return nil, err
					}
					return &__.ReviewGetResp{}, nil
				}
			}
			return nil, err
		}
		//点赞服务获取点赞信息，这先临时忽略
		for i := 0; i < len(*ans); i++ {
			key, err := json.Marshal((*ans)[i])
			if err != nil {
				return nil, err
			}
			pair := redis.Pair{
				Key:   string(key),
				Score: (*ans)[i].CreateAt,
			}
			redisData = append(redisData, pair)
		}
		_, err = l.svcCtx.ReviewRedis.Zadds(rediskey.ReviewKey(in.GetBaseId()), redisData...)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ReviewRedis.Expire(rediskey.ReviewKey(in.GetBaseId()), 10)
		if err != nil {
			return nil, err
		}
		redisData = redisData[min(max(in.GetOffset(), 0), int64(len(redisData)-1)+1):min(len(redisData), int(in.GetOffset()+in.GetLimit()))]
	}
	ans := &__.ReviewGetResp{Review: make([]*__.Review, 0)}
	for i := 0; i < len(redisData); i++ {
		var data model.Review
		err = json.Unmarshal([]byte(redisData[i].Key), &data)
		if err != nil {
			return nil, err
		}
		ans.Review = append(ans.Review, __.ModelReviewToReview(&data))
	}
	return ans, nil
}
