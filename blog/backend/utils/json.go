package utils

import (
	"encoding/json"
)

// ToJSON 将对象转换为JSON字符串
func ToJSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ParseJSON 解析JSON字符串到对象
func ParseJSON(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}
