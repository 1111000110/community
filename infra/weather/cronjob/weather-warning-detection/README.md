# 天气预警检测系统

这是一个基于Go语言开发的天气预警检测系统，能够自动获取天气信息，使用AI生成温馨提示，并通过飞书机器人发送通知。

## 功能特性

- 🌤️ **实时天气获取**: 使用高德地图API获取实时天气和预报信息
- 🤖 **AI智能提示**: 集成Moonshot AI生成个性化天气关怀提示
- 📱 **飞书通知**: 支持通过飞书Webhook发送天气信息
- ⚠️ **预警检测**: 自动识别恶劣天气并发送预警通知
- ⚙️ **配置灵活**: 支持JSON配置文件，易于部署和管理

## 配置说明

### config.json 配置文件

```json
{
  "weather_api": {
    "key": "YOUR_AMAP_API_KEY",
    "city_code": "110105",
    "city_name": "北京市朝阳区",
    "send_name": "张璇",
    "open_id": "ou_a4d2c263956ef162bd72bb3030500769"
  },
  "feishu": {
    "webhook_url": "YOUR_FEISHU_WEBHOOK_URL"
  }
}
```

### 配置项说明

- `weather_api.key`: 高德地图API密钥
- `weather_api.city_code`: 城市编码（如北京朝阳区：110105）
- `weather_api.city_name`: 城市名称
- `weather_api.send_name`: 发送者姓名
- `weather_api.open_id`: 飞书用户OpenID
- `feishu.webhook_url`: 飞书机器人Webhook地址

## 使用方法

### 1. 编译程序

```bash
cd /Users/zhangxuan/Data/work/xuan/community/infra/weather/cronjob/weather-warning-detection
go build -o weather-warning-detection main.go
```

### 2. 运行程序

```bash
# 使用默认配置文件 config.json
./weather-warning-detection

# 指定配置文件路径
./weather-warning-detection -config /path/to/your/config.json
```

### 3. 定时任务设置

可以使用crontab设置定时任务，例如每天早上8点执行：

```bash
# 编辑crontab
crontab -e

# 添加定时任务（每天8:00执行）
0 8 * * * /path/to/weather-warning-detection -config /path/to/config.json
```

## 预警关键词

系统会自动检测以下天气关键词并发送预警：

- 降水类：暴雨、大雨、雷雨、雷电
- 风类：台风、大风、沙尘暴
- 能见度：雾霾、大雾
- 温度类：霜冻、寒潮、高温、酷热
- 降雪类：暴雪、大雪、冰雹
- 其他：极端、预警、警报

## 依赖说明

- 高德地图开放平台API
- Moonshot AI API
- 飞书机器人Webhook

## 注意事项

1. 请确保高德地图API密钥有效且有足够的调用次数
2. Moonshot AI API密钥已在代码中配置，如需更换请修改源码
3. 飞书Webhook URL需要在飞书群中创建机器人获取
4. 程序会自动生成温馨的天气提示，字数控制在50字以内

## 错误处理

程序包含完善的错误处理机制：
- 配置文件读取失败
- API调用失败
- 网络连接问题
- JSON解析错误

所有错误都会输出详细的错误信息，便于调试和排查问题。