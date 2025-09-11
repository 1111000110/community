package Public

import (
	"community/service/post/api/internal/svc"
	"community/service/post/api/internal/types"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type PostListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostListLogic {
	return &PostListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostListLogic) PostList(req *types.PostListReq) (resp *types.PostListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
