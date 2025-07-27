# Gin-FataMorgana

一个基于Gin框架的简化Go Web服务，提供用户认证、钱包管理、健康监控等功能。

## 🚀 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 安装依赖
```bash
go mod tidy
```

### 配置
复制配置文件并修改：
```bash
cp config/config.example.yaml config/config.yaml
```

### 启动服务
```bash
# 方式1：直接运行
go run main.go

# 方式2：使用Makefile
make run

# 方式3：使用部署脚本
./deploy.sh
```

## 🗄️ 数据库管理

### 自动迁移
项目支持自动数据库迁移，会在启动时自动创建和更新表结构：

```bash
# 手动执行迁移
make db-migrate

# 或者使用迁移工具
go run cmd/migrate/main.go
```

### 初始化数据
```bash
# 初始化管理员账户和邀请码
make db-seed
```

### 迁移测试
```bash
# 测试迁移功能
./test_migration.sh
```

### 数据库表结构
项目包含以下核心表：

| 表名 | 说明 | 主要字段 |
|------|------|----------|
| `users` | 用户表 | uid, username, email, password, bank_card_info, status |
| `wallets` | 钱包表 | uid, balance, frozen_balance, total_income, total_expense |
| `wallet_transactions` | 交易流水表 | transaction_no, uid, type, amount, status |
| `admin_users` | 邀请码管理表 | admin_id, username, my_invite_code, role, status |
| `user_login_logs` | 登录日志表 | uid, login_time, login_ip, status |

## 📁 项目结构

```
gin-fataMorgana/
├── config/                 # 配置管理
│   ├── config.go          # 配置结构定义
│   └── config.example.yaml # 配置文件示例
├── controllers/           # 控制器层
│   ├── auth_controller.go     # 认证控制器
│   ├── health_controller.go   # 健康检查控制器
│   ├── session_controller.go  # 会话控制器
│   └── wallet_controller.go   # 钱包控制器
├── database/             # 数据库层
│   ├── mysql.go          # MySQL连接
│   ├── redis.go          # Redis连接
│   └── repository.go     # 数据访问层
├── middleware/           # 中间件
│   ├── auth.go           # 认证中间件
│   └── session.go        # 会话中间件
├── models/               # 数据模型
│   ├── user.go           # 用户模型
│   ├── wallet.go         # 钱包模型
│   └── wallet_transaction.go # 交易模型
├── services/             # 业务逻辑层
│   ├── user_service.go   # 用户服务
│   └── wallet_service.go # 钱包服务
├── utils/                # 工具函数
│   ├── jwt.go            # JWT工具
│   ├── response.go       # 响应工具
│   └── validator.go      # 验证工具
├── main.go               # 主程序
├── go.mod               # Go模块文件
└── README.md            # 项目说明
```

## 🔧 核心功能

### 用户认证
- 用户注册/登录
- JWT令牌管理
- 会话管理
- 银行卡绑定

### 钱包管理
- 钱包创建
- 余额查询
- 充值/提现
- 交易记录

### 系统监控
- 健康检查
- 数据库状态
- Redis状态

## 📡 API接口

### 认证接口
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新令牌
- `POST /api/v1/auth/profile` - 获取用户信息

### 钱包接口
- `GET /api/v1/wallet/info` - 获取钱包信息
- `GET /api/v1/wallet/transactions` - 获取交易记录
- `POST /api/v1/wallet/withdraw` - 申请提现

### 健康检查
- `GET /health` - 系统健康检查
- `GET /api/v1/health/check` - 系统健康检查
- `GET /api/v1/health/database` - 数据库健康检查
- `GET /api/v1/health/redis` - Redis健康检查

## 🛠️ 部署

### Docker部署
```bash
# 构建镜像
docker build -t gin-fataMorgana .

# 启动服务
docker-compose up -d
```

### 手动部署
```bash
# 编译
go build -o gin-fataMorgana main.go

# 运行
./gin-fataMorgana
```

## 🔍 简化特性

本项目经过简化优化，主要特点：

1. **简化配置** - 只保留核心配置项
2. **简化验证** - 银行卡验证只保留基本Luhn算法
3. **简化错误处理** - 统一的错误码和响应格式
4. **简化中间件** - 合并重复功能
5. **简化模型** - 移除业务逻辑，只保留数据结构

## �� 许可证

MIT License 