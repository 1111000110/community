package chat

import (
	"community.com/chat/api/internal/svc"
	"community.com/chat/api/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"log" // 用于日志记录
	"time"
)

type ChatClientLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	conn   *websocket.Conn    // WebSocket 连接
	Send   chan types.Message // 发送消息的缓冲通道
	UserId int64
}

func (c *ChatClientLogic) GetClientId() int64 {
	return c.UserId
}
func (c *ChatClientLogic) GetSendBuffer() chan types.Message {
	return c.Send
}

// readPump 从 WebSocket 连接读取消息并发送到 Hub
func (c *ChatClientLogic) ReadPump() {
	defer func() {
		// 当读取结束时，从 Hub 中注销客户端，并关闭连接
		c.svcCtx.Hub.Unregister <- c
		c.conn.Close()
	}()
	// 循环读取消息
	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			// 如果读取失败，记录错误并退出
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Println(time.Now().Unix())                     //zhangxuan 1
		err = c.svcCtx.Pusher.Push(c.ctx, string(message)) //读取消息立马发kafka
		if err != nil {
			log.Println(err.Error())
		}
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}

// writePump 从 Hub 接收消息并发送到 WebSocket 连接
func (c *ChatClientLogic) WritePump() {
	// 创建一个定时器，用于定期发送 Ping 消息
	ticker := time.NewTicker(time.Duration(c.svcCtx.Config.WebSocket.PingPeriod) * time.Second)
	defer func() {
		// 停止定时器并关闭连接
		ticker.Stop()
		c.conn.Close()
	}()

	// 循环发送消息
	for {
		select {
		case message, ok := <-c.Send:
			// 设置写入超时时间
			c.conn.SetWriteDeadline(time.Now().Add(time.Duration(c.svcCtx.Config.WebSocket.WriteWait) * time.Second))
			if !ok {
				// 如果 Hub 关闭了发送通道，发送关闭消息并退出
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 获取一个写入器，并写入消息
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			n := len(c.Send)
			msg := make([]types.Message, n+1)
			msg[0] = message
			for i := 0; i < n; i++ {
				msg[i+1] = <-c.Send
			}
			data, err := json.Marshal(msg)
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			} else {
				w.Write(data)
			}
			// 关闭写入器
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// 定期发送 Ping 消息
			c.conn.SetWriteDeadline(time.Now().Add(time.Duration(c.svcCtx.Config.WebSocket.PingPeriod) * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func NewChatClientLogic(ctx context.Context, svcCtx *svc.ServiceContext, conn *websocket.Conn, userId int64) *ChatClientLogic {
	conn.SetReadLimit(svcCtx.Config.WebSocket.MaxMessageSize)
	err := conn.SetReadDeadline(time.Now().Add(time.Duration(svcCtx.Config.WebSocket.PongWait) * time.Second))
	if err != nil {
		log.Fatalf("SetReadDeadline failed: %s", err.Error())
	}
	conn.SetPongHandler(func(string) error {
		// 当收到 Pong 消息时，重置读取超时时间
		err = conn.SetReadDeadline(time.Now().Add(time.Duration(svcCtx.Config.WebSocket.PongWait) * time.Second))
		if err != nil {
			log.Fatalf("SetReadDeadline failed: %s", err.Error())
		}
		return nil
	})
	return &ChatClientLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		conn:   conn,
		svcCtx: svcCtx,
		Send:   make(chan types.Message, svcCtx.Config.WebSocket.BufSize),
		UserId: userId,
	}
}
func (l *ChatClientLogic) ChatClient() error {
	go l.WritePump()
	go l.ReadPump()
	return nil
}
