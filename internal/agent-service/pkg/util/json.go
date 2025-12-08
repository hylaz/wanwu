package util

import (
	"encoding/json"
	"fmt"
)

/*func json2Str(i interface{}) (string, error) {
	// 序列化为JSON
	jsonData, err := json.Marshal(i)
	if err != nil {
		return "", fmt.Errorf("JSON Marshaling failed: %w", err)
	}
	// 将JSON字节切片转换为字符串并打印
	jsonString := string(jsonData)
	return jsonString, nil
}*/

func JSONParse[T any](jsonStr string, target *T) error {
	// 解析JSON到目标类型
	if err := json.Unmarshal([]byte(jsonStr), target); err != nil {
		return fmt.Errorf("JSON unmarshaling failed: %w", err)
	}
	return nil
}
