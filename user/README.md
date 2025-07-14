# user 用户模块

## 概述

user模块是社区系统的核心用户管理模块，负责用户的注册、登录、信息管理以及用户关系管理等功能。该模块采用微服务架构，提供HTTP API和gRPC服务两种接口形式。

## 模块结构

```
user/
├── api/                    # HTTP API服务
│   ├── user.api        # API接口定义文件
│   ├── user.go         # API服务启动入口
│   ├── etc/               # 配置文件
│   └── internal/          # 内部实现
│       ├── config/        # 配置结构定义
│       ├── handler/       # HTTP处理器
│       ├── logic/         # 业务逻辑层
│       ├── middleware/    # 中间件
│       ├── svc/          # 服务上下文
│       └── types/        # 类型定义
├── job/                   # 定时任务
│   └── weather/          # 天气相关任务
├── model/                 # 数据模型
│   ├── mongo/            # MongoDB模型
│   └── mysql/            # MySQL模型
│       └── user/      # 用户账户模型
└── rpc/                   # gRPC服务
    ├── user.proto     # Protobuf定义文件
    ├── user.go        # RPC服务启动入口
    ├── client/           # RPC客户端
    ├── etc/              # 配置文件
    ├── internal/         # 内部实现
    └── pb/               # 生成的protobuf代码
```

## 功能特性

### 🔐 用户认证
- 用户注册：支持手机号注册
- 用户登录：手机号+密码登录
- JWT令牌认证：安全的身份验证机制

### 👤 用户管理
- 用户信息查询：获取用户基本信息和私密信息
- 用户信息更新：修改用户资料
- 用户删除：支持用户账户删除

### 👥 用户关系
- 关系建立：用户之间建立各种关系（好友、关注等）
- 关系查询：查询用户间的关系状态
- 关系更新：修改用户关系类型

### 📊 数据模型
- **用户基本信息**：用户ID、昵称、头像、性别、生日等
- **用户私密信息**：手机号、邮箱、密码、角色、状态等
- **用户关系**：关系ID、用户ID、关系类型、创建时间等

## API接口

### 用户认证接口
- `POST /user/login` - 用户登录
- `POST /user/register` - 用户注册

### 用户管理接口（需JWT认证）
- `POST /user/delete` - 删除用户
- `POST /user/update` - 更新用户信息
- `POST /user/query` - 查询用户信息

### 用户关系接口（需JWT认证）
- `POST /userRelations/update` - 更新用户关系
- `POST /userRelations/get` - 获取用户关系

## gRPC服务

### UserService 用户服务
- `UserLogin` - 用户登录
- `UserRegister` - 用户注册
- `UserDelete` - 删除用户
- `UserUpdate` - 更新用户
- `UserQuery` - 查询用户

### RelationsService 关系服务
- `UserRelationsUpdate` - 更新用户关系
- `UserRelationsGet` - 获取用户关系

## 配置说明

### API服务配置 (api/etc/user.yaml)
```yaml
Name: user.api
Host: 0.0.0.0
Port: 8888
Auth:
  AccessSecret: "your_jwt_secret"
  AccessExpire: 1000
```

### RPC服务配置 (rpc/etc/user.yaml)
```yaml
Name: user.rpc
ListenOn: 0.0.0.0:8081
Etcd:
  Hosts:
    - etcd:2379
  Key: user.rpc
```

## 部署运行

### 启动API服务
```bash
cd user/api
go run user.go -f etc/user.yaml
```

### 启动RPC服务
```bash
cd user/rpc
go run user.go -f etc/user.yaml
```

## 数据库设计

### user表结构
```sql
CREATE TABLE user (
    user_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    password VARCHAR(255) NOT NULL COMMENT '用户密码',
    phone VARCHAR(255) NOT NULL COMMENT '手机号码',
    gender CHAR(10) DEFAULT 'unknown' COMMENT '性别',
    nickname VARCHAR(255) DEFAULT '默认名称' COMMENT '昵称',
    avatar VARCHAR(255) DEFAULT '' COMMENT '头像URL',
    birth_date BIGINT DEFAULT 0 COMMENT '出生日期（时间戳）',
    role VARCHAR(50) DEFAULT '' COMMENT '用户角色',
    status BIGINT DEFAULT 0 COMMENT '用户状态',
    email VARCHAR(255) DEFAULT '' COMMENT '电子邮箱',
    create_at BIGINT DEFAULT 0 COMMENT '创建时间',
    update_at BIGINT DEFAULT 0 COMMENT '更新时间'
);
```

## 开发指南

### 添加新接口
1. 在 `user.api` 中定义新的接口
2. 运行 `goctl api go` 生成代码
3. 在 `logic` 目录下实现业务逻辑
4. 在 `handler` 目录下处理HTTP请求

### 添加新的gRPC服务
1. 在 `user.proto` 中定义新的服务和消息
2. 运行 `goctl rpc protoc` 生成代码
3. 在 `logic` 目录下实现业务逻辑
4. 在 `server` 目录下实现gRPC服务

## 依赖项

- go-zero: 微服务框架
- gorm: ORM框架
- jwt-go: JWT认证
- etcd: 服务发现
- mysql: 数据库

## 注意事项

1. 所有涉及用户隐私的接口都需要JWT认证
2. 密码存储使用bcrypt加密
3. 手机号作为用户唯一标识
4. 用户状态字段用于软删除和账户状态管理
5. 关系类型需要根据业务需求定义具体含义

## 联系方式

- 作者：张璇
- 邮箱：xatuzx2025@163.com
- 版本：1.0