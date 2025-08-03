package snowflakes

import (
	"sync"
	"testing"
	"time"
)

// TestNewSnowflake 测试创建雪花算法实例
func TestNewSnowflake(t *testing.T) {
	// 正常情况
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if sf == nil {
		t.Error("Expected snowflake instance, got nil")
	}

	// 数据中心ID超出范围
	_, err = NewSnowflake(32, 1)
	if err == nil {
		t.Error("Expected error for datacenterID > 31")
	}

	// 机器ID超出范围
	_, err = NewSnowflake(1, 32)
	if err == nil {
		t.Error("Expected error for machineID > 31")
	}

	// 负数ID
	_, err = NewSnowflake(-1, 1)
	if err == nil {
		t.Error("Expected error for negative datacenterID")
	}
}

// TestNextID 测试ID生成
func TestNextID(t *testing.T) {
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		t.Fatalf("Failed to create snowflake: %v", err)
	}

	// 生成多个ID，确保唯一性
	ids := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id, err := sf.NextID()
		if err != nil {
			t.Errorf("Failed to generate ID: %v", err)
		}
		if ids[id] {
			t.Errorf("Duplicate ID generated: %d", id)
		}
		ids[id] = true
		if id <= 0 {
			t.Errorf("Generated negative or zero ID: %d", id)
		}
	}
}

// TestConcurrentGeneration 测试并发生成ID
func TestConcurrentGeneration(t *testing.T) {
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		t.Fatalf("Failed to create snowflake: %v", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	ids := make(map[int64]bool)
	goroutines := 10
	idsPerGoroutine := 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := sf.NextID()
				if err != nil {
					t.Errorf("Failed to generate ID: %v", err)
					return
				}
				mu.Lock()
				if ids[id] {
					t.Errorf("Duplicate ID generated: %d", id)
				}
				ids[id] = true
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	expectedCount := goroutines * idsPerGoroutine
	if len(ids) != expectedCount {
		t.Errorf("Expected %d unique IDs, got %d", expectedCount, len(ids))
	}
}

// TestParseID 测试ID解析
func TestParseID(t *testing.T) {
	sf, err := NewSnowflake(5, 10)
	if err != nil {
		t.Fatalf("Failed to create snowflake: %v", err)
	}

	id, err := sf.NextID()
	if err != nil {
		t.Fatalf("Failed to generate ID: %v", err)
	}

	timestamp, datacenterID, machineID, sequence := sf.ParseID(id)

	if datacenterID != 5 {
		t.Errorf("Expected datacenterID 5, got %d", datacenterID)
	}
	if machineID != 10 {
		t.Errorf("Expected machineID 10, got %d", machineID)
	}
	if timestamp <= sf.epoch {
		t.Errorf("Expected timestamp > epoch, got %d", timestamp)
	}
	if sequence < 0 || sequence > maxSequence {
		t.Errorf("Expected sequence in range [0, %d], got %d", maxSequence, sequence)
	}
}

// TestGetMethods 测试各种Get方法
func TestGetMethods(t *testing.T) {
	sf, err := NewSnowflake(3, 7)
	if err != nil {
		t.Fatalf("Failed to create snowflake: %v", err)
	}

	id, err := sf.NextID()
	if err != nil {
		t.Fatalf("Failed to generate ID: %v", err)
	}

	if sf.GetDatacenterID(id) != 3 {
		t.Errorf("Expected datacenterID 3, got %d", sf.GetDatacenterID(id))
	}
	if sf.GetMachineID(id) != 7 {
		t.Errorf("Expected machineID 7, got %d", sf.GetMachineID(id))
	}

	timestamp := sf.GetTimestamp(id)
	if timestamp.Before(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Error("Timestamp should be after epoch")
	}

	sequence := sf.GetSequence(id)
	if sequence < 0 || sequence > maxSequence {
		t.Errorf("Expected sequence in range [0, %d], got %d", maxSequence, sequence)
	}
}

// TestDefaultInstance 测试默认实例
func TestDefaultInstance(t *testing.T) {
	// 未初始化时应该返回错误
	_, err := NextID()
	if err == nil {
		t.Error("Expected error when default instance not initialized")
	}

	// 初始化默认实例
	err = InitDefault(1, 1)
	if err != nil {
		t.Fatalf("Failed to initialize default instance: %v", err)
	}

	// 生成ID
	id, err := NextID()
	if err != nil {
		t.Errorf("Failed to generate ID with default instance: %v", err)
	}
	if id <= 0 {
		t.Errorf("Generated invalid ID: %d", id)
	}

	// 解析ID
	timestamp, datacenterID, machineID, sequence := ParseID(id)
	if datacenterID != 1 {
		t.Errorf("Expected datacenterID 1, got %d", datacenterID)
	}
	if machineID != 1 {
		t.Errorf("Expected machineID 1, got %d", machineID)
	}
	if timestamp <= defaultEpoch {
		t.Errorf("Expected timestamp > epoch, got %d", timestamp)
	}
	if sequence < 0 {
		t.Errorf("Expected non-negative sequence, got %d", sequence)
	}

	// 获取时间戳
	ts := GetTimestamp(id)
	if ts.IsZero() {
		t.Error("Expected valid timestamp")
	}
}

// TestCustomEpoch 测试自定义起始时间
func TestCustomEpoch(t *testing.T) {
	customEpoch := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	sf, err := NewSnowflakeWithEpoch(1, 1, customEpoch)
	if err != nil {
		t.Fatalf("Failed to create snowflake with custom epoch: %v", err)
	}

	id, err := sf.NextID()
	if err != nil {
		t.Fatalf("Failed to generate ID: %v", err)
	}

	timestamp := sf.GetTimestamp(id)
	if timestamp.Before(customEpoch) {
		t.Error("Timestamp should be after custom epoch")
	}
}

// BenchmarkNextID 性能测试
func BenchmarkNextID(b *testing.B) {
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		b.Fatalf("Failed to create snowflake: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sf.NextID()
		if err != nil {
			b.Fatalf("Failed to generate ID: %v", err)
		}
	}
}

// BenchmarkConcurrentNextID 并发性能测试
func BenchmarkConcurrentNextID(b *testing.B) {
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		b.Fatalf("Failed to create snowflake: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := sf.NextID()
			if err != nil {
				b.Fatalf("Failed to generate ID: %v", err)
			}
		}
	})
}