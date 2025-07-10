package hub

import "community.com/chat/internal/types"

type Client interface {
	GetClientId() int64
	WritePump()
	ReadPump()
	GetSendBuffer() chan types.Message
}
