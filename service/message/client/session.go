package client

import (
	"community/xuan/xstring"
	"strings"
)

const (
	ChatGroup = "chatGroup" // 群聊
	Private   = "private"   // 私聊
)
const SplitByte = "-" // 分隔符

func GetSessionIdByPrivateIds(privateIds []int64) string {
	resp := Private
	for _, id := range privateIds {
		resp += SplitByte
		resp += xstring.IntToString(id)
	}
	return resp
}

// GetPrivateIdsBySessionIds 从sessionIds解析出用私聊id列表，去掉deleteIds的内容
func GetPrivateIdsBySessionIds(sessionIds string, deleteIds ...int64) []int64 {
	resp := make([]int64, 0)
	deleteMap := make(map[int64]struct{})
	for _, id := range deleteIds {
		deleteMap[id] = struct{}{}
	}
	for i, id := range strings.Split(sessionIds, SplitByte) {
		if i > 0 {
			privateId := xstring.StringToIntOrZero[int64](id)
			if _, ok := deleteMap[privateId]; !ok {
				resp = append(resp, privateId)
			}
		}
	}
	return resp
}

func GetSessionIdByChatGroupId(ChatGroupId int64) string {
	resp := ChatGroup + xstring.IntToString(ChatGroupId)
	return resp
}
