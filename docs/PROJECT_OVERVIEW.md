# Gin-FataMorgana 项目概览

## 项目简介

Gin-FataMorgana 是一个基于 Gin 框架的 Go Web 服务，提供用户认证、钱包管理、订单管理、邀请码系统等功能。项目采用分层架构，包含控制器层、服务层、数据访问层和模型层。

## 数据库表结构

### 1. 用户表 (users)
**作用**: 存储用户基本信息、认证信息、银行卡信息等

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 主键ID |
| uid | string(8) | 用户唯一ID |
| username | string(50) | 用户名 |
| email | string(100) | 邮箱地址 |
| password | string(255) | 密码哈希 |
| phone | string(20) | 手机号 |
| bank_card_info | json | 银行卡信息JSON |
| experience | int | 用户经验值 |
| credit_score | int | 用户信用分 |
| status | int | 用户状态 1:正常 0:禁用 |
| invited_by | string(6) | 注册时填写的邀请码 |
| has_group_buy_qualification | bool | 是否有拼单资格 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除时间 |

### 2. 钱包表 (wallets)
**作用**: 存储用户钱包信息，包括余额、冻结余额、总收入、总支出等

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 主键ID |
| uid | string(8) | 用户唯一ID |
| balance | decimal(15,2) | 钱包余额 |
| frozen_balance | decimal(15,2) | 冻结余额 |
| total_income | decimal(15,2) | 总收入 |
| total_expense | decimal(15,2) | 总支出 |
| status | int | 钱包状态 1:正常 0:冻结 |
| currency | string(3) | 货币类型 |
| last_active_at | datetime | 最后活跃时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 3. 钱包交易流水表 (wallet_transactions)
**作用**: 记录所有钱包交易明细，包括充值、提现、收入、支出、冻结、解冻等操作

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 主键ID |
| transaction_no | string(32) | 交易流水号 |
| uid | string(8) | 用户唯一ID |
| type | string(20) | 交易类型 |
| amount | decimal(15,2) | 交易金额 |
| balance_before | decimal(15,2) | 交易前余额 |
| balance_after | decimal(15,2) | 交易后余额 |
| frozen_before | decimal(15,2) | 交易前冻结余额 |
| frozen_after | decimal(15,2) | 交易后冻结余额 |
| status | string(20) | 交易状态 |
| description | string(200) | 交易描述 |
| remark | string(500) | 备注信息 |
| related_order_no | string(32) | 关联订单号 |

| operator_uid | string(8) | 操作员ID |
| ip_address | string(45) | 操作IP地址 |
| user_agent | string(500) | 用户代理 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

**交易类型**:
- recharge: 充值
- withdraw: 提现
- income: 收入（返现、奖励等）
- expense: 支出（消费、服务费等）
- freeze: 冻结
- unfreeze: 解冻
- refund: 退款
- transfer: 转账
- adjustment: 调整

**交易状态**:
- pending: 待处理
- success: 成功
- failed: 失败
- cancelled: 已取消

### 4. 订单表 (orders)
**作用**: 记录用户订单信息，包括买入金额、利润金额等

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 主键ID |
| order_no | string(32) | 订单编号 |
| uid | string(8) | 用户唯一ID |
| buy_amount | decimal(15,2) | 买入金额 |
| profit_amount | decimal(15,2) | 利润金额 |
| status | string(20) | 订单状态 |
| description | string(200) | 订单描述 |
| remark | string(500) | 备注信息 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

**订单状态**:
- pending: 待处理
- success: 成功
- failed: 失败
- cancelled: 已取消

### 5. 邀请码管理表 (admin_users)
**作用**: 存储邀请码信息，用于用户注册时的邀请码校验

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 主键ID |
| admin_id | uint | 管理员唯一ID |
| username | string(50) | 用户名 |
| password | string(255) | 密码哈希 |
| remark | string(500) | 备注 |
| status | int64 | 账户状态 1:正常 0:禁用 |
| avatar | string(255) | 头像URL |
| role | int64 | 身份角色 |
| my_invite_code | string(6) | 我的邀请码 |
| parent_id | uint | 上级用户ID |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除时间 |

**角色类型**:
- 1: 超级管理员
- 2: 经理
- 3: 主管
- 4: 业务员

### 6. 用户登录日志表 (user_login_logs)
**作用**: 记录用户登录历史，包括登录时间、IP地址、设备信息、登录状态等

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 主键ID |
| uid | string(8) | 用户UID |
| username | string(50) | 用户名 |
| email | string(100) | 邮箱 |
| login_ip | string(45) | 登录IP地址 |
| user_agent | string(500) | 用户代理 |
| login_time | datetime | 登录时间 |
| status | int | 登录状态 1:成功 0:失败 |
| fail_reason | string(200) | 失败原因 |
| device_info | string(200) | 设备信息 |
| location | string(100) | 登录地点 |
| created_at | datetime | 创建时间 |

### 7. 金额配置表 (amount_config)
**作用**: 存储充值、提现等操作的金额配置

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int64 | 主键ID |
| type | string(20) | 配置类型 recharge/withdraw |
| amount | decimal(10,2) | 金额 |
| description | string(100) | 描述 |
| is_active | bool | 是否激活 |
| sort_order | int | 排序 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

## API接口列表

### 1. 系统健康检查
- `GET /health` - 系统健康检查
- `GET /api/v2/health/check` - 系统健康检查
- `GET /api/v2/health/database` - 数据库健康检查
- `GET /api/v2/health/redis` - Redis健康检查

### 2. 认证相关接口
- `POST /api/v2/auth/register` - 用户注册
- `POST /api/v2/auth/login` - 用户登录
- `POST /api/v2/auth/refresh` - 刷新令牌
- `POST /api/v2/auth/logout` - 用户登出
- `POST /api/v2/auth/profile` - 获取用户信息
- `POST /api/v2/auth/bind-bank-card` - 绑定银行卡
- `POST /api/v2/auth/bank-card` - 获取银行卡信息

### 3. 会话管理接口
- `POST /api/v2/session/status` - 检查登录状态
- `POST /api/v2/session/user` - 获取当前用户信息
- `POST /api/v2/session/logout` - 用户登出
- `POST /api/v2/session/refresh` - 刷新会话

### 4. 钱包相关接口
- `POST /api/v2/wallet/info` - 获取钱包信息
- `POST /api/v2/wallet/transactions` - 获取资金记录
- `POST /api/v2/wallet/withdraw` - 申请提现
- `POST /api/v2/wallet/withdraw-summary` - 获取提现汇总
- `POST /api/v2/wallet/recharge-apply` - 充值申请
- `POST /api/v2/wallet/recharge-confirm` - 充值确认

### 5. 管理员接口
- `POST /api/v2/admin/withdraw/confirm` - 确认提现
- `POST /api/v2/admin/withdraw/cancel` - 取消提现

### 6. 订单相关接口
- `POST /api/v2/order/list` - 获取订单列表
- `POST /api/v2/order/create` - 创建订单
- `POST /api/v2/order/detail` - 获取订单详情
- `POST /api/v2/order/stats` - 获取订单统计
- `POST /api/v2/order/by-status` - 根据状态获取订单
- `POST /api/v2/order/by-date` - 根据日期范围获取订单

### 7. 热榜相关接口
- `POST /api/v2/leaderboard/ranking` - 获取任务热榜

### 8. 金额配置相关接口
- `POST /api/v2/amount-config/list` - 根据类型获取金额配置列表
- `GET /api/v2/amount-config/:id` - 根据ID获取金额配置详情

### 9. 假数据接口
- `POST /api/v2/fake/activities` - 获取假数据实时动态

## 核心功能模块

### 1. 用户认证模块
- 用户注册（需要邀请码）
- 用户登录（邮箱+密码）
- JWT令牌管理
- 会话管理
- 银行卡绑定

### 2. 钱包管理模块
- 钱包余额管理
- 充值功能（双阶段：申请+确认）
- 提现功能（申请+管理员确认）
- 交易流水记录
- 余额冻结/解冻

### 3. 订单管理模块
- 订单创建
- 订单状态管理
- 订单查询和统计
- 订单详情查看

### 4. 邀请码系统
- 邀请码校验
- 角色权限管理
- 层级关系管理

### 5. 金额配置管理
- 充值金额配置
- 提现金额配置
- 配置激活状态管理

### 6. 热榜系统
- 用户排行榜
- 个人数据统计

### 7. 日志记录
- 用户登录日志
- 操作审计

## 技术特性

### 1. 安全特性
- JWT认证
- 密码加密存储
- 敏感信息脱敏
- 接口限流防刷
- 幂等性校验

### 2. 数据库特性
- 连接池管理
- 事务处理
- 软删除
- 索引优化

### 3. 中间件
- CORS跨域处理
- 认证中间件
- 限流中间件
- 会话中间件
- 日志中间件

### 4. 文档
- Swagger API文档
- 自动生成文档
- 接口测试脚本

## 部署信息

- **服务端口**: 9002
- **数据库**: MySQL
- **缓存**: Redis
- **文档地址**: `/swagger/*`
- **健康检查**: `/health`

## 开发工具

- **编译**: `go build -o main main.go`
- **运行**: `./main`
- **测试**: `make test`
- **迁移**: `make migrate`
- **文档生成**: `swag init` 