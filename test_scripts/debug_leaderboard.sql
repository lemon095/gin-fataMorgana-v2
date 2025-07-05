-- 排行榜调试SQL脚本
-- 用于排查为什么排行榜只显示3条数据

-- 1. 检查当前时间和水单时间的周数差异
SELECT CURDATE() as current_date, 
       YEARWEEK(NOW(), 1) as current_week, 
       YEARWEEK('2025-07-05', 1) as water_order_week;

-- 2. 检查本周成功的水单数量
SELECT COUNT(*) as success_water_orders_this_week
FROM orders 
WHERE is_system_order = 1 
  AND status = 'success'
  AND YEARWEEK(updated_at, 1) = YEARWEEK(NOW(), 1);

-- 3. 检查水单的完成时间分布
SELECT DATE(updated_at) as completion_date,
       COUNT(*) as order_count
FROM orders 
WHERE is_system_order = 1 
  AND status = 'success'
GROUP BY DATE(updated_at)
ORDER BY completion_date DESC;

-- 4. 检查排行榜查询的具体条件（本周开始到结束）
SELECT o.uid, u.username, COUNT(*) as order_count, SUM(o.amount) as total_amount
FROM orders o
JOIN users u ON o.uid = u.uid
WHERE o.status = 'success' 
  AND o.updated_at >= DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY)
  AND o.updated_at < DATE_ADD(DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY), INTERVAL 7 DAY)
GROUP BY o.uid, u.username
ORDER BY order_count DESC, total_amount DESC
LIMIT 10;

-- 5. 检查本周时间范围
SELECT 
    DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY) as week_start,
    DATE_ADD(DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY), INTERVAL 7 DAY) as week_end;

-- 6. 检查水单完成时间是否在本周范围内
SELECT 
    id, order_no, uid, amount, status, updated_at,
    CASE 
        WHEN updated_at >= DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY) 
         AND updated_at < DATE_ADD(DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY), INTERVAL 7 DAY)
        THEN 'YES'
        ELSE 'NO'
    END as in_current_week
FROM orders 
WHERE is_system_order = 1 
  AND status = 'success'
ORDER BY updated_at DESC
LIMIT 10;

-- 7. 统计不同周的水单数量
SELECT 
    YEARWEEK(updated_at, 1) as week_number,
    COUNT(*) as order_count,
    SUM(amount) as total_amount
FROM orders 
WHERE is_system_order = 1 
  AND status = 'success'
GROUP BY YEARWEEK(updated_at, 1)
ORDER BY week_number DESC; 