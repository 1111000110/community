package logic

import (
	"context"

	"community.com/infra/review/rpc/internal/svc"
	__ "community.com/infra/review/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReviewCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReviewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewCreateLogic {
	return &ReviewCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReviewCreateLogic) ReviewCreate(in *__.ReviewCreateReq) (*__.ReviewCreateResp, error) {
	// todo: add your logic here and delete this line
	// var ans int64 = 0
	// var n = int64(len(l.svcCtx.ReviewModel))
	// rootIds := map[int64]int64{}
	// baseIds := map[int64]int64{}
	// snowflake, err := tool.NewSnowflake(1, 1)
	// if err != nil {
	// 	return nil, err
	// }
	// data := make(map[int64]*[]model.Review)
	// for i := 0; i < len(in.Review); i++ {
	// 	if in.Review[i].BaseId == 0 {
	// 		continue
	// 	}
	// 	in.Review[i].ReviewId, err = snowflake.Generate()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if data[in.Review[i].HeadId%n] == nil {
	// 		data[in.Review[i].HeadId%n] = &[]model.Review{*__.ReviewToModelReview(in.Review[i])}
	// 	} else {
	// 		*data[in.Review[i].HeadId%n] = append(*data[in.Review[i].HeadId%n], *__.ReviewToModelReview(in.Review[i]))
	// 	}
	// 	if in.Review[i].RootId == 0 {
	// 		baseIds[in.Review[i].BaseId]++
	// 	} else {
	// 		rootIds[in.Review[i].RootId]++
	// 	}
	// }
	// for i := 0; i < int(n); i++ {
	// 	if _, ok := data[int64(i)]; ok {
	// 		count, err := l.svcCtx.ReviewModel[i].InsertMany(l.ctx, data[int64(i)])
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		ans += count
	// 	}
	// }
	// for key, _ := range rootIds {
	// 	_, err = l.svcCtx.ReviewRedis.Del(rediskey.ReviewKey(key)) //删缓存
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// for key, _ := range baseIds {
	// 	_, err = l.svcCtx.ReviewRedis.Del(rediskey.ReviewKey(key)) //删缓存
	// }
	// return &__.ReviewCreateResp{
	// 	Count: ans,
	// }, nil
	return &__.ReviewCreateResp{
		Count: 0,
	}, nil
}
