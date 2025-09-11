package Public

import (
	"community/service/post/rpc/postservice"
	"community/service/user/rpc/client/userservice"
	"context"

	"community/service/post/api/internal/svc"
	"community/service/post/api/internal/types"

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
	data, err := l.svcCtx.RpcClient.PostClient.PostDetail(l.ctx, &postservice.PostDetailReq{
		PostId: req.PostId,
	})
	if err != nil {
		return resp, err
	}
	userData, err := l.svcCtx.RpcClient.UserClient.UserQuery(l.ctx, &userservice.UserQueryReq{
		UserId: data.GetPost().GetUserId(),
		Type:   userservice.GetPublicInfo,
	})
	if err != nil {
		return nil, err
	}
	return &types.PostDetailResp{
		Post: types.PostDetail{
			PostBase: types.PostBase{
				PostId:     data.GetPost().GetPostId(),
				UserId:     data.GetPost().GetUserId(),
				Title:      data.GetPost().GetTitle(),
				Content:    data.GetPost().GetContent(),
				Images:     data.GetPost().GetImages(),
				Theme:      data.GetPost().GetTheme(),
				Tags:       data.GetPost().GetTags(),
				Status:     data.GetPost().GetStatus(),
				CreateTime: data.GetPost().GetCreateTime(),
				UpdateTime: data.GetPost().GetUpdateTime(),
			},
			Author: types.PostAuthor{
				UserId:   userData.GetUser().GetUserId(),
				NickName: userData.GetUser().GetNickname(),
				Avatar:   userData.GetUser().GetAvatar(),
			},
		},
	}, nil
}
