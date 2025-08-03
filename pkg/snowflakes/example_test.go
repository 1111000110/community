package snowflakes

import (
	"fmt"
	"log"
	"time"
)

// ExampleSnowflake_NextID 演示基本的ID生成
func ExampleSnowflake_NextID() {
	// 创建雪花算法实例 (数据中心ID: 1, 机器ID: 1)
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		log.Fatal(err)
	}

	// 生成唯一ID
	id, err := sf.NextID()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated ID: %d\n", id)
	fmt.Printf("ID length: %d bits\n", 64)
}

// ExampleSnowflake_ParseID 演示ID解析
func ExampleSnowflake_ParseID() {
	sf, err := NewSnowflake(5, 10)
	if err != nil {
		log.Fatal(err)
	}

	id, err := sf.NextID()
	if err != nil {
		log.Fatal(err)
	}

	// 解析ID的各个组成部分
	timestamp, datacenterID, machineID, sequence := sf.ParseID(id)

	fmt.Printf("ID: %d\n", id)
	fmt.Printf("Timestamp: %d\n", timestamp)
	fmt.Printf("Datacenter ID: %d\n", datacenterID)
	fmt.Printf("Machine ID: %d\n", machineID)
	fmt.Printf("Sequence: %d\n", sequence)

	// 转换时间戳为可读时间
	readableTime := sf.GetTimestamp(id)
	fmt.Printf("Readable Time: %s\n", readableTime.Format("2006-01-02 15:04:05.000"))
}

// ExampleInitDefault 演示使用默认实例
func ExampleInitDefault() {
	// 初始化默认实例
	err := InitDefault(1, 1)
	if err != nil {
		log.Fatal(err)
	}

	// 使用默认实例生成ID
	id1, err := NextID()
	if err != nil {
		log.Fatal(err)
	}

	id2, err := NextID()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID 1: %d\n", id1)
	fmt.Printf("ID 2: %d\n", id2)
	fmt.Printf("IDs are unique: %t\n", id1 != id2)
}

// ExampleNewSnowflakeWithEpoch 演示自定义起始时间
func ExampleNewSnowflakeWithEpoch() {
	// 使用自定义起始时间 (2020-01-01)
	customEpoch := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	sf, err := NewSnowflakeWithEpoch(1, 1, customEpoch)
	if err != nil {
		log.Fatal(err)
	}

	id, err := sf.NextID()
	if err != nil {
		log.Fatal(err)
	}

	timestamp := sf.GetTimestamp(id)
	fmt.Printf("Generated ID: %d\n", id)
	fmt.Printf("Timestamp: %s\n", timestamp.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("Custom epoch: %s\n", customEpoch.Format("2006-01-02 15:04:05.000"))
}

// ExampleSnowflake_performance 演示性能特性
func ExampleSnowflake_performance() {
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		log.Fatal(err)
	}

	// 快速生成多个ID
	start := time.Now()
	count := 10000
	ids := make([]int64, count)

	for i := 0; i < count; i++ {
		ids[i], err = sf.NextID()
		if err != nil {
			log.Fatal(err)
		}
	}

	duration := time.Since(start)

	fmt.Printf("Generated %d IDs in %v\n", count, duration)
	fmt.Printf("Average: %.2f ns/ID\n", float64(duration.Nanoseconds())/float64(count))
	fmt.Printf("Rate: %.0f IDs/second\n", float64(count)/duration.Seconds())

	// 验证唯一性
	idSet := make(map[int64]bool)
	for _, id := range ids {
		if idSet[id] {
			fmt.Printf("Duplicate ID found: %d\n", id)
			return
		}
		idSet[id] = true
	}
	fmt.Printf("All %d IDs are unique\n", count)
}

// ExampleSnowflake_distributed 演示分布式场景
func ExampleSnowflake_distributed() {
	// 模拟不同数据中心和机器的ID生成器
	generators := []*Snowflake{}

	// 创建多个生成器 (不同的数据中心和机器ID)
	for datacenter := 0; datacenter < 3; datacenter++ {
		for machine := 0; machine < 3; machine++ {
			sf, err := NewSnowflake(int64(datacenter), int64(machine))
			if err != nil {
				log.Fatal(err)
			}
			generators = append(generators, sf)
		}
	}

	// 每个生成器生成ID
	allIDs := make(map[int64]bool)
	for i, sf := range generators {
		id, err := sf.NextID()
		if err != nil {
			log.Fatal(err)
		}

		datacenterID := sf.GetDatacenterID(id)
		machineID := sf.GetMachineID(id)

		fmt.Printf("Generator %d: ID=%d, DC=%d, Machine=%d\n", i, id, datacenterID, machineID)

		if allIDs[id] {
			fmt.Printf("Duplicate ID detected: %d\n", id)
		} else {
			allIDs[id] = true
		}
	}

	fmt.Printf("Generated %d unique IDs from %d generators\n", len(allIDs), len(generators))
}