package Post

import (
	"community.com/pkg/message"
	"community.com/service/post/api/internal/logic/Post"
	"community.com/service/post/api/internal/svc"
	"community.com/service/post/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func PostDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := Post.NewPostDeleteLogic(r.Context(), svcCtx)
		resp, err := l.PostDelete(&req)
		message.ResponseHandler(w, resp, err)
	}
}
