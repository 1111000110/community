package Public

import (
	"context"

	"community.com/service/post/api/internal/svc"
	"community.com/service/post/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostDetailLogic {
	return &PostDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostDetailLogic) PostDetail(req *types.PostDetailReq) (resp *types.PostDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
