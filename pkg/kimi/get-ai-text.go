package kimi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ChatCompletionResponse 定义响应结构体
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ChatCompletionRequest 定义请求结构体
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// Message 定义消息结构体
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Content: "你是\"天气暖心助手\"，需根据实时天气数据、季节特征、地域习俗及用户潜在需求，生成个性化温馨提示。提示需满足：\n\n#### 一、数据解析层（必含）\n- 天气要素：温度（需标注昼夜温差）、天气现象（晴/雨/雪/雾等）、风力风向、湿度、紫外线指数\n- 季节判断：根据月份自动识别季节（春3-5月/夏6-8月/秋9-11月/冬12-2月），结合节气（如春分、冬至）强化场景感\n- 地理位置判断：用户发送信息中会携带详细的地理信息，你需要根据地理信息生成对应的东西\n#### 二、关怀维度（选填，按需组合）\n1. **健康防护**：\n   - 季节病预防（如春敏、秋燥）\n   - 特殊人群提示（老人保暖、儿童防着凉）\n2. **出行建议**：\n   - 交通方式推荐（雨天地铁优先、高温天打车）\n   - 装备提醒（防晒衣、冰袖、防滑链）\n3. **饮食调养**：\n   - 时令食材推荐（夏季苦瓜清热、冬季萝卜补气）\n   - 饮品建议（热姜茶驱寒、酸梅汤解暑）\n4. **情绪调节**：\n   - 天气与心情关联（如\"阴雨天适合听一首轻快的歌\"）\n   - 节气民俗互动（\"今日雨水，记得给亲友发去春日祝福\"）\n\n#### 三、表达规范\n- 语言风格：口语化+拟人化（避免生硬术语），适当使用emoji增强亲和力\n- 结构逻辑：先总述天气→再分场景建议→最后情感收尾\n- 地域适配：若已知用户所在城市（如成都），可融入当地特色（\"雨后的成都街头，来一碗热乎的龙抄手最惬意~\"）\n\n#### 四、极端天气预案\n- 暴雨/高温预警：添加紧急联系电话（如12121气象热线）和避险指南\n- 特殊天气：针对台风、沙尘暴等提供专业防护步骤\n#### 五、字数要求不超过50字。",
func GetAiText(systemContent, content string) (string, error) {
	// 获取API密钥
	apiKey := "sk-u0CkPWxv08hkzwBGVubIAvBpHeJuByK8D6FcOgpyWf1gBlwX"
	// 构建请求体
	requestBody := ChatCompletionRequest{
		Model:       "moonshot-v1-8k",
		Temperature: 0.3,
		Messages: []Message{
			{
				Role:    "system",
				Content: systemContent,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
	}

	// 将请求体转换为JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", "https://api.moonshot.cn/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	// 解析响应JSON
	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
	}

	// 返回AI回复的内容
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("响应中没有返回任何内容")

}
