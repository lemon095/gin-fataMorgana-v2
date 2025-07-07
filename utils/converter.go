package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Converter 数据转换器
type Converter struct{}

// NewConverter 创建转换器
func NewConverter() *Converter {
	return &Converter{}
}

// ToInt64 转换为int64
func (c *Converter) ToInt64(value interface{}) (int64, error) {
	if value == nil {
		return 0, fmt.Errorf("value is nil")
	}

	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case []byte:
		return strconv.ParseInt(string(v), 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to int64", value)
	}
}

// ToFloat64 转换为float64
func (c *Converter) ToFloat64(value interface{}) (float64, error) {
	if value == nil {
		return 0, fmt.Errorf("value is nil")
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to float64", value)
	}
}

// ToString 转换为string
func (c *Converter) ToString(value interface{}) (string, error) {
	if value == nil {
		return "", fmt.Errorf("value is nil")
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case int:
		return strconv.Itoa(v), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case bool:
		return strconv.FormatBool(v), nil
	case time.Time:
		return v.Format(time.RFC3339), nil
	default:
		// 尝试JSON序列化
		if data, err := json.Marshal(v); err == nil {
			return string(data), nil
		}
		return fmt.Sprintf("%v", v), nil
	}
}

// ToBool 转换为bool
func (c *Converter) ToBool(value interface{}) (bool, error) {
	if value == nil {
		return false, fmt.Errorf("value is nil")
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case int64:
		return v != 0, nil
	case int:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int16:
		return v != 0, nil
	case int8:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case uint:
		return v != 0, nil
	case uint32:
		return v != 0, nil
	case uint16:
		return v != 0, nil
	case uint8:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case float32:
		return v != 0, nil
	case string:
		return strconv.ParseBool(v)
	case []byte:
		return strconv.ParseBool(string(v))
	default:
		return false, fmt.Errorf("cannot convert %v to bool", value)
	}
}

// ToTime 转换为time.Time
func (c *Converter) ToTime(value interface{}) (time.Time, error) {
	if value == nil {
		return time.Time{}, fmt.Errorf("value is nil")
	}

	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		// 尝试多种时间格式
		formats := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000Z",
			"2006-01-02",
		}
		
		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("cannot parse time string: %s", v)
	case []byte:
		return c.ToTime(string(v))
	case int64:
		// 假设是Unix时间戳
		return time.Unix(v, 0), nil
	case int:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("cannot convert %v to time.Time", value)
	}
}

// ToStruct 转换为结构体
func (c *Converter) ToStruct(value interface{}, dest interface{}) error {
	if value == nil {
		return fmt.Errorf("value is nil")
	}

	if dest == nil {
		return fmt.Errorf("destination is nil")
	}

	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, dest)
	case string:
		return json.Unmarshal([]byte(v), dest)
	case map[string]interface{}:
		// 将map转换为JSON，再解析为结构体
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, dest)
	default:
		// 尝试JSON序列化再反序列化
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("cannot marshal value: %v", err)
		}
		return json.Unmarshal(data, dest)
	}
}

// ToMap 转换为map[string]interface{}
func (c *Converter) ToMap(value interface{}) (map[string]interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("value is nil")
	}

	switch v := value.(type) {
	case map[string]interface{}:
		return v, nil
	case []byte:
		var result map[string]interface{}
		if err := json.Unmarshal(v, &result); err != nil {
			return nil, err
		}
		return result, nil
	case string:
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(v), &result); err != nil {
			return nil, err
		}
		return result, nil
	default:
		// 尝试JSON序列化再反序列化
		data, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal value: %v", err)
		}
		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		return result, nil
	}
}

// ToSlice 转换为[]interface{}
func (c *Converter) ToSlice(value interface{}) ([]interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("value is nil")
	}

	switch v := value.(type) {
	case []interface{}:
		return v, nil
	case []byte:
		var result []interface{}
		if err := json.Unmarshal(v, &result); err != nil {
			return nil, err
		}
		return result, nil
	case string:
		var result []interface{}
		if err := json.Unmarshal([]byte(v), &result); err != nil {
			return nil, err
		}
		return result, nil
	default:
		// 尝试JSON序列化再反序列化
		data, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal value: %v", err)
		}
		var result []interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		return result, nil
	}
}

// SafeToInt64 安全转换为int64，失败时返回默认值
func (c *Converter) SafeToInt64(value interface{}, defaultValue int64) int64 {
	if result, err := c.ToInt64(value); err == nil {
		return result
	}
	return defaultValue
}

// SafeToFloat64 安全转换为float64，失败时返回默认值
func (c *Converter) SafeToFloat64(value interface{}, defaultValue float64) float64 {
	if result, err := c.ToFloat64(value); err == nil {
		return result
	}
	return defaultValue
}

// SafeToString 安全转换为string，失败时返回默认值
func (c *Converter) SafeToString(value interface{}, defaultValue string) string {
	if result, err := c.ToString(value); err == nil {
		return result
	}
	return defaultValue
}

// SafeToBool 安全转换为bool，失败时返回默认值
func (c *Converter) SafeToBool(value interface{}, defaultValue bool) bool {
	if result, err := c.ToBool(value); err == nil {
		return result
	}
	return defaultValue
}

// 全局转换器实例
var GlobalConverter = NewConverter() 