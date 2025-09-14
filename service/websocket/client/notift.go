package client

import (
	"community/pkg/xstring"
	"fmt"
	"strings"
)

const SplitByte = ","

func GetRedisKeyName(key string) string {
	return fmt.Sprintf("websocket_client:{%s}", key)
}

func GetIdsByKey(key string) []int64 {
	if key == "" {
		return []int64{}
	}
	info := strings.Split(key, SplitByte)
	resp := make([]int64, 0)
	for _, d := range info {
		resp = append(resp, xstring.StringToIntOrZero[int64](d))
	}
	return resp
}

func GetKeyByIds(ids []int64) string {
	resp := ""
	for _, id := range ids {
		resp += xstring.IntToString(id)
	}
	return resp
}
