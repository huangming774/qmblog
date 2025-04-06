package utils

import (
	"strconv"
)

// StringToUint 将字符串转换为uint，如果转换失败则返回0
func StringToUint(str string) uint {
	if str == "" {
		return 0
	}

	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}

	return uint(val)
}
