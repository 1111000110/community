package logic

import (
	"context"

	"community.com/service/message/rpc/internal/svc"
	"community.com/service/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMessageByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMessageByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMessageByIdLogic {
	return &UpdateMessageByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMessageByIdLogic) UpdateMessageById(in *__.UpdateMessageByIdReq) (*__.UpdateMessageByIdResp, error) {
	// todo: add your logic here and delete this line

	return &__.UpdateMessageByIdResp{}, nil
}
