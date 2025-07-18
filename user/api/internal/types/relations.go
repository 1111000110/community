// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3

package types

type UserRelationsGetReq struct {
	UserId  int64 `json:"user_id"`   // 查询发起用户ID
	OUserId int64 `json:"o_user_id"` // 查询目标用户ID
}

type UserRelationsGetResp struct {
	Relations Relations `json:"relations"` // 用户关系详细信息，如果无关系则为空
}

type UserRelationsUpdateReq struct {
	UserId           int64 `json:"user_id"`   // 发起操作的用户ID
	OUserId          int64 `json:"o_user_id"` // 目标用户ID
	RelationshipType int64 `json:"relations"` // 关系类型：1-好友 2-关注 3-拉黑 4-屏蔽 0-删除关系
}

type UserRelationsUpdateResp struct {
	Relations Relations `json:"relations"` // 更新后的关系详细信息
}
