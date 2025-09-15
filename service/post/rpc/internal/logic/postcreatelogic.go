package logic

import (
	"community/pkg/snowflakes"
	model "community/service/post/model/mongo/post"
	"context"
	"log"

	"community/service/post/rpc/internal/svc"
	"community/service/post/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	snowflakes *snowflakes.Snowflake
}

func NewPostCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostCreateLogic {
	snowflake, err := snowflakes.NewSnowflake(1, 1)
	if err != nil {
		panic(err)
	}
	return &PostCreateLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		snowflakes: snowflake,
	}
}

func (l *PostCreateLogic) PostCreate(in *__.PostCreateReq) (*__.PostCreateResp, error) {
	postId, err := l.snowflakes.NextID()
	if err != nil {
		log.Fatal(err)
	}
	in.GetPost().PostId = postId
	err = l.svcCtx.ModelClient.MongoPost.Insert(l.ctx, model.RpcPostToPost(in.GetPost()))
	return &__.PostCreateResp{
		PostId: postId,
	}, err
}
