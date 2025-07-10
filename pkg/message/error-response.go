package message

type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 实现error的接口  然后CodeError继承一下Error方法  CodeError就为error类型的返回值
func (e *CodeError) Error() string {
	return e.Message
}

// 返回给前端的数据
func (e *CodeError) ErrorResponse() CodeError {
	return CodeError{
		Code:    e.Code,
		Message: e.Message,
	}
}

// 用来自定义抛出的异常，支持调用传递
func NewCodeError(code int, msg string) error {
	return &CodeError{
		Code:    code,
		Message: msg,
	}
}

// 提供new方法，任意地方传递参数返回CodeError类型的数据
func New(code int, msg string) CodeError {
	return CodeError{
		Code:    code,
		Message: msg,
	}
}

// 默认异常状态码函数，只需传递错误信息即可，默认返回code-10001
func DefaultErrHandler(message string) error {
	return &CodeError{
		Code:    ServerError.Code,
		Message: message,
	}
}

// 自定义错误返回函数 错误函数主入口
func ErrHandler(err error) interface{} {
	switch err.(type) {
	// 如果错误类型为CodeError，就返回错误类型的结构体
	case *CodeError:
		return err
	default:
		// 系统错误，500 错误提示
		return CodeError{Code: ServerError.Code, Message: err.Error()}
	}
}
