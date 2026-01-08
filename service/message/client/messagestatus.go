package messageclient

const (
	MsgStatusNormal   int64 = 0 // 正常：消息成功发送/接收，状态正常
	MsgStatusAbnormal int64 = 1 // 异常：消息发送失败、格式错误、解析异常等
)
