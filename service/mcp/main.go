package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/mcp"
)

func main() {
	// 加载配置
	var c mcp.McpConf
	if err := conf.Load("config.yaml", &c); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置日志
	logx.DisableStat()

	// 创建 MCP 服务器
	server := mcp.NewMcpServer(c)
	defer server.Stop()

	// 注册一个简单的回显工具
	echoTool := mcp.Tool{
		Name:        "echo",
		Description: "Echoes back the message provided by the user",
		InputSchema: mcp.InputSchema{
			Properties: map[string]any{
				"message": map[string]any{
					"type":        "string",
					"description": "The message to echo back",
				},
				"prefix": map[string]any{
					"type":        "string",
					"description": "Optional prefix to add to the echoed message",
					"default":     "Echo: ",
				},
			},
			Required: []string{"message"},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			var req struct {
				Message string `json:"message"`
				Prefix  string `json:"prefix,optional"`
			}

			if err := mcp.ParseArguments(params, &req); err != nil {
				return nil, fmt.Errorf("failed to parse args: %w", err)
			}

			prefix := "Echo: "
			if len(req.Prefix) > 0 {
				prefix = req.Prefix
			}

			return prefix + req.Message, nil
		},
	}
	err := server.RegisterTool(echoTool)
	if err != nil {
		return
	}

	// 注册一个静态 prompt
	server.RegisterPrompt(mcp.Prompt{
		Name:        "greeting",
		Description: "A simple greeting prompt",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "name",
				Description: "The name to greet",
				Required:    true,
			},
		},
		Content: "Hello {{name}}! How can I assist you today?",
	})

	// 注册一个动态 prompt
	server.RegisterPrompt(mcp.Prompt{
		Name:        "dynamic-prompt",
		Description: "A prompt that uses a handler to generate dynamic content",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "username",
				Description: "User's name for personalized greeting",
				Required:    true,
			},
			{
				Name:        "topic",
				Description: "Topic of expertise",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, args map[string]string) ([]mcp.PromptMessage, error) {
			var req struct {
				Username string `json:"username"`
				Topic    string `json:"topic"`
			}

			if err := mcp.ParseArguments(args, &req); err != nil {
				return nil, fmt.Errorf("failed to parse args: %w", err)
			}

			// 创建包含当前时间的消息
			currentTime := time.Now().Format(time.RFC1123)
			return []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Text: fmt.Sprintf("Hello, I'm %s and I'd like to learn about %s.", req.Username, req.Topic),
					},
				},
				{
					Role: mcp.RoleAssistant,
					Content: mcp.TextContent{
						Text: fmt.Sprintf("Hello %s! I'm an AI assistant and I'll help you learn about %s. The current time is %s.",
							req.Username, req.Topic, currentTime),
					},
				},
			}, nil
		},
	})

	// 注册一个资源
	server.RegisterResource(mcp.Resource{
		Name:        "example-doc",
		URI:         "file:///example/doc.txt",
		Description: "An example document",
		MimeType:    "text/plain",
		Handler: func(ctx context.Context) (mcp.ResourceContent, error) {
			return mcp.ResourceContent{
				URI:      "file:///example/doc.txt",
				MimeType: "text/plain",
				Text:     "This is the content of the example document.",
			}, nil
		},
	})

	// 启动服务器
	fmt.Printf("Starting MCP server on %s:%d\n", c.Host, c.Port)
	server.Start()
}
