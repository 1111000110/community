package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 飞书消息结构
type FeishuMessage struct {
	MsgType string                 `json:"msg_type"`
	Content map[string]interface{} `json:"content"`
} // 发送消息到飞书
func SendStringToFeishu(url string, text string) error {
	// 构建消息内容
	message := FeishuMessage{
		MsgType: "text",
		Content: map[string]interface{}{
			"text": text,
		},
	}

	// 将消息转换为JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息JSON编码失败: %v", err)
	}

	// 发送POST请求
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageJSON))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送消息到飞书失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取飞书响应失败: %v", err)
	}

	// 打印响应
	fmt.Printf("飞书响应: %s\n", string(body))

	return nil
}
