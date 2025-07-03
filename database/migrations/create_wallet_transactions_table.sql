-- 创建钱包交易流水表
CREATE TABLE IF NOT EXISTS `wallet_transactions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `transaction_no` varchar(32) NOT NULL COMMENT '交易流水号',
  `uid` varchar(8) NOT NULL COMMENT '用户唯一ID',
  `type` varchar(20) NOT NULL COMMENT '交易类型',
  `amount` decimal(15,2) NOT NULL COMMENT '交易金额',
  `balance_before` decimal(15,2) NOT NULL COMMENT '交易前余额',
  `balance_after` decimal(15,2) NOT NULL COMMENT '交易后余额',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '交易状态',
  `description` varchar(255) DEFAULT NULL COMMENT '交易描述',
  `related_order_no` varchar(32) DEFAULT NULL COMMENT '关联订单编号',
  `operator_uid` varchar(8) DEFAULT NULL COMMENT '操作员ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transaction_no` (`transaction_no`),
  KEY `idx_uid` (`uid`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_related_order_no` (`related_order_no`),
  KEY `idx_operator_uid` (`operator_uid`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='钱包交易流水表 - 记录所有钱包交易流水'; 