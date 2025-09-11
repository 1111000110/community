package logic

import (
	"context"

	"community/infra/review/rpc/internal/svc"
	"community/infra/review/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReviewDeleteByHeadIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReviewDeleteByHeadIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewDeleteByHeadIdLogic {
	return &ReviewDeleteByHeadIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReviewDeleteByHeadIdLogic) ReviewDeleteByHeadId(in *__.ReviewDeleteByHeadIdReq) (*__.ReviewDeleteResp, error) {
	// todo: add your logic here and delete this line
	deleteIds := make(map[int64][]int64)
	for i := 0; i < len(in.GetHeadId()); i++ {
		deleteIds[in.GetHeadId()[i]%int64(len(l.svcCtx.ReviewModel))] = append(deleteIds[in.GetHeadId()[i]%int64(len(l.svcCtx.ReviewModel))], in.GetHeadId()[i])
	}
	var Count int64 = 0
	for key, val := range deleteIds {
		count, err := l.svcCtx.ReviewModel[key].DeleteByHeadIds(l.ctx, val)
		if err != nil {
			return nil, err
		}
		Count += count
	}
	return &__.ReviewDeleteResp{
		Count: Count,
	}, nil
}
