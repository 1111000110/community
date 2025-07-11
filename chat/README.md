# Chat 聊天模块

## 概述

Chat模块是社区系统的实时通讯模块，负责用户之间的消息传递、群聊管理以及聊天记录的存储和检索。该模块采用WebSocket技术实现实时通讯，支持单聊和群聊功能。

## 模块结构

```
chat/
├── api/                    # HTTP API服务
│   ├── chat.api           # API接口定义文件
│   ├── chat.go            # API服务启动入口
│   ├── etc/               # 配置文件
│   │   └── chat.yaml      # 服务配置
│   └── internal/          # 内部实现
│       ├── config/        # 配置结构定义
│       ├── handler/       # HTTP处理器
│       ├── logic/         # 业务逻辑层
│       ├── middleware/    # 中间件
│       ├── svc/          # 服务上下文
│       └── types/        # 类型定义
└── etc/                   # 全局配置
    └── chat.yaml         # 聊天服务配置
```

## 功能特性

### 💬 实时通讯
- WebSocket连接：支持实时双向通讯
- 消息推送：即时消息传递
- 连接管理：自动重连和心跳检测

### 📝 消息管理
- 消息发送：支持文本消息发送
- 消息接收：实时接收新消息
- 消息历史：支持历史消息查询和分页
- 上拉刷新：获取更多历史消息

### 👥 聊天类型
- 单聊：用户一对一私聊
- 群聊：多用户群组聊天
- 消息分组：按群组ID组织消息

### 📊 数据模型
- **消息结构**：群组ID、消息ID、发送者、接收者、消息内容、创建时间
- **消息分组**：按群组ID对消息进行分类管理
- **用户会话**：维护用户的聊天会话状态

## API接口

### WebSocket连接
- `GET /chat` - 建立WebSocket连接
  - 参数：`userId` - 用户ID
  - 功能：建立实时通讯连接

### 消息接口
- 消息接收：通过WebSocket实时接收
- 消息发送：通过WebSocket实时发送
- 历史消息：支持分页查询历史记录

## 数据结构

### Message 消息结构
```go
type Message struct {
    GroupId    int64  `json:"groupId"`    // 群组ID（0表示私聊）
    MessageId  int64  `json:"messageId"`  // 消息唯一ID
    FromUserId int64  `json:"fromUserId"` // 发送者用户ID
    ToUserId   int64  `json:"toUserId"`   // 接收者用户ID
    Text       string `json:"text"`       // 消息文本内容
    CreateDate int64  `json:"createDate"` // 创建时间戳
}
```

### MessageReciveReq 消息接收请求
```go
type MessageReciveReq struct {
    MessageId int64 `json:"messageId"` // 起始消息ID（用于分页）
    UserId    int64 `json:"userId"`    // 用户ID
    Limit     int64 `json:"limit"`     // 获取消息数量限制
}
```

### MessageReciveResp 消息接收响应
```go
type MessageReciveResp struct {
    Message map[int64][]Message `json:"message"` // 按群组ID分组的消息列表
}
```

## WebSocket通讯协议

### 连接建立
```javascript
// 前端连接示例
const ws = new WebSocket('ws://localhost:8888/chat?userId=123');
```

### 消息格式
```json
{
  "groupId": 0,
  "messageId": 1001,
  "fromUserId": 123,
  "toUserId": 456,
  "text": "Hello, World!",
  "createDate": 1640995200000
}
```

### 消息类型
- **私聊消息**：`groupId = 0`，指定 `toUserId`
- **群聊消息**：`groupId > 0`，`toUserId` 可为0
- **系统消息**：特殊的消息类型，用于通知

## 配置说明

### 服务配置 (etc/chat.yaml)
```yaml
Name: chat.api
Host: 0.0.0.0
Port: 8888
Timeout: 3s
MaxBytes: 1048576

# WebSocket配置
WebSocket:
  ReadBufferSize: 1024
  WriteBufferSize: 1024
  CheckOrigin: true

# 中间件配置
Middleware:
  - Cors
  - Auth
```

## 部署运行

### 启动聊天服务
```bash
cd chat/api
go run chat.go -f etc/chat.yaml
```

### Docker部署
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o chat chat.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/chat .
COPY --from=builder /app/etc ./etc
EXPOSE 8888
CMD ["./chat", "-f", "etc/chat.yaml"]
```

## 开发指南

### 添加新的消息类型
1. 在 `chat.api` 中定义新的消息结构
2. 运行 `goctl api go` 生成代码
3. 在 `logic` 目录下实现消息处理逻辑
4. 在 `handler` 目录下处理WebSocket消息

### 扩展群聊功能
1. 定义群组管理相关的API接口
2. 实现群组创建、加入、退出逻辑
3. 添加群组权限管理
4. 实现群组消息广播

### 消息持久化
1. 设计消息存储数据库表结构
2. 实现消息的增删改查操作
3. 添加消息索引优化查询性能
4. 实现消息的分页和搜索功能

## 性能优化

### 连接管理
- 连接池：复用WebSocket连接
- 心跳检测：定期检查连接状态
- 自动重连：网络断开时自动重连

### 消息处理
- 消息队列：使用Redis或RabbitMQ缓冲消息
- 批量处理：批量发送和接收消息
- 消息压缩：大消息内容压缩传输

### 数据存储
- 分库分表：按用户或时间分片存储
- 缓存策略：热点消息Redis缓存
- 归档策略：历史消息定期归档

## 安全考虑

### 身份验证
- JWT令牌：WebSocket连接需要有效令牌
- 用户权限：验证用户发送消息的权限
- 频率限制：防止消息轰炸攻击

### 内容安全
- 消息过滤：敏感词汇过滤
- 内容审核：自动或人工审核机制
- 消息加密：敏感消息端到端加密

## 监控告警

### 关键指标
- 在线用户数：实时连接数统计
- 消息吞吐量：每秒消息发送接收量
- 连接稳定性：连接断开重连率
- 响应时间：消息传递延迟

### 告警规则
- 连接数异常：超过阈值告警
- 消息堆积：消息队列积压告警
- 服务异常：服务不可用告警

## 依赖项

- go-zero: 微服务框架
- gorilla/websocket: WebSocket库
- redis: 消息缓存和会话存储
- mysql: 消息持久化存储

## 注意事项

1. WebSocket连接需要处理网络断开和重连
2. 消息ID需要保证全局唯一性
3. 群聊消息需要考虑广播性能
4. 历史消息查询需要分页避免大量数据传输
5. 需要实现消息的已读未读状态管理

## 联系方式

- 作者：张璇
- 邮箱：xatuzx2025@163.com
- 版本：1.0