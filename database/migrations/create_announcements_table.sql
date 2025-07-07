-- 创建公告表
CREATE TABLE IF NOT EXISTS `announcements` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `title` varchar(100) NOT NULL COMMENT '标题',
  `content` text NOT NULL COMMENT '纯文本内容（用于摘要、搜索等）',
  `rich_content` longtext DEFAULT NULL COMMENT '富文本内容（HTML格式）',
  `tag` varchar(20) NOT NULL COMMENT '标签',
  `status` int NOT NULL DEFAULT 0 COMMENT '状态 0-草稿 1-已发布',
  `is_publish` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否发布',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_status` (`status`),
  KEY `idx_tag` (`tag`),
  KEY `idx_status_deleted_at_created_at` (`status`, `deleted_at`, `created_at`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告表 - 存储系统公告信息，支持富文本内容，包括标题、纯文本内容、富文本内容、标签、状态等'; 