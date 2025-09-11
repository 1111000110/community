package relations

import (
	"context"

	"community/service/user/api/internal/svc"
	"community/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRelationsGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRelationsGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRelationsGetLogic {
	return &UserRelationsGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRelationsGetLogic) UserRelationsGet(req *types.UserRelationsGetReq) (resp *types.UserRelationsGetResp, err error) {
	// todo: add your logic here and delete this line

	return
}
