-- 创建钱包表
CREATE TABLE IF NOT EXISTS `wallets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` varchar(8) NOT NULL COMMENT '用户唯一ID',
  `balance` decimal(15,2) NOT NULL DEFAULT 0.00 COMMENT '钱包余额',
  `status` int NOT NULL DEFAULT 1 COMMENT '钱包状态 1:正常 0:冻结',
  `currency` varchar(3) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
  `last_active_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='钱包表 - 存储用户钱包信息，包括余额等'; 