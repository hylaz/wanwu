package util

import (
	"math"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func UUID() string {
	return uuid.NewV4().String()
}

func Reverse[T any](lst []T) {
	// 反转
	for i, j := 0, len(lst)-1; i < j; i, j = i+1, j-1 {
		lst[i], lst[j] = lst[j], lst[i]
	}
}

func InLst(item string, lst []string) bool {
	item = strings.TrimSpace(item)
	for _, l := range lst {
		if item == strings.TrimSpace(l) {
			return true
		}
	}

	return false
}

// int型除法。第一返回值为商，第二为余
func DivideAndRemainder(dividend, divisor int) (quotient, remainder int) {
	quotient = int(math.Floor(float64(dividend) / float64(divisor)))
	remainder = dividend - quotient*divisor
	return quotient, remainder
}

// len  一共101组数据，每页8个，那么需要 12+1=13页
func GetPageTotal(len, pageSize int) int {
	q, r := DivideAndRemainder(len, pageSize)
	if r == 0 {
		return q
	}
	return q + 1
}

// Intersection 函数用于找出两个可比较类型切片的交集。
func Intersection[T comparable](slice1, slice2 []T) []T {
	m := make(map[T]struct{}) // 使用空结构体来节省空间
	var result []T

	// 将第一个切片的所有元素放入 map 中
	for _, item := range slice1 {
		m[item] = struct{}{}
	}

	// 检查第二个切片中的元素是否存在于 map 中
	for _, item := range slice2 {
		if _, ok := m[item]; ok {
			// 如果存在，并且结果集中还没有这个元素，则添加到结果集中
			if !contains(result, item) {
				result = append(result, item)
			}
		}
	}

	return result
}

// contains 函数检查一个切片中是否包含某个元素。
func contains[T comparable](slice []T, elem T) bool {
	for _, item := range slice {
		if item == elem {
			return true
		}
	}
	return false
}
