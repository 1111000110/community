package user

import (
	"community.com/pkg/message"
	"community.com/user/api/internal/logic/user"
	"community.com/user/api/internal/svc"
	"community.com/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func UserLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUserLoginLogic(r.Context(), svcCtx)
		resp, err := l.UserLogin(&req)
		message.ResponseHandler(w, resp, err)
	}
}
