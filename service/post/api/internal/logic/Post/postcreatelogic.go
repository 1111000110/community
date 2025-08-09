package Post

import (
	"community.com/pkg/tool"
	"community.com/service/post/rpc/postservice"
	"context"
	"time"

	"community.com/service/post/api/internal/svc"
	"community.com/service/post/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostCreateLogic {
	return &PostCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostCreateLogic) PostCreate(req *types.PostCreateReq) (resp *types.PostCreateResp, err error) {
	userId, err := tool.GetUserId(l.ctx)
	postCreateResp, err := l.svcCtx.RpcClient.PostClient.PostCreate(l.ctx, &postservice.PostCreateReq{
		Post: &postservice.Post{
			PostId:     0,
			UserId:     userId,
			Title:      req.Title,
			Content:    req.Content,
			Images:     req.Images,
			Theme:      req.Theme,
			Tags:       req.Tags,
			Status:     req.Status,
			CreateTime: time.Now().Unix(),
			UpdateTime: 0,
		},
	})
	if err != nil {
		return nil, err
	}
	resp.PostId = postCreateResp.PostId
	return
}
