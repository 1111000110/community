package deepseek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 支持的模型类型
const (
	ModelDeepSeekChat     = "deepseek-chat"     // 标准对话模型
	ModelDeepSeekCoder    = "deepseek-coder"    // 代码专用模型
	ModelDeepSeekReasoner = "deepseek-reasoner" // 推理模型
)

type ResponseFormatType string // 响应格式类型
const (
	ResponseFormatText ResponseFormatType = "text"        // 文本格式传输
	ResponseFormatJSON ResponseFormatType = "json_object" // json格式传输
)

type ToolChoice string // 工具选择类型
const (
	ToolChoiceNone     ToolChoice = "none"
	ToolChoiceAuto     ToolChoice = "auto"
	ToolChoiceRequired ToolChoice = "required"
)

const (
	User      string = "user"
	System    string = "system"
	Assistant string = "assistant"
)

type Message struct { // 消息结构体
	Role    string `json:"role"`    // 角色: system, user, assistant
	Content string `json:"content"` // 消息内容
}

type ResponseFormat struct { // 响应格式配置
	Type ResponseFormatType `json:"type"` // 响应格式类型
}

type Tool struct { // 工具定义
	Type     string                 `json:"type"`     // 工具类型
	Function map[string]interface{} `json:"function"` // 函数定义
}

type StreamOptions struct { // 流式选项
	IncludeUsage bool `json:"include_usage"` // 是否包含使用统计
}

// ChatRequest 聊天请求配置结构体
type ChatRequest struct {
	// Model 指定使用的AI模型名称
	// 可选值: "deepseek-chat"(标准对话), "deepseek-coder"(代码专用), "deepseek-reasoner"(推理模型)
	// 影响: 不同模型在特定任务上的表现不同，选择合适的模型能获得更好的效果
	Model string `json:"model"`

	// Messages 对话消息列表，包含完整的对话历史
	// 格式: [{"role": "system", "content": "系统提示"}, {"role": "user", "content": "用户消息"}]
	// 影响: 提供上下文信息，帮助AI理解对话背景和用户意图
	Messages []Message `json:"messages"`

	// MaxTokens 最大生成token数量限制
	// 范围: 1-4096 (具体上限取决于模型)
	// 影响: 控制回复长度，设置过小可能导致回复被截断，过大可能浪费资源
	MaxTokens int `json:"max_tokens,omitempty"`

	// Temperature 温度参数，控制输出的随机性和创造性
	// 范围: 0.0-2.0
	// 影响: 值越高输出越随机和创造性，值越低输出越确定和保守
	// 建议: 0.7-0.9用于创意写作，0.1-0.3用于事实性回答
	Temperature float64 `json:"temperature,omitempty"`

	// TopP 核采样参数，控制词汇选择的多样性
	// 范围: 0.0-1.0
	// 影响: 值越高词汇选择越多样，值越低选择越集中
	// 建议: 0.9用于创意任务，0.5用于精确任务
	TopP float64 `json:"top_p,omitempty"`

	// FrequencyPenalty 频率惩罚，减少重复内容的生成
	// 范围: -2.0到2.0
	// 影响: 正值减少重复，负值增加重复
	// 建议: 0.1-0.6用于减少重复，-0.1到-0.6用于增加重复
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	// PresencePenalty 存在惩罚，鼓励模型谈论新话题
	// 范围: -2.0到2.0
	// 影响: 正值鼓励谈论新话题，负值鼓励重复已有话题
	// 建议: 0.1-0.6用于增加话题多样性，-0.1到-0.6用于保持话题集中
	PresencePenalty float64 `json:"presence_penalty,omitempty"`

	// Stop 停止词列表，遇到这些词时停止生成
	// 格式: ["END", "STOP", "用户说"]
	// 影响: 控制生成何时停止，避免生成不需要的内容
	// 建议: 设置明确的结束标记，如["\n\n用户:", "END"]
	Stop []string `json:"stop,omitempty"`

	// Stream 是否使用流式输出
	// 影响: true时实时返回生成内容，false时等待完整生成后返回
	// 建议: true用于实时对话，false用于批量处理
	Stream bool `json:"stream"`

	// StreamOptions 流式输出选项配置
	// 影响: 控制流式输出的行为，如是否包含使用统计信息
	StreamOptions *StreamOptions `json:"stream_options,omitempty"`

	// ResponseFormat 响应格式配置
	// 可选值: "text"(纯文本), "json_object"(JSON格式)
	// 影响: 控制返回内容的格式，JSON格式便于程序解析
	// 建议: 需要结构化数据时使用json_object
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`

	// Tools 工具函数列表，用于Function Calling
	// 格式: [{"type": "function", "function": {"name": "get_weather", "description": "获取天气"}}]
	// 影响: 允许AI调用外部工具，扩展功能范围
	// 建议: 根据需求定义具体的工具函数
	Tools []Tool `json:"tools,omitempty"`

	// ToolChoice 工具选择策略
	// 可选值: "none"(不使用工具), "auto"(自动选择), "required"(必须使用工具)
	// 影响: 控制AI是否以及如何使用工具
	// 建议: "auto"用于灵活使用，"required"用于强制使用特定工具
	ToolChoice ToolChoice `json:"tool_choice,omitempty"`

	// Logprobs 是否返回每个token的对数概率
	// 影响: 提供生成过程的概率信息，用于分析和调试
	// 建议: 仅在需要分析生成质量时启用
	Logprobs bool `json:"logprobs,omitempty"`

	// TopLogprobs 返回的top logprobs数量
	// 范围: 0-5
	// 影响: 控制返回多少个最可能的token及其概率
	// 建议: 1-3用于一般分析，5用于详细分析
	TopLogprobs int `json:"top_logprobs,omitempty"`
}

// StreamResponse 流式响应结构体
type StreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Delta        Message     `json:"delta"`
		Logprobs     interface{} `json:"logprobs,omitempty"`
		FinishReason string      `json:"finish_reason,omitempty"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

// ChatResponse 非流式响应结构体
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      Message     `json:"message"`
		Logprobs     interface{} `json:"logprobs,omitempty"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ErrorResponse 错误响应结构体
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Client DeepSeek API客户端
type Client struct {
	APIKey     string       // API密钥
	BaseURL    string       // API基础URL
	HTTPClient *http.Client // HTTP客户端
	Request    *ChatRequest // 请求配置
}

// NewClient 创建新的DeepSeek客户端
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		BaseURL:    "https://api.deepseek.com",
		HTTPClient: &http.Client{},
		Request: &ChatRequest{
			Model:            ModelDeepSeekChat,
			Messages:         nil,
			MaxTokens:        4096,
			Temperature:      1.0,
			TopP:             1.0,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
			Stop:             nil,
			Stream:           false,
			StreamOptions:    nil,
			ResponseFormat:   nil,
			Tools:            nil,
			ToolChoice:       ToolChoiceNone,
			Logprobs:         false,
			TopLogprobs:      0,
		},
	}
}

// NewClientWithRequest 创建带有预配置请求的DeepSeek客户端
func NewClientWithRequest(apiKey string, req *ChatRequest) *Client {
	return &Client{
		APIKey:     apiKey,
		BaseURL:    "https://api.deepseek.com",
		HTTPClient: &http.Client{},
		Request:    req,
	}
}

// SetModel 设置模型
func (c *Client) SetModel(model string) *Client {
	c.Request.Model = model
	return c
}

// SetMessages 设置消息列表
func (c *Client) SetMessages(messages []Message) *Client {
	c.Request.Messages = messages
	return c
}

// SetMaxTokens 设置最大token数
func (c *Client) SetMaxTokens(maxTokens int) *Client {
	c.Request.MaxTokens = maxTokens
	return c
}

// SetTemperature 设置温度参数
func (c *Client) SetTemperature(temperature float64) *Client {
	c.Request.Temperature = temperature
	return c
}

// SetTopP 设置TopP参数
func (c *Client) SetTopP(topP float64) *Client {
	c.Request.TopP = topP
	return c
}

// SetFrequencyPenalty 设置频率惩罚
func (c *Client) SetFrequencyPenalty(penalty float64) *Client {
	c.Request.FrequencyPenalty = penalty
	return c
}

// SetPresencePenalty 设置存在惩罚
func (c *Client) SetPresencePenalty(penalty float64) *Client {
	c.Request.PresencePenalty = penalty
	return c
}

// SetStop 设置停止词
func (c *Client) SetStop(stop []string) *Client {
	c.Request.Stop = stop
	return c
}

// SetStream 设置是否流式输出
func (c *Client) SetStream(stream bool) *Client {
	c.Request.Stream = stream
	return c
}

// SetResponseFormat 设置响应格式
func (c *Client) SetResponseFormat(format *ResponseFormat) *Client {
	c.Request.ResponseFormat = format
	return c
}

// SetTools 设置工具列表
func (c *Client) SetTools(tools []Tool) *Client {
	c.Request.Tools = tools
	return c
}

// SetToolChoice 设置工具选择策略
func (c *Client) SetToolChoice(choice ToolChoice) *Client {
	c.Request.ToolChoice = choice
	return c
}

// AddMessage 添加消息到对话历史
func (c *Client) AddMessage(role, content string) *Client {
	c.Request.Messages = append(c.Request.Messages, Message{
		Role:    role,
		Content: content,
	})
	return c
}

// AddDialogue 创建一组对话
func (c *Client) AddDialogue(systemPrompt, userMessage string) *Client {
	return c.AddAssistantMessage(systemPrompt).AddUserMessage(userMessage)
}

// AddSystemMessage 添加系统消息
func (c *Client) AddSystemMessage(content string) *Client {
	return c.AddMessage(System, content)
}

// AddUserMessage 添加用户消息
func (c *Client) AddUserMessage(content string) *Client {
	return c.AddMessage(User, content)
}

// AddAssistantMessage 添加助手消息
func (c *Client) AddAssistantMessage(content string) *Client {
	return c.AddMessage(Assistant, content)
}

// ClearMessages 清空消息历史
func (c *Client) ClearMessages() *Client {
	c.Request.Messages = []Message{}
	return c
}

// Reset 重置客户端配置到默认值
func (c *Client) Reset() *Client {
	c.Request = &ChatRequest{}
	return c
}

// Send 使用客户端配置发送聊天请求（非流式）
func (c *Client) Send() (*ChatResponse, error) {
	jsonData, err := json.Marshal(c.Request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("API错误: %s", errorResp.Error.Message)
		}
		return nil, fmt.Errorf("HTTP错误: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &chatResp, nil
}

// SendStream 使用客户端配置发送流式聊天请求
func (c *Client) SendStream(callback func(*StreamResponse)) error {
	// 强制设置为流式
	c.Request.Stream = true

	// 序列化请求体
	jsonData, err := json.Marshal(c.Request)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	// 发送请求
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP错误: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 读取流式数据
	buffer := make([]byte, 1024)
	var lineBuffer strings.Builder

	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			lineBuffer.Write(buffer[:n])
			content := lineBuffer.String()

			// 处理SSE格式的数据
			lines := strings.Split(content, "\n")
			for i, line := range lines {
				if i == len(lines)-1 {
					// 最后一行可能不完整，保留在buffer中
					lineBuffer.Reset()
					lineBuffer.WriteString(line)
					continue
				}

				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "data: ") {
					data := strings.TrimPrefix(line, "data: ")
					if data == "[DONE]" {
						return nil
					}

					// 解析JSON数据
					var streamResp StreamResponse
					if err := json.Unmarshal([]byte(data), &streamResp); err == nil {
						callback(&streamResp)
					}
				}
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取流式数据失败: %v", err)
		}
	}

	return nil
}
