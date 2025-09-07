package websocket

import (
	"log"
	"net/http"

	"community.com/service/websocket/api/internal/logic/websocket"
	"community.com/service/websocket/api/internal/svc"
	"community.com/service/websocket/api/internal/types"
	websocketclient "github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 定义 WebSocket 升级器
var upGrader = websocketclient.Upgrader{
	ReadBufferSize:  1024, // 设置读取缓冲区大小
	WriteBufferSize: 1024, // 设置写入缓冲区大小
	CheckOrigin: func(r *http.Request) bool {
		return true
	}, // 是否允许跨域
}

func WebSocketClientHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WebSocketReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		conn, err := upGrader.Upgrade(w, r, nil) //升级Http为websocket
		if err != nil {
			log.Println(err)
			return
		}
		l := websocket.NewWebSocketClientLogic(r.Context(), svcCtx, &req, conn) // 获取一个消息连接
		svcCtx.MessageHub.AddClient(l)                                          // 增加进消息中心
		err = l.RunSocketClient()                                               // 启动消息处理
		if err != nil {
			panic(err)
		}
	}
}
