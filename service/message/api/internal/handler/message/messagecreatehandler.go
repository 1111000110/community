package message

import (
	"net/http"

	"community/service/message/api/internal/logic/message"
	"community/service/message/api/internal/svc"
	"community/service/message/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func MessageCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MessageCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := message.NewMessageCreateLogic(r.Context(), svcCtx)
		resp, err := l.MessageCreate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
