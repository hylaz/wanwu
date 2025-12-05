package util

import (
	"strings"
)

func UrlNameFilter(old string) string {
	var str = strings.ReplaceAll(old, "\\r\\n", "")
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\r\n", "")
	return str
}
