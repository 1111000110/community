# DeepSeek API Go客户端

这是一个功能完整的DeepSeek API Go客户端，支持多种模型、流式和非流式调用、上下文保留等高级功能。

## 功能特性

- ✅ 支持多种DeepSeek模型（chat、coder、reasoner）
- ✅ 流式和非流式两种调用方式
- ✅ 完整的参数控制（温度、top_p、惩罚参数等）
- ✅ 上下文保留和历史对话管理
- ✅ 工具调用支持（Function Calling）
- ✅ JSON格式输出支持
- ✅ 详细的错误处理和日志
- ✅ 类型安全的结构体定义
- ✅ 面向对象设计，支持链式调用
- ✅ 详细的参数注释和使用建议

## 支持的模型

- `deepseek-chat`: 标准对话模型
- `deepseek-coder`: 代码专用模型
- `deepseek-reasoner`: 推理模型

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 设置API密钥

```bash
export DEEPSEEK_API_KEY="your-api-key-here"
```

### 3. 基础使用

#### 方式一：面向对象链式调用（推荐）

```go
package main

import (
    "fmt"
    "log"
    "community/pkg/aiapi/deepseek"
)

func main() {
    // 创建客户端
    client := deepseek.NewClient("your-api-key-here")
    
    // 链式配置并发送
    response, err := client.
        SetModel(deepseek.ModelDeepSeekChat).
        AddSystemMessage("你是一个有用的助手").
        AddUserMessage("你好！").
        SetTemperature(0.7).
        SetMaxTokens(1000).
        Send()
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(response.Choices[0].Message.Content)
}
```

#### 方式二：传统方式

```go
package main

import (
    "fmt"
    "log"
    "community/pkg/aiapi/deepseek"
)

func main() {
    // 创建客户端
    client := deepseek.NewClient("your-api-key-here")
    
    // 简单对话
    messages := []deepseek.Message{
        {Role: "system", Content: "你是一个有用的助手"},
        {Role: "user", Content: "你好！"},
    }
    
    response, err := client.QuickChat(messages, deepseek.ModelDeepSeekChat, false)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(response)
}
```

## 详细使用说明

### 1. 创建客户端

```go
// 基础创建
client := deepseek.NewClient("your-api-key-here")

// 带预配置创建
req := &deepseek.ChatRequest{
    Model:       deepseek.ModelDeepSeekChat,
    Temperature: 0.7,
    MaxTokens:   1000,
}
client := deepseek.NewClientWithRequest("your-api-key-here", req)
```

### 2. 参数详细说明

#### 核心参数

| 参数 | 类型 | 范围 | 默认值 | 说明 |
|------|------|------|--------|------|
| `Model` | string | - | "deepseek-chat" | 指定使用的AI模型 |
| `Messages` | []Message | - | [] | 对话消息列表 |
| `MaxTokens` | int | 1-4096 | 4096 | 最大生成token数 |
| `Temperature` | float64 | 0.0-2.0 | 1.0 | 控制输出随机性 |
| `TopP` | float64 | 0.0-1.0 | 1.0 | 控制词汇选择多样性 |

#### 惩罚参数

| 参数 | 类型 | 范围 | 默认值 | 说明 |
|------|------|------|--------|------|
| `FrequencyPenalty` | float64 | -2.0到2.0 | 0.0 | 频率惩罚，减少重复内容 |
| `PresencePenalty` | float64 | -2.0到2.0 | 0.0 | 存在惩罚，鼓励新话题 |

#### 控制参数

| 参数 | 类型 | 说明 |
|------|------|------|
| `Stop` | []string | 停止词列表，遇到时停止生成 |
| `Stream` | bool | 是否使用流式输出 |
| `ResponseFormat` | *ResponseFormat | 响应格式（text/json_object） |

#### 高级参数

| 参数 | 类型 | 说明 |
|------|------|------|
| `Tools` | []Tool | 工具函数列表，用于Function Calling |
| `ToolChoice` | ToolChoice | 工具选择策略（none/auto/required） |
| `Logprobs` | bool | 是否返回token概率信息 |
| `TopLogprobs` | int | 返回的top logprobs数量 |

### 3. 参数使用建议

#### Temperature 温度参数
- **0.1-0.3**: 事实性回答，代码生成
- **0.7-0.9**: 创意写作，对话生成
- **1.0-1.5**: 高度创造性内容

#### TopP 核采样
- **0.1-0.5**: 精确任务，需要确定性
- **0.7-0.9**: 平衡创造性和准确性
- **0.9-1.0**: 高度创造性

#### 惩罚参数
- **FrequencyPenalty**: 正值减少重复，负值增加重复
- **PresencePenalty**: 正值鼓励新话题，负值保持话题集中

## JSON 示例

### 基础请求示例

```json
{
  "model": "deepseek-chat",
  "messages": [
    {
      "role": "system",
      "content": "你是一个有用的AI助手"
    },
    {
      "role": "user",
      "content": "请解释什么是机器学习"
    }
  ],
  "max_tokens": 1000,
  "temperature": 0.7,
  "top_p": 0.9,
  "stream": false
}
```

### 流式请求示例

```json
{
  "model": "deepseek-coder",
  "messages": [
    {
      "role": "system",
      "content": "你是一个专业的Go语言编程助手"
    },
    {
      "role": "user",
      "content": "请写一个HTTP服务器"
    }
  ],
  "max_tokens": 2000,
  "temperature": 0.3,
  "top_p": 0.8,
  "stream": true,
  "stream_options": {
    "include_usage": true
  }
}
```

### 高级配置示例

```json
{
  "model": "deepseek-chat",
  "messages": [
    {
      "role": "system",
      "content": "你是一个创意写作助手"
    },
    {
      "role": "user",
      "content": "写一个关于机器人的短故事"
    }
  ],
  "max_tokens": 1500,
  "temperature": 0.8,
  "top_p": 0.9,
  "frequency_penalty": 0.1,
  "presence_penalty": 0.1,
  "stop": ["END", "故事结束"],
  "response_format": {
    "type": "text"
  },
  "stream": false
}
```

### 工具调用示例

```json
{
  "model": "deepseek-chat",
  "messages": [
    {
      "role": "system",
      "content": "你是一个天气助手，可以查询天气信息"
    },
    {
      "role": "user",
      "content": "北京今天天气怎么样？"
    }
  ],
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_weather",
        "description": "获取指定城市的天气信息",
        "parameters": {
          "type": "object",
          "properties": {
            "city": {
              "type": "string",
              "description": "城市名称"
            }
          },
          "required": ["city"]
        }
      }
    }
  ],
  "tool_choice": "auto",
  "max_tokens": 500,
  "temperature": 0.3
}
```

### 响应示例

#### 非流式响应

```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "deepseek-chat",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "机器学习是人工智能的一个分支..."
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 20,
    "completion_tokens": 150,
    "total_tokens": 170
  }
}
```

#### 流式响应

```json
{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"deepseek-chat","choices":[{"index":0,"delta":{"role":"assistant","content":"机器"},"finish_reason":null}]}
{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"deepseek-chat","choices":[{"index":0,"delta":{"content":"学习"},"finish_reason":null}]}
{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"deepseek-chat","choices":[{"index":0,"delta":{"content":"是"},"finish_reason":null}]}
{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"deepseek-chat","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}
```

### 2. 面向对象链式调用

#### 基础对话

```go
// 链式配置并发送
response, err := client.
    SetModel(deepseek.ModelDeepSeekChat).
    AddSystemMessage("你是一个有用的助手").
    AddUserMessage("请解释什么是机器学习").
    SetTemperature(0.7).
    SetMaxTokens(1000).
    Send()

if err != nil {
    log.Fatal(err)
}

fmt.Println(response.Choices[0].Message.Content)
```

#### 流式对话

```go
// 流式配置并发送
err := client.
    SetModel(deepseek.ModelDeepSeekCoder).
    AddSystemMessage("你是一个Go语言专家").
    AddUserMessage("请写一个HTTP服务器").
    SetTemperature(0.3).
    SetMaxTokens(2000).
    SetStream(true).
    SendStream(func(resp *deepseek.StreamResponse) {
        if len(resp.Choices) > 0 {
            fmt.Print(resp.Choices[0].Delta.Content)
        }
    })
```

#### 多轮对话
- system：全局系统消息
- assistant：ai曾经回复的消息
- user：用户发送的消息
必须保证顺序，且最后必须为user消息。
```go
// 构建对话历史
client.
    SetModel(deepseek.ModelDeepSeekChat).
    AddSystemMessage("你是一个编程导师").
    AddUserMessage("我想学习Go语言").
    AddAssistantMessage("很好！Go是一门优秀的语言，你想从哪个方面开始？").
    AddUserMessage("请给我一个Hello World示例").
    SetTemperature(0.5).
    SetMaxTokens(500)

// 发送请求
response, err := client.Send()
[]Message{
{Role: "system", Content: "你是技术助手，只回答编程相关问题"}, // 系统指令
{Role: "user", Content: "Go语言的切片怎么扩容？"},          // 用户第一次提问
{Role: "assistant", Content: "Go切片扩容会先检查容量是否足够，不足则按规则扩容..."}, // AI第一次回复
{Role: "user", Content: "那扩容规则具体是什么？"},          // 用户第二次提问（依赖上一次对话）
}
```

#### 高级配置

```go
// 创意写作配置
response, err := client.
    SetModel(deepseek.ModelDeepSeekChat).
    AddSystemMessage("你是一个创意写作助手").
    AddUserMessage("写一个关于机器人的短故事").
    SetTemperature(0.8).                    // 高创造性
    SetTopP(0.9).                          // 词汇多样性
    SetFrequencyPenalty(0.1).              // 减少重复
    SetPresencePenalty(0.1).               // 鼓励新话题
    SetStop([]string{"END", "故事结束"}).   // 停止词
    SetMaxTokens(1500).
    Send()
```

### 3. 传统方式调用

#### 非流式调用

```go
req := &deepseek.ChatRequest{
    Model:    deepseek.ModelDeepSeekChat,
    Messages: messages,
    MaxTokens: 1000,
    Temperature: 0.7,
    Stream:   false,
}

resp, err := client.Chat(req)
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.Choices[0].Message.Content)
```

#### 流式调用

```go
err := client.ChatStream(req, func(resp *deepseek.StreamResponse) {
    if len(resp.Choices) > 0 {
        fmt.Print(resp.Choices[0].Delta.Content)
    }
})
```

#### 带上下文的对话

```go
// 历史对话
history := []deepseek.Message{
    {Role: "user", Content: "我想学习Go"},
    {Role: "assistant", Content: "很好！Go是一门优秀的语言"},
}

// 创建带上下文的消息
messages := deepseek.CreateMessages(
    "你是一个编程导师",
    "请给我一个Go的Hello World示例",
    history,
)

response, err := client.QuickChat(messages, deepseek.ModelDeepSeekCoder, false)
```

### 4. JSON Output 功能示例

JSON Output 功能确保模型严格按照 JSON 格式输出，便于后续解析。

#### 基础 JSON 输出

```go
// 配置 JSON 输出格式
response, err := client.
    SetModel(deepseek.ModelDeepSeekChat).
    AddSystemMessage(`你是一个数据解析助手。请将用户提供的问题和答案解析成JSON格式。

示例输入：世界上最高的山是什么？珠穆朗玛峰。

示例JSON输出：
{
    "question": "世界上最高的山是什么？",
    "answer": "珠穆朗玛峰"
}`).
    AddUserMessage("世界上最长的河流是什么？尼罗河。").
    SetResponseFormat(&deepseek.ResponseFormat{
        Type: deepseek.ResponseFormatJSON,
    }).
    SetMaxTokens(500).
    SetTemperature(0.3).
    Send()

if err != nil {
    log.Fatal(err)
}

// 解析 JSON 输出
var result map[string]string
json.Unmarshal([]byte(response.Choices[0].Message.Content), &result)
fmt.Printf("问题: %s\n", result["question"])
fmt.Printf("答案: %s\n", result["answer"])
```

#### 复杂 JSON 结构输出

```go
// 定义复杂的数据结构
type ExamResult struct {
    StudentName string `json:"student_name"`
    Subject     string `json:"subject"`
    Score       int    `json:"score"`
    Grade       string `json:"grade"`
    Comments    string `json:"comments"`
}

// 配置复杂 JSON 输出
response, err := client.
    SetModel(deepseek.ModelDeepSeekChat).
    AddSystemMessage(`你是一个成绩分析助手。请将学生的考试信息解析成JSON格式。

示例输入：张三，数学考试，85分

示例JSON输出：
{
    "student_name": "张三",
    "subject": "数学",
    "score": 85,
    "grade": "B",
    "comments": "成绩良好，需要加强练习"
}`).
    AddUserMessage("李四，英语考试，92分").
    SetResponseFormat(&deepseek.ResponseFormat{
        Type: deepseek.ResponseFormatJSON,
    }).
    SetMaxTokens(300).
    SetTemperature(0.1).
    Send()

if err != nil {
    log.Fatal(err)
}

// 解析为结构体
var result ExamResult
json.Unmarshal([]byte(response.Choices[0].Message.Content), &result)
fmt.Printf("学生: %s, 科目: %s, 分数: %d, 等级: %s\n", 
    result.StudentName, result.Subject, result.Score, result.Grade)
```

#### 数组格式 JSON 输出

```go
// 配置数组格式输出
response, err := client.
    SetModel(deepseek.ModelDeepSeekChat).
    AddSystemMessage(`你是一个商品信息提取助手。请将商品列表解析成JSON数组格式。

示例输入：苹果 5元，香蕉 3元，橙子 4元

示例JSON输出：
[
    {"name": "苹果", "price": 5},
    {"name": "香蕉", "price": 3},
    {"name": "橙子", "price": 4}
]`).
    AddUserMessage("牛奶 8元，面包 6元，鸡蛋 12元").
    SetResponseFormat(&deepseek.ResponseFormat{
        Type: deepseek.ResponseFormatJSON,
    }).
    SetMaxTokens(400).
    SetTemperature(0.2).
    Send()

if err != nil {
    log.Fatal(err)
}

// 解析为数组
var products []map[string]interface{}
json.Unmarshal([]byte(response.Choices[0].Message.Content), &products)
for _, product := range products {
    fmt.Printf("商品: %s, 价格: %.0f元\n", product["name"], product["price"])
}
```

### 5. 工具调用示例

```go
// 定义工具
tools := []deepseek.Tool{
    {
        Type: "function",
        Function: map[string]interface{}{
            "name": "get_weather",
            "description": "获取指定城市的天气信息",
            "parameters": map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "city": map[string]interface{}{
                        "type":        "string",
                        "description": "城市名称",
                    },
                },
                "required": []string{"city"},
            },
        },
    },
}

// 配置并发送
response, err := client.
    SetModel(deepseek.ModelDeepSeekChat).
    AddSystemMessage("你是一个天气助手").
    AddUserMessage("北京今天天气怎么样？").
    SetTools(tools).
    SetToolChoice(deepseek.ToolChoiceAuto).
    SetTemperature(0.3).
    Send()
```

## JSON Output 功能详解

### 功能说明

JSON Output 功能确保模型严格按照 JSON 格式输出，便于后续程序解析和处理。

### 使用要点

1. **设置响应格式**：必须设置 `ResponseFormat` 为 `json_object`
2. **提示词要求**：system 或 user 提示词中必须包含 "json" 字样
3. **提供示例**：在提示词中给出期望的 JSON 格式示例
4. **合理设置 max_tokens**：防止 JSON 字符串被截断
5. **低温度设置**：建议使用 0.1-0.3 的温度值确保输出稳定性

### 最佳实践

```go
// ✅ 正确的 JSON Output 配置
response, err := client.
    SetModel(deepseek.ModelDeepSeekChat).
    AddSystemMessage(`请将以下信息解析成JSON格式。必须严格按照示例格式输出。

示例：
输入：姓名：张三，年龄：25
输出：{"name": "张三", "age": 25}`).
    AddUserMessage("姓名：李四，年龄：30").
    SetResponseFormat(&deepseek.ResponseFormat{
        Type: deepseek.ResponseFormatJSON,
    }).
    SetMaxTokens(200).        // 合理设置token限制
    SetTemperature(0.1).      // 低温度确保稳定性
    Send()

// ❌ 错误的配置 - 缺少json关键字
response, err := client.
    AddSystemMessage("请输出结构化数据").  // 缺少"json"关键字
    SetResponseFormat(&deepseek.ResponseFormat{
        Type: deepseek.ResponseFormatJSON,
    }).
    Send()
```

### 常见问题

1. **空内容返回**：如果API返回空的content，尝试修改提示词，增加更明确的JSON格式要求
2. **JSON格式错误**：确保提示词中包含完整的JSON示例
3. **内容截断**：适当增加max_tokens参数值

## API参考

### 主要结构体

#### ChatRequest
```go
type ChatRequest struct {
    Model            string          `json:"model"`                       // 模型名称
    Messages         []Message       `json:"messages"`                    // 对话消息列表
    MaxTokens        int             `json:"max_tokens,omitempty"`        // 最大生成token数
    Temperature      float64         `json:"temperature,omitempty"`       // 温度参数 (0-2)
    TopP             float64         `json:"top_p,omitempty"`             // Top-p参数 (0-1)
    FrequencyPenalty float64         `json:"frequency_penalty,omitempty"` // 频率惩罚 (-2到2)
    PresencePenalty  float64         `json:"presence_penalty,omitempty"`  // 存在惩罚 (-2到2)
    Stop             []string        `json:"stop,omitempty"`              // 停止词列表
    Stream           bool            `json:"stream"`                      // 是否流式输出
    StreamOptions    *StreamOptions  `json:"stream_options,omitempty"`    // 流式选项
    ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`   // 响应格式
    Tools            []Tool          `json:"tools,omitempty"`             // 工具列表
    ToolChoice       ToolChoice      `json:"tool_choice,omitempty"`       // 工具选择策略
    Logprobs         bool            `json:"logprobs,omitempty"`          // 是否返回logprobs
    TopLogprobs      int             `json:"top_logprobs,omitempty"`      // 返回的top logprobs数量
}
```

#### Message
```go
type Message struct {
    Role    string `json:"role"`    // 角色: system, user, assistant
    Content string `json:"content"` // 消息内容
}
```

### 主要方法

#### NewClient
```go
func NewClient(apiKey string) *Client
```
创建新的DeepSeek客户端实例。

#### Chat
```go
func (c *Client) Chat(req *ChatRequest) (*ChatResponse, error)
```
发送非流式聊天请求。

#### ChatStream
```go
func (c *Client) ChatStream(req *ChatRequest, callback func(*StreamResponse)) error
```
发送流式聊天请求。

#### QuickChat
```go
func (c *Client) QuickChat(messages []Message, model string, stream bool) (string, error)
```
快速聊天方法，简化调用。

#### CreateMessages
```go
func CreateMessages(systemPrompt, userMessage string, history []Message) []Message
```
创建带上下文的对话消息。


## 错误处理

客户端提供了详细的错误信息：

```go
resp, err := client.Chat(req)
if err != nil {
    // 错误可能包括：
    // - 网络错误
    // - API错误（认证失败、配额超限等）
    // - 解析错误
    log.Printf("错误: %v", err)
    return
}
```
