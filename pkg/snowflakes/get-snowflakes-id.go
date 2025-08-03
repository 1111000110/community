package snowflakes

import (
	"errors"
	"sync"
	"time"
)

// Snowflake 雪花算法ID生成器
// 64位ID结构: 1位符号位(0) + 41位时间戳 + 5位数据中心ID + 5位机器ID + 12位序列号
type Snowflake struct {
	mutex         sync.Mutex // 互斥锁，保证线程安全
	timestamp     int64      // 上次生成ID的时间戳
	datacenterID  int64      // 数据中心ID (0-31)
	machineID     int64      // 机器ID (0-31)
	sequence      int64      // 序列号 (0-4095)
	epoch         int64      // 起始时间戳 (2023-01-01 00:00:00 UTC)
}

// 位数定义
const (
	// 时间戳位数 (41位，可用69年)
	timestampBits = 41
	// 数据中心ID位数 (5位，支持32个数据中心)
	datacenterIDBits = 5
	// 机器ID位数 (5位，支持32台机器)
	machineIDBits = 5
	// 序列号位数 (12位，每毫秒支持4096个ID)
	sequenceBits = 12

	// 最大值计算
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits) // 31
	maxMachineID    = -1 ^ (-1 << machineIDBits)    // 31
	maxSequence     = -1 ^ (-1 << sequenceBits)     // 4095

	// 位移量
	machineIDShift     = sequenceBits                            // 12
	datacenterIDShift  = sequenceBits + machineIDBits            // 17
	timestampLeftShift = sequenceBits + machineIDBits + datacenterIDBits // 22
)

// 默认起始时间: 2023-01-01 00:00:00 UTC
var defaultEpoch = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

// NewSnowflake 创建新的雪花算法ID生成器
// datacenterID: 数据中心ID (0-31)
// machineID: 机器ID (0-31)
func NewSnowflake(datacenterID, machineID int64) (*Snowflake, error) {
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errors.New("datacenterID must be between 0 and 31")
	}
	if machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("machineID must be between 0 and 31")
	}

	return &Snowflake{
		datacenterID: datacenterID,
		machineID:    machineID,
		sequence:     0,
		epoch:        defaultEpoch,
	}, nil
}

// NewSnowflakeWithEpoch 创建带自定义起始时间的雪花算法ID生成器
func NewSnowflakeWithEpoch(datacenterID, machineID int64, epoch time.Time) (*Snowflake, error) {
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errors.New("datacenterID must be between 0 and 31")
	}
	if machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("machineID must be between 0 and 31")
	}

	return &Snowflake{
		datacenterID: datacenterID,
		machineID:    machineID,
		sequence:     0,
		epoch:        epoch.UnixMilli(),
	}, nil
}

// NextID 生成下一个唯一ID
func (s *Snowflake) NextID() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().UnixMilli()

	// 时钟回拨检测
	if now < s.timestamp {
		return 0, errors.New("clock moved backwards, refusing to generate id")
	}

	// 同一毫秒内生成ID
	if now == s.timestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		// 序列号溢出，等待下一毫秒
		if s.sequence == 0 {
			now = s.waitNextMillis(s.timestamp)
		}
	} else {
		// 新的毫秒，序列号重置为0
		s.sequence = 0
	}

	s.timestamp = now

	// 组装64位ID
	id := ((now - s.epoch) << timestampLeftShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.machineID << machineIDShift) |
		s.sequence

	return id, nil
}

// waitNextMillis 等待下一毫秒
func (s *Snowflake) waitNextMillis(lastTimestamp int64) int64 {
	timestamp := time.Now().UnixMilli()
	for timestamp <= lastTimestamp {
		time.Sleep(time.Millisecond)
		timestamp = time.Now().UnixMilli()
	}
	return timestamp
}

// ParseID 解析雪花算法ID，返回时间戳、数据中心ID、机器ID、序列号
func (s *Snowflake) ParseID(id int64) (timestamp int64, datacenterID int64, machineID int64, sequence int64) {
	timestamp = (id >> timestampLeftShift) + s.epoch
	datacenterID = (id >> datacenterIDShift) & maxDatacenterID
	machineID = (id >> machineIDShift) & maxMachineID
	sequence = id & maxSequence
	return
}

// GetTimestamp 从ID中提取时间戳
func (s *Snowflake) GetTimestamp(id int64) time.Time {
	timestamp := (id >> timestampLeftShift) + s.epoch
	return time.UnixMilli(timestamp)
}

// GetDatacenterID 从ID中提取数据中心ID
func (s *Snowflake) GetDatacenterID(id int64) int64 {
	return (id >> datacenterIDShift) & maxDatacenterID
}

// GetMachineID 从ID中提取机器ID
func (s *Snowflake) GetMachineID(id int64) int64 {
	return (id >> machineIDShift) & maxMachineID
}

// GetSequence 从ID中提取序列号
func (s *Snowflake) GetSequence(id int64) int64 {
	return id & maxSequence
}

// 全局默认实例
var defaultSnowflake *Snowflake
var once sync.Once

// InitDefault 初始化默认雪花算法实例
func InitDefault(datacenterID, machineID int64) error {
	var err error
	once.Do(func() {
		defaultSnowflake, err = NewSnowflake(datacenterID, machineID)
	})
	return err
}

// NextID 使用默认实例生成ID
func NextID() (int64, error) {
	if defaultSnowflake == nil {
		return 0, errors.New("default snowflake not initialized, call InitDefault first")
	}
	return defaultSnowflake.NextID()
}

// ParseID 使用默认实例解析ID
func ParseID(id int64) (timestamp int64, datacenterID int64, machineID int64, sequence int64) {
	if defaultSnowflake == nil {
		return 0, 0, 0, 0
	}
	return defaultSnowflake.ParseID(id)
}

// GetTimestamp 使用默认实例从ID中提取时间戳
func GetTimestamp(id int64) time.Time {
	if defaultSnowflake == nil {
		return time.Time{}
	}
	return defaultSnowflake.GetTimestamp(id)
}