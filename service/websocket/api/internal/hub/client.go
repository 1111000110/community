package hub

import (
	"time"
)

const (
	ConnTypeWeb  = 0
	ConnTypeApp  = 1
	ConnTypeMini = 2
)

const (
	Timeout = 30 * time.Second
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

type NotifyList struct {
	Key string // 发送的连接列表
	Val string // 发送的消息
}

func NewNotifyList(key, val string) *NotifyList {
	return &NotifyList{
		Key: key,
		Val: val,
	}
}

type Notify struct {
	ConnId int64
	Val    string
}

func NewNotify(connId int64, val string) *Notify {
	return &Notify{
		ConnId: connId,
		Val:    val,
	}
}

type Client interface {
	GetType() int64             // 获取客户端类型,注册和注销的时候需要去对应的map处理
	GetClientId() int64         // 获取连接id
	GetSendBuffer() chan string // 自己的消息队列
	Close()                     // 关闭
	Start() error               // 启动
}
