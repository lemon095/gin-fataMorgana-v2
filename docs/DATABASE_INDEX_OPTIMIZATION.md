# 数据库索引优化文档

## 概述

本文档详细分析了整个项目的数据库索引需求，基于实际的查询模式和业务逻辑，为每个表提供了完整的索引优化方案。

## 索引优化原则

1. **查询频率优先**：优先为高频查询字段创建索引
2. **选择性原则**：优先为选择性高的字段创建索引
3. **复合索引优化**：合理使用复合索引减少索引数量
4. **覆盖索引**：尽可能使用覆盖索引减少回表查询
5. **避免过度索引**：避免创建不必要的索引影响写入性能

## 表索引详细分析

### 1. users 表（用户表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
UNIQUE KEY `idx_users_uid` (`uid`)
UNIQUE KEY `idx_users_email` (`email`)
KEY `idx_users_username` (`username`)
KEY `idx_users_phone` (`phone`)
KEY `idx_users_invited_by` (`invited_by`)
KEY `idx_users_deleted_at` (`deleted_at`)
```

**建议优化：**
```sql
-- 添加复合索引优化用户状态查询
ALTER TABLE `users` ADD INDEX `idx_status_deleted_at` (`status`, `deleted_at`);

-- 添加经验值索引（用于等级计算）
ALTER TABLE `users` ADD INDEX `idx_experience` (`experience`);

-- 添加信用分索引（用于用户筛选）
ALTER TABLE `users` ADD INDEX `idx_credit_score` (`credit_score`);

-- 添加创建时间索引（用于用户统计）
ALTER TABLE `users` ADD INDEX `idx_created_at` (`created_at`);

-- 添加复合索引优化用户搜索
ALTER TABLE `users` ADD INDEX `idx_username_email` (`username`, `email`);
```

**优化理由：**
- `status + deleted_at`：高频查询活跃用户
- `experience`：用户等级计算和升级查询
- `credit_score`：用户信用筛选
- `created_at`：用户注册统计和分页查询
- `username + email`：用户搜索功能

### 2. admin_users 表（管理员用户表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
UNIQUE KEY `idx_admin_users_admin_id` (`admin_id`)
UNIQUE KEY `idx_admin_users_username` (`username`)
UNIQUE KEY `idx_admin_users_my_invite_code` (`my_invite_code`)
KEY `idx_admin_users_parent_id` (`parent_id`)
KEY `idx_admin_users_deleted_at` (`deleted_at`)
```

**建议优化：**
```sql
-- 添加角色索引（用于权限管理）
ALTER TABLE `admin_users` ADD INDEX `idx_role` (`role`);

-- 添加状态索引（用于活跃管理员查询）
ALTER TABLE `admin_users` ADD INDEX `idx_status` (`status`);

-- 添加复合索引优化管理员列表查询
ALTER TABLE `admin_users` ADD INDEX `idx_role_status_deleted_at` (`role`, `status`, `deleted_at`);

-- 添加创建时间索引（用于管理员统计）
ALTER TABLE `admin_users` ADD INDEX `idx_created_at` (`created_at`);
```

**优化理由：**
- `role`：按角色筛选管理员
- `status`：查询活跃管理员
- `role + status + deleted_at`：管理员列表分页查询
- `created_at`：管理员注册统计

### 3. wallets 表（钱包表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
UNIQUE KEY `uk_uid` (`uid`)
```

**建议优化：**
```sql
-- 添加状态索引（用于钱包状态查询）
ALTER TABLE `wallets` ADD INDEX `idx_status` (`status`);

-- 添加余额索引（用于余额统计）
ALTER TABLE `wallets` ADD INDEX `idx_balance` (`balance`);

-- 添加最后活跃时间索引（用于活跃钱包查询）
ALTER TABLE `wallets` ADD INDEX `idx_last_active_at` (`last_active_at`);

-- 添加复合索引优化钱包查询
ALTER TABLE `wallets` ADD INDEX `idx_status_balance` (`status`, `balance`);
```

**优化理由：**
- `status`：查询正常/冻结钱包
- `balance`：余额统计和筛选
- `last_active_at`：活跃钱包查询
- `status + balance`：钱包状态和余额组合查询

### 4. wallet_transactions 表（钱包交易流水表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
UNIQUE KEY `uk_transaction_no` (`transaction_no`)
KEY `idx_uid` (`uid`)
KEY `idx_type` (`type`)
KEY `idx_status` (`status`)
KEY `idx_related_order_no` (`related_order_no`)
KEY `idx_operator_uid` (`operator_uid`)
KEY `idx_created_at` (`created_at`)
```

**建议优化：**
```sql
-- 添加复合索引优化用户交易查询
ALTER TABLE `wallet_transactions` ADD INDEX `idx_uid_created_at` (`uid`, `created_at`);

-- 添加复合索引优化交易类型查询
ALTER TABLE `wallet_transactions` ADD INDEX `idx_uid_type_created_at` (`uid`, `type`, `created_at`);

-- 添加复合索引优化交易状态查询
ALTER TABLE `wallet_transactions` ADD INDEX `idx_uid_status_created_at` (`uid`, `status`, `created_at`);

-- 添加金额索引（用于金额统计）
ALTER TABLE `wallet_transactions` ADD INDEX `idx_amount` (`amount`);

-- 添加复合索引优化交易统计
ALTER TABLE `wallet_transactions` ADD INDEX `idx_type_status_created_at` (`type`, `status`, `created_at`);
```

**优化理由：**
- `uid + created_at`：用户交易历史分页查询
- `uid + type + created_at`：用户特定类型交易查询
- `uid + status + created_at`：用户交易状态查询
- `amount`：交易金额统计
- `type + status + created_at`：交易类型统计

### 5. orders 表（订单表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
UNIQUE KEY `uk_order_no` (`order_no`)
UNIQUE KEY `uk_uid_period_number` (`uid`, `period_number`)
KEY `idx_uid` (`uid`)
KEY `idx_status` (`status`)
KEY `idx_expire_time` (`expire_time`)
KEY `idx_auditor_uid` (`auditor_uid`)
KEY `idx_created_at` (`created_at`)
```

**建议优化：**
```sql
-- 添加复合索引优化用户订单查询
ALTER TABLE `orders` ADD INDEX `idx_uid_status_created_at` (`uid`, `status`, `created_at`);

-- 添加期号索引（用于期数相关查询）
ALTER TABLE `orders` ADD INDEX `idx_period_number` (`period_number`);

-- 添加复合索引优化期数订单查询
ALTER TABLE `orders` ADD INDEX `idx_period_number_status` (`period_number`, `status`);

-- 添加金额索引（用于金额统计）
ALTER TABLE `orders` ADD INDEX `idx_amount` (`amount`);

-- 添加利润索引（用于利润统计）
ALTER TABLE `orders` ADD INDEX `idx_profit_amount` (`profit_amount`);

-- 添加复合索引优化过期订单查询
ALTER TABLE `orders` ADD INDEX `idx_status_expire_time` (`status`, `expire_time`);

-- 添加更新时间索引（用于排行榜查询）
ALTER TABLE `orders` ADD INDEX `idx_updated_at` (`updated_at`);

-- 添加复合索引优化排行榜查询
ALTER TABLE `orders` ADD INDEX `idx_status_updated_at` (`status`, `updated_at`);
```

**优化理由：**
- `uid + status + created_at`：用户订单列表分页查询
- `period_number`：期数相关订单查询
- `period_number + status`：期数订单状态统计
- `amount`：订单金额统计
- `profit_amount`：订单利润统计
- `status + expire_time`：过期订单处理
- `updated_at`：订单更新时间查询
- `status + updated_at`：排行榜数据查询

### 6. group_buys 表（拼单表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
UNIQUE KEY `uk_group_buy_no` (`group_buy_no`)
KEY `idx_order_no` (`order_no`)
KEY `idx_creator_uid` (`creator_uid`)
KEY `idx_uid` (`uid`)
KEY `idx_group_buy_type` (`group_buy_type`)
KEY `idx_deadline` (`deadline`)
KEY `idx_status` (`status`)
KEY `idx_created_at` (`created_at`)
```

**建议优化：**
```sql
-- 添加复合索引优化活跃拼单查询
ALTER TABLE `group_buys` ADD INDEX `idx_deadline_status` (`deadline`, `status`);

-- 添加复合索引优化用户拼单查询
ALTER TABLE `group_buys` ADD INDEX `idx_uid_deadline` (`uid`, `deadline`);

-- 添加复合索引优化拼单类型查询
ALTER TABLE `group_buys` ADD INDEX `idx_type_status_deadline` (`group_buy_type`, `status`, `deadline`);

-- 添加复合索引优化创建者查询
ALTER TABLE `group_buys` ADD INDEX `idx_creator_uid_status` (`creator_uid`, `status`);
```

**优化理由：**
- `deadline + status`：活跃拼单查询（高频）
- `uid + deadline`：用户参与的拼单查询
- `type + status + deadline`：特定类型活跃拼单查询
- `creator_uid + status`：创建者拼单查询

### 7. user_login_logs 表（用户登录日志表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
KEY `idx_uid` (`uid`)
KEY `idx_username` (`username`)
KEY `idx_email` (`email`)
KEY `idx_login_ip` (`login_ip`)
KEY `idx_login_time` (`login_time`)
KEY `idx_status` (`status`)
```

**建议优化：**
```sql
-- 添加复合索引优化用户登录历史查询
ALTER TABLE `user_login_logs` ADD INDEX `idx_uid_login_time` (`uid`, `login_time`);

-- 添加复合索引优化登录状态查询
ALTER TABLE `user_login_logs` ADD INDEX `idx_uid_status_login_time` (`uid`, `status`, `login_time`);

-- 添加复合索引优化IP查询
ALTER TABLE `user_login_logs` ADD INDEX `idx_uid_login_ip` (`uid`, `login_ip`);

-- 添加复合索引优化失败登录查询
ALTER TABLE `user_login_logs` ADD INDEX `idx_uid_status` (`uid`, `status`);

-- 添加创建时间索引（用于日志清理）
ALTER TABLE `user_login_logs` ADD INDEX `idx_created_at` (`created_at`);
```

**优化理由：**
- `uid + login_time`：用户登录历史分页查询
- `uid + status + login_time`：用户登录状态统计
- `uid + login_ip`：用户IP登录记录查询
- `uid + status`：用户登录失败次数统计
- `created_at`：日志清理和统计

### 8. amount_config 表（金额配置表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
KEY `idx_type` (`type`)
KEY `idx_is_active` (`is_active`)
KEY `idx_sort_order` (`sort_order`)
```

**建议优化：**
```sql
-- 添加复合索引优化配置查询
ALTER TABLE `amount_config` ADD INDEX `idx_type_is_active_sort` (`type`, `is_active`, `sort_order`);

-- 添加金额索引（用于金额范围查询）
ALTER TABLE `amount_config` ADD INDEX `idx_amount` (`amount`);
```

**优化理由：**
- `type + is_active + sort_order`：获取激活配置列表（高频查询）
- `amount`：金额范围查询和统计

### 9. announcements 表（公告表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
KEY `idx_deleted_at` (`deleted_at`)
```

**建议优化：**
```sql
-- 添加状态索引（用于发布状态查询）
ALTER TABLE `announcements` ADD INDEX `idx_status` (`status`);

-- 添加标签索引（用于标签筛选）
ALTER TABLE `announcements` ADD INDEX `idx_tag` (`tag`);

-- 添加复合索引优化公告列表查询
ALTER TABLE `announcements` ADD INDEX `idx_status_deleted_at_created_at` (`status`, `deleted_at`, `created_at`);

-- 添加创建时间索引（用于时间排序）
ALTER TABLE `announcements` ADD INDEX `idx_created_at` (`created_at`);
```

**优化理由：**
- `status`：查询已发布公告
- `tag`：按标签筛选公告
- `status + deleted_at + created_at`：公告列表分页查询
- `created_at`：公告时间排序

### 10. lottery_periods 表（期数开奖表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
UNIQUE KEY `uk_period_number` (`period_number`)
KEY `idx_status` (`status`)
KEY `idx_order_start_time` (`order_start_time`)
KEY `idx_order_end_time` (`order_end_time`)
KEY `idx_created_at` (`created_at`)
```

**建议优化：**
```sql
-- 添加复合索引优化当前期数查询
ALTER TABLE `lottery_periods` ADD INDEX `idx_status_order_end_time` (`status`, `order_end_time`);

-- 添加复合索引优化期数时间范围查询
ALTER TABLE `lottery_periods` ADD INDEX `idx_order_start_time_order_end_time` (`order_start_time`, `order_end_time`);

-- 添加开奖结果索引（用于开奖查询）
ALTER TABLE `lottery_periods` ADD INDEX `idx_lottery_result` (`lottery_result`);
```

**优化理由：**
- `status + order_end_time`：查询当前活跃期数（高频）
- `order_start_time + order_end_time`：期数时间范围查询
- `lottery_result`：开奖结果查询

### 11. member_level 表（用户等级配置表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
UNIQUE KEY `uniq_level` (`level`)
KEY `idx_deleted_at` (`deleted_at`)
```

**建议优化：**
```sql
-- 添加等级索引（用于等级范围查询）
ALTER TABLE `member_level` ADD INDEX `idx_level_deleted_at` (`level`, `deleted_at`);

-- 添加返现比例索引（用于返现统计）
ALTER TABLE `member_level` ADD INDEX `idx_cashback_ratio` (`cashback_ratio`);
```

**优化理由：**
- `level + deleted_at`：等级配置查询
- `cashback_ratio`：返现比例统计

### 12. announcement_banners 表（公告图片表）

**现有索引：**
```sql
PRIMARY KEY (`id`)
KEY `idx_announcement_id` (`announcement_id`)
KEY `idx_deleted_at` (`deleted_at`)
```

**建议优化：**
```sql
-- 添加复合索引优化图片查询
ALTER TABLE `announcement_banners` ADD INDEX `idx_announcement_id_deleted_at_sort` (`announcement_id`, `deleted_at`, `sort`);

-- 添加排序索引（用于图片排序）
ALTER TABLE `announcement_banners` ADD INDEX `idx_sort` (`sort`);
```

**优化理由：**
- `announcement_id + deleted_at + sort`：公告图片列表查询
- `sort`：图片排序查询

## 性能监控建议

### 1. 慢查询监控
```sql
-- 开启慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;
SET GLOBAL log_queries_not_using_indexes = 'ON';
```

### 2. 索引使用情况监控
```sql
-- 查看索引使用情况
SELECT 
    table_name,
    index_name,
    cardinality,
    sub_part,
    packed,
    nullable,
    index_type
FROM information_schema.statistics 
WHERE table_schema = 'gin_fataMorgana'
ORDER BY table_name, index_name;
```

### 3. 查询性能分析
```sql
-- 分析查询执行计划
EXPLAIN SELECT * FROM users WHERE status = 1 AND deleted_at IS NULL;

-- 查看表统计信息
ANALYZE TABLE users;
```

## 索引维护建议

### 1. 定期重建索引
```sql
-- 重建索引（建议在低峰期执行）
OPTIMIZE TABLE users;
OPTIMIZE TABLE orders;
OPTIMIZE TABLE wallet_transactions;
```

### 2. 监控索引大小
```sql
-- 查看索引大小
SELECT 
    table_name,
    index_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)'
FROM information_schema.tables 
WHERE table_schema = 'gin_fataMorgana'
ORDER BY (data_length + index_length) DESC;
```

### 3. 删除无用索引
```sql
-- 查看未使用的索引
SELECT 
    object_schema,
    object_name,
    index_name
FROM performance_schema.table_io_waits_summary_by_index_usage
WHERE index_name IS NOT NULL
AND count_star = 0;
```

## 实施计划

### 第一阶段：核心表索引优化
1. users 表复合索引
2. orders 表复合索引
3. wallet_transactions 表复合索引

### 第二阶段：业务表索引优化
1. group_buys 表复合索引
2. user_login_logs 表复合索引
3. lottery_periods 表复合索引

### 第三阶段：配置表索引优化
1. amount_config 表复合索引
2. announcements 表复合索引
3. member_level 表索引

### 第四阶段：监控和调优
1. 实施性能监控
2. 分析慢查询
3. 根据实际使用情况调整索引

## 注意事项

1. **分批执行**：索引创建会锁表，建议在低峰期分批执行
2. **监控影响**：创建索引期间监控系统性能
3. **备份数据**：执行前务必备份数据库
4. **测试验证**：在生产环境执行前先在测试环境验证
5. **回滚方案**：准备索引删除的回滚方案

## 预期效果

1. **查询性能提升**：高频查询响应时间减少50-80%
2. **系统并发能力**：支持更高的并发用户数
3. **资源利用率**：减少CPU和I/O资源消耗
4. **用户体验**：页面加载速度明显提升
5. **系统稳定性**：减少数据库连接超时和死锁

通过以上索引优化方案，可以显著提升系统的查询性能和并发处理能力，为用户提供更好的使用体验。 