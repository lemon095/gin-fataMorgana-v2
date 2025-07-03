-- 初始化用户等级配置数据
-- 执行前请确保 member_level 表已创建

-- 清空现有数据（可选）
-- DELETE FROM member_level;

-- 插入等级配置数据
INSERT INTO member_level (level, name, min_experience, max_experience, cashback_ratio, status, created_at, updated_at) VALUES
(1, '青铜会员', 1, 99, 0.50, 1, NOW(), NOW()),
(2, '白银会员', 100, 299, 1.00, 1, NOW(), NOW()),
(3, '黄金会员', 300, 599, 1.50, 1, NOW(), NOW()),
(4, '铂金会员', 600, 999, 2.00, 1, NOW(), NOW()),
(5, '钻石会员', 1000, 1999, 2.50, 1, NOW(), NOW()),
(6, '皇冠会员', 2000, 4999, 3.00, 1, NOW(), NOW()),
(7, '至尊会员', 5000, 9999, 3.50, 1, NOW(), NOW()),
(8, '传奇会员', 10000, 99999, 4.00, 1, NOW(), NOW()),
(9, '神话会员', 100000, 999999, 4.50, 1, NOW(), NOW()),
(10, '永恒会员', 1000000, 9999999, 5.00, 1, NOW(), NOW());

-- 验证插入结果
SELECT * FROM member_level ORDER BY level ASC; 