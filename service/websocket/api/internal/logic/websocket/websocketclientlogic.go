package websocket

import (
	"community/pkg/snowflakes"
	"community/pkg/xstring"
	"community/service/websocket/client"
	"context"
	"encoding/json"
	"sync/atomic"
	"time"

	"community/service/websocket/api/internal/hub"
	"community/service/websocket/api/internal/svc"
	"community/service/websocket/api/internal/types"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type WebSocketClientLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	cancel       context.CancelFunc // 取消内部上下文（控制读写协程退出）
	userId       int64              // 客户端唯一标识
	connType     int64              // 连接类型
	sendBuffer   chan string        // 消息发送缓冲区
	bufferClosed atomic.Bool        // 标记sendBuffer是否已关闭（防止重复close panic）
	conn         *websocket.Conn    // WebSocket连接实例
}

func NewWebSocketClientLogic(_ context.Context, svcCtx *svc.ServiceContext, req *types.WebSocketReq, conn *websocket.Conn) (*WebSocketClientLogic, error) {
	// 1. 获取用户ID
	//userId, err := tool.GetUserId(ctx) todo
	snowflake, err := snowflakes.NewSnowflake(1, 1)
	if err != nil {
		return nil, err
	}
	userId, err := snowflake.NextID()
	if err != nil {
		logx.Errorf("failed to get user ID: %v", err)
		return nil, err
	}

	// 2. 创建可取消上下文（控制读写协程生命周期）
	innerCtx, cancel := context.WithCancel(context.Background())

	// 3. 初始化WebSocket连接参数（读限制、读超时）
	conn.SetReadLimit(svcCtx.Config.WebSocket.MaxMessageSize)
	if err = conn.SetReadDeadline(time.Now().Add(time.Duration(svcCtx.Config.WebSocket.PongWait) * time.Second)); err != nil {
		logx.Errorf("failed to set read deadline (userId: %d): %v", userId, err)
		cancel() // 失败时取消上下文，避免资源泄漏
		return nil, err
	}

	// 4. 返回客户端实例
	resp := &WebSocketClientLogic{
		Logger:       logx.WithContext(innerCtx),
		ctx:          innerCtx,
		cancel:       cancel,
		conn:         conn,
		svcCtx:       svcCtx,
		sendBuffer:   make(chan string, svcCtx.Config.WebSocket.BufSize), // 带缓冲通道，避免阻塞
		userId:       userId,
		connType:     hub.GetConnType(req.ClientType),
		bufferClosed: atomic.Bool{}, // 初始为false（未关闭）
	}
	// 5. 注册Pong回调（重置读超时，维持连接）
	resp.conn.SetPongHandler(func(_ string) error {
		if err = resp.conn.SetReadDeadline(time.Now().Add(time.Duration(resp.svcCtx.Config.WebSocket.PongWait) * time.Second)); err != nil {
			logx.Errorf("failed to reset read deadline in PongHandler (userId: %d): %v", userId, err)
			return err
		}
		err = resp.ResetRedis() // 重刷redis
		return err
	})
	return resp, nil
}

func (l *WebSocketClientLogic) GetClientId() int64 {
	return l.userId
}

func (l *WebSocketClientLogic) GetSendBuffer() chan string {
	return l.sendBuffer
}

func (l *WebSocketClientLogic) GetType() int64 {
	return l.connType
}

func (l *WebSocketClientLogic) ResetRedis() error {
	err := l.svcCtx.Model.RedisClient.Setex(
		client.GetRedisKeyName(xstring.IntToString(l.userId)),
		"node1", // todo 修改机器
		int(hub.Timeout.Seconds()),
	)
	if err != nil {
		logx.Errorf("reset redis failed: %v", err)
	}
	return err
}

func (l *WebSocketClientLogic) ReadPump() {
	for {
		select {
		case <-l.ctx.Done(): // 上下文取消（如Close()被调用）
			logx.Infof("ReadPump exited (ctx done), userId: %d,closeErr: %s ", l.userId, l.ctx.Err())
			return
		default:
			// 读取WebSocket消息（文本/二进制）
			_, message, err := l.conn.ReadMessage()
			if err != nil {
				// 判断是否为「致命错误」（无法恢复的连接问题）
				if isFatalReadErr(err) {
					logx.Errorf("fatal read error (userId: %d): %v", l.userId, err)
					// 调用Hub移除客户端，触发统一注销流程
					l.svcCtx.MessageHub.RemoveClient(l)
					return
				}
				// 非致命错误（如临时网络波动），短暂重试
				logx.Errorf("temporary read error (userId: %d): %v, retry after 100ms", l.userId, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			err = l.svcCtx.Model.KafkaMessageClient.PushWithKey(l.ctx, "", string(message))
			if err != nil {
				logx.Errorf(err.Error())
			}
		}
	}
}

func (l *WebSocketClientLogic) WritePump() {
	// 定时发送Ping消息（维持连接）
	pingTicker := time.NewTicker(time.Duration(l.svcCtx.Config.WebSocket.PingPeriod) * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		// 1. 从sendBuffer接收消息（Hub转发的消息）
		case msg, ok := <-l.sendBuffer:
			// 设置写入超时（避免阻塞）
			if err := l.conn.SetWriteDeadline(time.Now().Add(time.Duration(l.svcCtx.Config.WebSocket.WriteWait) * time.Second)); err != nil {
				logx.Errorf("failed to set write deadline (userId: %d): %v", l.userId, err)
				l.svcCtx.MessageHub.RemoveClient(l)
				return
			}

			// sendBuffer已关闭（致命错误）
			if !ok {
				logx.Infof("sendBuffer closed (userId: %d), sending close frame", l.userId)
				sendCloseFrame(l.conn, websocket.CloseNormalClosure, "sendBuffer closed")
				l.svcCtx.MessageHub.RemoveClient(l)
				return
			}

			// 获取WebSocket写入器（批量发送优化）
			writer, err := l.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logx.Errorf("failed to get next writer (userId: %d): %v", l.userId, err)
				l.svcCtx.MessageHub.RemoveClient(l)
				return
			}

			// 批量读取sendBuffer中所有消息（减少写入次数，提升性能）
			msgList := []string{msg}
			batchReadMessages(l.sendBuffer, &msgList)

			// 序列化消息（JSON格式）
			data, err := json.Marshal(msgList)
			if err != nil {
				logx.Errorf("failed to marshal message (userId: %d): %v", l.userId, err)
				_ = writer.Close() // 确保写入器关闭
				continue
			}

			// 写入消息到WebSocket
			if _, err := writer.Write(data); err != nil {
				logx.Errorf("failed to write message (userId: %d): %v", l.userId, err)
				_ = writer.Close()
				// 写入失败为致命错误，触发注销
				if isFatalWriteErr(err) {
					l.svcCtx.MessageHub.RemoveClient(l)
					return
				}
				continue
			}

			// 关闭写入器（完成一次批量发送）
			if err := writer.Close(); err != nil {
				logx.Errorf("failed to close writer (userId: %d): %v", l.userId, err)
				l.svcCtx.MessageHub.RemoveClient(l)
				return
			}

		// 2. 定时发送Ping消息（维持连接）
		case <-pingTicker.C:
			if err := l.conn.SetWriteDeadline(time.Now().Add(time.Duration(l.svcCtx.Config.WebSocket.PingPeriod) * time.Second)); err != nil {
				logx.Errorf("failed to set ping deadline (userId: %d): %v", l.userId, err)
				l.svcCtx.MessageHub.RemoveClient(l)
				return
			}
			// 发送Ping帧（空内容）
			if err := l.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logx.Errorf("failed to send ping (userId: %d): %v", l.userId, err)
				l.svcCtx.MessageHub.RemoveClient(l)
				return
			}

		// 3. 上下文取消（如Close()被调用）
		case <-l.ctx.Done():
			logx.Infof("WritePump exited (ctx done), userId: %d, closeErr: %s ", l.userId, l.ctx.Err())
			return
		}
	}
}

func (l *WebSocketClientLogic) Start() error {
	logx.Infof("starting websocket client (userId: %d, connType: %d)", l.userId, l.connType)
	go l.ReadPump()  // 启动读协程
	go l.WritePump() // 启动写协程
	return l.ResetRedis()
}

func (l *WebSocketClientLogic) Close() {
	logx.Infof("closing websocket client (userId: %d)", l.userId)
	// 1. 取消上下文（终止读写协程，幂等操作）
	l.cancel()

	// 2. 关闭sendBuffer（原子变量确保仅关闭一次，避免panic）
	if l.bufferClosed.CompareAndSwap(false, true) {
		close(l.sendBuffer)
		logx.Infof("sendBuffer closed (userId: %d)", l.userId)
	}

	// 3. 关闭WebSocket连接（gorilla/websocket.Conn.Close()幂等）
	if l.conn != nil {
		sendCloseFrame(l.conn, websocket.CloseNormalClosure, "client closed by hub")
		if err := l.conn.Close(); err != nil {
			logx.Errorf("failed to close websocket conn (userId: %d): %v", l.userId, err)
		} else {
			logx.Infof("websocket conn closed (userId: %d)", l.userId)
		}
	}
}

// isFatalReadErr 判断读错误是否为致命错误
func isFatalReadErr(err error) bool {
	return websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) ||
		websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway)
}

// isFatalWriteErr 判断写错误是否为致命错误
func isFatalWriteErr(err error) bool {
	return websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
}

// sendCloseFrame 发送WebSocket标准关闭帧
func sendCloseFrame(conn *websocket.Conn, code int, reason string) {
	closeMsg := websocket.FormatCloseMessage(code, reason)
	if err := conn.WriteMessage(websocket.CloseMessage, closeMsg); err != nil {
		logx.Errorf("failed to send close frame: %v (reason: %s)", err, reason)
	}
}

// batchReadMessages 批量读取sendBuffer中的消息（避免通道关闭时读取零值）
func batchReadMessages(buf chan string, msgList *[]string) {
	for {
		select {
		case msg, ok := <-buf:
			if !ok {
				return
			}
			*msgList = append(*msgList, msg)
		default:
			return
		}
	}
}
