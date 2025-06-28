# Gin-FataMorgana

一个基于Gin框架的Go Web服务项目，支持用户认证、钱包管理、银行卡绑定等功能。

## 🚀 功能特性

- 🚀 基于Gin框架的高性能Web服务
- 📡 RESTful API设计
- 🔐 JWT认证和授权
- 👤 用户注册、登录、登出
- 🔄 Token自动刷新
- 📱 邮箱注册和登录
- 🔒 密码确认验证
- 🤖 自动生成用户名
- 🆔 雪花算法生成八位数用户ID
- 💳 银行卡信息管理和验证
- 💰 钱包管理和交易记录
- 📊 用户经验值和信用分系统
- 📝 请求日志记录
- 🛡️ 错误恢复机制
- 🐳 Docker容器化部署
- 🔧 一键部署脚本

## 🚀 一键部署

### 前置要求

- Docker 20.10+
- Docker Compose 2.0+
- 4GB+ 可用内存

### 快速开始

1. **克隆项目**
```bash
git clone <repository-url>
cd gin-fataMorgana
```

2. **一键部署**
```bash
# 开发环境部署
./deploy.sh dev

# 生产环境部署
./deploy.sh prod
```

3. **使用Makefile（推荐）**
```bash
# 查看所有可用命令
make help

# 快速启动开发环境
make dev

# 快速启动生产环境
make prod

# 查看服务状态
make status

# 查看日志
make logs

# 停止服务
make stop
```

### 部署后访问

- **应用地址**: http://localhost:8080
- **Nginx地址**: http://localhost:80
- **健康检查**: http://localhost:8080/health

### 默认账户

- **管理员邮箱**: admin@example.com
- **管理员密码**: admin123
- **邀请码**: ADMIN1

## 🔧 手动部署

### 前置要求

- Go 1.21 或更高版本
- MySQL 8.0+
- Redis 7.0+

### 安装依赖

```bash
go mod tidy
```

### 配置数据库

1. 创建MySQL数据库
2. 复制配置文件
```bash
cp config/config.example.yaml config/config.yaml
```

3. 修改配置文件中的数据库连接信息

### 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

### 构建可执行文件

```bash
go build -o gin-fataMorgana main.go
```

然后运行：
```bash
./gin-fataMorgana
```

## 📋 API接口

### 基础接口

- `GET /` - 首页，返回服务状态
- `GET /health` - 健康检查
- `GET /health/check` - 系统健康检查
- `GET /health/database` - 数据库健康检查
- `GET /health/redis` - Redis健康检查

### 认证接口

- `POST /auth/register` - 用户注册
- `POST /auth/login` - 用户登录
- `POST /auth/refresh` - 刷新访问令牌
- `POST /auth/logout` - 用户登出
- `GET /auth/profile` - 获取用户信息
- `POST /auth/bind-bank-card` - 绑定银行卡
- `GET /auth/bank-card` - 获取银行卡信息

### 会话管理

- `GET /session/status` - 检查登录状态
- `GET /session/user` - 获取当前用户信息
- `POST /session/logout` - 用户登出
- `POST /session/refresh` - 刷新会话

### 钱包接口

- `GET /wallet/info` - 获取钱包信息
- `GET /wallet/transactions` - 获取资金记录
- `POST /wallet/withdraw` - 申请提现
- `GET /wallet/withdraw-summary` - 获取提现汇总

### 管理员接口

- `POST /admin/withdraw/confirm` - 确认提现
- `POST /admin/withdraw/cancel` - 取消提现

## 🗄️ 数据库结构

项目包含以下数据表：

- `users` - 用户表
- `wallets` - 钱包表
- `wallet_transactions` - 钱包交易记录表
- `user_login_logs` - 用户登录日志表
- `admin_users` - 管理员用户表

详细的数据库设计请参考 [README_DATABASE.md](README_DATABASE.md)

## 🐳 Docker部署

### 使用Docker Compose

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 单独构建镜像

```bash
# 构建应用镜像
docker build -t gin-fataMorgana .

# 运行容器
docker run -d -p 8080:8080 --name gin-app gin-fataMorgana
```

## 🔧 开发工具

### 测试脚本

项目包含多个测试脚本：

```bash
# 认证测试
./test_auth.sh

# 银行卡测试
./test_bank_card.sh

# 钱包测试
./test_wallet.sh

# 性能测试
./test_performance.sh

# Bug修复测试
./test_bug_fixes.sh
```

### 数据库管理

```bash
# 初始化管理员账户
./init_admin.sh

# 数据库备份
make backup

# 数据库恢复
make restore file=backups/backup_20240101_120000.sql
```

## 📊 监控和管理

### 健康检查

```bash
# 应用健康检查
curl http://localhost:8080/health

# 数据库健康检查
curl http://localhost:8080/health/database

# Redis健康检查
curl http://localhost:8080/health/redis
```

### 日志管理

```bash
# 查看应用日志
docker-compose logs -f app

# 查看数据库日志
docker-compose logs -f mysql

# 查看Redis日志
docker-compose logs -f redis
```

### 性能监控

```bash
# 查看数据库统计
curl http://localhost:8080/health/stats

# 查看查询统计
curl http://localhost:8080/health/query-stats

# 性能优化建议
curl http://localhost:8080/health/optimization
```

## 🛡️ 安全特性

- JWT令牌认证
- 密码加密存储
- 银行卡信息验证
- 请求频率限制
- SQL注入防护
- XSS防护
- CSRF防护

## 📝 项目结构

```
gin-fataMorgana/
├── main.go                    # 主程序文件
├── go.mod                     # Go模块文件
├── go.sum                     # 依赖校验文件
├── Dockerfile                 # Docker镜像构建文件
├── docker-compose.yml         # Docker Compose配置
├── deploy.sh                  # 一键部署脚本
├── Makefile                   # 项目管理工具
├── .dockerignore              # Docker忽略文件
├── config/                    # 配置文件目录
│   ├── config.go             # 配置结构定义
│   ├── config.example.yaml   # 配置示例文件
│   └── config.yaml           # 实际配置文件
├── models/                    # 数据模型
│   ├── user.go               # 用户模型
│   ├── wallet.go             # 钱包模型
│   └── wallet_transaction.go # 交易模型
├── controllers/               # 控制器
│   ├── auth_controller.go    # 认证控制器
│   ├── wallet_controller.go  # 钱包控制器
│   └── health_controller.go  # 健康检查控制器
├── services/                  # 业务服务
│   ├── user_service.go       # 用户服务
│   └── wallet_service.go     # 钱包服务
├── database/                  # 数据库相关
│   ├── mysql.go              # MySQL连接
│   ├── redis.go              # Redis连接
│   └── repository.go         # 数据仓库
├── middleware/                # 中间件
│   ├── auth.go               # 认证中间件
│   └── session.go            # 会话中间件
├── utils/                     # 工具函数
│   ├── jwt.go                # JWT工具
│   ├── snowflake.go          # 雪花算法
│   └── bank_card_validator.go # 银行卡验证
├── docker/                    # Docker相关文件
│   ├── mysql/                # MySQL配置
│   └── nginx/                # Nginx配置
├── docs/                      # 文档
├── logs/                      # 日志文件
└── test_*.sh                  # 测试脚本
```

## 🔄 更新和升级

### 更新代码

```bash
# 拉取最新代码
git pull

# 重新部署
./deploy.sh prod
```

### 数据备份

```bash
# 备份数据库
make backup

# 备份配置文件
cp config/config.yaml config/config.yaml.backup
```

## 🐛 故障排除

### 常见问题

1. **端口冲突**
   ```bash
   # 检查端口占用
   lsof -i :8080
   
   # 修改配置文件中的端口
   ```

2. **数据库连接失败**
   ```bash
   # 检查数据库状态
   docker-compose logs mysql
   
   # 检查网络连接
   docker-compose exec app ping mysql
   ```

3. **Redis连接失败**
   ```bash
   # 检查Redis状态
   docker-compose logs redis
   
   # 测试Redis连接
   docker-compose exec redis redis-cli ping
   ```

### 日志分析

```bash
# 查看应用错误日志
docker-compose logs app | grep ERROR

# 查看数据库慢查询
docker-compose exec mysql tail -f /var/log/mysql/slow.log
```

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📞 支持

如有问题，请查看：
- [Bug分析文档](docs/BUG_ANALYSIS.md)
- [数据库设计文档](README_DATABASE.md)
- [银行卡API文档](docs/BANK_CARD_API.md) 