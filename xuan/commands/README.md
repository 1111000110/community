# Xuan 扩展性命令系统

## 概述

Xuan 命令系统采用模块化设计，支持命令的层级结构和扩展性。每个命令对应一个目录，便于管理和维护。

## 目录结构

```
commands/
├── gen/          # 代码生成命令
├── run/          # 运行/重启命令
├── end/          # 停止命令
├── test/         # 测试命令（示例）
└── README.md     # 本文档
```

## 如何添加新命令

### 1. 创建命令目录

```bash
mkdir commands/your-command
```

### 2. 实现命令结构

在 `commands/your-command/command.go` 中实现：

```go
package yourcommand

import (
    "community/xuan/types"
    "fmt"
)

// YourCommand 你的命令
type YourCommand struct {
    *types.BaseCommand
}

func NewYourCommand() *YourCommand {
    cmd := &YourCommand{
        BaseCommand: types.NewBaseCommand("your-command", "命令描述", 1, nil),
    }
    
    // 添加子命令（可选）
    subCmd := NewSubCommand(cmd)
    cmd.AddSon(subCmd)
    
    return cmd
}

// 如果需要支持参数，实现 ExecuteWithArgs 方法
func (y *YourCommand) ExecuteWithArgs(args []string) error {
    // 处理参数逻辑
    return nil
}
```

### 3. 注册命令

在 `main.go` 中注册新命令：

```go
import (
    "community/xuan/commands/your-command"
    // ... 其他导入
)

func initCommands() {
    // ... 现有命令
    
    yourCmd := yourcommand.NewYourCommand()
    types.AddCommand(yourCmd)
}
```

## 命令接口说明

### Command 接口

```go
type Command interface {
    Exec() error             // 执行命令
    Father() Command         // 获取父命令
    Son() map[string]Command // 获取子命令列表
    Help()                   // 显示帮助信息
    IsEnd() bool             // 是否为最终命令
    Level() int64            // 命令级别
    String() string          // 命令名称
}
```

### BaseCommand 基础实现

`BaseCommand` 提供了 `Command` 接口的基础实现，包含：

- `name`: 命令名称
- `description`: 命令描述
- `level`: 命令级别
- `father`: 父命令
- `sons`: 子命令映射
- `isEnd`: 是否为最终命令

### 关键方法

- `NewBaseCommand(name, description, level, father)`: 创建基础命令
- `AddSon(cmd)`: 添加子命令
- `SetEnd(isEnd)`: 设置是否为最终命令

## 命令执行流程

1. 解析命令行参数
2. 从根命令开始查找
3. 递归查找子命令
4. 执行最终命令或显示帮助

## 示例：test 命令

`test` 命令展示了如何添加新命令：

```bash
# 查看 test 命令帮助
xuan test

# 运行单元测试
xuan test unit

# 运行单元测试（带参数）
xuan test unit ./service/user

# 运行集成测试
xuan test integration
```

## 最佳实践

1. **目录结构清晰**: 每个命令对应一个目录
2. **命令描述完整**: 提供清晰的命令描述和帮助信息
3. **参数处理**: 对于需要参数的最终命令，实现 `ExecuteWithArgs` 方法
4. **错误处理**: 适当的错误处理和用户友好的错误信息
5. **扩展性**: 设计时考虑未来扩展需求

## 现有命令说明

### gen 命令
- `xuan gen api <目录>`: 生成API代码
- `xuan gen rpc <目录>`: 生成RPC代码

### run 命令
- `xuan run all`: 启动所有服务
- `xuan run docker`: 重启Docker服务
- `xuan run api <目录>`: 重启API服务
- `xuan run rpc <目录>`: 重启RPC服务

### end 命令
- `xuan end all`: 停止所有服务
- `xuan end all --clean-data`: 停止所有服务并清理数据
- `xuan end docker`: 停止Docker服务
- `xuan end api <目录>`: 停止API服务
- `xuan end rpc <目录>`: 停止RPC服务
- `xuan end etcd`: 停止etcd服务
- `xuan end clean logs`: 清理日志文件
- `xuan end clean data`: 清理数据目录
- `xuan end clean binaries`: 清理二进制文件

### test 命令（示例）
- `xuan test unit`: 运行单元测试
- `xuan test integration`: 运行集成测试

