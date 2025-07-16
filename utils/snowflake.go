package utils

import (
	"fmt"
	"sync"
	"time"
)

// SnowflakeUID 雪花算法简化版UID生成器
type SnowflakeUID struct {
	mutex     sync.Mutex
	lastTime  int64
	sequence  int64
	machineID int64
}

// NewSnowflakeUID 创建新的UID生成器
func NewSnowflakeUID(machineID int64) *SnowflakeUID {
	return &SnowflakeUID{
		lastTime:  0,
		sequence:  0,
		machineID: machineID % 100, // 确保机器ID在0-99范围内
	}
}

// GenerateUID 生成8位UID
func (s *SnowflakeUID) GenerateUID() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 获取当前时间戳（毫秒）
	currentTime := time.Now().UnixNano() / 1e6

	// 处理时钟回退
	if currentTime < s.lastTime {
		// 如果时钟回退超过1秒，记录警告
		if s.lastTime-currentTime > 1000 {

		}

		// 等待到下一个毫秒
		time.Sleep(time.Millisecond)
		currentTime = time.Now().UnixNano() / 1e6

		// 如果仍然回退，使用上次时间
		if currentTime < s.lastTime {
			currentTime = s.lastTime
		}
	}

	// 如果是同一毫秒内，序列号递增
	if currentTime == s.lastTime {
		s.sequence = (s.sequence + 1) % 100 // 序列号范围0-99
	} else {
		// 不同毫秒，序列号重置
		s.sequence = 0
	}

	s.lastTime = currentTime

	// 生成UID：时间戳(4位) + 机器ID(2位) + 序列号(2位)
	timestamp := currentTime % 10000 // 取时间戳后4位
	return fmt.Sprintf("%04d%02d%02d", timestamp, s.machineID, s.sequence)
}

// 全局UID生成器实例
var globalUIDGenerator *SnowflakeUID

// InitSnowflake 初始化雪花算法
func InitSnowflake(workerID int64) {
	// 使用配置中的worker_id作为机器ID
	machineID := workerID
	globalUIDGenerator = NewSnowflakeUID(machineID)

}

// GenerateUID 全局UID生成函数
func GenerateUID() string {
	if globalUIDGenerator == nil {

		// 备用方案：使用时间戳生成UID
		timestamp := time.Now().UnixNano() / 1e6
		return fmt.Sprintf("%08d", timestamp%100000000)
	}
	return globalUIDGenerator.GenerateUID()
}

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	return "ORD" + GenerateUID()
}
