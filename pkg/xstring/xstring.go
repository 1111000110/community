package xstring

import "strconv"

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// IntToString 泛型函数：将任意整数类型转换为 string
// 支持 int/int8/int16/int32/int64/uint 等所有整数类型
func IntToString[T Integer](info T) string {
	// 根据类型分别处理，避免大整数溢出
	switch any(info).(type) {
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(int64(info), 10)
	case uint, uint8, uint16, uint32, uint64, uintptr:
		return strconv.FormatUint(uint64(info), 10)
	default:
		// 理论上不会走到这里（受 Integer 约束限制）
		return ""
	}
}

// StringToInt 将字符串转换为指定的整数类型 T
// 返回转换后的值和可能的错误（格式错误、溢出等）
func StringToInt[T Integer](s string) (T, error) {
	var zero T // 用于返回零值（类型兼容）

	// 根据目标类型的底层类型选择不同的解析方法
	switch any(zero).(type) {
	// 有符号整数类型：int、int8、int16、int32、int64
	case int, int8, int16, int32, int64:
		// 解析为int64（最大有符号整数类型，避免中间溢出）
		// 第三个参数0表示自动推断基数（0-9为10进制，0x开头为16进制等）
		// 第四个参数指定目标类型的位大小（如int8为8，int为系统位数）
		val, err := strconv.ParseInt(s, 0, strconv.IntSize)
		if err != nil {
			return zero, err
		}
		// 转换为目标类型（因约束限制，此处不会溢出）
		return T(val), nil

	// 无符号整数类型：uint、uint8、uint16、uint32、uint64、uintptr
	case uint, uint8, uint16, uint32, uint64, uintptr:
		// 解析为uint64（最大无符号整数类型，避免中间溢出）
		val, err := strconv.ParseUint(s, 0, strconv.IntSize)
		if err != nil {
			return zero, err
		}
		// 转换为目标类型（因约束限制，此处不会溢出）
		return T(val), nil

	default:
		// 理论上不会走到这里（受Integer约束限制）
		return zero, strconv.ErrSyntax
	}
}

// StringToIntOrZero 将字符串转换为指定整数类型 T
// 转换成功返回对应值，失败（格式错误/溢出等）返回 T 的默认值（零值）
func StringToIntOrZero[T Integer](s string) T {
	// 调用带错误处理的版本
	val, err := StringToInt[T](s)
	if err != nil {
		var zero T // 声明 T 类型的零值（如 int 为 0，uint8 为 0 等）
		return zero
	}
	return val
}
