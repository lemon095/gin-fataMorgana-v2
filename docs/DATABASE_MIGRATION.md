# 数据库迁移文档

## 概述

本文档记录了项目中所有数据库表的自动迁移配置和手动迁移脚本。

## 自动迁移配置

### 已配置自动迁移的表

在 `database/mysql.go` 的 `AutoMigrate()` 函数中，以下表已配置自动迁移：

1. **users** - 用户表
2. **wallets** - 钱包表
3. **wallet_transactions** - 钱包交易流水表
4. **admin_users** - 管理员用户表
5. **user_login_logs** - 用户登录日志表
6. **orders** - 订单表
7. **amount_config** - 金额配置表
8. **announcements** - 公告表
9. **announcement_banners** - 公告图片表
10. **group_buys** - 拼单表
11. **member_level** - 用户等级配置表 ⭐ **新增**
12. **lottery_periods** - 期数开奖表 ⭐ **新增**

### 表注释配置

所有表都配置了相应的注释，在 `addTableComments()` 函数中定义。

## 新增表详情

### 1. member_level (用户等级配置表)

**用途：** 存储用户等级配置信息，包括等级、经验值范围、返现比例等

**主要字段：**
- `level` - 等级
- `name` - 等级名称
- `min_experience` - 最小经验值
- `max_experience` - 最大经验值
- `cashback_ratio` - 返现比例（百分比）
- `status` - 状态（1:启用 0:禁用）

**默认数据：**
- 青铜会员 (1-99经验值, 0.5%返现)
- 白银会员 (100-299经验值, 1.0%返现)
- 黄金会员 (300-599经验值, 1.5%返现)
- 铂金会员 (600-999经验值, 2.0%返现)
- 钻石会员 (1000-1999经验值, 2.5%返现)
- 皇冠会员 (2000-4999经验值, 3.0%返现)
- 至尊会员 (5000-9999经验值, 3.5%返现)
- 传奇会员 (10000-99999经验值, 4.0%返现)
- 神话会员 (100000-999999经验值, 4.5%返现)
- 永恒会员 (1000000-9999999经验值, 5.0%返现)

### 2. lottery_periods (期数开奖表)

**用途：** 记录每期开奖信息和订单统计

**主要字段：**
- `period_number` - 期数
- `lottery_result` - 开奖结果存储字符串
- `total_order_amount` - 本期购买订单金额
- `is_drawn` - 是否开奖
- `is_paid` - 是否发放
- `winning_order_count` - 发奖订单数
- `order_start_time` - 订单开始时间
- `order_end_time` - 订单结束时间

## 手动迁移脚本

### SQL迁移文件

- `database/migrations/create_lottery_periods_table.sql` - 创建期数开奖表和用户等级配置表

### 执行手动迁移

```bash
# 执行SQL迁移文件
mysql -u root -p gin_fataMorgana < database/migrations/create_lottery_periods_table.sql
```

## 测试脚本

### 迁移测试

```bash
# 测试数据库迁移
chmod +x test_scripts/test_migration.sh
./test_scripts/test_migration.sh
```

### 利润计算测试

```bash
# 测试利润计算功能
chmod +x test_scripts/test_profit_calculation.sh
./test_scripts/test_profit_calculation.sh
```

## 注意事项

1. **自动迁移优先级：** GORM的自动迁移会优先执行，如果表已存在且结构不匹配，可能会报错。

2. **手动迁移：** 如果自动迁移失败，可以使用手动SQL脚本创建表。

3. **数据备份：** 在生产环境中执行迁移前，请务必备份数据库。

4. **索引优化：** 所有表都配置了必要的索引，确保查询性能。

5. **字符集：** 所有表使用 `utf8mb4` 字符集和 `utf8mb4_unicode_ci` 排序规则。

## 验证迁移结果

### 检查表是否存在

```sql
-- 检查新表
SHOW TABLES LIKE 'member_level';
SHOW TABLES LIKE 'lottery_periods';

-- 检查表结构
DESCRIBE member_level;
DESCRIBE lottery_periods;

-- 检查索引
SHOW INDEX FROM member_level;
SHOW INDEX FROM lottery_periods;
```

### 检查默认数据

```sql
-- 检查等级配置数据
SELECT * FROM member_level ORDER BY level;
```

## 相关文件

- `models/member_level.go` - 用户等级配置模型
- `models/lottery_period.go` - 期数开奖模型
- `database/member_level_repository.go` - 用户等级配置Repository
- `database/lottery_period_repository.go` - 期数开奖Repository
- `database/mysql.go` - 数据库连接和迁移配置
- `services/order_service.go` - 订单服务（包含利润计算）
- `services/group_buy_service.go` - 拼单服务（包含利润计算） 