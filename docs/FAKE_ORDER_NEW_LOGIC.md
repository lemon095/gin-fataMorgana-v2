# 水单生成逻辑更新

## 更新概述

水单生成逻辑和用户拼单逻辑已更新，主要变化如下：

### 1. 任务类型选择逻辑
- **旧逻辑**: 每个任务类型随机生成1-10个数量（水单）/ 1-8000个数量（用户拼单）
- **新逻辑**: 从四个类型（点赞、分享、关注、收藏）中随机选择1-4个，每个选中的类型数量固定为1

### 2. 金额计算逻辑
- **旧逻辑**: 基于缓存的价格配置计算总金额（单价×数量）
- **新逻辑**: 直接生成10万到1000万之间的随机总金额，不再依赖缓存的价格配置

### 3. 拼单金额逻辑
- **旧逻辑**: 基于价格配置计算人均金额
- **新逻辑**: 随机选择1-4个类型（每个类型数量为1），随机生成单价，计算总金额 = 单价 × 总任务数量，人均金额 = 总金额 ÷ 目标人数

### 4. 用户拼单逻辑
- **旧逻辑**: 用户参与拼单时，每个任务类型随机生成1-8000个数量
- **新逻辑**: 用户参与拼单时，随机选择1-4个类型，每个类型数量为1

## 技术实现

### 购买单生成逻辑
```go
// 随机选择1-4个类型，每个类型数量为1
likeCount := 0
shareCount := 0
followCount := 0
favoriteCount := 0

// 随机选择类型数量（1-4个）
typeCount := rand.Intn(4) + 1

// 创建类型数组并随机打乱
types := []string{"like", "share", "follow", "favorite"}
rand.Shuffle(len(types), func(i, j int) {
    types[i], types[j] = types[j], types[i]
})

// 选择前typeCount个类型，数量设为1
for i := 0; i < typeCount; i++ {
    switch types[i] {
    case "like":
        likeCount = 1
    case "share":
        shareCount = 1
    case "follow":
        followCount = 1
    case "favorite":
        favoriteCount = 1
    }
}

// 生成总金额（10万到1000万之间）
totalAmount := float64(rand.Intn(9900000)+100000) // 100000-10000000
```

### 拼单生成逻辑
```go
// 随机选择1-4个类型，每个类型数量为1
likeCount := 0
shareCount := 0
followCount := 0
favoriteCount := 0

// 随机选择类型数量（1-4个）
typeCount := rand.Intn(4) + 1

// 创建类型数组并随机打乱
types := []string{"like", "share", "follow", "favorite"}
rand.Shuffle(len(types), func(i, j int) {
    types[i], types[j] = types[j], types[i]
})

// 选择前typeCount个类型，数量设为1
for i := 0; i < typeCount; i++ {
    switch types[i] {
    case "like":
        likeCount = 1
    case "share":
        shareCount = 1
    case "follow":
        followCount = 1
    case "favorite":
        favoriteCount = 1
    }
}

// 随机生成单价（1万到10万之间）
unitPrice := float64(rand.Intn(90000)+10000) // 10000-100000

// 计算总任务数量
totalTaskCount := likeCount + shareCount + followCount + favoriteCount

// 计算总金额：单价 × 总任务数量
totalAmount := unitPrice * float64(totalTaskCount)

// 计算人均金额：总金额 ÷ 目标人数
perPersonAmount := totalAmount / float64(targetParticipants)
```

### 用户拼单逻辑
```go
// 用户参与拼单时，随机选择1-4个类型，每个类型数量为1
likeCount := 0
shareCount := 0
followCount := 0
favoriteCount := 0

// 随机选择类型数量（1-4个）
typeCount := rand.Intn(4) + 1

// 创建类型数组并随机打乱
types := []string{"like", "share", "follow", "favorite"}
rand.Shuffle(len(types), func(i, j int) {
    types[i], types[j] = types[j], types[i]
})

// 选择前typeCount个类型，数量设为1
for i := 0; i < typeCount; i++ {
    switch types[i] {
    case "like":
        likeCount = 1
    case "share":
        shareCount = 1
    case "follow":
        followCount = 1
    case "favorite":
        favoriteCount = 1
    }
}

// 创建订单时使用新的任务数量
order := &models.Order{
    LikeCount:      likeCount,
    ShareCount:     shareCount,
    FollowCount:    followCount,
    FavoriteCount:  favoriteCount,
    // ... 其他字段
}

## 配置变更

### 废弃的方法
- `getPurchaseConfig()`: 不再从Redis获取价格配置
- `calculateGroupBuyAmount()`: 不再基于价格配置计算

### 保留的功能
- 时间生成逻辑（前后10分钟随机）
- 状态分布逻辑（60%进行中，30%已完成，10%已关闭）
- 利润计算逻辑（基于用户等级）
- 期号分配逻辑

## 测试验证

使用测试脚本验证新逻辑：
```bash
./test_scripts/test_fake_order_new_logic.sh
```

### 验证要点
1. 任务类型数量是否正确（每个类型最多1个）
2. 购买单总金额是否在10万-1000万范围内
3. 拼单总金额是否合理（单价×任务数量）
4. 拼单人均金额是否合理（总金额÷目标人数）
5. 用户拼单任务数量是否正确（每个类型最多1个）
6. 生成的水单和用户拼单是否能正常查询

## 影响分析

### 正面影响
- 简化了价格计算逻辑，不再依赖外部缓存
- 提高了生成效率，减少了Redis查询
- 金额范围更符合业务需求

### 注意事项
- 水单金额不再基于实际价格配置
- 需要确保前端显示逻辑能正确处理新的金额范围
- 排行榜统计逻辑不受影响，仍基于实际金额

## 部署说明

1. 更新代码后重启服务
2. 运行测试脚本验证功能
3. 监控水单生成日志确认新逻辑生效
4. 检查排行榜数据是否正常 