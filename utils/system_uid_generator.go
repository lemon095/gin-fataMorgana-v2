package utils

import (
	"fmt"
	"sync"
	"time"
)

// SystemUIDGenerator 系统UID生成器
type SystemUIDGenerator struct {
	mutex           sync.Mutex
	sequence        int64
	lastTime        int64
	machineID       int64
	groupBuyCounter int64 // 拼单计数器，确保拼单号唯一性
	orderCounter    int64 // 订单计数器，确保订单号唯一性
}

// NewSystemUIDGenerator 创建新的系统UID生成器
func NewSystemUIDGenerator(machineID int64) *SystemUIDGenerator {
	return &SystemUIDGenerator{
		lastTime:     0,
		sequence:     0,
		machineID:    machineID % 100, // 确保机器ID在0-99范围内
		orderCounter: 0,
	}
}

// GenerateSystemUID 生成7位系统UID
func (s *SystemUIDGenerator) GenerateSystemUID() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 获取当前时间戳（纳秒级精度）
	currentTime := time.Now().UnixNano()

	// 处理时钟回退
	if currentTime < s.lastTime {
		// 等待到下一个微秒
		time.Sleep(time.Microsecond)
		currentTime = time.Now().UnixNano()

		// 如果仍然回退，使用上次时间
		if currentTime < s.lastTime {
			currentTime = s.lastTime
		}
	}

	// 如果是同一微秒内，序列号递增
	if currentTime == s.lastTime {
		s.sequence = (s.sequence + 1) % 1000 // 序列号范围0-999
	} else {
		// 不同微秒，序列号重置
		s.sequence = 0
	}

	s.lastTime = currentTime

	// 生成7位UID：时间戳(4位) + 机器ID(2位) + 序列号(1位)
	timestamp := (currentTime / 1000000) % 10000 // 取纳秒时间戳后4位
	return fmt.Sprintf("%04d%02d%01d", timestamp, s.machineID, s.sequence%10)
}

// GenerateSystemOrderNo 生成系统订单号
func (s *SystemUIDGenerator) GenerateSystemOrderNo() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 订单计数器递增
	s.orderCounter = (s.orderCounter + 1) % 10000 // 计数器范围0-9999

	// 获取当前时间戳（纳秒级精度）
	currentTime := time.Now().UnixNano()
	
	// 处理时钟回退
	if currentTime < s.lastTime {
		// 等待到下一个微秒
		time.Sleep(time.Microsecond)
		currentTime = time.Now().UnixNano()

		// 如果仍然回退，使用上次时间
		if currentTime < s.lastTime {
			currentTime = s.lastTime
		}
	}

	// 如果是同一微秒内，序列号递增
	if currentTime == s.lastTime {
		s.sequence = (s.sequence + 1) % 1000 // 序列号范围0-999
	} else {
		// 不同微秒，序列号重置
		s.sequence = 0
	}

	s.lastTime = currentTime

	// 生成订单号：ORD + 时间戳后4位 + 机器ID2位 + 计数器4位
	// 使用纳秒时间戳的后4位，提高唯一性
	timestamp := (currentTime / 1000000) % 10000 // 取纳秒时间戳后4位
	return fmt.Sprintf("ORD%04d%02d%04d", timestamp, s.machineID, s.orderCounter)
}

// GenerateSystemGroupBuyNo 生成系统拼单号
func (s *SystemUIDGenerator) GenerateSystemGroupBuyNo() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 拼单计数器递增
	s.groupBuyCounter = (s.groupBuyCounter + 1) % 10000 // 计数器范围0-9999

	// 获取当前时间戳（纳秒级精度）
	currentTime := time.Now().UnixNano()
	
	// 处理时钟回退
	if currentTime < s.lastTime {
		// 等待到下一个微秒
		time.Sleep(time.Microsecond)
		currentTime = time.Now().UnixNano()

		// 如果仍然回退，使用上次时间
		if currentTime < s.lastTime {
			currentTime = s.lastTime
		}
	}

	// 如果是同一微秒内，序列号递增
	if currentTime == s.lastTime {
		s.sequence = (s.sequence + 1) % 1000 // 序列号范围0-999
	} else {
		// 不同微秒，序列号重置
		s.sequence = 0
	}

	s.lastTime = currentTime

	// 生成拼单号：GB + 时间戳后4位 + 机器ID2位 + 计数器4位
	// 使用纳秒时间戳的后4位，提高唯一性
	timestamp := (currentTime / 1000000) % 10000 // 取纳秒时间戳后4位
	return fmt.Sprintf("GB%04d%02d%04d", timestamp, s.machineID, s.groupBuyCounter)
}

// 全局系统UID生成器实例
var globalSystemUIDGenerator *SystemUIDGenerator

// InitSystemUIDGenerator 初始化系统UID生成器
func InitSystemUIDGenerator(workerID int64) {
	globalSystemUIDGenerator = NewSystemUIDGenerator(workerID)
}

// GenerateSystemUID 全局系统UID生成函数
func GenerateSystemUID() string {
	if globalSystemUIDGenerator == nil {
		// 备用方案：使用时间戳生成7位UID
		timestamp := time.Now().UnixNano() / 1e6
		return fmt.Sprintf("%07d", timestamp%10000000)
	}
	return globalSystemUIDGenerator.GenerateSystemUID()
}

// GenerateSystemOrderNo 全局系统订单号生成函数
func GenerateSystemOrderNo() string {
	if globalSystemUIDGenerator == nil {
		return "ORD" + GenerateSystemUID()
	}
	return globalSystemUIDGenerator.GenerateSystemOrderNo()
}

// GenerateSystemGroupBuyNo 全局系统拼单号生成函数
func GenerateSystemGroupBuyNo() string {
	if globalSystemUIDGenerator == nil {
		return "GB" + GenerateSystemUID()
	}
	return globalSystemUIDGenerator.GenerateSystemGroupBuyNo()
} 