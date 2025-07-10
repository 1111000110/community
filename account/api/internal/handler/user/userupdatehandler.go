package user

import (
	"community.com/account/api/internal/logic/user"
	"community.com/account/api/internal/svc"
	"community.com/account/api/internal/types"
	"community.com/pkg/message"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func UserUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUserUpdateLogic(r.Context(), svcCtx)
		resp, err := l.UserUpdate(&req)
		message.ResponseHandler(w, resp, err)
	}
}
