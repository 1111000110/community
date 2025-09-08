package logic

import (
	model "community.com/service/post/model/mongo/post"
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
	data, err := l.svcCtx.Model.PostMongoClient.FindAllByPostIds(l.ctx, in.PostIds)
	if err != nil {
		return nil, err
	}
	return &__.PostListResp{
		Posts: model.PostsToRpcPosts(data),
	}, nil
}
