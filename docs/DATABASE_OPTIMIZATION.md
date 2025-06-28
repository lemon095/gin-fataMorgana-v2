# 数据库优化配置指南

## 概述

本文档详细说明了项目中数据库优化的配置和最佳实践。

## 1. MySQL连接优化

### 1.1 连接字符串参数

```go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true&cachePrepStmts=true&prepStmtCacheSize=256&prepStmtCacheSqlLimit=2048&rejectReadOnly=true&timeout=10s&readTimeout=30s&writeTimeout=30s&multiStatements=true&autocommit=true")
```

**参数说明：**
- `interpolateParams=true`: 启用参数插值，减少预处理语句数量
- `cachePrepStmts=true`: 启用预处理语句缓存
- `prepStmtCacheSize=256`: 预处理语句缓存大小（256条）
- `prepStmtCacheSqlLimit=2048`: 预处理语句SQL长度限制（2KB）
- `rejectReadOnly=true`: 拒绝只读连接
- `timeout=10s`: 连接超时时间
- `readTimeout=30s`: 读取超时时间
- `writeTimeout=30s`: 写入超时时间
- `multiStatements=true`: 支持多语句执行
- `autocommit=true`: 自动提交事务

### 1.2 GORM配置优化

```go
DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
    DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
    PrepareStmt: true, // 启用预处理语句缓存
    SkipDefaultTransaction: true, // 跳过默认事务，提升性能
    DryRun: false, // 生产环境禁用DryRun
})
```

## 2. 连接池配置

### 2.1 推荐配置

```yaml
database:
  max_idle_conns: 10      # 最大空闲连接数
  max_open_conns: 100     # 最大打开连接数
  conn_max_lifetime: 3600 # 连接最大生存时间（秒）
  conn_max_idle_time: 1800 # 空闲连接最大生存时间（秒）
```

### 2.2 配置建议

- **max_open_conns**: 根据并发量调整，建议为CPU核心数的2-4倍
- **max_idle_conns**: 建议为max_open_conns的10-20%
- **conn_max_lifetime**: 建议1小时，避免连接长时间占用
- **conn_max_idle_time**: 建议30分钟，及时释放空闲连接

## 3. Redis缓存优化

### 3.1 缓存策略

```go
// 邮箱存在性检查缓存
CacheEmailExists(ctx, email, exists) // 缓存5分钟

// 邀请码存在性检查缓存
CacheInviteCodeExists(ctx, inviteCode, exists) // 缓存10分钟
```

### 3.2 缓存键设计

```
邮箱缓存: email_exists:user@example.com
邀请码缓存: invite_code_exists:ABC123
用户信息缓存: user_info:12345
```

## 4. 查询优化

### 4.1 批量查询优化

```go
// 批量检查邮箱是否存在
func (r *UserRepository) BatchCheckEmails(ctx context.Context, emails []string) (map[string]bool, error)

// 批量检查邀请码是否存在
func (r *UserRepository) BatchCheckInviteCodes(ctx context.Context, inviteCodes []string) (map[string]bool, error)
```

### 4.2 查询缓存中间件

```go
// 带缓存的查询
func (qc *QueryCache) WithCache(ctx context.Context, key string, expiration time.Duration, fn func() (interface{}, error)) (interface{}, error)

// 清除缓存
func (qc *QueryCache) InvalidateCache(ctx context.Context, pattern string) error
```

## 5. 性能监控

### 5.1 监控接口

```
GET /health/system          # 系统健康检查
GET /health/database        # 数据库健康检查
GET /health/db-stats        # 数据库统计信息
GET /health/query-stats     # 查询统计信息
GET /health/optimization    # 性能优化建议
```

### 5.2 监控指标

- **连接池状态**: 最大连接数、当前连接数、使用中连接数、空闲连接数
- **等待统计**: 等待连接数、等待时间
- **缓存命中率**: 预处理语句缓存命中率
- **连接利用率**: 当前连接使用率

## 6. 最佳实践

### 6.1 查询优化

1. **使用索引**: 为常用查询字段创建索引
2. **避免SELECT ***: 只查询需要的字段
3. **使用批量操作**: 减少数据库往返次数
4. **合理使用缓存**: 缓存热点数据
5. **使用预处理语句**: 提升查询性能

### 6.2 连接管理

1. **合理设置连接池参数**: 根据实际负载调整
2. **监控连接状态**: 定期检查连接池健康状态
3. **及时释放连接**: 避免连接泄漏
4. **使用连接超时**: 设置合理的超时时间

### 6.3 缓存策略

1. **缓存热点数据**: 缓存频繁访问的数据
2. **设置合理的TTL**: 根据数据更新频率设置缓存时间
3. **及时清除缓存**: 数据更新时清除相关缓存
4. **监控缓存命中率**: 优化缓存策略

## 7. 故障排查

### 7.1 常见问题

1. **连接池耗尽**: 检查max_open_conns设置
2. **查询超时**: 检查readTimeout和writeTimeout设置
3. **缓存未命中**: 检查缓存键设计和TTL设置
4. **性能下降**: 检查索引和查询语句优化

### 7.2 排查工具

```bash
# 查看数据库统计信息
curl http://localhost:8080/health/db-stats

# 查看查询统计信息
curl http://localhost:8080/health/query-stats

# 查看性能优化建议
curl http://localhost:8080/health/optimization
```

## 8. 配置示例

### 8.1 生产环境配置

```yaml
database:
  max_idle_conns: 20
  max_open_conns: 200
  conn_max_lifetime: 3600
  conn_max_idle_time: 1800

redis:
  pool_size: 20
  min_idle_conns: 10
  max_retries: 3
  dial_timeout: 5
  read_timeout: 3
  write_timeout: 3
```

### 8.2 开发环境配置

```yaml
database:
  max_idle_conns: 5
  max_open_conns: 50
  conn_max_lifetime: 3600
  conn_max_idle_time: 1800

redis:
  pool_size: 10
  min_idle_conns: 5
  max_retries: 3
  dial_timeout: 5
  read_timeout: 3
  write_timeout: 3
``` 