package Post

import (
	"context"

	"community.com/service/post/api/internal/svc"
	"community.com/service/post/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostUpdateLogic {
	return &PostUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostUpdateLogic) PostUpdate(req *types.PostUpdateReq) (resp *types.PostUpdateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
