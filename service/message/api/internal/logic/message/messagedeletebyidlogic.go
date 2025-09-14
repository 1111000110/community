package message

import (
	"context"

	"community/service/message/api/internal/svc"
	"community/service/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageDeleteByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageDeleteByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageDeleteByIdLogic {
	return &MessageDeleteByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageDeleteByIdLogic) MessageDeleteById(req *types.MessageDeleteByIdReq) (resp *types.MessageDeleteByIdResp, err error) {
	// todo: add your logic here and delete this line

	return
}
