# Snowflake ID Generator

一个高性能、线程安全的雪花算法ID生成器，用于生成全局唯一的64位整数ID。

## 特性

- **全局唯一**: 在分布式环境中保证ID的全局唯一性
- **高性能**: 单机每秒可生成数百万个ID
- **线程安全**: 使用互斥锁保证并发安全
- **时间有序**: 生成的ID按时间递增
- **可解析**: 可以从ID中提取时间戳、数据中心ID、机器ID等信息
- **时钟回拨检测**: 自动检测和处理系统时钟回拨问题

## ID结构

64位ID的组成结构：

```
+----------+----------+----------+----------+
|  1 bit   |  41 bit  |  5 bit   |  5 bit   |  12 bit  |
| 符号位(0) |  时间戳   | 数据中心ID | 机器ID   |  序列号   |
+----------+----------+----------+----------+
```

- **符号位**: 1位，固定为0，保证生成的ID为正数
- **时间戳**: 41位，精确到毫秒，可使用69年
- **数据中心ID**: 5位，支持32个数据中心 (0-31)
- **机器ID**: 5位，支持32台机器 (0-31)
- **序列号**: 12位，同一毫秒内支持4096个ID (0-4095)

## 安装

```bash
go get community.com/pkg/snowflakes
```

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "log"
    "community.com/pkg/snowflakes"
)

func main() {
    // 创建雪花算法实例 (数据中心ID: 1, 机器ID: 1)
    sf, err := snowflakes.NewSnowflake(1, 1)
    if err != nil {
        log.Fatal(err)
    }

    // 生成唯一ID
    id, err := sf.NextID()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Generated ID: %d\n", id)
}
```

### 使用默认实例

```go
package main

import (
    "fmt"
    "log"
    "community.com/pkg/snowflakes"
)

func main() {
    // 初始化默认实例
    err := snowflakes.InitDefault(1, 1)
    if err != nil {
        log.Fatal(err)
    }

    // 使用默认实例生成ID
    id, err := snowflakes.NextID()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Generated ID: %d\n", id)
}
```

### 自定义起始时间

```go
package main

import (
    "fmt"
    "log"
    "time"
    "community.com/pkg/snowflakes"
)

func main() {
    // 使用自定义起始时间
    customEpoch := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
    sf, err := snowflakes.NewSnowflakeWithEpoch(1, 1, customEpoch)
    if err != nil {
        log.Fatal(err)
    }

    id, err := sf.NextID()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Generated ID: %d\n", id)
}
```

## API 文档

### 类型定义

#### Snowflake

```go
type Snowflake struct {
    // 私有字段
}
```

雪花算法ID生成器的主要结构体。

### 构造函数

#### NewSnowflake

```go
func NewSnowflake(datacenterID, machineID int64) (*Snowflake, error)
```

创建新的雪花算法ID生成器。

**参数:**
- `datacenterID`: 数据中心ID (0-31)
- `machineID`: 机器ID (0-31)

**返回:**
- `*Snowflake`: 雪花算法实例
- `error`: 错误信息

#### NewSnowflakeWithEpoch

```go
func NewSnowflakeWithEpoch(datacenterID, machineID int64, epoch time.Time) (*Snowflake, error)
```

创建带自定义起始时间的雪花算法ID生成器。

### 实例方法

#### NextID

```go
func (s *Snowflake) NextID() (int64, error)
```

生成下一个唯一ID。

#### ParseID

```go
func (s *Snowflake) ParseID(id int64) (timestamp int64, datacenterID int64, machineID int64, sequence int64)
```

解析雪花算法ID，返回各个组成部分。

#### GetTimestamp

```go
func (s *Snowflake) GetTimestamp(id int64) time.Time
```

从ID中提取时间戳。

#### GetDatacenterID

```go
func (s *Snowflake) GetDatacenterID(id int64) int64
```

从ID中提取数据中心ID。

#### GetMachineID

```go
func (s *Snowflake) GetMachineID(id int64) int64
```

从ID中提取机器ID。

#### GetSequence

```go
func (s *Snowflake) GetSequence(id int64) int64
```

从ID中提取序列号。

### 全局函数

#### InitDefault

```go
func InitDefault(datacenterID, machineID int64) error
```

初始化默认雪花算法实例。

#### NextID

```go
func NextID() (int64, error)
```

使用默认实例生成ID。

#### ParseID

```go
func ParseID(id int64) (timestamp int64, datacenterID int64, machineID int64, sequence int64)
```

使用默认实例解析ID。

#### GetTimestamp

```go
func GetTimestamp(id int64) time.Time
```

使用默认实例从ID中提取时间戳。

## 性能

在 Apple M4 Pro 上的性能测试结果：

- **单线程**: ~314 ns/op (约320万ID/秒)
- **并发**: ~335 ns/op (约300万ID/秒)

## 分布式部署

在分布式环境中，需要为每个节点分配唯一的数据中心ID和机器ID组合：

```go
// 节点1: 数据中心0, 机器0
sf1, _ := snowflakes.NewSnowflake(0, 0)

// 节点2: 数据中心0, 机器1
sf2, _ := snowflakes.NewSnowflake(0, 1)

// 节点3: 数据中心1, 机器0
sf3, _ := snowflakes.NewSnowflake(1, 0)
```

## 注意事项

1. **时钟同步**: 确保所有节点的系统时钟同步，避免时钟回拨
2. **ID分配**: 在分布式环境中，确保每个节点的数据中心ID和机器ID组合唯一
3. **时钟回拨**: 当检测到时钟回拨时，会返回错误，需要等待时钟恢复正常
4. **序列号溢出**: 当同一毫秒内生成超过4096个ID时，会自动等待下一毫秒

## 错误处理

常见错误及处理方法：

- `datacenterID must be between 0 and 31`: 数据中心ID超出范围
- `machineID must be between 0 and 31`: 机器ID超出范围
- `clock moved backwards, refusing to generate id`: 检测到时钟回拨
- `default snowflake not initialized, call InitDefault first`: 默认实例未初始化

## 测试

运行测试：

```bash
go test ./pkg/snowflakes -v
```

运行性能测试：

```bash
go test ./pkg/snowflakes -bench=.
```

## 许可证

MIT License