package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoBaseModel struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`                  // MongoDB 内置主键
	UpdateTime int64         `bson:"update_time,omitempty" json:"update_time,omitempty"` // 最近更新时间（Unix 时间戳）
	CreateTime int64         `bson:"create_time,omitempty" json:"create_time,omitempty"` // 创建时间（Unix 时间戳）
}

type Agent struct {
	MongoBaseModel
	AgentId    int64           `bson:"agent_id" json:"agent_id"`             // 业务层使用的 Agent 唯一 ID（对应 proto 中的 agent_id）
	ApiKey     string          `bson:"api_key" json:"api_key"`               // 调用该 Agent 的 API Key（用于鉴权）
	Name       string          `bson:"name" json:"name"`                     // Agent 名称
	Desc       string          `bson:"desc,omitempty" json:"desc,omitempty"` // Agent 描述信息
	Icon       string          `bson:"icon,omitempty" json:"icon,omitempty"` // Agent 头像 / 图标地址
	Status     int64           `bson:"status" json:"status"`                 // Agent 状态（如：1=启用，0=禁用）
	ChatConfig AgentChatConfig `bson:"chat_config" json:"chat_config"`       // 智能体聊天配置（核心参数）
}

// AgentChatConfig 定义 Agent 的聊天行为与推理参数
type AgentChatConfig struct {
	SystemPrompt string  `bson:"system_prompt" json:"system_prompt"` // 系统提示词（System Prompt）
	IsStream     int32   `bson:"is_stream" json:"is_stream"`         // 是否启用流式输出：1=是，0=否
	ChatType     int32   `bson:"chat_type" json:"chat_type"`         // 会话模式：0=新会话，1=session 会话，2=永久会话
	ChatRound    int32   `bson:"chat_round" json:"chat_round"`       // 携带的历史对话轮数上限
	Temperature  float32 `bson:"temperature" json:"temperature"`     // 采样温度（0~1，越大越随机）
	MaxTokens    int32   `bson:"max_tokens" json:"max_tokens"`       // 单次生成的最大 Token 数
	TopP         float32 `bson:"top_p" json:"top_p"`                 // Top-P 核采样参数（0~1）

	EnableTools     []string `bson:"enable_tools,omitempty" json:"enable_tools,omitempty"`         // 启用的工具列表（如：web、kb）
	EnableFunctions []string `bson:"enable_functions,omitempty" json:"enable_functions,omitempty"` // 启用的业务函数白名单

	ToolConfig ToolConfig `bson:"tool_config" json:"tool_config"` // 工具级配置（不同工具的详细参数）
	LlmId      int64      `bson:"llm_id" json:"llm_id"`           // 模型id
}

// ToolConfig 定义 Agent 可使用的外部工具配置
type ToolConfig struct {
	KB  *KBSearchConfig  `bson:"kb,omitempty" json:"kb,omitempty"`   // 知识库检索工具配置（可选）
	Web *WebSearchConfig `bson:"web,omitempty" json:"web,omitempty"` // 联网搜索工具配置（可选）
}

// KBSearchConfig 知识库检索配置
type KBSearchConfig struct {
	KnowledgeList  []string `bson:"knowledge_list,omitempty" json:"knowledge_list,omitempty"` // 允许挂载的知识库名称列表
	TopK           int32    `bson:"top_k" json:"top_k"`                                       // 向量检索返回的最大文档数量
	ScoreThreshold float32  `bson:"score_threshold" json:"score_threshold"`                   // 相似度阈值（0~1，低于该值的文档将被过滤）
}

// WebSearchConfig 联网搜索配置
type WebSearchConfig struct {
	TopK            int32    `bson:"top_k" json:"top_k"`                                     // 单次搜索返回的最大结果条数
	RecencyDays     int32    `bson:"recency_days" json:"recency_days"`                       // 搜索结果时间范围（单位：天）
	AllowDomains    []string `bson:"allow_domains,omitempty" json:"allow_domains,omitempty"` // 允许访问的域名白名单
	BlockDomains    []string `bson:"block_domains,omitempty" json:"block_domains,omitempty"` // 禁止访问的域名黑名单
	MaxCallsPerTurn int32    `bson:"max_calls_per_turn" json:"max_calls_per_turn"`           // 每轮对话中允许的最大联网调用次数
}
