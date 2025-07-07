-- 创建金额配置表
CREATE TABLE IF NOT EXISTS `amount_config` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `type` varchar(20) NOT NULL COMMENT '配置类型: recharge-充值, withdraw-提现',
  `amount` decimal(10,2) NOT NULL COMMENT '金额',
  `description` varchar(100) DEFAULT NULL COMMENT '描述',
  `is_active` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否激活',
  `sort_order` int NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_is_active` (`is_active`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_type_is_active_sort` (`type`, `is_active`, `sort_order`),
  KEY `idx_amount` (`amount`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='金额配置表 - 存储充值提现金额配置'; 