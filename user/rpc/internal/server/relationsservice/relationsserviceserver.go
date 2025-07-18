// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3
// Source: user.proto

package server

import (
	"context"

	"community.com/user/rpc/internal/logic/relationsservice"
	"community.com/user/rpc/internal/svc"
	"community.com/user/rpc/pb"
)

type RelationsServiceServer struct {
	svcCtx *svc.ServiceContext
	__.UnimplementedRelationsServiceServer
}

func NewRelationsServiceServer(svcCtx *svc.ServiceContext) *RelationsServiceServer {
	return &RelationsServiceServer{
		svcCtx: svcCtx,
	}
}

// UserRelationsUpdate 用户关系更新
func (s *RelationsServiceServer) UserRelationsUpdate(ctx context.Context, in *__.UserRelationsUpdateReq) (*__.UserRelationsUpdateResp, error) {
	l := relationsservicelogic.NewUserRelationsUpdateLogic(ctx, s.svcCtx)
	return l.UserRelationsUpdate(in)
}

// UserRelationsGet 用户关系查询
func (s *RelationsServiceServer) UserRelationsGet(ctx context.Context, in *__.UserRelationsGetReq) (*__.UserRelationsGetResp, error) {
	l := relationsservicelogic.NewUserRelationsGetLogic(ctx, s.svcCtx)
	return l.UserRelationsGet(in)
}
