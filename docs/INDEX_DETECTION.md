# 数据库索引自动检测功能

## 概述

本项目实现了智能的数据库索引检测和创建功能，确保数据库性能优化索引始终存在。

## 迁移执行顺序

### 自动迁移流程
项目启动时的数据库迁移按照以下顺序执行：

1. **第一步：创建/更新表结构**
   - 使用GORM AutoMigrate创建表
   - 自动创建主键、唯一索引、普通索引等基础索引
   - 更新表结构（新增字段、修改字段类型等）

2. **第二步：添加表注释**
   - 为所有表添加中文注释
   - 说明表的用途和字段含义

3. **第三步：检测和创建优化索引**
   - 检测现有索引，避免重复创建
   - 创建复合索引和性能优化索引
   - 提供详细的创建统计信息

4. **第四步：添加特殊表注释**
   - 为特定表添加额外的注释信息

### 执行顺序的重要性

**为什么必须先迁移表再创建索引？**

1. **表必须先存在**：索引只能创建在已存在的表上
2. **GORM基础索引**：AutoMigrate会自动创建模型定义的基础索引
3. **优化索引补充**：我们创建的索引是对基础索引的补充和优化
4. **避免冲突**：确保不会与GORM创建的基础索引冲突

### 启动日志示例

```
🚀 开始数据库迁移...
📋 第一步：创建/更新表结构...
✅ 表结构迁移完成
📝 第二步：添加表注释...
✅ 表注释添加完成
🔍 第三步：检测和创建优化索引...
索引 idx_users_uid_status_deleted_at 已存在，跳过创建
✅ 索引创建成功: idx_orders_uid_status_created_at (uid, status, created_at)
📊 索引检测完成 - 创建: 15, 跳过: 8, 失败: 0
✅ 索引检测和创建完成
📝 第四步：添加特殊表注释...
✅ 特殊表注释添加完成
🎉 数据库迁移全部完成！
```

## 功能特性

### 1. 自动索引检测
- 项目启动时自动检测缺失的索引
- 只创建不存在的索引，避免重复创建
- 详细的检测日志和统计信息

### 2. 手动索引管理
- 支持手动检测和创建索引
- 支持查看当前数据库所有索引
- 提供命令行工具和Makefile命令

### 3. 智能索引创建
- 包含所有核心表的优化索引
- 复合索引支持多字段查询优化
- 软删除字段索引支持

## 支持的索引类型

### 用户相关索引
- `users` 表：用户ID、状态、软删除、创建时间等复合索引
- `user_login_logs` 表：登录时间、状态、IP地址等索引

### 钱包相关索引
- `wallets` 表：状态、余额、活跃时间等索引
- `wallet_transactions` 表：用户ID、类型、状态、时间等复合索引

### 订单相关索引
- `orders` 表：用户ID、状态、期号、过期时间等复合索引
- `group_buys` 表：截止时间、状态、类型等复合索引

### 系统配置索引
- `admin_users` 表：角色、状态、软删除等索引
- `amount_config` 表：类型、激活状态、排序等索引
- `announcements` 表：状态、标签、软删除等索引
- `member_level` 表：等级、软删除、返现比例等索引

## 使用方法

### 1. 自动检测（项目启动时）
项目启动时会自动执行索引检测和创建：

```bash
go run main.go
```

启动日志会显示：
```
🔍 开始检测数据库索引...
索引 idx_users_uid_status_deleted_at 已存在，跳过创建
✅ 索引创建成功: idx_orders_uid_status_created_at (uid, status, created_at)
📊 索引检测完成 - 创建: 15, 跳过: 8, 失败: 0
```

### 2. 手动检测索引
```bash
# 使用Makefile命令
make db-check-index

# 或直接使用Go命令
go run cmd/migrate/main.go -check-index
```

### 3. 查看所有索引
```bash
# 使用Makefile命令
make db-show-index

# 或直接使用Go命令
go run cmd/migrate/main.go -show-index
```

### 4. 完整数据库迁移
```bash
# 使用Makefile命令
make db-migrate

# 或直接使用Go命令
go run cmd/migrate/main.go
```

## 索引检测原理

### 1. 检测机制
使用 `information_schema.statistics` 表查询当前数据库的索引信息：

```sql
SELECT COUNT(*) 
FROM information_schema.statistics 
WHERE table_schema = DATABASE() 
AND table_name = ? 
AND index_name = ?
```

### 2. 创建策略
- 只创建不存在的索引
- 单个索引创建失败不影响其他索引
- 提供详细的创建日志和统计

### 3. 索引定义
每个索引包含以下信息：
- 表名（TableName）
- 索引名（IndexName）
- 字段列表（Columns）
- 创建SQL（SQL）

## 索引列表

### users 表索引
| 索引名 | 字段 | 用途 |
|--------|------|------|
| idx_users_uid_status_deleted_at | uid, status, deleted_at | 用户查询优化 |
| idx_users_username_deleted_at | username, deleted_at | 用户名查询 |
| idx_users_email_deleted_at | email, deleted_at | 邮箱查询 |
| idx_users_phone_deleted_at | phone, deleted_at | 手机号查询 |
| idx_users_invited_by_deleted_at | invited_by, deleted_at | 邀请码查询 |
| idx_users_status_deleted_at | status, deleted_at | 状态查询 |
| idx_users_created_at | created_at | 创建时间查询 |

### orders 表索引
| 索引名 | 字段 | 用途 |
|--------|------|------|
| idx_orders_uid_status_created_at | uid, status, created_at | 用户订单查询 |
| idx_orders_status_updated_at | status, updated_at | 状态更新查询 |
| idx_orders_period_number_status | period_number, status | 期号状态查询 |
| idx_orders_expire_time_status | expire_time, status | 过期时间查询 |
| idx_orders_auditor_uid_status | auditor_uid, status | 审核员查询 |
| idx_orders_is_system_order_status | is_system_order, status | 系统订单查询 |

### wallet_transactions 表索引
| 索引名 | 字段 | 用途 |
|--------|------|------|
| idx_wallet_transactions_uid_type_status | uid, type, status | 用户交易查询 |
| idx_wallet_transactions_uid_created_at | uid, created_at | 用户交易时间查询 |
| idx_wallet_transactions_type_status_created_at | type, status, created_at | 交易类型查询 |
| idx_wallet_transactions_transaction_no | transaction_no | 交易号查询 |
| idx_wallet_transactions_amount | amount | 金额查询 |

## 性能优化建议

### 1. 定期监控
- 定期运行索引检测确保索引完整性
- 监控查询性能，根据实际使用情况调整索引

### 2. 索引维护
- 定期分析索引使用情况
- 删除不必要的索引减少维护开销

### 3. 查询优化
- 充分利用现有索引优化查询
- 避免在索引字段上使用函数

## 故障排除

### 1. 索引创建失败
- 检查数据库连接权限
- 确认表结构是否正确
- 查看详细错误日志

### 2. 性能问题
- 检查索引是否被正确使用
- 分析慢查询日志
- 考虑添加缺失的索引

### 3. 存储空间
- 监控索引占用的存储空间
- 定期清理不必要的索引

## 测试

运行索引检测测试：

```bash
./test_scripts/test_index_detection.sh
```

测试包括：
- 显示当前索引
- 检测并创建缺失索引
- 验证创建结果

## 总结

自动索引检测功能确保了：
- 数据库性能优化索引的完整性
- 减少手动维护索引的工作量
- 提高查询性能和系统稳定性
- 支持灵活的索引管理策略

## 迁移执行顺序

### 自动迁移流程
项目启动时的数据库迁移按照以下顺序执行：

1. **第一步：创建/更新表结构**
   - 使用GORM AutoMigrate创建表
   - 自动创建主键、唯一索引、普通索引等基础索引
   - 更新表结构（新增字段、修改字段类型等）

2. **第二步：添加表注释**
   - 为所有表添加中文注释
   - 说明表的用途和字段含义

3. **第三步：检测和创建优化索引**
   - 检测现有索引，避免重复创建
   - 创建复合索引和性能优化索引
   - 提供详细的创建统计信息

4. **第四步：添加特殊表注释**
   - 为特定表添加额外的注释信息

### 执行顺序的重要性

**为什么必须先迁移表再创建索引？**

1. **表必须先存在**：索引只能创建在已存在的表上
2. **GORM基础索引**：AutoMigrate会自动创建模型定义的基础索引
3. **优化索引补充**：我们创建的索引是对基础索引的补充和优化
4. **避免冲突**：确保不会与GORM创建的基础索引冲突

### 启动日志示例

```
🚀 开始数据库迁移...
📋 第一步：创建/更新表结构...
✅ 表结构迁移完成
📝 第二步：添加表注释...
✅ 表注释添加完成
🔍 第三步：检测和创建优化索引...
索引 idx_users_uid_status_deleted_at 已存在，跳过创建
✅ 索引创建成功: idx_orders_uid_status_created_at (uid, status, created_at)
📊 索引检测完成 - 创建: 15, 跳过: 8, 失败: 0
✅ 索引检测和创建完成
📝 第四步：添加特殊表注释...
✅ 特殊表注释添加完成
🎉 数据库迁移全部完成！
``` 