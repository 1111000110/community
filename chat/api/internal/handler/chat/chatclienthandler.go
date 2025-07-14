package chat

import (
	"community.com/chat/api/internal/types"
	"community.com/pkg/message"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	"community.com/chat/api/internal/logic/chat"
	"community.com/chat/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 定义 WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // 设置读取缓冲区大小
	WriteBufferSize: 1024, // 设置写入缓冲区大小
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatClientHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MessageClientReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil) //升级Http为websocket
		if err != nil {
			log.Println(err)
			return
		}
		l := chat.NewChatClientLogic(r.Context(), svcCtx, conn, req.UserId)
		svcCtx.Hub.Register <- l
		err = l.ChatClient()
		message.ResponseHandler(w, r.Context(), err)
	}
}
