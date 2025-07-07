-- 创建用户等级配置表
CREATE TABLE IF NOT EXISTS `member_level` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `level` int NOT NULL COMMENT '等级数值',
  `name` varchar(50) NOT NULL COMMENT '等级名称',
  `logo` varchar(255) DEFAULT NULL COMMENT '等级logo',
  `cashback_ratio` decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '返现比例（百分比）',
  `single_amount` int DEFAULT 1 COMMENT '单数字额',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_level` (`level`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_level_deleted_at` (`level`, `deleted_at`),
  KEY `idx_cashback_ratio` (`cashback_ratio`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户等级配置表 - 存储用户等级配置信息,包括等级、名称、logo、返现比例、单数字额等'; 