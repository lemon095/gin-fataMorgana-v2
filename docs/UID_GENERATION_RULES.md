# UID生成规则文档

## 概述

本项目使用基于雪花算法的简化版UID生成器来生成8位用户唯一标识符（UID）。该算法确保在分布式环境下的唯一性和有序性。

## 算法实现

### 核心文件
- **实现文件**: `utils/snowflake.go`
- **配置**: `config.yaml` 中的 `snowflake.worker_id`
- **初始化**: `main.go` 中的 `utils.InitSnowflake()`

### 算法原理

基于雪花算法（Snowflake）的简化版本，生成8位数字UID：

```
UID = 时间戳(4位) + 机器ID(2位) + 序列号(2位)
```

#### 组成部分

1. **时间戳部分** (4位)
   - 来源：当前时间戳（毫秒）的后4位
   - 范围：0000-9999
   - 作用：确保时间顺序性

2. **机器ID部分** (2位)
   - 来源：配置文件中的 `worker_id`
   - 范围：00-99
   - 作用：区分不同服务器节点

3. **序列号部分** (2位)
   - 来源：同一毫秒内的递增序列
   - 范围：00-99
   - 作用：处理同一毫秒内的并发

## 配置说明

### 配置文件
```yaml
# config.yaml
snowflake:
  worker_id: 1        # 工作节点ID，范围0-99
  datacenter_id: 1    # 数据中心ID（当前未使用）
```

### 初始化代码
```go
// main.go
utils.InitSnowflake(config.GlobalConfig.Snowflake.WorkerID)
```

## 生成规则详解

### 1. 时间戳处理
```go
// 获取当前时间戳（毫秒）
currentTime := time.Now().UnixNano() / 1e6

// 取后4位作为时间戳部分
timestamp := currentTime % 10000
```

### 2. 机器ID处理
```go
// 确保机器ID在0-99范围内
machineID := workerID % 100
```

### 3. 序列号处理
```go
// 同一毫秒内序列号递增
if currentTime == s.lastTime {
    s.sequence = (s.sequence + 1) % 100 // 序列号范围0-99
} else {
    s.sequence = 0 // 不同毫秒重置序列号
}
```

### 4. 时钟回退处理
```go
// 处理时钟回退情况
if currentTime < s.lastTime {
    // 等待到下一个毫秒
    time.Sleep(time.Millisecond)
    currentTime = time.Now().UnixNano() / 1e6
    
    // 如果仍然回退，使用上次时间
    if currentTime < s.lastTime {
        currentTime = s.lastTime
    }
}
```

## 使用示例

### 生成UID
```go
// 生成8位UID
uid := utils.GenerateUID()
// 示例输出: "12345678"
```

### 生成订单号
```go
// 生成订单号（前缀 + UID）
orderNo := utils.GenerateOrderNo()
// 示例输出: "ORD12345678"
```

## 性能特点

### 优势
1. **唯一性**: 通过时间戳+机器ID+序列号确保唯一性
2. **有序性**: 时间戳确保时间顺序
3. **高性能**: 本地生成，无需数据库交互
4. **分布式友好**: 支持多服务器部署

### 限制
1. **时钟依赖**: 依赖系统时钟准确性
2. **容量限制**: 每毫秒最多生成100个UID
3. **时间精度**: 依赖毫秒级时间戳

## 容量计算

### 理论容量
- **时间戳**: 10000个不同值
- **机器ID**: 100个不同值
- **序列号**: 100个不同值
- **总容量**: 10000 × 100 × 100 = 100,000,000

### 实际容量
- **每毫秒**: 最多100个UID
- **每秒**: 最多100,000个UID
- **每天**: 最多8,640,000,000个UID

## 使用场景

### 1. 用户注册
```go
// services/user_service.go
userID := utils.GenerateUID()
user := &models.User{
    Uid: userID,
    // ... 其他字段
}
```

### 2. 订单创建
```go
// utils/snowflake.go
func GenerateOrderNo() string {
    return "ORD" + GenerateUID()
}
```

### 3. 其他业务场景
- 交易流水号
- 系统内部ID
- 临时标识符

## 故障处理

### 1. 时钟回退
- 自动检测时钟回退
- 等待到下一个毫秒
- 记录警告日志

### 2. 初始化失败
```go
// 备用方案：使用时间戳生成UID
if globalUIDGenerator == nil {
    timestamp := time.Now().UnixNano() / 1e6
    return fmt.Sprintf("%08d", timestamp%100000000)
}
```

### 3. 序列号溢出
- 序列号范围：0-99
- 溢出后自动重置为0
- 等待下一个毫秒

## 部署建议

### 1. 服务器配置
- 确保每台服务器的 `worker_id` 唯一
- 配置NTP服务保持时钟同步
- 监控系统时钟准确性

### 2. 高可用部署
```yaml
# 服务器1
snowflake:
  worker_id: 1

# 服务器2  
snowflake:
  worker_id: 2

# 服务器3
snowflake:
  worker_id: 3
```

### 3. 监控指标
- UID生成速率
- 时钟偏差
- 序列号使用情况
- 错误日志

## 测试验证

### 测试脚本
```bash
# 测试UID生成
go test -v ./utils -run TestGenerateUID

# 测试并发生成
go test -v ./utils -run TestConcurrentUIDGeneration
```

### 验证要点
1. 唯一性验证
2. 有序性验证
3. 并发安全性
4. 时钟回退处理
5. 性能测试

## 总结

本项目使用的UID生成算法具有以下特点：

1. **简单高效**: 8位数字，易于使用和记忆
2. **分布式友好**: 支持多服务器部署
3. **性能优秀**: 本地生成，无网络开销
4. **容错性强**: 具备时钟回退处理机制
5. **扩展性好**: 可根据需要调整位数和组成

该算法适合中小型项目的UID生成需求，在保证唯一性和性能的同时，提供了良好的可维护性。 