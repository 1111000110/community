package Post

import (
	"community/pkg/tool"
	"community/service/post/rpc/postservice"
	"context"
	"errors"

	"community/service/post/api/internal/svc"
	"community/service/post/api/internal/types"

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
	userId, err := tool.GetUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	postDetail, err := l.svcCtx.RpcClient.PostClient.PostDetail(l.ctx, &postservice.PostDetailReq{
		PostId: req.PostId,
	})
	if err != nil {
		return nil, err
	}
	if postDetail.GetPost().GetUserId() != userId {
		return nil, errors.New("不可以删除其他人的作品哦~")
	}
	_, err = l.svcCtx.RpcClient.PostClient.PostDelete(l.ctx, &postservice.PostDeleteReq{
		PostId: req.PostId,
	})
	if err != nil {
		return nil, err
	}
	return
}
