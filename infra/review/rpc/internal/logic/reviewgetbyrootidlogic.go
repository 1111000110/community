package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"math"
	model "community/review/model"
	rediskey "community/review/model/redisKey"

	"community/review/rpc/generateFileV1/internal/svc"
	"community/review/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReviewGetByRootIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReviewGetByRootIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewGetByRootIdLogic {
	return &ReviewGetByRootIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReviewGetByRootIdLogic) ReviewGetByRootId(in *__.ReviewGetByRootIdReq) (*__.ReviewGetResp, error) {
	// todo: add your logic here and delete this line
	if in.GetRootId() == 0 || in.GetHeadId() == 0 {
		return &__.ReviewGetResp{}, errors.New("RootId or HeadId is null")
	}
	redisData := []redis.Pair{}
	redisData, err := l.svcCtx.ReviewRedis.ZrangebyscoreWithScoresAndLimitCtx(l.ctx, rediskey.ReviewKey(in.GetRootId()), -1, math.MaxInt64, int(in.GetOffset()), int(in.GetLimit()))
	if err != nil {
		return nil, err
	}
	if len(redisData) == 1 && redisData[0].Key == "" && redisData[0].Score == 0 {
		return &__.ReviewGetResp{}, nil
	}
	if len(redisData) == 0 {
		var ans *[]model.Review
		ans, err = l.svcCtx.ReviewModel[int(in.GetHeadId()%(int64(len(l.svcCtx.ReviewModel))))].FindByRootId(l.ctx, in.GetRootId())
		if err != nil || ans == nil { //空值解决缓存穿透
			if errors.Is(err, model.ErrNotFound) {
				ok, err := l.svcCtx.ReviewRedis.ZaddnxCtx(l.ctx, rediskey.ReviewKey(in.GetRootId()), 0, "")
				if err != nil {
					return nil, err
				}
				if ok {
					err := l.svcCtx.ReviewRedis.Expire(rediskey.ReviewKey(in.GetRootId()), 3)
					if err != nil {
						return nil, err
					}
					return &__.ReviewGetResp{}, nil
				}
			}
			return nil, err
		}
		_, err := l.svcCtx.ReviewModel[int(in.GetHeadId()%(int64(len(l.svcCtx.ReviewModel))))].UpdataRootcommitCountByReviewId(l.ctx, in.GetRootId(), int64(len(*ans)))
		if err != nil {
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
		_, err = l.svcCtx.ReviewRedis.Zadds(rediskey.ReviewKey(in.GetRootId()), redisData...)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ReviewRedis.Expire(rediskey.ReviewKey(in.GetRootId()), 10)
		if err != nil {
			return nil, err
		}
		redisData = redisData[min(max(in.GetOffset(), 0), int64(len(redisData)-1)+1):min(len(redisData), int(in.GetOffset()+in.GetLimit()))]
	}
	ans := &__.ReviewGetResp{Review: make([]*__.Review, 0)}
	for i := 0; i < len(redisData); i++ {
		var data __.Review
		err = json.Unmarshal([]byte(redisData[i].Key), &data)
		if err != nil {
			return nil, err
		}
		ans.Review = append(ans.Review, &data)
	}
	return ans, nil
}
