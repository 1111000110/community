package logic

import (
	"context"

	"community/service/post/rpc/internal/svc"
	"community/service/post/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostUpdateLogic {
	return &PostUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostUpdateLogic) PostUpdate(in *__.PostUpdateReq) (*__.PostUpdateResp, error) {
	// todo: add your logic here and delete this line

	return &__.PostUpdateResp{}, nil
}
