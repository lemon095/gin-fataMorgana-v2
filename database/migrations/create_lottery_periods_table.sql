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
  KEY `idx_status` (`status`),
  KEY `idx_order_start_time` (`order_start_time`),
  KEY `idx_order_end_time` (`order_end_time`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_order_start_time_order_end_time` (`order_start_time`, `order_end_time`),
  KEY `idx_status_order_start_time` (`status`, `order_start_time`),
  KEY `idx_total_order_amount` (`total_order_amount`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='游戏期数表'; 