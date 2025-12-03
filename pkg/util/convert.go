package util

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func I64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func MustI64(s string) int64 {
	i64, _ := I64(s)
	return i64
}

func I32(s string) (int32, error) {
	i64, err := I64(s)
	if err != nil {
		return 0, err
	}
	return int32(i64), nil
}

func MustI32(s string) int32 {
	i64, _ := I64(s)
	return int32(i64)
}

func U32(s string) (uint32, error) {
	i64, err := I64(s)
	if err != nil {
		return 0, err
	}
	return uint32(i64), nil
}

func MustU32(s string) uint32 {
	i64, _ := I64(s)
	return uint32(i64)
}

func Int2Str[T ~int | ~int32 | ~uint32 | ~int64](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

// 将 map[string]interface{} 转换为 map[string]string
func ConvertMapToString(req map[string]interface{}) map[string]string {
	formData := make(map[string]string)
	for key, value := range req {
		switch v := value.(type) {
		case string:
			formData[key] = v
		case int, int32, int64, float32, float64, bool:
			formData[key] = fmt.Sprintf("%v", v)
		default:
			// 对于其他类型，可以序列化为JSON字符串或根据需求处理
			if jsonStr, err := json.Marshal(v); err == nil {
				formData[key] = string(jsonStr)
			}
		}
	}
	return formData
}
