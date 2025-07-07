# 数据库迁移文件

## 文件说明

- `01_create_all_tables.sql` - 完整的数据库初始化脚本
  - 包含所有表的创建语句
  - 包含所有必要的索引优化
  - 包含表统计信息更新
  - 包含创建结果查询

## 表结构概览

### 1. 用户相关表
- `users` - 用户表（用户基本信息、认证信息、银行卡信息等）
- `wallets` - 钱包表（用户钱包信息，包括余额、冻结余额等）
- `wallet_transactions` - 钱包交易流水表（所有钱包交易明细）
- `user_login_logs` - 用户登录日志表（用户登录历史）

### 2. 订单相关表
- `orders` - 订单表（用户订单信息，包括订单金额、状态等）
- `group_buys` - 拼单表（拼单信息，支持多人拼单功能）
- `lottery_periods` - 游戏期数表（每期的编号、订单金额、状态等）

### 3. 管理相关表
- `admin_users` - 管理员用户表（管理员信息，用于系统管理）

### 4. 配置相关表
- `amount_config` - 金额配置表（充值、提现等操作的金额配置）
- `announcements` - 公告表（系统公告信息）
- `announcement_banners` - 公告图片表（公告相关的图片信息）
- `member_level` - 用户等级配置表（用户等级配置信息）

## 使用说明

1. **一键初始化**：执行 `01_create_all_tables.sql` 即可完成所有表的创建和索引优化
2. **字符集**：所有表使用 `utf8mb4` 字符集和 `utf8mb4_unicode_ci` 排序规则
3. **存储引擎**：所有表使用 `InnoDB` 存储引擎
4. **索引优化**：每个表都包含了必要的索引，包括主键、唯一键和普通索引

## 注意事项

- 所有表都包含 `created_at` 和 `updated_at` 时间戳字段
- 支持软删除的表包含 `deleted_at` 字段
- 外键关系通过业务逻辑维护，未使用数据库外键约束
- 所有金额字段使用 `decimal(15,2)` 类型确保精度
- 状态字段使用字符串枚举，便于扩展

## 执行方式

```bash
# 方式1：直接执行SQL文件
mysql -u root -p gin_fataMorgana < database/migrations/01_create_all_tables.sql

# 方式2：使用项目迁移工具
go run cmd/migrate/main.go
``` 