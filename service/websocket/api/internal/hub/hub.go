package hub

import (
	"community/service/websocket/client"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	HashCount  = 32    // 节点数量（分片数）
	BufferSize = 10000 // 通道缓冲区大小
)

type Hub struct {
	NodeCount    int64             // 节点数量（分片数）
	HashFunction func(int64) int64 // 哈希函数
	Node         []*Node           // 节点数组（每个节点是一个分片）
	stopped      atomic.Bool       // 停止状态标记
}

type Node struct {
	Clients    []sync.Map   // 按连接类型分组的客户端映射
	Register   chan Client  // 注册通道
	Unregister chan Client  // 注销通道
	Notify     chan *Notify // 消息通道
}

// NewHub 创建一个新的 Hub 实例
func NewHub() *Hub {
	nodeSlice := make([]*Node, HashCount)
	for i := 0; i < HashCount; i++ {
		// 初始化每种连接类型的客户端映射
		clients := make([]sync.Map, ConnTypeCount)
		for j := 0; j < ConnTypeCount; j++ {
			clients[j] = sync.Map{}
		}

		nodeSlice[i] = &Node{
			Clients:    clients,
			Register:   make(chan Client, BufferSize),
			Unregister: make(chan Client, BufferSize),
			Notify:     make(chan *Notify, BufferSize*10),
		}
	}

	return &Hub{
		NodeCount: HashCount,
		Node:      nodeSlice,
		HashFunction: func(i int64) int64 {
			return i % HashCount
		},
	}
}

// Start 启动 Hub 的所有分片节点
func (h *Hub) Start() {
	for i := 0; i < int(h.NodeCount); i++ {
		go func(node *Node, nodeIndex int) {
			defer func() {
				if err := recover(); err != nil {
					logx.Errorf("Hub node %d panic: %v", nodeIndex, err)
				}
			}()

			logx.Infof("Starting hub node %d", nodeIndex)

			for {
				select {
				case client, ok := <-node.Register:
					if !ok {
						logx.Infof("Node %d register channel closed", nodeIndex)
						return
					}
					h.handleRegister(node, client, nodeIndex)

				case client, ok := <-node.Unregister:
					if !ok {
						logx.Infof("Node %d unregister channel closed", nodeIndex)
						return
					}
					h.handleUnregister(node, client, nodeIndex)

				case notify, ok := <-node.Notify:
					if !ok {
						logx.Infof("Node %d message channel closed", nodeIndex)
						return
					}
					h.handleMessage(node, notify, nodeIndex)
				}
			}
		}(h.Node[i], i)
	}
}

// handleRegister 处理客户端注册
func (h *Hub) handleRegister(node *Node, client Client, nodeIndex int) {
	connType := client.GetType()
	if int(connType) < len(node.Clients) {
		node.Clients[connType].Store(client.GetClientId(), client)
		logx.Debugf("Client %d registered in node %d, type %d",
			client.GetClientId(), nodeIndex, connType)
		if err := client.Start(); err != nil {
			client.Close()
		}
	} else {
		logx.Errorf("Invalid connection type %d for client %d", connType, client.GetClientId())
		client.Close()
	}
}

// handleUnregister 处理客户端注销
func (h *Hub) handleUnregister(node *Node, client Client, nodeIndex int) {
	connType := client.GetType()
	if int(connType) < len(node.Clients) {
		node.Clients[connType].Delete(client.GetClientId())
		logx.Debugf("Client %d unregistered from node %d, type %d",
			client.GetClientId(), nodeIndex, connType)
	} else {
		logx.Errorf("Invalid connection type %d for client %d during unregister",
			connType, client.GetClientId())
	}
	client.Close()
}

// handleMessage 处理消息发送
func (h *Hub) handleMessage(node *Node, notify *Notify, nodeIndex int) {
	if notify == nil {
		return
	}
	// 在所有连接类型中查找目标客户端
	for connType := 0; connType < len(node.Clients); connType++ {
		if client, exists := node.Clients[connType].Load(notify.ConnId); exists {
			// 异步发送避免阻塞消息处理循环
			go func(c Client, msg string, nodeIdx int) {
				select {
				case c.GetSendBuffer() <- msg:
					logx.Debugf("Message sent to client %d in node %d",
						c.GetClientId(), nodeIdx)
				case <-time.After(100 * time.Millisecond):
					logx.Errorf("Send timeout for client %d in node %d, message dropped",
						c.GetClientId(), nodeIdx)
				}
			}(client.(Client), notify.Val, nodeIndex)
			return // 找到客户端后立即返回
		}
	}

	logx.Debugf("Client %d not found in node %d for message", notify.ConnId, nodeIndex)
}

// Stop 停止 Hub 并清理资源
func (h *Hub) Stop() {
	if h.stopped.Swap(true) {
		return
	}

	logx.Info("Stopping hub...")

	// 先关闭所有通道，停止新的操作
	for i := 0; i < int(h.NodeCount); i++ {
		close(h.Node[i].Register)
		close(h.Node[i].Unregister)
		close(h.Node[i].Notify)
	}

	// 等待一段时间让节点 goroutine 处理完剩余消息
	time.Sleep(100 * time.Millisecond)

	// 安全关闭所有客户端
	var wg sync.WaitGroup
	for i := 0; i < int(h.NodeCount); i++ {
		wg.Add(1)
		go func(node *Node, nodeIndex int) {
			defer wg.Done()

			clientCount := 0
			for connType := 0; connType < len(node.Clients); connType++ {
				node.Clients[connType].Range(func(k, v interface{}) bool {
					v.(Client).Close()
					node.Clients[connType].Delete(k)
					return true
				})
			}
			logx.Infof("Node %d closed %d clients", nodeIndex, clientCount)
		}(h.Node[i], i)
	}

	wg.Wait()
	logx.Info("Hub stopped completely")
}

// AddClient 添加客户端到对应的分片节点
func (h *Hub) AddClient(client Client) {
	if h.stopped.Load() {
		logx.Info("Hub stopped")
		client.Close()
		return
	}

	nodeIdx := h.HashFunction(client.GetClientId())
	if nodeIdx >= 0 && nodeIdx < h.NodeCount {
		select {
		case h.Node[nodeIdx].Register <- client:
			logx.Debugf("Client %d queued for registration in node %d",
				client.GetClientId(), nodeIdx)
		case <-time.After(50 * time.Millisecond):
			logx.Errorf("Register channel full for node %d, client %d closed",
				nodeIdx, client.GetClientId())
			client.Close()
		}
	} else {
		logx.Errorf("Invalid node index %d for client %d", nodeIdx, client.GetClientId())
		client.Close()
	}
}

// RemoveClient 从对应的分片节点移除客户端
func (h *Hub) RemoveClient(client Client) {
	if h.stopped.Load() {
		return
	}

	nodeIdx := h.HashFunction(client.GetClientId())
	if nodeIdx >= 0 && nodeIdx < h.NodeCount {
		select {
		case h.Node[nodeIdx].Unregister <- client:
			logx.Debugf("Client %d queued for unregistration from node %d",
				client.GetClientId(), nodeIdx)
		case <-time.After(50 * time.Millisecond):
			logx.Errorf("Unregister channel full for node %d, client %d force closed",
				nodeIdx, client.GetClientId())
			client.Close()
		}
	} else {
		logx.Errorf("Invalid node index %d for client %d", nodeIdx, client.GetClientId())
		client.Close()
	}
}

// AddMessage 添加消息到对应的分片节点
func (h *Hub) AddMessage(notifyList *NotifyList) {
	if h.stopped.Load() || notifyList == nil {
		return
	}
	connList := client.GetIdsByKey(notifyList.Key)
	if len(connList) == 0 {
		h.Broadcast(notifyList)
		return
	}
	for _, conn := range connList {
		nodeIdx := h.HashFunction(conn)
		notify := NewNotify(conn, notifyList.Val)
		if nodeIdx >= 0 && nodeIdx < h.NodeCount {
			select {
			case h.Node[nodeIdx].Notify <- notify:
				logx.Debugf("Message queued for client %d in node %d",
					conn, nodeIdx)
			default:
				logx.Errorf("Message channel full for node %d, message for client %d dropped",
					nodeIdx, conn)
			}
		} else {
			logx.Errorf("Invalid node index %d for client %d", nodeIdx, conn)
		}
	}
}

// Consume 消费 Kafka 消息
func (h *Hub) Consume(_ context.Context, key, val string) error {
	h.AddMessage(NewNotifyList(key, val))
	return nil
}

// GetClientCount 获取总客户端数量
func (h *Hub) GetClientCount() int {
	total := 0
	for i := 0; i < int(h.NodeCount); i++ {
		for connType := 0; connType < len(h.Node[i].Clients); connType++ {
			h.Node[i].Clients[connType].Range(func(key, value interface{}) bool {
				total++
				return true
			})
		}
	}
	return total
}

// GetNodeStats 获取节点统计信息
func (h *Hub) GetNodeStats() map[int]int {
	stats := make(map[int]int)
	for i := 0; i < int(h.NodeCount); i++ {
		nodeTotal := 0
		for connType := 0; connType < len(h.Node[i].Clients); connType++ {
			h.Node[i].Clients[connType].Range(func(key, value interface{}) bool {
				nodeTotal++
				return true
			})
		}
		stats[i] = nodeTotal
	}
	return stats
}

// IsStopped 检查 Hub 是否已停止
func (h *Hub) IsStopped() bool {
	return h.stopped.Load()
}

// Broadcast 广播消息给所有客户端
func (h *Hub) Broadcast(conn *NotifyList) {
	if h.stopped.Load() || conn == nil {
		return
	}
	for i := 0; i < int(h.NodeCount); i++ {
		for connType := 0; connType < len(h.Node[i].Clients); connType++ {
			h.Node[i].Clients[connType].Range(func(key, value interface{}) bool {
				go func(c Client, conn *NotifyList, nodeIdx int) {
					select {
					case c.GetSendBuffer() <- conn.Val:
					case <-time.After(100 * time.Millisecond):
						logx.Errorf("Broadcast timeout for client %d in node %d",
							c.GetClientId(), nodeIdx)
					}
				}(value.(Client), conn, i)
				return true
			})
		}
	}
}
