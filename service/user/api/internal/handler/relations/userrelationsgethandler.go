package relations

import (
	"community/pkg/message"
	"community/service/user/api/internal/logic/relations"
	"community/service/user/api/internal/svc"
	"community/service/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func UserRelationsGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRelationsGetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := relations.NewUserRelationsGetLogic(r.Context(), svcCtx)
		resp, err := l.UserRelationsGet(&req)
		message.ResponseHandler(w, resp, err)
	}
}
