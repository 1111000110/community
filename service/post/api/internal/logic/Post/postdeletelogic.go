package Post

import (
	"context"

	"community.com/service/post/api/internal/svc"
	"community.com/service/post/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostDeleteLogic {
	return &PostDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostDeleteLogic) PostDelete(req *types.PostDeleteReq) (resp *types.PostDeleteResp, err error) {
	// todo: add your logic here and delete this line

	return
}
