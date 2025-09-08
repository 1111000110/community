package logic

import (
	model "community.com/service/post/model/mongo/post"
	"context"

	"community.com/service/post/rpc/internal/svc"
	"community.com/service/post/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostDetailLogic {
	return &PostDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostDetailLogic) PostDetail(in *__.PostDetailReq) (*__.PostDetailResp, error) {
	post, err := l.svcCtx.Model.PostMongoClient.FindOneByPostId(l.ctx, in.GetPostId())
	return &__.PostDetailResp{
		Post: model.PostToRpcPost(post),
	}, err
}
