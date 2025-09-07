package hub

import (
	"fmt"
	"sync"
)

const HashCount = 10

type Hub struct {
	NodeCount    int64             // 多少个节点
	HashFunction func(int64) int64 // 哈希函数
	Node         []*Node
}
type Node struct {
	Clients    []map[int64]Client
	Register   chan Client        // 注册请求的通道，用于接收客户端的注册请求
	Unregister chan Client        // 注销请求的通道，用于接收客户端的注销请求
	Message    chan *Notification // 消息
}

// NewHub 创建一个新的 Hub 实例。
func NewHub() *Hub {
	nodeSlice := make([]*Node, HashCount)
	for i := 0; i < HashCount; i++ {
		clients := make([]map[int64]Client, ConnTypeCount)
		for j := 0; j < ConnTypeCount; j++ {
			clients[j] = make(map[int64]Client)
		}
		nodeSlice[i] = &Node{
			Clients:    clients,
			Register:   make(chan Client),
			Unregister: make(chan Client),
			Message:    make(chan *Notification),
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

// Run 启动 Hub 的主循环，处理注册、注销和消息发送。
func (h *Hub) Run() {
	wg := sync.WaitGroup{}
	for i := 0; i < int(h.NodeCount); i++ {
		wg.Add(1)
		go func(node *Node) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("panic for Hub:", err)
				}
			}()
			for {
				select {
				case client := <-node.Register: // 处理客户端注册请求
					node.Clients[client.GetType()][client.GetClientId()] = client
				case client := <-node.Unregister: // 处理客户端注销请求
					delete(node.Clients[client.GetType()], client.GetClientId())
					close(client.GetSendBuffer()) // 关闭客户端的发送通道
				case message := <-node.Message: //有消息
					for _, clientType := range node.Clients {
						if info, ok := clientType[message.ConnId]; ok {
							if client, ok := info.(Client); ok && message != nil {
								client.GetSendBuffer() <- *message
							}
						}
					}
				}
			}
		}(h.Node[i])
	}
	wg.Wait()
}

func (h *Hub) AddClient(client Client) {
	nodeIdx := h.HashFunction(client.GetClientId())
	if nodeIdx >= 0 && nodeIdx < h.NodeCount {
		h.Node[nodeIdx].Register <- client
	}
}
func (h *Hub) RemoveClient(client Client) {
	nodeIdx := h.HashFunction(client.GetClientId())
	if nodeIdx >= 0 && nodeIdx < h.NodeCount {
		h.Node[nodeIdx].Unregister <- client
	}
}
