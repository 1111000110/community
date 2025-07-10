package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"community.com/pkg/feishu"
	"community.com/pkg/kimi"
)

// 用户配置结构
type UserConfig struct {
	SendName string `json:"send_name"`
	OpenID   string `json:"open_id"`
	CityCode string `json:"city_code"`
	CityName string `json:"city_name"`
	Hour     int    `json:"hour"`
	Minute   int    `json:"minute"`
}

type Config struct {
	WeatherAPI struct {
		Key string `json:"key"`
	} `json:"weather_api"`
	Users  []UserConfig `json:"users"`
	Feishu []string     `json:"feishu_webhooks"`
}

// 高德天气API响应结构
type WeatherResponse struct {
	Status    string        `json:"status"`
	Count     string        `json:"count"`
	Info      string        `json:"info"`
	InfoCode  string        `json:"infocode"`
	Lives     []WeatherInfo `json:"lives,omitempty"`
	Forecasts []Forecast    `json:"forecasts,omitempty"`
}

type WeatherInfo struct {
	Province      string `json:"province"`
	City          string `json:"city"`
	Adcode        string `json:"adcode"`
	Weather       string `json:"weather"`
	Temperature   string `json:"temperature"`
	WindDirection string `json:"winddirection"`
	WindPower     string `json:"windpower"`
	Humidity      string `json:"humidity"`
	ReportTime    string `json:"reporttime"`
}

type Forecast struct {
	City       string        `json:"city"`
	Adcode     string        `json:"adcode"`
	Province   string        `json:"province"`
	ReportTime string        `json:"reporttime"`
	Casts      []WeatherCast `json:"casts"`
}

type WeatherCast struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	DayWeather   string `json:"dayweather"`
	NightWeather string `json:"nightweather"`
	DayTemp      string `json:"daytemp"`
	NightTemp    string `json:"nighttemp"`
	DayWind      string `json:"daywind"`
	NightWind    string `json:"nightwind"`
	DayPower     string `json:"daypower"`
	NightPower   string `json:"nightpower"`
}

// 全局配置变量
var config Config

// 加载配置文件
func loadConfig(configPath string) error {
	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON配置
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证必要的配置项
	if config.WeatherAPI.Key == "YOUR_AMAP_API_KEY" {
		fmt.Println("警告: 您尚未设置高德API密钥，请在config.json中设置有效的API密钥")
	}

	return nil
}

// 使用高德开放平台API获取天气
func getAmapWeather(user UserConfig) (string, bool, error) {
	fmt.Printf("开始获取%s的天气信息...\n", user.CityName)

	// 构建请求URL
	url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?key=%s&city=%s&extensions=all",
		config.WeatherAPI.Key, user.CityCode)

	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		return "", false, fmt.Errorf("获取天气信息失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 解析JSON响应
	var weatherResp WeatherResponse
	if err = json.Unmarshal(body, &weatherResp); err != nil {
		return "", false, fmt.Errorf("解析天气数据失败: %v", err)
	}

	// 检查API响应状态
	if weatherResp.Status != "1" {
		return "", false, fmt.Errorf("天气API返回错误: %s", weatherResp.Info)
	}

	// 获取实时天气
	url = fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?key=%s&city=%s&extensions=base",
		config.WeatherAPI.Key, user.CityCode)

	resp, err = http.Get(url)
	if err != nil {
		return "", false, fmt.Errorf("获取实时天气信息失败: %v", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("读取实时天气响应内容失败: %v", err)
	}

	var liveWeatherResp WeatherResponse
	if err := json.Unmarshal(body, &liveWeatherResp); err != nil {
		return "", false, fmt.Errorf("解析实时天气数据失败: %v", err)
	}

	// 格式化天气信息
	var weatherText string

	// 添加用户@信息
	weatherText += fmt.Sprintf("<at user_id=\"%s\">%s</at> 今日%s天气预报\n\n", user.OpenID, user.SendName, user.CityName)

	// 实时天气信息
	if len(liveWeatherResp.Lives) > 0 {
		live := liveWeatherResp.Lives[0]
		weatherText += "【实时天气】\n"
		weatherText += fmt.Sprintf("天气: %s\n", live.Weather)
		weatherText += fmt.Sprintf("温度: %s℃\n", live.Temperature)
		weatherText += fmt.Sprintf("风向: %s\n", live.WindDirection)
		weatherText += fmt.Sprintf("风力: %s级\n", live.WindPower)
		weatherText += fmt.Sprintf("湿度: %s%%\n", live.Humidity)
		weatherText += fmt.Sprintf("发布时间: %s\n\n", live.ReportTime)
	}

	// 天气预报
	if len(weatherResp.Forecasts) > 0 && len(weatherResp.Forecasts[0].Casts) > 0 {
		weatherText += "【未来天气预报】\n"

		// 获取今天和明天的预报
		forecasts := weatherResp.Forecasts[0].Casts

		// 今天预报
		if len(forecasts) > 0 {
			today := forecasts[0]
			weatherText += fmt.Sprintf("今天 (%s):\n", today.Date)
			weatherText += fmt.Sprintf("白天: %s %s℃ %s风 %s级\n",
				today.DayWeather, today.DayTemp, today.DayWind, today.DayPower)
			weatherText += fmt.Sprintf("夜间: %s %s℃ %s风 %s级\n\n",
				today.NightWeather, today.NightTemp, today.NightWind, today.NightPower)
		}

		// 明天预报
		if len(forecasts) > 1 {
			tomorrow := forecasts[1]
			weatherText += fmt.Sprintf("明天 (%s):\n", tomorrow.Date)
			weatherText += fmt.Sprintf("白天: %s %s℃ %s风 %s级\n",
				tomorrow.DayWeather, tomorrow.DayTemp, tomorrow.DayWind, tomorrow.DayPower)
			weatherText += fmt.Sprintf("夜间: %s %s℃ %s风 %s级\n\n",
				tomorrow.NightWeather, tomorrow.NightTemp, tomorrow.NightWind, tomorrow.NightPower)
		}
	}

	// 检查是否有恶劣天气预警
	hasWarning := shouldSendWarning(weatherText)

	// 使用AI生成温馨提示
	promptText, err := SendMoonshotChatRequest(weatherText)
	if err != nil {
		return "", hasWarning, err
	}
	// 添加温馨提示
	weatherText += promptText
	return weatherText, hasWarning, nil
}

// 发送Moonshot聊天请求
func SendMoonshotChatRequest(weatherText string) (string, error) {
	systemContent := "你是\"天气暖心助手\"，需根据实时天气数据、季节特征、地域习俗及用户潜在需求，生成个性化温馨提示。提示需满足：\n\n#### 一、数据解析层（必含）\n- 天气要素：温度（需标注昼夜温差）、天气现象（晴/雨/雪/雾等）、风力风向、湿度、紫外线指数\n- 季节判断：根据月份自动识别季节（春3-5月/夏6-8月/秋9-11月/冬12-2月），结合节气（如春分、冬至）强化场景感\n- 地理位置判断：用户发送信息中会携带详细的地理信息，你需要根据地理信息生成对应的东西\n#### 二、关怀维度（选填，按需组合）\n1. **健康防护**：\n   - 季节病预防（如春敏、秋燥）\n   - 特殊人群提示（老人保暖、儿童防着凉）\n2. **出行建议**：\n   - 交通方式推荐（雨天地铁优先、高温天打车）\n   - 装备提醒（防晒衣、冰袖、防滑链）\n3. **饮食调养**：\n   - 时令食材推荐（夏季苦瓜清热、冬季萝卜补气）\n   - 饮品建议（热姜茶驱寒、酸梅汤解暑）\n4. **情绪调节**：\n   - 天气与心情关联（如\"阴雨天适合听一首轻快的歌\"）\n   - 节气民俗互动（\"今日雨水，记得给亲友发去春日祝福\"）\n\n#### 三、表达规范\n- 语言风格：口语化+拟人化（避免生硬术语），适当使用emoji增强亲和力\n- 结构逻辑：先总述天气→再分场景建议→最后情感收尾\n- 地域适配：若已知用户所在城市（如成都），可融入当地特色（\"雨后的成都街头，来一碗热乎的龙抄手最惬意~\"）\n\n#### 四、极端天气预案\n- 暴雨/高温预警：添加紧急联系电话（如12121气象热线）和避险指南\n- 特殊天气：针对台风、沙尘暴等提供专业防护步骤\n#### 五、字数要求不超过50字。"

	return kimi.GetAiText(systemContent, weatherText)
}

// 判断字符串是否包含指定字符串列表中的任意一个
func contains(s string, substrs []string) bool {
	for _, substr := range substrs {
		if bytes.Contains([]byte(s), []byte(substr)) {
			return true
		}
	}
	return false
}

// 检查是否需要发送天气预警
func shouldSendWarning(weatherText string) bool {
	// 定义需要预警的天气关键词
	warningKeywords := []string{
		"暴雨", "大雨", "雷雨", "雷电", "台风", "大风", "沙尘暴",
		"雾霾", "大雾", "霜冻", "冰雹", "暴雪", "大雪", "寒潮",
		"高温", "酷热", "极端", "预警", "警报", "雨", "雪",
	}

	return contains(weatherText, warningKeywords)
}

// 发送消息到所有配置的飞书webhook
func sendToAllFeishuWebhooks(message string) error {
	if len(config.Feishu) == 0 {
		fmt.Println("未配置飞书Webhook URL，跳过消息发送")
		return nil
	}

	var errors []string
	for i, webhookURL := range config.Feishu {
		if webhookURL == "" {
			continue
		}
		if err := feishu.SendStringToFeishu(webhookURL, message); err != nil {
			errorMsg := fmt.Sprintf("发送到第%d个飞书Webhook失败: %v", i+1, err)
			fmt.Println(errorMsg)
			errors = append(errors, errorMsg)
		} else {
			fmt.Printf("消息已成功发送到第%d个飞书Webhook\n", i+1)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("部分消息发送失败: %s", strings.Join(errors, "; "))
	}
	return nil
}

func main() {
	// 解析命令行参数
	allWeather := flag.Bool("all", false, "强制发送天气信息，忽略时间检查和预警检测")
	configFile := flag.String("config", "config.json", "配置文件路径")
	flag.Parse()

	// 读取配置文件
	if err := loadConfig(*configFile); err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 获取当前时间
	now := time.Now()
	currentHour := now.Hour()
	currentMinute := now.Minute()

	// 为每个用户检查时间并处理天气信息
	for _, user := range config.Users {
		// 如果不是强制执行，检查该用户的时间
		if !*allWeather && !(user.Hour == currentHour && user.Minute == currentMinute) {
			// 跳过不匹配时间的用户
			continue
		}

		fmt.Printf("开始为%s获取天气信息...\n", user.SendName)

		weatherText, hasWarning, err := getAmapWeather(user)
		if err != nil {
			fmt.Printf("获取%s的天气信息失败: %v\n", user.CityName, err)
			continue
		}

		fmt.Printf("%s的天气信息获取成功:\n%s\n", user.CityName, weatherText)

		// 检查是否需要发送预警或强制发送
		if *allWeather {
			fmt.Printf("使用-all参数，强制发送%s的天气信息...\n", user.SendName)
			sendWeather := true
			if sendWeather {
				if err := sendToAllFeishuWebhooks(weatherText); err != nil {
					fmt.Printf("发送%s的消息到飞书失败: %v\n", user.SendName, err)
				} else {
					fmt.Printf("%s的天气信息发送成功\n", user.SendName)
				}
			}
		} else if hasWarning {
			fmt.Printf("检测到%s有恶劣天气，发送预警消息...\n", user.SendName)
			if err := sendToAllFeishuWebhooks(weatherText); err != nil {
				fmt.Printf("发送%s的预警消息到飞书失败: %v\n", user.SendName, err)
			} else {
				fmt.Printf("%s的预警消息发送成功\n", user.SendName)
			}
		} else {
			fmt.Printf("%s天气正常，发送日常天气信息...\n", user.SendName)
			if err := sendToAllFeishuWebhooks(weatherText); err != nil {
				fmt.Printf("发送%s的消息到飞书失败: %v\n", user.SendName, err)
			} else {
				fmt.Printf("%s的天气信息发送成功\n", user.SendName)
			}
		}
	}

	fmt.Println("天气信息处理完成")
}
