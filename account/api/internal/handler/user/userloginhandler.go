package user

import (
	"community.com/account/api/internal/logic/user"
	"community.com/account/api/internal/svc"
	"community.com/account/api/internal/types"
	"community.com/pkg/message"
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
