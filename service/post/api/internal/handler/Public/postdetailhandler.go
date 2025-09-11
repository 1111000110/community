package Public

import (
	"community/pkg/message"
	"community/service/post/api/internal/logic/Public"
	"community/service/post/api/internal/svc"
	"community/service/post/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func PostDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := Public.NewPostDetailLogic(r.Context(), svcCtx)
		resp, err := l.PostDetail(&req)
		message.ResponseHandler(w, resp, err)
	}
}
