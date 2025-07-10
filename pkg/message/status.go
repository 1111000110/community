package message

var (
	OK          = New(0, "操作成功")
	ReidsError  = New(6379, "redis连接失败")
	JwtError    = New(500, "JWT解密失败")
	TokenError  = New(401, "Token已过期")
	ServerError = New(500, "服务器内部错误")
)
