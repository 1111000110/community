package logic

import (
	"community.com/pkg/snowflakes"
	model "community.com/service/post/model/mongo/post"
	"context"
	"log"

	"community.com/service/post/rpc/internal/svc"
	"community.com/service/post/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostCreateLogic {
	return &PostCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostCreateLogic) PostCreate(in *__.PostCreateReq) (*__.PostCreateResp, error) {
	sf, err := snowflakes.NewSnowflake(1, 1)
	if err != nil {
		log.Fatal(err)
	}
	postId, err := sf.NextID()
	if err != nil {
		log.Fatal(err)
	}
	in.GetPost().PostId = postId
	err = l.svcCtx.PostMongoClient.Insert(l.ctx, model.PostToModelPost(in.GetPost()))
	return &__.PostCreateResp{
		PostId: postId,
	}, err
}
