# Weather 天气服务

## 概述

Weather服务是社区平台的天气信息和预警管理系统，提供实时天气查询、智能天气预警、个性化订阅等功能。通过集成高德地图API获取天气数据，使用AI生成个性化关怀提示，并通过飞书机器人发送通知。

## 主要功能

### 🌤️ 天气信息服务
- **实时天气**：获取当前天气状况、温度、湿度等信息
- **天气预报**：提供多日天气预报数据
- **城市支持**：支持全国各大城市的天气查询
- **数据缓存**：优化API调用频率，提高响应速度

### ⚠️ 智能预警系统
- **自动检测**：智能识别恶劣天气条件
- **预警关键词**：暴雨、大风、雾霾、高温、台风等预警识别
- **实时通知**：发现预警天气立即推送通知
- **多渠道通知**：支持飞书群组等多渠道预警通知

### 🤖 AI智能提示
- **个性化关怀**：基于用户个人信息生成温馨提示
- **智能建议**：根据天气情况提供出行、穿衣等生活建议
- **内容优化**：控制提示长度，适配移动端显示

### 👥 用户订阅管理
- **个性化配置**：每个用户可独立配置城市和通知时间
- **多地点支持**：不同用户可订阅不同城市的天气
- **定时推送**：支持精确到分钟的个性化定时推送
- **订阅数据**：MongoDB存储用户订阅和偏好数据

## 技术架构

- **数据获取**：高德地图天气API提供实时数据
- **AI服务**：Moonshot AI生成个性化天气提示
- **消息通知**：飞书机器人Webhook发送通知
- **数据存储**：MongoDB存储用户订阅数据
- **定时任务**：支持cron表达式和精确时间控制

## 定时任务系统

### Weather Warning Detection
独立运行的定时任务系统，负责：
- 遍历所有用户配置
- 检查执行时间匹配
- 获取对应城市天气数据
- 检测恶劣天气预警
- 生成AI关怀提示
- 发送飞书通知

### 执行模式
- **定时模式**：按用户配置的时间精确执行
- **强制模式**：使用`-all`参数忽略时间检查，立即执行所有用户
- **热重载**：配置文件修改立即生效，无需重启

## 数据模型

### 用户订阅数据 (MongoDB)
```javascript
{
  userId: String,        // 用户ID
  cityCode: String,      // 城市编码
  cityName: String,      // 城市名称
  sendName: String,      // 用户姓名
  openId: String,        // 飞书OpenID
  hour: Number,          // 执行小时
  minute: Number,        // 执行分钟
  isActive: Boolean,     // 是否激活
  createTime: Date,      // 创建时间
  updateTime: Date       // 更新时间
}
```

## 预警检测机制

系统自动检测以下天气关键词并触发预警：
- **降水类**：暴雨、大雨、雷雨、雷电
- **风类**：台风、大风、沙尘暴
- **能见度**：雾霾、大雾
- **温度类**：霜冻、寒潮、高温、酷热
- **降雪类**：暴雪、大雪、冰雹
- **其他**：极端、预警、警报

## 配置管理

### 服务配置
```yaml
Name: weather.rpc
Host: 0.0.0.0
Port: 8891
WeatherAPI:
  Key: "your_amap_api_key"
MongoDB:
  URL: "mongodb://localhost:27017"
  Database: "weather"
```

### 定时任务配置
```json
{
  "weather_api": {
    "key": "your_amap_api_key"
  },
  "users": [
    {
      "send_name": "张璇",
      "open_id": "ou_xxx",
      "city_code": "110105",
      "city_name": "北京市朝阳区",
      "hour": 8,
      "minute": 0
    }
  ],
  "feishu_webhooks": [
    "https://open.feishu.cn/open-apis/bot/v2/hook/xxx"
  ]
}
```

## 部署方式

### RPC服务部署
```bash
# 启动天气RPC服务
go run weather.go
```

### 定时任务部署
```bash
# 编译定时任务
go build -o weather-warning-detection cronjob/weather-warning-detection/main.go

# 设置crontab（每分钟执行一次）
* * * * * /path/to/weather-warning-detection -config /path/to/config.json
```

## 应用场景

- **社区天气提醒**：为社区用户提供天气关怀服务
- **智能生活助手**：结合AI提供个性化生活建议
- **应急预警系统**：及时发现恶劣天气并推送预警
- **个性化订阅**：满足不同用户对不同城市天气的需求
