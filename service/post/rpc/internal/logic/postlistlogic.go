package logic

import (
	"context"

	"community.com/service/post/rpc/internal/svc"
	"community.com/service/post/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostListLogic {
	return &PostListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostListLogic) PostList(in *__.PostListReq) (*__.PostListResp, error) {
	// todo: add your logic here and delete this line

	return &__.PostListResp{}, nil
}
