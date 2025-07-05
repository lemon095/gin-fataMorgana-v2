-- 创建数据库
CREATE DATABASE IF NOT EXISTS gin_fataMorgana CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE gin_fataMorgana;

-- 创建管理员用户表（邀请码管理表）
CREATE TABLE IF NOT EXISTS admin_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    admin_id BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '管理员唯一ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码哈希',
    remark VARCHAR(500) COMMENT '备注',
    status BIGINT DEFAULT 1 COMMENT '账户状态 1:正常 0:禁用',
    avatar VARCHAR(255) COMMENT '头像URL',
    role BIGINT NOT NULL DEFAULT 4 COMMENT '身份角色 1:超级管理员 2:经理 3:主管 4:业务员（默认业务员）',
    my_invite_code VARCHAR(6) NOT NULL UNIQUE COMMENT '我的邀请码',
    parent_id INT UNSIGNED NULL COMMENT '上级用户ID',
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at TIMESTAMP(3) NULL COMMENT '软删除时间',
    INDEX idx_admin_users_admin_id (admin_id),
    INDEX idx_admin_users_username (username),
    INDEX idx_admin_users_my_invite_code (my_invite_code),
    INDEX idx_admin_users_parent_id (parent_id),
    INDEX idx_admin_users_status (status),
    INDEX idx_admin_users_role (role),
    INDEX idx_admin_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认管理员用户（密码：admin123）
INSERT INTO admin_users (admin_id, username, password, remark, status, role, my_invite_code) VALUES 
(1, 'admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '系统管理员', 1, 1, 'ADMIN1')
ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP(3);

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid VARCHAR(8) NOT NULL UNIQUE COMMENT '用户唯一ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱地址',
    password VARCHAR(255) NOT NULL COMMENT '密码哈希',
    phone VARCHAR(20) COMMENT '手机号',
    bank_card_info JSON COMMENT '银行卡信息JSON',
    experience INT DEFAULT 0 COMMENT '用户经验值',
    credit_score INT DEFAULT 100 COMMENT '用户信用分',
    status TINYINT DEFAULT 1 COMMENT '用户状态 1:正常 0:禁用',
    invited_by VARCHAR(6) COMMENT '注册时填写的邀请码',
    has_group_buy_qualification BOOLEAN DEFAULT FALSE COMMENT '是否有拼单资格',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_uid (uid),
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_phone (phone),
    INDEX idx_invited_by (invited_by),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建钱包表
CREATE TABLE IF NOT EXISTS wallets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid VARCHAR(8) NOT NULL UNIQUE COMMENT '用户唯一ID',
    balance DECIMAL(15,2) DEFAULT 0.00 COMMENT '总余额',
    frozen_balance DECIMAL(15,2) DEFAULT 0.00 COMMENT '冻结余额',
    total_income DECIMAL(15,2) DEFAULT 0.00 COMMENT '总收入',
    total_expense DECIMAL(15,2) DEFAULT 0.00 COMMENT '总支出',
    status TINYINT DEFAULT 1 COMMENT '钱包状态 1:正常 0:冻结',
    currency VARCHAR(3) DEFAULT 'PHP' COMMENT '货币类型',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_uid (uid),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建钱包交易记录表
CREATE TABLE IF NOT EXISTS wallet_transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    transaction_no VARCHAR(32) NOT NULL UNIQUE COMMENT '交易流水号',
    uid VARCHAR(8) NOT NULL COMMENT '用户唯一ID',
    type VARCHAR(20) NOT NULL COMMENT '交易类型',
    amount DECIMAL(15,2) NOT NULL COMMENT '交易金额',
    balance_before DECIMAL(15,2) NOT NULL COMMENT '交易前余额',
    balance_after DECIMAL(15,2) NOT NULL COMMENT '交易后余额',
    frozen_before DECIMAL(15,2) DEFAULT 0.00 COMMENT '交易前冻结余额',
    frozen_after DECIMAL(15,2) DEFAULT 0.00 COMMENT '交易后冻结余额',
    status VARCHAR(20) DEFAULT 'success' COMMENT '交易状态',
    description TEXT COMMENT '交易描述',
    operator_uid VARCHAR(8) COMMENT '操作员UID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_transaction_no (transaction_no),
    INDEX idx_uid (uid),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建用户登录日志表
CREATE TABLE IF NOT EXISTS user_login_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid VARCHAR(8) NOT NULL COMMENT '用户唯一ID',
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    login_ip VARCHAR(45) COMMENT '登录IP',
    user_agent TEXT COMMENT '用户代理',
    status TINYINT DEFAULT 1 COMMENT '登录状态 1:成功 0:失败',
    failure_reason VARCHAR(255) COMMENT '失败原因',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_uid (uid),
    INDEX idx_login_time (login_time),
    INDEX idx_status (status),
    INDEX idx_login_ip (login_ip)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建索引
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_wallets_uid ON wallets(uid);
CREATE INDEX idx_transactions_uid_created ON wallet_transactions(uid, created_at);
CREATE INDEX idx_login_logs_uid_time ON user_login_logs(uid, login_time); 