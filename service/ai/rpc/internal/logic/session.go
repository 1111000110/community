package logic

type Role string

const (
	User      Role = "user"      // 用户
	Assistant Role = "assistant" // 历史对话
	System    Role = "system"    // 系统，prompt
)

type Message struct {
	MessageId string `json:"messageId"` // 消息id
	Role      Role   `json:"role"`      // 角色
	Content   string `json:"content"`   // 文本
}
type Session struct {
	SessionId int64     `json:"session_id"` // 会话id
	Messages  []Message `json:"messages"`   // 消息
	AgentId   int64     `json:"agentId"`    // 使用的agentId
}
