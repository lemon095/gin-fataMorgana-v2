# 数据库索引优化总结

## 优化概述

基于对项目代码中SQL查询模式的分析，我为所有数据库表添加了必要的索引，以优化查询性能。主要优化了以下查询场景：

1. **用户相关查询** - 按状态、删除时间等条件查询
2. **订单相关查询** - 按用户、状态、时间等条件查询
3. **钱包相关查询** - 按用户、类型、状态等条件查询
4. **拼单相关查询** - 按用户、截止时间、状态等条件查询
5. **排行榜查询** - 按时间范围、状态等条件查询

## 索引优化详情

### 1. orders表索引优化

**原有索引：**
- `uk_order_no` (order_no)
- `uk_uid_period_number` (uid, period_number)
- `idx_uid` (uid)
- `idx_status` (status)
- `idx_expire_time` (expire_time)
- `idx_auditor_uid` (auditor_uid)
- `idx_is_system_order` (is_system_order)
- `idx_created_at` (created_at)

**新增索引：**
- `idx_uid_status_created_at` (uid, status, created_at) - 优化用户订单查询
- `idx_status_updated_at` (status, updated_at) - 优化排行榜查询
- `idx_status_expire_time` (status, expire_time) - 优化过期订单查询
- `idx_uid_created_at` (uid, created_at) - 优化用户订单时间查询
- `idx_period_number_status` (period_number, status) - 优化期数订单查询
- `idx_amount` (amount) - 优化金额统计查询
- `idx_profit_amount` (profit_amount) - 优化利润统计查询
- `idx_updated_at` (updated_at) - 优化更新时间查询

### 2. wallet_transactions表索引优化

**原有索引：**
- `uk_transaction_no` (transaction_no)
- `idx_uid` (uid)
- `idx_type` (type)
- `idx_status` (status)
- `idx_related_order_no` (related_order_no)
- `idx_operator_uid` (operator_uid)
- `idx_created_at` (created_at)

**新增索引：**
- `idx_uid_created_at` (uid, created_at) - 优化用户交易时间查询
- `idx_uid_type_created_at` (uid, type, created_at) - 优化用户交易类型查询
- `idx_uid_status_created_at` (uid, status, created_at) - 优化用户交易状态查询
- `idx_type_status_created_at` (type, status, created_at) - 优化交易类型状态查询
- `idx_amount` (amount) - 优化金额统计查询

### 3. group_buys表索引优化

**原有索引：**
- `uk_group_buy_no` (group_buy_no)
- `idx_order_no` (order_no)
- `idx_creator_uid` (creator_uid)
- `idx_uid` (uid)
- `idx_group_buy_type` (group_buy_type)
- `idx_deadline` (deadline)
- `idx_status` (status)
- `idx_created_at` (created_at)

**新增索引：**
- `idx_uid_deadline_created_at` (uid, deadline, created_at) - 优化用户拼单查询
- `idx_deadline_status` (deadline, status) - 优化活跃拼单查询
- `idx_status_deadline` (status, deadline) - 优化拼单状态查询
- `idx_uid_status` (uid, status) - 优化用户拼单状态查询
- `idx_total_amount` (total_amount) - 优化金额统计查询
- `idx_per_person_amount` (per_person_amount) - 优化人均金额查询

### 4. users表索引优化

**原有索引：**
- `idx_users_uid` (uid) - UNIQUE
- `idx_users_email` (email) - UNIQUE
- `idx_users_username` (username)
- `idx_users_phone` (phone)
- `idx_users_invited_by` (invited_by)
- `idx_users_deleted_at` (deleted_at)

**新增索引：**
- `idx_status_deleted_at` (status, deleted_at) - 优化用户状态查询
- `idx_experience` (experience) - 优化用户等级查询
- `idx_credit_score` (credit_score) - 优化信用分查询
- `idx_created_at` (created_at) - 优化创建时间查询
- `idx_username_email` (username, email) - 优化用户名邮箱查询

### 5. wallets表索引优化

**原有索引：**
- `uk_uid` (uid) - UNIQUE

**新增索引：**
- `idx_status` (status) - 优化钱包状态查询
- `idx_balance` (balance) - 优化余额查询
- `idx_last_active_at` (last_active_at) - 优化活跃时间查询
- `idx_created_at` (created_at) - 优化创建时间查询
- `idx_updated_at` (updated_at) - 优化更新时间查询

### 6. lottery_periods表索引优化

**原有索引：**
- `uk_period_number` (period_number) - UNIQUE
- `idx_status` (status)
- `idx_order_start_time` (order_start_time)
- `idx_order_end_time` (order_end_time)
- `idx_created_at` (created_at)

**新增索引：**
- `idx_order_start_time_order_end_time` (order_start_time, order_end_time) - 优化时间范围查询
- `idx_status_order_start_time` (status, order_start_time) - 优化状态时间查询
- `idx_total_order_amount` (total_order_amount) - 优化金额统计查询

### 7. admin_users表索引优化

**原有索引：**
- `idx_admin_users_admin_id` (admin_id) - UNIQUE
- `idx_admin_users_username` (username) - UNIQUE
- `idx_admin_users_my_invite_code` (my_invite_code) - UNIQUE
- `idx_admin_users_parent_id` (parent_id)
- `idx_admin_users_deleted_at` (deleted_at)

**新增索引：**
- `idx_role_deleted_at` (role, deleted_at) - 优化角色查询
- `idx_status_deleted_at` (status, deleted_at) - 优化状态查询
- `idx_parent_id_deleted_at` (parent_id, deleted_at) - 优化上级查询
- `idx_created_at` (created_at) - 优化创建时间查询

### 8. user_login_logs表索引优化

**原有索引：**
- `idx_uid` (uid)
- `idx_username` (username)
- `idx_email` (email)
- `idx_login_ip` (login_ip)
- `idx_login_time` (login_time)
- `idx_status` (status)

**新增索引：**
- `idx_uid_login_time` (uid, login_time) - 优化用户登录时间查询
- `idx_uid_status` (uid, status) - 优化用户登录状态查询
- `idx_login_ip_login_time` (login_ip, login_time) - 优化IP登录时间查询
- `idx_status_login_time` (status, login_time) - 优化状态登录时间查询

### 9. amount_config表索引优化

**原有索引：**
- `idx_type` (type)
- `idx_is_active` (is_active)
- `idx_sort_order` (sort_order)

**新增索引：**
- `idx_type_is_active_sort` (type, is_active, sort_order) - 优化配置查询
- `idx_amount` (amount) - 优化金额查询

### 10. announcements表索引优化

**原有索引：**
- `idx_deleted_at` (deleted_at)

**新增索引：**
- `idx_status` (status) - 优化状态查询
- `idx_tag` (tag) - 优化标签查询
- `idx_status_deleted_at_created_at` (status, deleted_at, created_at) - 优化发布查询
- `idx_created_at` (created_at) - 优化创建时间查询

### 11. announcement_banners表索引优化

**原有索引：**
- `idx_announcement_id` (announcement_id)
- `idx_deleted_at` (deleted_at)

**新增索引：**
- `idx_announcement_id_deleted_at_sort` (announcement_id, deleted_at, sort) - 优化排序查询
- `idx_sort` (sort) - 优化排序查询

### 12. member_level表索引优化

**原有索引：**
- `uniq_level` (level) - UNIQUE
- `idx_deleted_at` (deleted_at)

**新增索引：**
- `idx_level_deleted_at` (level, deleted_at) - 优化等级查询
- `idx_cashback_ratio` (cashback_ratio) - 优化返现比例查询

## 查询性能优化效果

### 1. 订单查询优化
```sql
-- 优化前：需要扫描大量数据
SELECT * FROM orders WHERE uid = 'user123' AND status = 'success' AND created_at <= NOW();

-- 优化后：使用复合索引 idx_uid_status_created_at
-- 查询时间从 O(n) 降低到 O(log n)
```

### 2. 排行榜查询优化
```sql
-- 优化前：需要扫描所有订单
SELECT uid, COUNT(*) as order_count, SUM(amount) as total_amount 
FROM orders 
WHERE status = 'success' AND updated_at >= ? AND updated_at <= ?
GROUP BY uid;

-- 优化后：使用复合索引 idx_status_updated_at
-- 查询时间大幅减少
```

### 3. 钱包交易查询优化
```sql
-- 优化前：需要扫描用户所有交易
SELECT * FROM wallet_transactions 
WHERE uid = 'user123' AND type = 'recharge' AND created_at >= ? AND created_at <= ?;

-- 优化后：使用复合索引 idx_uid_type_created_at
-- 查询效率显著提升
```

### 4. 拼单查询优化
```sql
-- 优化前：需要扫描所有拼单
SELECT * FROM group_buys 
WHERE uid = 'user123' AND deadline > NOW() AND created_at <= NOW();

-- 优化后：使用复合索引 idx_uid_deadline_created_at
-- 查询性能大幅提升
```

## 索引维护建议

### 1. 定期分析索引使用情况
```sql
-- 查看索引使用统计
SELECT 
    table_name,
    index_name,
    cardinality,
    sub_part,
    packed,
    null,
    index_type
FROM information_schema.statistics 
WHERE table_schema = 'gin_fataMorgana'
ORDER BY table_name, index_name;
```

### 2. 监控慢查询
```sql
-- 开启慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;

-- 查看慢查询
SELECT * FROM mysql.slow_log ORDER BY start_time DESC LIMIT 10;
```

### 3. 定期更新表统计信息
```sql
-- 更新表统计信息
ANALYZE TABLE orders;
ANALYZE TABLE wallet_transactions;
ANALYZE TABLE group_buys;
ANALYZE TABLE users;
ANALYZE TABLE wallets;
```

### 4. 监控索引大小
```sql
-- 查看索引大小
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)',
    ROUND((data_length / 1024 / 1024), 2) AS 'Data (MB)',
    ROUND((index_length / 1024 / 1024), 2) AS 'Index (MB)'
FROM information_schema.tables 
WHERE table_schema = 'gin_fataMorgana'
ORDER BY (data_length + index_length) DESC;
```

## 注意事项

### 1. 索引选择原则
- **最左前缀原则**：复合索引的列顺序很重要
- **选择性原则**：选择区分度高的列作为索引
- **覆盖索引**：尽量使用覆盖索引减少回表查询

### 2. 索引维护成本
- **写入性能**：索引会增加INSERT、UPDATE、DELETE的开销
- **存储空间**：索引会占用额外的存储空间
- **维护成本**：需要定期维护和优化索引

### 3. 监控指标
- **查询响应时间**：监控关键查询的响应时间
- **索引命中率**：监控索引的使用情况
- **慢查询数量**：监控慢查询的数量和分布

## 总结

通过这次索引优化，我们为项目中的所有主要表添加了必要的索引，特别是复合索引，这将显著提升以下场景的查询性能：

1. **用户订单查询** - 提升用户体验
2. **排行榜查询** - 提升系统响应速度
3. **钱包交易查询** - 提升财务查询效率
4. **拼单查询** - 提升拼单功能性能
5. **管理后台查询** - 提升管理效率

这些优化将显著提升系统的整体性能，特别是在数据量增长的情况下，效果会更加明显。 