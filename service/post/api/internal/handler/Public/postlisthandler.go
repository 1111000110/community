package Public

import (
	"community/pkg/message"
	"community/service/post/api/internal/logic/Public"
	"community/service/post/api/internal/svc"
	"community/service/post/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func PostListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := Public.NewPostListLogic(r.Context(), svcCtx)
		resp, err := l.PostList(&req)
		message.ResponseHandler(w, resp, err)
	}
}
