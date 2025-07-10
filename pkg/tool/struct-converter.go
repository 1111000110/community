package tool

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// StructConverter 通用结构体转换器
func StructConverter(src any, dest any) error {
	// 1. 检查源和目标是否为指针类型
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)
	if srcValue.Kind() != reflect.Ptr || destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("参数必须为指针类型")
	}

	// 2. 将源结构体序列化为 JSON 字节流
	srcJSON, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("序列化源结构体失败: %v", err)
	}

	// 3. 将 JSON 反序列化为目标结构体（自动处理字段匹配和默认值）
	if err := json.Unmarshal(srcJSON, dest); err != nil {
		return fmt.Errorf("反序列化到目标结构体失败: %v", err)
	}

	return nil
}
