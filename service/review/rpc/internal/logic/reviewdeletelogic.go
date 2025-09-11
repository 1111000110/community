package logic

import (
	"context"

	"community/infra/review/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReviewDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReviewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewDeleteLogic {
	return &ReviewDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReviewDeleteLogic) ReviewDelete(in *pb.ReviewDeleteReq) (*pb.ReviewDeleteResp, error) {
	// todo: add your logic here and delete this line
	// deleteIds := make(map[int][]int64)
	// rootDeleteIds := make(map[int][]int64)
	// rootIds := make([]int64, 0)
	// var Count int64 = 0
	// for i := 0; i < len(in.GetReview()); i++ {
	// 	if in.GetReview()[i].GetReviewId() == 0 || in.GetReview()[i].GetHeadId() == 0 {
	// 		continue
	// 	}
	// 	if in.GetReview()[i].GetRootId() == 0 {
	// 		rootDeleteIds[int(in.GetReview()[i].GetHeadId()%int64(len(l.svcCtx.ReviewModel)))] = append(rootDeleteIds[int(in.GetReview()[i].GetHeadId()%int64(len(l.svcCtx.ReviewModel)))], in.GetReview()[i].GetReviewId())
	// 	} else {
	// 		deleteIds[int(in.GetReview()[i].GetHeadId()%int64(len(l.svcCtx.ReviewModel)))] = append(deleteIds[int(in.GetReview()[i].GetHeadId()%int64(len(l.svcCtx.ReviewModel)))], in.GetReview()[i].GetReviewId())
	// 		rootIds = append(rootIds, in.GetReview()[i].GetRootId())
	// 	}
	// }
	// for key, val := range deleteIds {
	// 	count, err := l.svcCtx.ReviewModel[key].DeleteByIds(l.ctx, val)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	Count += count
	// }
	// for key, val := range rootDeleteIds {
	// 	count, err := l.svcCtx.ReviewModel[key].DeleteByRootIds(l.ctx, val)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	Count += count
	// }
	// for i := 0; i < len(rootIds); i++ {
	// 	_, err := l.svcCtx.ReviewRedis.Del(rediskey.ReviewKey(rootIds[i]))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// return &__.ReviewDeleteResp{
	// 	Count: Count,
	// }, nil
	return nil, nil
}
