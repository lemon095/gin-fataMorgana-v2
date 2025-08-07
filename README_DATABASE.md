# 数据库和Redis配置说明

## 概述

本项目已集成MySQL数据库和Redis缓存，支持连接池配置和自动表迁移。

## 配置文件

### 主配置文件
- `config/config.yaml` - 主配置文件
- `config/config.example.yaml` - 示例配置文件

### 配置项说明

#### 服务器配置
```yaml
server:
  host: "0.0.0.0"      # 服务器监听地址
  port: 9002           # 服务器端口
domain: "localhost:9002"  # 域名
  mode: "debug"        # 运行模式：debug/release
```

#### 数据库配置
```yaml
database:
  driver: "mysql"           # 数据库驱动
  host: "localhost"         # 数据库主机
  port: 3306               # 数据库端口
  username: "root"         # 数据库用户名
  password: "123456"       # 数据库密码
  dbname: "gin_fata_morgana"  # 数据库名
  charset: "utf8mb4"       # 字符集
  parse_time: true         # 解析时间
  loc: "Local"             # 时区
  # 连接池配置
  max_idle_conns: 10       # 最大空闲连接数
  max_open_conns: 100      # 最大打开连接数
  conn_max_lifetime: 3600  # 连接最大生命周期（秒）
```

#### Redis配置
```yaml
redis:
  host: "localhost"         # Redis主机
  port: 6379               # Redis端口
  password: ""             # Redis密码
  db: 0                    # Redis数据库编号
  # 连接池配置
  pool_size: 10            # 连接池大小
  min_idle_conns: 5        # 最小空闲连接数
  max_retries: 3           # 最大重试次数
  dial_timeout: 5          # 连接超时（秒）
  read_timeout: 3          # 读取超时（秒）
  write_timeout: 3         # 写入超时（秒）
```

#### JWT配置
```yaml
jwt:
  secret: "your-secret-key-here"  # JWT密钥
  access_token_expire: 3600       # 访问令牌过期时间（秒）
  refresh_token_expire: 604800    # 刷新令牌过期时间（秒）
```

#### 认证配置
```yaml
auth:
  session_timeout: 3600    # 会话超时时间（秒）
  max_login_attempts: 5    # 最大登录尝试次数
  lockout_duration: 1800   # 锁定持续时间（秒）
```

#### 日志配置
```yaml
log:
  level: "info"            # 日志级别：debug/info/warn/error
  format: "json"           # 日志格式：json/text
  output: "stdout"         # 输出方式：stdout/file
  file:
    path: "logs/app.log"   # 日志文件路径
    max_size: 100          # 最大文件大小（MB）
    max_age: 30            # 最大保留天数
    max_backups: 10        # 最大备份文件数
    compress: true         # 是否压缩
```

#### 雪花算法配置
```yaml
snowflake:
  worker_id: 1             # 工作节点ID
  datacenter_id: 1         # 数据中心ID
```

## 数据库初始化

### 1. 创建数据库
```sql
CREATE DATABASE gin_fata_morgana CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 配置数据库连接
编辑 `config/config.yaml` 文件，修改数据库配置：
```yaml
database:
  host: "localhost"
  port: 3306
  username: "your_username"
  password: "your_password"
  dbname: "gin_fata_morgana"
```

### 3. 自动迁移
项目启动时会自动创建和迁移数据库表结构，包括：
- `users` 表 - 用户信息表

## Redis初始化

### 1. 安装Redis
```bash
# macOS
brew install redis

# Ubuntu/Debian
sudo apt-get install redis-server

# CentOS/RHEL
sudo yum install redis
```

### 2. 启动Redis服务
```bash
# macOS
brew services start redis

# Linux
sudo systemctl start redis
```

### 3. 配置Redis连接
编辑 `config/config.yaml` 文件，修改Redis配置：
```yaml
redis:
  host: "localhost"
  port: 6379
  password: "your_redis_password"  # 如果设置了密码
```

## 项目结构

```
gin-fataMorgana/
├── config/
│   ├── config.go           # 配置管理
│   ├── config.yaml         # 主配置文件
│   └── config.example.yaml # 示例配置文件
├── database/
│   ├── mysql.go            # MySQL连接和迁移
│   └── redis.go            # Redis连接和操作
├── models/
│   └── user.go             # 用户模型（已更新为GORM）
├── services/
│   └── user_service.go     # 用户服务（已更新为使用数据库）
├── utils/
│   ├── jwt.go              # JWT工具（已更新为使用配置）
│   └── snowflake.go        # 雪花算法（已更新为使用配置）
└── main.go                 # 主程序（已集成数据库和Redis）
```

## 使用说明

### 1. 启动项目
```bash
go run main.go
```

### 2. 检查连接状态
项目启动时会显示：
- MySQL连接状态
- Redis连接状态
- 数据库迁移状态

### 3. 数据库操作
- 用户注册：数据自动保存到MySQL
- 用户登录：从MySQL验证用户信息
- 会话管理：使用Redis存储会话信息

### 4. Redis操作
项目提供了以下Redis操作函数：
```go
// 设置键值对
database.SetKey(ctx, key, value, expiration)

// 获取键值
database.GetKey(ctx, key)

// 删除键
database.DelKey(ctx, key)

// 检查键是否存在
database.ExistsKey(ctx, key)

// 设置过期时间
database.SetExpire(ctx, key, expiration)
```

## 注意事项

1. **生产环境配置**
   - 修改JWT密钥
   - 设置强密码
   - 配置Redis密码
   - 调整连接池参数

2. **数据库备份**
   - 定期备份MySQL数据
   - 备份Redis数据（如需要）

3. **性能优化**
   - 根据实际负载调整连接池大小
   - 监控数据库和Redis性能
   - 配置适当的超时时间

4. **安全考虑**
   - 不要将敏感配置提交到版本控制
   - 使用环境变量覆盖敏感配置
   - 定期更新依赖包

## 故障排除

### 数据库连接失败
1. 检查数据库服务是否启动
2. 验证连接参数是否正确
3. 确认数据库用户权限

### Redis连接失败
1. 检查Redis服务是否启动
2. 验证连接参数是否正确
3. 确认Redis密码设置

### 表迁移失败
1. 检查数据库权限
2. 确认数据库字符集设置
3. 查看错误日志 