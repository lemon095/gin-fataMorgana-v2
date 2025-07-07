# 服务版本管理策略

## 问题现状

项目中存在多个版本的服务实现，这导致了以下问题：
1. **代码冗余**：相同功能有多个实现
2. **维护困难**：需要同时维护多个版本
3. **使用混乱**：不同地方使用不同版本
4. **学习成本高**：新开发者需要理解多个版本

## 版本演进历史

### 钱包服务演进
```
V1 (wallet_service.go) 
    ↓ 添加缓存
V2 (wallet_cache_service.go) 
    ↓ 优化缓存策略
V3 (wallet_cache_service_v2.go) 
    ↓ 添加版本控制
V4 (wallet_cache_service_v3.go) 
    ↓ 简化实现
V5 (wallet_cache_service_v4.go) 
    ↓ 添加并发安全
V6 (wallet_concurrent_service.go) 
    ↓ 支持分布式
V7 (wallet_distributed_lock_service.go)
```

## 版本清理策略

### 1. 确定当前推荐版本

| 服务类型 | 推荐版本 | 原因 | 状态 |
|---------|---------|------|------|
| 钱包服务 | `WalletDistributedLockService` | 支持跨进程、统一Key管理 | ✅ 推荐 |
| 钱包缓存 | `WalletCacheServiceV4` | 简化实现、统一Key管理 | ✅ 推荐 |
| 订单服务 | `OrderService` | 基础功能完整 | ✅ 推荐 |
| 订单缓存 | `OrderCacheService` | 缓存策略合理 | ✅ 推荐 |

### 2. 版本迁移计划

#### 第一阶段：标记废弃版本
```go
// 在废弃版本文件开头添加
// DEPRECATED: 此版本已废弃，请使用 WalletDistributedLockService
// 迁移时间：2024年12月
// 删除时间：2025年3月
```

#### 第二阶段：逐步迁移
1. **更新依赖**：将所有使用旧版本的地方改为使用推荐版本
2. **测试验证**：确保功能正常
3. **性能对比**：验证新版本性能

#### 第三阶段：删除废弃版本
1. **确认无依赖**：确保没有代码使用废弃版本
2. **备份代码**：保留废弃版本到备份分支
3. **删除文件**：清理废弃版本文件

## 具体迁移步骤

### 1. 钱包服务迁移

#### 当前使用情况
```go
// 需要迁移的服务
- group_buy_service.go 使用 WalletCacheServiceV3
- wallet_concurrent_service.go 使用 WalletCacheServiceV3
- wallet_distributed_lock_service.go 使用 WalletCacheServiceV4
```

#### 迁移计划
```go
// 1. 统一使用 WalletDistributedLockService
// 2. 删除其他钱包服务版本
// 3. 保留 wallet_service.go 作为基础服务（如果还有用）
```

### 2. 缓存服务迁移

#### 当前使用情况
```go
// 钱包缓存服务
- wallet_cache_service.go (V1) - 废弃
- wallet_cache_service_v2.go (V2) - 废弃  
- wallet_cache_service_v3.go (V3) - 部分使用
- wallet_cache_service_v4.go (V4) - 推荐版本
```

#### 迁移计划
```go
// 1. 所有服务统一使用 WalletCacheServiceV4
// 2. 删除 V1、V2、V3 版本
// 3. 重命名 V4 为 wallet_cache_service.go
```

## 版本管理最佳实践

### 1. 版本命名规范
```go
// 推荐：使用描述性名称而不是版本号
type WalletService struct{}           // 当前版本
type WalletServiceLegacy struct{}     // 旧版本（如果需要保留）

// 避免：使用版本号
type WalletServiceV1 struct{}         // 不推荐
type WalletServiceV2 struct{}         // 不推荐
```

### 2. 接口设计
```go
// 定义统一接口
type WalletServiceInterface interface {
    GetWalletBalance(ctx context.Context, uid string) (*models.Wallet, error)
    WithdrawBalance(ctx context.Context, uid string, amount float64) error
    AddBalance(ctx context.Context, uid string, amount float64) error
    // ... 其他方法
}

// 实现接口
type WalletService struct {
    // 当前推荐实现
}

// 工厂方法
func NewWalletService() WalletServiceInterface {
    return &WalletService{}
}
```

### 3. 配置化选择
```go
// 通过配置选择服务实现
type Config struct {
    WalletServiceType string `yaml:"wallet_service_type"` // "distributed", "concurrent", "basic"
}

func NewWalletService(config Config) WalletServiceInterface {
    switch config.WalletServiceType {
    case "distributed":
        return NewWalletDistributedLockService()
    case "concurrent":
        return NewWalletConcurrentService()
    default:
        return NewWalletService()
    }
}
```

## 立即行动项

### 1. 标记废弃版本
```bash
# 在以下文件开头添加废弃标记
services/wallet_cache_service.go
services/wallet_cache_service_v2.go
services/wallet_cache_service_v3.go
services/wallet_concurrent_service.go
```

### 2. 更新依赖
```bash
# 更新以下文件使用推荐版本
services/group_buy_service.go
```

### 3. 重命名推荐版本
```bash
# 将推荐版本重命名为标准名称
mv services/wallet_distributed_lock_service.go services/wallet_service_new.go
mv services/wallet_cache_service_v4.go services/wallet_cache_service.go
```

### 4. 删除废弃版本
```bash
# 确认无依赖后删除
rm services/wallet_cache_service_v2.go
rm services/wallet_cache_service_v3.go
rm services/wallet_concurrent_service.go
```

## 长期维护策略

### 1. 版本控制
- 使用Git分支管理不同版本
- 主分支只保留当前推荐版本
- 废弃版本移到legacy分支

### 2. 文档管理
- 维护版本演进文档
- 记录每个版本的优缺点
- 提供迁移指南

### 3. 测试策略
- 为新版本编写完整测试
- 保留旧版本的测试用例
- 定期进行性能对比测试

## 总结

通过统一的版本管理策略，我们可以：
1. **减少代码冗余**：只保留必要的版本
2. **提高维护效率**：专注于推荐版本
3. **降低学习成本**：新开发者只需要学习当前版本
4. **提升代码质量**：持续优化推荐版本

建议立即开始版本清理工作，优先处理钱包相关服务，然后逐步扩展到其他服务。 