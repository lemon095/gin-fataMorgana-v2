-- 公告状态修复脚本
-- 使用方法: mysql -u root -p123456 gin_fata_morgana < fix_announcement_status.sql

-- 1. 查看当前公告状态
SELECT '=== 当前公告状态 ===' as info;
SELECT id, title, status, is_publish, created_at FROM announcements ORDER BY created_at DESC;

-- 2. 统计各状态公告数量
SELECT '=== 公告状态统计 ===' as info;
SELECT 
    COUNT(*) as total_count,
    SUM(CASE WHEN status = 1 AND is_publish = 1 THEN 1 ELSE 0 END) as published_count,
    SUM(CASE WHEN status = 0 OR is_publish = 0 THEN 1 ELSE 0 END) as draft_count
FROM announcements;

-- 3. 更新所有公告为已发布状态
UPDATE announcements SET is_publish = 1, status = 1;

-- 4. 查看更新后的状态
SELECT '=== 更新后的公告状态 ===' as info;
SELECT id, title, status, is_publish, created_at FROM announcements ORDER BY created_at DESC;

-- 5. 验证更新结果
SELECT '=== 更新结果验证 ===' as info;
SELECT 
    COUNT(*) as total_count,
    SUM(CASE WHEN status = 1 AND is_publish = 1 THEN 1 ELSE 0 END) as published_count
FROM announcements; 