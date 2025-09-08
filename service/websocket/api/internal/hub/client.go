package hub

import (
	"encoding/json"
)

const (
	ConnTypeWeb  = 0
	ConnTypeApp  = 1
	ConnTypeMini = 2
)

func GetConnType(clientType string) int64 {
	switch clientType {
	case "web":
		return ConnTypeWeb
	case "app":
		return ConnTypeApp
	case "mini":
		return ConnTypeMini
	default:
		return ConnTypeWeb
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
	GetSendBuffer() chan Notification // 自己的消息队列
	Close()                           // 关闭
	Start()                           //启动
}
