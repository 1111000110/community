# End 命令使用说明

## 概述

`end` 命令用于停止服务、清理数据和资源。采用动态服务发现机制，无需硬编码服务列表，支持灵活的清理选项。

## 命令结构

```
xuan end [子命令] [参数]
```

## 子命令

### 1. all - 停止所有服务

```bash
# 停止所有服务（保留数据）
xuan end all

# 停止所有服务并清理数据
xuan end all --clean-data
xuan end all -c
```

**功能：**
- 停止 Docker Compose 服务
- 停止 etcd 服务
- 动态发现并停止所有业务服务
- 可选清理所有数据

### 2. docker - Docker 服务管理

```bash
# 停止所有 Docker 服务
xuan end docker

# 停止特定 Docker 服务
xuan end docker <service-name>
```

### 3. api - API 服务管理

```bash
# 停止指定目录的 API 服务
xuan end api <service-path>
```

### 4. rpc - RPC 服务管理

```bash
# 停止指定目录的 RPC 服务
xuan end rpc <service-path>
```

### 5. etcd - Etcd 服务管理

```bash
# 停止 etcd 服务
xuan end etcd
```

### 6. clean - 数据清理

```bash
# 清理所有数据
xuan end clean

# 清理特定类型的数据
xuan end clean logs      # 清理日志文件
xuan end clean data      # 清理数据目录
xuan end clean binaries  # 清理二进制文件
xuan end clean all       # 清理所有数据
```

## 动态服务发现

### 服务发现机制

`end` 命令会自动扫描 `service/` 目录，发现所有 `api` 和 `rpc` 子目录，无需手动配置服务列表。

**目录结构示例：**
```
service/
├── user/
│   ├── api/     # 自动发现为 user-api
│   └── rpc/     # 自动发现为 user-rpc
├── post/
│   ├── api/     # 自动发现为 post-api
│   └── rpc/     # 自动发现为 post-rpc
└── message/
    ├── api/     # 自动发现为 message-api
    └── rpc/     # 自动发现为 message-rpc
```

### 服务命名规则

服务名称遵循 `{parent-dir}-{sub-dir}` 的命名规则：
- `service/user/api` → `user-api`
- `service/user/rpc` → `user-rpc`
- `service/post/api` → `post-api`

## 清理功能详解

### 日志清理 (logs)

清理以下日志文件：
- 根目录的 `*.log` 文件
- `etcd.log` 文件
- `service/` 目录下的所有 `*.log` 文件

### 数据目录清理 (data)

清理以下数据目录：
- `etcd-data/` - etcd 数据目录
- `docker-data/` - Docker 数据目录

### 二进制文件清理 (binaries)

清理以下二进制文件：
- `service/` 目录下所有可执行的 `*-api` 和 `*-rpc` 文件

## 使用示例

### 开发环境清理

```bash
# 停止所有服务，保留数据（用于快速重启）
xuan end all

# 完全清理开发环境
xuan end all --clean-data
```

### 选择性清理

```bash
# 只清理日志
xuan end clean logs

# 只清理数据目录
xuan end clean data

# 只清理二进制文件
xuan end clean binaries
```

### 服务管理

```bash
# 停止特定服务的 API
xuan end api ./service/user

# 停止特定服务的 RPC
xuan end rpc ./service/user

# 停止 etcd
xuan end etcd
```

## 优势特性

### 1. 动态发现
- 无需硬编码服务列表
- 自动适应项目结构变化
- 支持任意数量的服务

### 2. 灵活清理
- 支持选择性清理
- 参数化控制清理范围
- 安全的清理操作

### 3. 错误容忍
- 服务不存在时不会报错
- 清理失败时继续执行其他操作
- 详细的执行日志

### 4. 扩展性
- 易于添加新的清理类型
- 支持自定义清理逻辑
- 模块化设计

## 注意事项

1. **数据安全**: 使用 `--clean-data` 参数会删除所有数据，请谨慎使用
2. **权限要求**: 某些操作可能需要适当的文件系统权限
3. **服务状态**: 确保在停止服务前保存重要数据
4. **依赖关系**: 停止服务时注意服务间的依赖关系

## 故障排除

### 常见问题

1. **权限不足**: 确保有足够的权限删除文件和目录
2. **服务仍在运行**: 某些服务可能需要强制停止
3. **目录不存在**: 这是正常情况，命令会跳过不存在的目录

### 调试模式

命令会输出详细的执行信息，包括：
- 正在停止的服务
- 正在清理的文件
- 操作结果和错误信息
