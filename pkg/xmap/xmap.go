package xmap

// GetKeys 从map中获取所有键，返回键的切片
// K：键类型（必须可比较，符合map键的要求）
// V：值类型（任意类型）
func GetKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GetVals 从map中获取所有值，返回值的切片
// K：键类型（必须可比较）
// V：值类型（任意类型）
func GetVals[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}
