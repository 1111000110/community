package message

import (
	"context"

	"community/service/message/api/internal/svc"
	"community/service/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageUpdateByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageUpdateByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageUpdateByIdLogic {
	return &MessageUpdateByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageUpdateByIdLogic) MessageUpdateById(req *types.MessageUpdateByIdReq) (resp *types.MessageUpdateByIdResp, err error) {
	// todo: add your logic here and delete this line

	return
}
