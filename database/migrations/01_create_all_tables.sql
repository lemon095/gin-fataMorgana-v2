-- ========================================
-- Gin-FataMorgana 数据库初始化脚本
-- 包含所有表的创建和索引优化
-- ========================================

USE future;

-- ========================================
-- 1. 用户相关表
-- ========================================

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(8) NOT NULL COMMENT '用户唯一ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `password_hash` varchar(255) NOT NULL COMMENT '密码哈希',
  `status` varchar(20) NOT NULL DEFAULT 'active' COMMENT '用户状态: active-正常, inactive-禁用',
  `experience` int NOT NULL DEFAULT 0 COMMENT '经验值',
  `credit_score` int NOT NULL DEFAULT 100 COMMENT '信用分',
  `bank_card_number` varchar(20) DEFAULT NULL COMMENT '银行卡号',
  `bank_card_holder` varchar(50) DEFAULT NULL COMMENT '持卡人姓名',
  `bank_name` varchar(50) DEFAULT NULL COMMENT '银行名称',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uid` (`uid`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_status_deleted_at` (`status`, `deleted_at`),
  KEY `idx_experience` (`experience`),
  KEY `idx_credit_score` (`credit_score`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_username_email` (`username`, `email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表 - 存储用户基本信息、认证信息、银行卡信息、经验值、信用分等';

-- 钱包表
CREATE TABLE IF NOT EXISTS `wallets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(8) NOT NULL COMMENT '用户唯一ID',
  `balance` decimal(15,2) NOT NULL DEFAULT '0.00' COMMENT '余额',
  `frozen_balance` decimal(15,2) NOT NULL DEFAULT '0.00' COMMENT '冻结余额',
  `total_income` decimal(15,2) NOT NULL DEFAULT '0.00' COMMENT '总收入',
  `total_expense` decimal(15,2) NOT NULL DEFAULT '0.00' COMMENT '总支出',
  `status` varchar(20) NOT NULL DEFAULT 'active' COMMENT '钱包状态: active-正常, frozen-冻结',
  `last_active_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后活跃时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uid` (`uid`),
  KEY `idx_status` (`status`),
  KEY `idx_balance` (`balance`),
  KEY `idx_last_active_at` (`last_active_at`),
  KEY `idx_status_balance` (`status`, `balance`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='钱包表 - 存储用户钱包信息，包括余额、冻结余额、总收入、总支出等';

-- 钱包交易流水表
CREATE TABLE IF NOT EXISTS `wallet_transactions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `transaction_no` varchar(32) NOT NULL COMMENT '交易流水号',
  `uid` varchar(8) NOT NULL COMMENT '用户唯一ID',
  `type` varchar(20) NOT NULL COMMENT '交易类型: recharge-充值, withdraw-提现, order_buy-购买, group_buy-拼单, profit-收益',
  `amount` decimal(15,2) NOT NULL COMMENT '交易金额',
  `balance_before` decimal(15,2) NOT NULL COMMENT '交易前余额',
  `balance_after` decimal(15,2) NOT NULL COMMENT '交易后余额',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '交易状态: pending-处理中, success-成功, failed-失败',
  `description` varchar(255) DEFAULT NULL COMMENT '交易描述',
  `related_order_no` varchar(32) DEFAULT NULL COMMENT '关联订单号',
  `operator_uid` varchar(8) DEFAULT NULL COMMENT '操作员ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transaction_no` (`transaction_no`),
  KEY `idx_uid_created_at` (`uid`, `created_at`),
  KEY `idx_uid_type_created_at` (`uid`, `type`, `created_at`),
  KEY `idx_uid_status_created_at` (`uid`, `status`, `created_at`),
  KEY `idx_amount` (`amount`),
  KEY `idx_type_status_created_at` (`type`, `status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='钱包交易流水表 - 记录所有钱包交易明细，包括充值、提现、购买、拼单等操作';

-- 用户登录日志表
CREATE TABLE IF NOT EXISTS `user_login_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(8) NOT NULL COMMENT '用户唯一ID',
  `login_time` timestamp NOT NULL COMMENT '登录时间',
  `login_ip` varchar(45) NOT NULL COMMENT '登录IP地址',
  `user_agent` varchar(500) DEFAULT NULL COMMENT '用户代理',
  `device_info` varchar(200) DEFAULT NULL COMMENT '设备信息',
  `status` varchar(20) NOT NULL DEFAULT 'success' COMMENT '登录状态: success-成功, failed-失败',
  `failure_reason` varchar(200) DEFAULT NULL COMMENT '失败原因',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid_login_time` (`uid`, `login_time`),
  KEY `idx_uid_status_login_time` (`uid`, `status`, `login_time`),
  KEY `idx_uid_login_ip` (`uid`, `login_ip`),
  KEY `idx_uid_status` (`uid`, `status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户登录日志表 - 记录用户登录历史，包括登录时间、IP地址、设备信息、登录状态等';

-- ========================================
-- 2. 订单相关表
-- ========================================

-- 订单表
CREATE TABLE IF NOT EXISTS `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `order_no` varchar(32) NOT NULL COMMENT '订单编号',
  `uid` varchar(8) NOT NULL COMMENT '用户唯一ID',
  `period_number` varchar(32) NOT NULL COMMENT '期号',
  `amount` decimal(15,2) NOT NULL COMMENT '订单金额',
  `profit_amount` decimal(15,2) NOT NULL COMMENT '利润金额',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '订单状态',
  `expire_time` datetime NOT NULL COMMENT '订单剩余时间',
  `like_count` int NOT NULL DEFAULT 0 COMMENT '点赞数',
  `share_count` int NOT NULL DEFAULT 0 COMMENT '转发数',
  `follow_count` int NOT NULL DEFAULT 0 COMMENT '关注数',
  `favorite_count` int NOT NULL DEFAULT 0 COMMENT '收藏数',
  `like_status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '点赞完成状态',
  `share_status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '转发完成状态',
  `follow_status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '关注完成状态',
  `favorite_status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '收藏完成状态',
  `auditor_uid` varchar(8) DEFAULT NULL COMMENT '审核员ID',
  `is_system_order` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否系统订单 0-否 1-是',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  UNIQUE KEY `uk_uid_period_number` (`uid`, `period_number`),
  KEY `idx_uid_status_created_at` (`uid`, `status`, `created_at`),
  KEY `idx_period_number` (`period_number`),
  KEY `idx_period_number_status` (`period_number`, `status`),
  KEY `idx_amount` (`amount`),
  KEY `idx_profit_amount` (`profit_amount`),
  KEY `idx_status_expire_time` (`status`, `expire_time`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_status_updated_at` (`status`, `updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表 - 存储用户订单信息，包括订单金额、状态、任务完成情况等';

-- 拼单表
CREATE TABLE IF NOT EXISTS `group_buys` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `group_buy_no` varchar(32) NOT NULL COMMENT '拼单编号',
  `order_no` varchar(32) DEFAULT NULL COMMENT '关联订单编号',
  `uid` varchar(8) NOT NULL COMMENT '参与用户ID',
  `creator_uid` varchar(8) NOT NULL COMMENT '创建用户ID',
  `current_participants` int NOT NULL DEFAULT 1 COMMENT '当前参与人数',
  `target_participants` int NOT NULL DEFAULT 2 COMMENT '目标参与人数',
  `group_buy_type` varchar(20) NOT NULL DEFAULT 'normal' COMMENT '拼单类型: normal-普通, flash-限时, vip-VIP',
  `total_amount` decimal(15,2) NOT NULL COMMENT '拼单总金额',
  `paid_amount` decimal(15,2) NOT NULL DEFAULT '0.00' COMMENT '已付款金额',
  `per_person_amount` decimal(15,2) NOT NULL COMMENT '每人需要付款金额',
  `profit_margin` decimal(5,4) NOT NULL DEFAULT '0.0000' COMMENT '利润比例（小数）',
  `deadline` timestamp NOT NULL COMMENT '拼单截止时间',
  `status` varchar(20) NOT NULL DEFAULT 'not_started' COMMENT '拼单状态: not_started-未开启, pending-进行中, success-已完成',
  `description` text COMMENT '拼单描述',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_buy_no` (`group_buy_no`),
  KEY `idx_deadline_status` (`deadline`, `status`),
  KEY `idx_uid_deadline` (`uid`, `deadline`),
  KEY `idx_type_status_deadline` (`group_buy_type`, `status`, `deadline`),
  KEY `idx_creator_uid_status` (`creator_uid`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='拼单表 - 存储拼单信息，支持多人拼单功能';

-- 游戏期数表
CREATE TABLE IF NOT EXISTS `lottery_periods` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `period_number` varchar(20) NOT NULL COMMENT '期数编号',
  `total_order_amount` decimal(15,2) NOT NULL DEFAULT '0.00' COMMENT '本期购买订单金额',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '期数状态: pending-待开始, active-进行中, closed-已结束',
  `lottery_result` varchar(50) DEFAULT NULL COMMENT '开奖结果',
  `order_start_time` timestamp NOT NULL COMMENT '订单开始时间',
  `order_end_time` timestamp NOT NULL COMMENT '订单结束时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_period_number` (`period_number`),
  KEY `idx_status_order_end_time` (`status`, `order_end_time`),
  KEY `idx_order_start_time_order_end_time` (`order_start_time`, `order_end_time`),
  KEY `idx_lottery_result` (`lottery_result`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='游戏期数表 - 记录每期的编号、订单金额、状态和时间信息';

-- ========================================
-- 3. 管理相关表
-- ========================================

-- 管理员用户表
CREATE TABLE IF NOT EXISTS `admin_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(8) NOT NULL COMMENT '管理员唯一ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password_hash` varchar(255) NOT NULL COMMENT '密码哈希',
  `role` varchar(20) NOT NULL DEFAULT 'admin' COMMENT '角色: super_admin-超级管理员, admin-管理员, operator-操作员',
  `status` varchar(20) NOT NULL DEFAULT 'active' COMMENT '状态: active-正常, inactive-禁用',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '上级管理员ID',
  `last_login_at` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uid` (`uid`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_role` (`role`),
  KEY `idx_status` (`status`),
  KEY `idx_role_status_deleted_at` (`role`, `status`, `deleted_at`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邀请码管理表 - 存储邀请码信息，用于用户注册时的邀请码校验，默认角色为业务员(4)';

-- ========================================
-- 4. 配置相关表
-- ========================================

-- 金额配置表
CREATE TABLE IF NOT EXISTS `amount_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `type` varchar(20) NOT NULL COMMENT '配置类型: recharge-充值, withdraw-提现',
  `amount` decimal(15,2) NOT NULL COMMENT '金额',
  `is_active` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否激活: 0-否, 1-是',
  `sort_order` int NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type_is_active_sort` (`type`, `is_active`, `sort_order`),
  KEY `idx_amount` (`amount`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='金额配置表 - 存储充值、提现等操作的金额配置，支持排序和激活状态管理';

-- 公告表
CREATE TABLE IF NOT EXISTS `announcements` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `title` varchar(200) NOT NULL COMMENT '公告标题',
  `content` text NOT NULL COMMENT '公告内容',
  `plain_content` text DEFAULT NULL COMMENT '纯文本内容',
  `rich_content` longtext DEFAULT NULL COMMENT '富文本内容',
  `tag` varchar(50) DEFAULT NULL COMMENT '标签',
  `status` varchar(20) NOT NULL DEFAULT 'active' COMMENT '状态: active-正常, inactive-禁用',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_tag` (`tag`),
  KEY `idx_status_deleted_at_created_at` (`status`, `deleted_at`, `created_at`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告表 - 存储系统公告信息，支持富文本内容，包括标题、纯文本内容、富文本内容、标签、状态等';

-- 公告图片表
CREATE TABLE IF NOT EXISTS `announcement_banners` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `announcement_id` bigint unsigned NOT NULL COMMENT '公告ID',
  `image_url` varchar(500) NOT NULL COMMENT '图片URL',
  `link_url` varchar(500) DEFAULT NULL COMMENT '跳转链接',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_announcement_id_deleted_at_sort` (`announcement_id`, `deleted_at`, `sort`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告图片表 - 存储公告相关的图片信息，支持排序和跳转链接';

-- 用户等级配置表
CREATE TABLE IF NOT EXISTS `member_level` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `level` int NOT NULL COMMENT '等级',
  `level_name` varchar(50) NOT NULL COMMENT '等级名称',
  `min_experience` int NOT NULL COMMENT '最小经验值',
  `max_experience` int NOT NULL COMMENT '最大经验值',
  `cashback_ratio` decimal(5,4) NOT NULL DEFAULT '0.0000' COMMENT '返现比例',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_level` (`level`),
  KEY `idx_level_deleted_at` (`level`, `deleted_at`),
  KEY `idx_cashback_ratio` (`cashback_ratio`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户等级配置表 - 存储用户等级配置信息，包括等级、经验值范围、返现比例等';

-- ========================================
-- 表统计信息更新
-- ========================================

-- 更新表统计信息
ANALYZE TABLE users;
ANALYZE TABLE wallets;
ANALYZE TABLE wallet_transactions;
ANALYZE TABLE user_login_logs;
ANALYZE TABLE orders;
ANALYZE TABLE group_buys;
ANALYZE TABLE lottery_periods;
ANALYZE TABLE admin_users;
ANALYZE TABLE amount_config;
ANALYZE TABLE announcements;
ANALYZE TABLE announcement_banners;
ANALYZE TABLE member_level;

-- ========================================
-- 显示创建结果
-- ========================================

-- 显示所有表
SELECT 
    table_name,
    table_comment,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)',
    table_rows
FROM information_schema.tables 
WHERE table_schema = 'gin_fataMorgana'
ORDER BY table_name;

-- 显示索引统计
SELECT 
    table_name,
    COUNT(*) as index_count,
    SUM(CASE WHEN index_name = 'PRIMARY' THEN 1 ELSE 0 END) as primary_keys,
    SUM(CASE WHEN non_unique = 0 AND index_name != 'PRIMARY' THEN 1 ELSE 0 END) as unique_keys,
    SUM(CASE WHEN non_unique = 1 THEN 1 ELSE 0 END) as regular_indexes
FROM information_schema.statistics 
WHERE table_schema = 'gin_fataMorgana'
GROUP BY table_name
ORDER BY table_name; 