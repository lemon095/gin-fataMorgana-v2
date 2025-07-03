package utils

import (
	"encoding/json"
)

// StructToJSON 将结构体转换为JSON字符串
func StructToJSON(v interface{}) (string, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// JSONToStruct 将JSON字符串转换为结构体
func JSONToStruct(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}
