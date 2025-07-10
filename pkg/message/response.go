package message

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"data"`
}

// 统一返回入口，
func ResponseHandler(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		httpx.OkJson(w, ErrHandler(err))
		// 如果err不为空的话，走错误处理函数，将err传递过去
	} else {
		// 没有错误信息，返回相应内容
		httpx.OkJson(w, Body{
			Code:    OK.Code,
			Message: OK.Message,
			Result:  resp,
		})
	}
}
