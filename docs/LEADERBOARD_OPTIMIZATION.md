# 热榜功能优化文档

## 概述

本文档描述了热榜功能的优化方案，主要目标是移除窗口函数的使用，采用更简单直接的SQL查询方式，提高代码的可读性和维护性。

## 优化前的问题

### 1. 窗口函数的使用
- 使用了 `ROW_NUMBER() OVER` 窗口函数
- 查询逻辑复杂，难以理解和维护
- 在某些数据库版本中可能存在兼容性问题

### 2. 排名计算复杂
- 使用了大量的子查询来计算用户排名
- 查询性能可能受到影响
- 代码逻辑复杂，容易出错

## 优化后的方案

### 1. 获取前10名用户数据

**优化前（使用窗口函数）：**
```sql
WITH user_stats AS (
    SELECT 
        o.uid,
        u.username,
        MAX(o.updated_at) as completed_at,
        COUNT(*) as order_count,
        SUM(o.amount) as total_amount,
        SUM(o.profit_amount) as total_profit,
        ROW_NUMBER() OVER (
            ORDER BY COUNT(*) DESC, SUM(o.amount) DESC, MAX(o.updated_at) ASC
        ) as rank
    FROM orders o
    JOIN users u ON o.uid = u.uid
    WHERE o.status = ? 
    AND o.updated_at >= ? 
    AND o.updated_at <= ?
    GROUP BY o.uid, u.username
)
SELECT 
    uid,
    username,
    completed_at,
    order_count,
    total_amount,
    total_profit
FROM user_stats
WHERE rank <= 10
ORDER BY rank
```

**优化后（不使用窗口函数）：**
```sql
SELECT 
    o.uid,
    u.username,
    MAX(o.updated_at) as completed_at,
    COUNT(*) as order_count,
    SUM(o.amount) as total_amount,
    SUM(o.profit_amount) as total_profit
FROM orders o
JOIN users u ON o.uid = u.uid
WHERE o.status = ? 
AND o.updated_at >= ? 
AND o.updated_at <= ?
GROUP BY o.uid, u.username
ORDER BY 
    order_count DESC,
    total_amount DESC,
    completed_at ASC
LIMIT 10
```

### 2. 计算用户排名

**优化前（复杂子查询）：**
```sql
SELECT COUNT(*) + 1 as rank
FROM (
    SELECT 
        o.uid,
        COUNT(*) as order_count,
        SUM(o.amount) as total_amount,
        MAX(o.updated_at) as completed_at
    FROM orders o
    WHERE o.status = ? 
    AND o.updated_at >= ? 
    AND o.updated_at <= ?
    GROUP BY o.uid
    HAVING (
        order_count > (SELECT COUNT(*) FROM orders WHERE status = ? AND uid = ? AND updated_at >= ? AND updated_at <= ?)
        OR (
            order_count = (SELECT COUNT(*) FROM orders WHERE status = ? AND uid = ? AND updated_at >= ? AND updated_at <= ?)
            AND total_amount > (SELECT SUM(amount) FROM orders WHERE status = ? AND uid = ? AND updated_at >= ? AND updated_at <= ?)
        )
        OR (
            order_count = (SELECT COUNT(*) FROM orders WHERE status = ? AND uid = ? AND updated_at >= ? AND updated_at <= ?)
            AND total_amount = (SELECT SUM(amount) FROM orders WHERE status = ? AND uid = ? AND updated_at >= ? AND updated_at <= ?)
            AND completed_at < (SELECT MAX(updated_at) FROM orders WHERE status = ? AND uid = ? AND updated_at >= ? AND updated_at <= ?)
        )
    )
) as ranking
```

**优化后（简化查询）：**
```sql
SELECT COUNT(*) + 1 as rank
FROM (
    SELECT 
        o.uid,
        COUNT(*) as order_count,
        SUM(o.amount) as total_amount,
        MAX(o.updated_at) as completed_at
    FROM orders o
    WHERE o.status = ? 
    AND o.updated_at >= ? 
    AND o.updated_at <= ?
    GROUP BY o.uid
    HAVING 
        order_count > ? 
        OR (order_count = ? AND total_amount > ?)
        OR (order_count = ? AND total_amount = ? AND completed_at < ?)
) as better_users
```

## 优化效果

### 1. 代码可读性提升
- 移除了复杂的窗口函数语法
- SQL查询更加直观易懂
- 减少了嵌套查询的复杂度

### 2. 性能优化
- 简化了查询逻辑，减少了数据库计算负担
- 减少了子查询的数量
- 提高了查询执行效率

### 3. 兼容性提升
- 不依赖窗口函数，兼容性更好
- 适用于更多数据库版本
- 降低了数据库升级的风险

### 4. 维护性提升
- 代码逻辑更清晰
- 更容易调试和修改
- 减少了出错的可能性

## 排名规则

热榜排名按照以下优先级进行排序：

1. **完成订单数量**（降序）：完成订单数量越多，排名越靠前
2. **总金额**（降序）：在订单数量相同的情况下，总金额越高排名越靠前
3. **最新完成时间**（升序）：在订单数量和总金额都相同的情况下，完成时间越早排名越靠前

## 缓存策略

- **缓存键**：`leaderboard:weekly:{week_start_date}`
- **缓存时间**：5分钟
- **缓存内容**：包含前10名用户数据和本周时间范围
- **用户排名**：每次请求时动态计算，确保实时性

## 错误处理

### 1. 用户无订单数据
- 返回默认排名信息（rank: 999）
- 显示用户基本信息（脱敏后的用户名）
- 所有数值字段设为0

### 2. 数据库查询失败
- 返回数据库错误码
- 提供友好的错误信息
- 记录详细的错误日志

### 3. 缓存异常
- 降级到数据库查询
- 不影响正常功能
- 记录缓存异常日志

## 测试验证

使用 `test_scripts/test_leaderboard_optimized.sh` 脚本进行测试：

```bash
./test_scripts/test_leaderboard_optimized.sh
```

测试内容包括：
1. 接口功能测试
2. 缓存功能验证
3. 性能测试
4. 数据结构验证

## 总结

通过移除窗口函数，我们实现了以下目标：

1. **简化了SQL查询逻辑**，提高了代码可读性
2. **提升了查询性能**，减少了数据库计算负担
3. **增强了兼容性**，适用于更多数据库环境
4. **改善了维护性**，降低了代码复杂度

新的实现方案更加稳定、高效，同时保持了原有的功能特性。

## 用户名脱敏规则

热榜中的用户名会进行脱敏处理，保护用户隐私：

### 脱敏规则
1. **用户名长度 = 1**: 不脱敏，直接显示
2. **用户名长度 = 2**: 在中间加 `*`
3. **用户名长度 3-4**: 显示首尾字符，中间用 `*` 替换
4. **用户名长度 ≥ 5**: 显示首尾各2个字符，中间用 `*` 替换

### 示例
- `张三` → `张*三`
- `张三丰` → `张*丰`
- `张三丰李` → `张*李`
- `张三丰李四` → `张三*李四`
- `张三丰李四王` → `张三*李四`
- `test_user_123` → `te*23`

### 优化效果
- 统一了脱敏规则，所有用户名都进行脱敏
- 脱敏长度大幅缩短
- 保留了更多可识别信息
- 提高了用户体验
- 仍然保护了用户隐私 