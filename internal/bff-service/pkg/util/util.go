// @Author wangxm 8/21/星期四 16:01:00
package util

// 拼接 UniqueId
func ConcatAssistantToolUniqueId(typeStr, IdStr string) string {
	if typeStr == "" || IdStr == "" {
		return ""
	}
	return typeStr + "_" + IdStr
}
