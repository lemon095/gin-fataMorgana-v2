-- 数据库索引优化脚本
-- 执行前请务必备份数据库
-- 建议在低峰期分批执行

USE gin_fataMorgana;

-- ========================================
-- 第一阶段：核心表索引优化
-- ========================================

-- 1. users 表索引优化
ALTER TABLE `users` ADD INDEX `idx_status_deleted_at` (`status`, `deleted_at`);
ALTER TABLE `users` ADD INDEX `idx_experience` (`experience`);
ALTER TABLE `users` ADD INDEX `idx_credit_score` (`credit_score`);
ALTER TABLE `users` ADD INDEX `idx_created_at` (`created_at`);
ALTER TABLE `users` ADD INDEX `idx_username_email` (`username`, `email`);

-- 2. orders 表索引优化
ALTER TABLE `orders` ADD INDEX `idx_uid_status_created_at` (`uid`, `status`, `created_at`);
ALTER TABLE `orders` ADD INDEX `idx_period_number` (`period_number`);
ALTER TABLE `orders` ADD INDEX `idx_period_number_status` (`period_number`, `status`);
ALTER TABLE `orders` ADD INDEX `idx_amount` (`amount`);
ALTER TABLE `orders` ADD INDEX `idx_profit_amount` (`profit_amount`);
ALTER TABLE `orders` ADD INDEX `idx_status_expire_time` (`status`, `expire_time`);
ALTER TABLE `orders` ADD INDEX `idx_updated_at` (`updated_at`);
ALTER TABLE `orders` ADD INDEX `idx_status_updated_at` (`status`, `updated_at`);

-- 3. wallet_transactions 表索引优化
ALTER TABLE `wallet_transactions` ADD INDEX `idx_uid_created_at` (`uid`, `created_at`);
ALTER TABLE `wallet_transactions` ADD INDEX `idx_uid_type_created_at` (`uid`, `type`, `created_at`);
ALTER TABLE `wallet_transactions` ADD INDEX `idx_uid_status_created_at` (`uid`, `status`, `created_at`);
ALTER TABLE `wallet_transactions` ADD INDEX `idx_amount` (`amount`);
ALTER TABLE `wallet_transactions` ADD INDEX `idx_type_status_created_at` (`type`, `status`, `created_at`);

-- ========================================
-- 第二阶段：业务表索引优化
-- ========================================

-- 4. group_buys 表索引优化
ALTER TABLE `group_buys` ADD INDEX `idx_deadline_status` (`deadline`, `status`);
ALTER TABLE `group_buys` ADD INDEX `idx_uid_deadline` (`uid`, `deadline`);
ALTER TABLE `group_buys` ADD INDEX `idx_type_status_deadline` (`group_buy_type`, `status`, `deadline`);
ALTER TABLE `group_buys` ADD INDEX `idx_creator_uid_status` (`creator_uid`, `status`);

-- 5. user_login_logs 表索引优化
ALTER TABLE `user_login_logs` ADD INDEX `idx_uid_login_time` (`uid`, `login_time`);
ALTER TABLE `user_login_logs` ADD INDEX `idx_uid_status_login_time` (`uid`, `status`, `login_time`);
ALTER TABLE `user_login_logs` ADD INDEX `idx_uid_login_ip` (`uid`, `login_ip`);
ALTER TABLE `user_login_logs` ADD INDEX `idx_uid_status` (`uid`, `status`);
ALTER TABLE `user_login_logs` ADD INDEX `idx_created_at` (`created_at`);

-- 6. lottery_periods 表索引优化
ALTER TABLE `lottery_periods` ADD INDEX `idx_status_order_end_time` (`status`, `order_end_time`);
ALTER TABLE `lottery_periods` ADD INDEX `idx_order_start_time_order_end_time` (`order_start_time`, `order_end_time`);
ALTER TABLE `lottery_periods` ADD INDEX `idx_lottery_result` (`lottery_result`);

-- ========================================
-- 第三阶段：配置表索引优化
-- ========================================

-- 7. admin_users 表索引优化
ALTER TABLE `admin_users` ADD INDEX `idx_role` (`role`);
ALTER TABLE `admin_users` ADD INDEX `idx_status` (`status`);
ALTER TABLE `admin_users` ADD INDEX `idx_role_status_deleted_at` (`role`, `status`, `deleted_at`);
ALTER TABLE `admin_users` ADD INDEX `idx_created_at` (`created_at`);

-- 8. wallets 表索引优化
ALTER TABLE `wallets` ADD INDEX `idx_status` (`status`);
ALTER TABLE `wallets` ADD INDEX `idx_balance` (`balance`);
ALTER TABLE `wallets` ADD INDEX `idx_last_active_at` (`last_active_at`);
ALTER TABLE `wallets` ADD INDEX `idx_status_balance` (`status`, `balance`);

-- 9. amount_config 表索引优化
ALTER TABLE `amount_config` ADD INDEX `idx_type_is_active_sort` (`type`, `is_active`, `sort_order`);
ALTER TABLE `amount_config` ADD INDEX `idx_amount` (`amount`);

-- 10. announcements 表索引优化
ALTER TABLE `announcements` ADD INDEX `idx_status` (`status`);
ALTER TABLE `announcements` ADD INDEX `idx_tag` (`tag`);
ALTER TABLE `announcements` ADD INDEX `idx_status_deleted_at_created_at` (`status`, `deleted_at`, `created_at`);
ALTER TABLE `announcements` ADD INDEX `idx_created_at` (`created_at`);

-- 11. member_level 表索引优化
ALTER TABLE `member_level` ADD INDEX `idx_level_deleted_at` (`level`, `deleted_at`);
ALTER TABLE `member_level` ADD INDEX `idx_cashback_ratio` (`cashback_ratio`);

-- 12. announcement_banners 表索引优化
ALTER TABLE `announcement_banners` ADD INDEX `idx_announcement_id_deleted_at_sort` (`announcement_id`, `deleted_at`, `sort`);
ALTER TABLE `announcement_banners` ADD INDEX `idx_sort` (`sort`);

-- ========================================
-- 索引创建完成后的优化操作
-- ========================================

-- 更新表统计信息
ANALYZE TABLE users;
ANALYZE TABLE orders;
ANALYZE TABLE wallet_transactions;
ANALYZE TABLE group_buys;
ANALYZE TABLE user_login_logs;
ANALYZE TABLE lottery_periods;
ANALYZE TABLE admin_users;
ANALYZE TABLE wallets;
ANALYZE TABLE amount_config;
ANALYZE TABLE announcements;
ANALYZE TABLE member_level;
ANALYZE TABLE announcement_banners;

-- 显示索引创建结果
SELECT 
    table_name,
    index_name,
    column_name,
    seq_in_index,
    cardinality
FROM information_schema.statistics 
WHERE table_schema = 'gin_fataMorgana'
AND index_name NOT LIKE 'PRIMARY'
ORDER BY table_name, index_name, seq_in_index;

-- 显示表大小统计
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)',
    ROUND((data_length / 1024 / 1024), 2) AS 'Data (MB)',
    ROUND((index_length / 1024 / 1024), 2) AS 'Index (MB)',
    table_rows
FROM information_schema.tables 
WHERE table_schema = 'gin_fataMorgana'
ORDER BY (data_length + index_length) DESC; 