package hub

import (
	"encoding/json"
)

const (
	ios = iota
	web
)

func GetConnType(connType string) int64 {
	switch connType {
	case "ios":
		return ios
	case "web":
		return web
	default:
		return web
	}
}

const ConnTypeCount = 2 // 客户端和web

type Notification struct { // 发送的消息
	ConnId int64           `json:"conn_id"`
	Type   string          `json:"type"` // 消息类型标识（如"message"、"group_status"、"settings"）
	Data   json.RawMessage `json:"data"` // 具体内容（使用json.RawMessage支持动态解析）
}

type Client interface {
	GetType() int64                   // 获取客户端类型
	GetClientId() int64               // 获取连接id
	WritePump()                       // 写如函数
	ReadPump()                        // 读取函数
	GetSendBuffer() chan Notification // 自己的消息队列
}
