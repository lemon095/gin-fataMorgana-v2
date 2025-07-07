-- 创建公告Banner图表
CREATE TABLE IF NOT EXISTS `announcement_banners` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `announcement_id` bigint unsigned NOT NULL COMMENT '公告ID',
  `image_url` varchar(255) NOT NULL COMMENT '图片URL',
  `title` varchar(100) DEFAULT NULL COMMENT '图片标题',
  `link` varchar(255) DEFAULT NULL COMMENT '跳转链接',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_announcement_id` (`announcement_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_announcement_id_deleted_at_sort` (`announcement_id`, `deleted_at`, `sort`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告Banner图表'; 