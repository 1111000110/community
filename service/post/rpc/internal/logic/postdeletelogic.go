package logic

import (
	"context"

	"community.com/service/post/rpc/internal/svc"
	"community.com/service/post/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostDeleteLogic {
	return &PostDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostDeleteLogic) PostDelete(in *__.PostDeleteReq) (*__.PostDeleteResp, error) {
	return &__.PostDeleteResp{}, l.svcCtx.PostMongoClient.DeleteOneByPostId(l.ctx, in.GetPostId())
}
