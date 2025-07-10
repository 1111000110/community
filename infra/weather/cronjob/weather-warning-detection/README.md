# 天气预警检测系统

这是一个基于Go语言开发的天气预警检测系统，能够自动获取天气信息，使用AI生成温馨提示，并通过飞书机器人发送通知。

## 功能特性

- 🌤️ **实时天气获取**: 使用高德地图API获取实时天气和预报信息
- 🤖 **AI智能提示**: 集成Moonshot AI生成个性化天气关怀提示
- 📱 **多群组通知**: 支持同时发送到多个飞书群组
- 👥 **多用户独立配置**: 每个用户可配置独立的城市、时间和个人信息
- 🌍 **多地点支持**: 不同用户可获取不同城市的天气信息
- ⚠️ **智能预警**: 自动识别恶劣天气并发送预警通知
- ⏰ **个性化定时**: 每个用户可设置独立的执行时间，精确到分钟
- 🔄 **配置热重载**: 修改配置文件立即生效，无需重启
- 🚀 **强制执行**: 支持-all参数强制发送，忽略时间和预警检测
- ⚙️ **配置灵活**: 支持JSON配置文件，易于部署和管理

## 配置说明

### config.json 配置文件

``` json
{
  "weather_api": {
    "key": "726915a2274707078bf0db65ee067322"
  },
  "users": [
    {
      "send_name": "张璇",
      "open_id": "ou_YOU——OPEN——ID",
      "city_code": "110105",
      "city_name": "北京市朝阳区",
      "hour": 8,
      "minute": 0
    }
  ],
  "feishu_webhooks": [
    "FEISHU_WEBHOOK"
  ]
}
```

### 配置项说明

- `weather_api.key`: 高德地图API密钥
- `users`: 用户列表，支持多个用户，每个用户可配置独立的城市和时间
  - `send_name`: 用户姓名
  - `open_id`: 飞书用户OpenID
  - `city_code`: 城市编码（如北京朝阳区：110105）
  - `city_name`: 城市名称
  - `hour`: 小时（0-23）
  - `minute`: 分钟（0-59）
- `feishu_webhooks`: 飞书机器人Webhook地址列表，支持多个群组

## 使用方法

### 1. 编译程序

```bash
cd /Users/zhangxuan/Data/work/xuan/community/infra/weather/cronjob/weather-warning-detection
go build -o weather-warning-detection main.go
```

### 2. 运行程序

```bash
# 使用默认配置文件 config.json（仅在用户配置的时间点执行）
./weather-warning-detection

# 指定配置文件路径
./weather-warning-detection -config /path/to/your/config.json

# 强制执行所有用户，忽略时间检查和预警检测
./weather-warning-detection -all

# 组合使用
./weather-warning-detection -config /path/to/config.json -all
```

### 3. 定时任务设置

**推荐配置：每分钟执行一次，程序内部为每个用户根据配置的时间点决定是否真正执行**

```bash
# 编辑crontab
crontab -e

# 添加定时任务（每分钟执行一次，程序内部为每个用户判断时间）
* * * * * /path/to/weather-warning-detection -config /path/to/config.json >/dev/null 2>&1

# 如果需要日志记录
* * * * * /path/to/weather-warning-detection -config /path/to/config.json >> /var/log/weather.log 2>&1
```

**传统配置：如果所有用户时间相同，可直接在crontab中设置具体时间**

```bash
# 每天8:00和18:00执行（适用于所有用户时间相同的情况）
0 8,18 * * * /path/to/weather-warning-detection -config /path/to/config.json -all
```

### 4. 工作原理

- 程序每次运行时会遍历所有用户配置
- 对于每个用户，检查当前时间是否匹配该用户的执行时间（hour:minute）
- 匹配的用户会获取其对应城市的天气信息并发送
- 使用`-all`参数时，忽略时间检查，为所有用户发送天气信息
- 每个用户可以设置不同的城市和执行时间，实现个性化配置

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