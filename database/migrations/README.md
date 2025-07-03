# 数据库迁移文件说明

本目录包含所有数据库表的完整建表SQL文件，每个表对应一个独立的SQL文件。

## 文件列表

### 1. 用户相关表
- `create_users_table.sql` - 用户表（users）
  - 存储用户基本信息、认证信息、银行卡信息、经验值、信用分等
  - 包含软删除功能

- `create_wallets_table.sql` - 钱包表（wallets）
  - 存储用户钱包信息，包括余额、状态等
  - 每个用户只能有一个钱包

- `create_wallet_transactions_table.sql` - 钱包交易流水表（wallet_transactions）
  - 记录所有钱包交易流水
  - 包括充值、提现、订单购买、拼单参与等交易类型

### 2. 订单相关表
- `create_orders_table.sql` - 订单表（orders）
  - 记录用户订单信息
  - 包含任务数量、状态、利润金额等
  - 支持期号关联（period_number字段）

- `create_group_buys_table.sql` - 拼单表（group_buys）
  - 记录拼单信息，包括参与人数、付款金额、截止时间等
  - 支持多种拼单类型（normal/flash/vip）

### 3. 系统配置表
- `create_member_level_table.sql` - 用户等级配置表（member_level）
  - 存储用户等级配置信息
  - 包括等级、名称、logo、返现比例、单数字额等
  - 包含软删除功能

- `create_amount_config_table.sql` - 金额配置表（amount_config）
  - 存储充值提现金额配置
  - 支持激活状态和排序

- `create_admin_users_table.sql` - 管理员用户表（admin_users）
  - 存储邀请码信息，用于用户注册时的邀请码校验
  - 支持多级角色管理
  - 包含软删除功能

### 4. 内容管理表
- `create_announcements_table.sql` - 公告表（announcements）
  - 存储系统公告信息，支持富文本内容
  - 包含公告图片表（announcement_banners）
  - 支持发布状态和软删除

### 5. 日志记录表
- `create_user_login_logs_table.sql` - 用户登录日志表（user_login_logs）
  - 记录用户登录历史
  - 包括登录时间、IP地址、设备信息、登录状态等

### 6. 业务功能表
- `create_lottery_periods_table.sql` - 期数开奖表（lottery_periods）
  - 记录每期开奖信息和订单统计
  - 支持开奖状态和奖金发放状态

## 使用说明

1. **按顺序执行**：建议按照文件名的字母顺序执行，确保外键依赖正确
2. **字符集**：所有表使用 `utf8mb4` 字符集和 `utf8mb4_unicode_ci` 排序规则
3. **存储引擎**：所有表使用 `InnoDB` 存储引擎
4. **索引优化**：每个表都包含了必要的索引，包括主键、唯一键和普通索引

## 注意事项

- 所有表都包含 `created_at` 和 `updated_at` 时间戳字段
- 支持软删除的表包含 `deleted_at` 字段
- 外键关系通过业务逻辑维护，未使用数据库外键约束
- 所有金额字段使用 `decimal(15,2)` 类型确保精度
- 状态字段使用整数或字符串枚举，便于扩展

## 数据初始化

某些表（如 `create_group_buys_table.sql`）包含了示例数据，可以直接用于测试。
生产环境部署时请根据实际需要修改或删除示例数据。 