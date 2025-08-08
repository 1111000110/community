# Infra 基础设施模块

## 概述

Infra模块是社区系统的基础设施服务集合，包含各种支撑性的微服务。这些服务主要提供gRPC接口，不直接对外提供HTTP API，专注于为其他业务模块提供基础功能支持。

## 模块结构

```
infra/
└── weather/                # 天气服务模块
    ├── cronjob/           # 定时任务
    │   └── weather-warning-detection/  # 天气预警检测任务
    │       ├── README.md  # 任务说明文档
    │       ├── config.json # 配置文件
    │       ├── main.go    # 主程序入口
    │       └── weather-warning-detection # 编译后的可执行文件
    ├── model/             # 数据模型
    │   └── mongo/         # MongoDB数据模型
    │       └── weather/   # 天气数据模型
    └── rpc/               # gRPC服务
        ├── weather.proto  # Protobuf定义文件
        ├── weather.go     # RPC服务启动入口
        ├── etc/           # 配置文件
        ├── internal/      # 内部实现
        └── pb/            # 生成的protobuf代码
```

## 服务列表

### 🌤️ Weather 天气服务

天气服务是一个综合性的气象信息管理系统，提供天气数据收集、存储、查询和预警功能。

#### 主要功能
- **天气数据收集**：从第三方API获取实时天气信息
- **用户订阅管理**：管理用户的天气关注城市
- **定时推送**：定时向用户推送天气预警信息
- **多平台通知**：支持飞书等多种通知渠道

#### 技术特性
- **多用户支持**：每个用户可配置独立的城市和推送时间
- **多城市支持**：支持全国各地城市的天气查询
- **个性化定时**：用户可自定义接收天气信息的时间
- **智能推送**：根据天气变化智能推送预警信息

#### 服务接口

##### WeatherService gRPC服务
- `WeatherAddData` - 添加天气数据
  - 请求：`WeatherAddDataReq`
  - 响应：`WeatherAddDataResp`
  - 功能：将用户的天气关注信息存储到数据库

#### 数据模型

##### WeatherAddDataReq 天气数据请求
```protobuf
message WeatherAddDataReq {
  string UserName = 1;        // 用户姓名
  string OpenId = 2;          // 用户OpenID（飞书等平台）
  string UserId = 3;          // 用户ID
  string City = 4;            // 关注城市
  string Time = 5;            // 推送时间
  string MaxTemperature = 6;  // 最高温度
  string MinTemperature = 7;  // 最低温度
  repeated string Weather = 8; // 天气状况列表
  int64 Status = 9;           // 状态标识
}
```

#### 定时任务

##### weather-warning-detection 天气预警检测

这是一个独立的定时任务程序，负责定期检查天气信息并向用户发送预警通知。

**主要特性：**
- **多用户独立配置**：每个用户可配置自己的城市、时间和个人信息
- **多地点支持**：支持不同城市的天气信息获取
- **个性化定时**：每个用户独立的执行时间安排
- **智能执行**：程序内部检查各用户时间，自动执行相应用户的天气信息获取和发送

**配置示例：**
```json
{
  "users": [
    {
      "open_id": "ou_xxx",
      "send_name": "张璇",
      "city_code": "110105",
      "city_name": "北京市朝阳区",
      "hour": 8,
      "minute": 0
    }
  ],
  "amap_key": "your_amap_api_key",
  "webhook_urls": ["https://your-webhook-url"]
}
```

**运行方式：**
```bash
# 检查所有用户的时间配置并执行
./weather-warning-detection

# 强制为所有用户执行（忽略时间检查）
./weather-warning-detection -all
```

**推荐部署：**
```bash
# 添加到crontab，每分钟执行一次
* * * * * /path/to/weather-warning-detection
```

## 部署运行

### Weather RPC服务
```bash
cd infra/weather/rpc
go run weather.go -f etc/weather.yaml
```

### 天气预警定时任务
```bash
cd infra/weather/cronjob/weather-warning-detection
# 编译
go build -o weather-warning-detection main.go
# 运行
./weather-warning-detection
```

## 配置说明

### Weather RPC配置
```yaml
Name: weather.rpc
ListenOn: 0.0.0.0:8082
Etcd:
  Hosts:
    - etcd:2379
  Key: weather.rpc

# MongoDB配置
Mongo:
  Url: mongodb://localhost:27017
  Database: community
```

### 定时任务配置
详见 `weather/cronjob/weather-warning-detection/README.md`

## 开发指南

### 添加新的基础服务
1. 在 `infra` 目录下创建新的服务目录
2. 定义 `.proto` 文件描述服务接口
3. 使用 `goctl rpc protoc` 生成代码框架
4. 实现业务逻辑
5. 配置服务注册和发现

### 扩展Weather服务
1. 在 `weather.proto` 中添加新的服务方法
2. 重新生成protobuf代码
3. 在 `logic` 目录下实现新的业务逻辑
4. 更新数据模型和数据库操作

## 监控告警

### 关键指标
- **服务可用性**：gRPC服务健康状态
- **任务执行状态**：定时任务执行成功率
- **数据准确性**：天气数据获取成功率
- **通知送达率**：消息推送成功率

### 告警规则
- 服务不可用告警
- 定时任务执行失败告警
- 第三方API调用失败告警
- 数据库连接异常告警

## 依赖项

- go-zero: 微服务框架
- MongoDB: 数据存储
- etcd: 服务发现
- 高德地图API: 天气数据源
- 飞书API: 消息推送

## 注意事项

1. **服务依赖**：基础服务应保持高可用性
2. **数据一致性**：注意跨服务的数据一致性
3. **API限制**：注意第三方API的调用频率限制
4. **错误处理**：完善的错误处理和重试机制
5. **配置管理**：敏感信息使用环境变量或配置中心

## 联系方式

- 作者：张璇
- 邮箱：xatuzx2025@163.com
- 版本：1.0