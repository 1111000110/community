package websocket

import (
	"community.com/pkg/tool"
	"community.com/service/websocket/api/internal/hub"
	"community.com/service/websocket/api/internal/svc"
	"community.com/service/websocket/api/internal/types"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"time"
)

type WebSocketClientLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	userId     int64
	sendBuffer chan hub.Notification
	conn       *websocket.Conn
	connType   int64
}

func NewWebSocketClientLogic(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WebSocketReq, conn *websocket.Conn) *WebSocketClientLogic {
	userId, err := tool.GetUserId(ctx)
	if err != nil {
		return nil
	}
	conn.SetReadLimit(svcCtx.Config.WebSocket.MaxMessageSize)
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(svcCtx.Config.WebSocket.PongWait) * time.Second))
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
	return &WebSocketClientLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		conn:       conn,
		svcCtx:     svcCtx,
		sendBuffer: make(chan hub.Notification, svcCtx.Config.WebSocket.BufSize),
		userId:     userId,
		connType:   hub.GetConnType(req.ClientType),
	}
}

func (l *WebSocketClientLogic) GetClientId() int64 {
	return l.userId
}
func (l *WebSocketClientLogic) GetSendBuffer() chan hub.Notification {
	return l.sendBuffer
}
func (l *WebSocketClientLogic) GetType() int64 {
	return l.connType
}

// readPump 从 WebSocket 连接读取消息并发送到 Hub
func (l *WebSocketClientLogic) ReadPump() {
	defer func() {
		l.svcCtx.MessageHub.RemoveClient(l)
		func() {
			err := l.conn.Close()
			if err != nil {
				logx.Errorf("close websocket client conn err: %s", err.Error())
			}
		}()
		close(l.sendBuffer)
	}()
	// 循环读取消息
	for {
		_, message, err := l.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logx.Errorf("read message err: %s", err.Error())
			}
			break
		}
		logx.Info(string(message))
		if err = l.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			logx.Errorf("write message err: %s", err.Error())
		}
	}
}

// writePump 从 Hub 接收消息并发送到 WebSocket 连接
func (l *WebSocketClientLogic) WritePump() {
	// 创建一个定时器，用于定期发送 Ping 消息
	ticker := time.NewTicker(time.Duration(l.svcCtx.Config.WebSocket.PingPeriod) * time.Second)
	defer func() {
		// 停止定时器并关闭连接
		ticker.Stop()
		func() {
			err := l.conn.Close()
			if err != nil {
				logx.Errorf("close websocket client conn err: %s", err.Error())
			}
		}()
	}()

	// 循环发送消息
	for {
		select {
		case message, ok := <-l.sendBuffer:
			// 设置写入超时时间
			err := l.conn.SetWriteDeadline(time.Now().Add(time.Duration(l.svcCtx.Config.WebSocket.WriteWait) * time.Second))
			if err != nil {
				logx.Errorf("set write deadline err: %s", err.Error())
			}
			if !ok {
				// 如果 Hub 关闭了发送通道，发送关闭消息并退出
				err = l.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logx.Errorf("close websocket client conn err: %s", err.Error())
				}
				return
			}

			// 获取一个写入器，并写入消息
			w, err := l.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			n := len(l.sendBuffer)
			msg := make([]hub.Notification, n+1)
			msg[0] = message
			for i := 0; i < n; i++ {
				msg[i+1] = <-l.sendBuffer
			}
			data, err := json.Marshal(msg)
			if err != nil {
				logx.Errorf("json marshal err: %s", err.Error())
			}
			if _, err = w.Write(data); err != nil {
				logx.Errorf("write message err: %s", err.Error())
			}
			// 关闭写入器
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// 定期发送 Ping 消息
			err := l.conn.SetWriteDeadline(time.Now().Add(time.Duration(l.svcCtx.Config.WebSocket.PingPeriod) * time.Second))
			if err != nil {
				logx.Errorf("set write deadline err: %s", err.Error())
			}
			if err := l.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-l.ctx.Done():
			return
		}
	}
}
func (l *WebSocketClientLogic) RunSocketClient() error {
	go l.WritePump()
	go l.ReadPump()
	return nil
}
