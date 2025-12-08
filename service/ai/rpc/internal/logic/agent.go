package logic

type ResponseFormatType string // 指定的输出格式

const (
	text       ResponseFormatType = "text"        // 默认值
	jsonObject ResponseFormatType = "json_object" // json结构
	markdown   ResponseFormatType = "markdown"    // markdown格式，例：{"name": "张三", "age": 30, "gender": "男"}
	csv        ResponseFormatType = "csv"         // csv结构，例：name,age,gender\n张三,30,男\n李四,25,女
)

type ResponseFormat struct {
	Type ResponseFormatType `json:"type"`
}
type ToolChoice string // 指定的工具选择模式

const (
	None     ToolChoice = "None"     // 意味着模型不会调用任何 tool，而是生成一条消息。
	Auto     ToolChoice = "Auto"     // 意味着模型可以选择生成一条消息或调用一个或多个 tool。
	Required ToolChoice = "Required" // 意味着模型必须调用一个或多个 tool。
)

type LLM struct { // 大模型，deepSeek，豆包等
	LLMId   int64  `json:"model_id"`
	LLMName string `json:"model_name"`
	LLMKey  string `json:"model_key"`
}

type Tool struct { // todo,可以使用的工具
	ToolId int64 `json:"toolId"`
}
type Agent struct {
	AgentId          int64          `json:"agent_id"`          // id
	AgentPrompt      string         `json:"agent_prompt"`      // 提示词
	Temperature      int64          `json:"temperature"`       // 温度，值越高生成的越有创造性，值越低生成的越保守。
	MaxToken         int64          `json:"max_token"`         // 最大使用token数
	Llm              LLM            `json:"llm"`               // 大模型实体
	Model            string         `json:"model"`             // 使用的模型
	PresencePenalty  float64        `json:"presence_penalty"`  // 是否出现过，只要出现一次后续就会降低很多。避免模型 “钻牛角尖”（如只围绕一个点展开，不拓展新信息）。
	FrequencyPenalty float64        `json:"frequency_penalty"` // 重复惩罚参数，一个词 / 短语在已有文本中出现的次数越多，后续被生成的概率越低（正值时）。避免 “车轱辘话”（如反复说 “总结来说”“具体而言”）。
	ResponseFormat   ResponseFormat `json:"response_format"`   // 返回格式定义
	Stop             []string       `json:"stop"`              // 停止词汇，在遇到这些词时，API 将停止生成更多的 token。
	Stream           bool           `json:"stream"`            // 是否流式回复，SSE形式发送消息增量
	Tools            []Tool         `json:"tools"`             // 可以使用的工具
	ToolChoice       ToolChoice     `json:"tool_choice"`       // 工具选择情况
}
