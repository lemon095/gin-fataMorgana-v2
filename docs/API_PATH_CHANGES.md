# API路径变更说明

## 概述

为了更好的API管理和版本控制，所有业务接口现在都使用统一的 `/api/v1` 前缀。

## 变更详情

### 变更前
- 认证接口：`/auth/*`
- 会话接口：`/session/*`
- 钱包接口：`/wallet/*`
- 管理员接口：`/admin/*`
- 健康检查：`/health/*`

### 变更后
- 认证接口：`/api/v1/auth/*`
- 会话接口：`/api/v1/session/*`
- 钱包接口：`/api/v1/wallet/*`
- 管理员接口：`/api/v1/admin/*`
- 健康检查：`/api/v1/health/*`

## 保持不变的路径

以下路径保持不变，以确保向后兼容：
- 根路径：`/` - 欢迎页面
- 健康检查：`/health` - 基础健康检查（用于监控）
- Swagger文档：`/swagger/*` - API文档

## 详细的API路径映射

### 认证相关接口
| 功能 | 旧路径 | 新路径 | 方法 |
|------|--------|--------|------|
| 用户注册 | `/auth/register` | `/api/v1/auth/register` | POST |
| 用户登录 | `/auth/login` | `/api/v1/auth/login` | POST |
| 刷新令牌 | `/auth/refresh` | `/api/v1/auth/refresh` | POST |
| 用户登出 | `/auth/logout` | `/api/v1/auth/logout` | POST |
| 获取用户信息 | `/auth/profile` | `/api/v1/auth/profile` | GET |
| 绑定银行卡 | `/auth/bind-bank-card` | `/api/v1/auth/bind-bank-card` | POST |
| 获取银行卡信息 | `/auth/bank-card` | `/api/v1/auth/bank-card` | GET |

### 会话管理接口
| 功能 | 旧路径 | 新路径 | 方法 |
|------|--------|--------|------|
| 检查登录状态 | `/session/status` | `/api/v1/session/status` | GET |
| 获取当前用户信息 | `/session/user` | `/api/v1/session/user` | GET |
| 用户登出 | `/session/logout` | `/api/v1/session/logout` | POST |
| 刷新会话 | `/session/refresh` | `/api/v1/session/refresh` | POST |

### 钱包相关接口
| 功能 | 旧路径 | 新路径 | 方法 |
|------|--------|--------|------|
| 获取钱包信息 | `/wallet/info` | `/api/v1/wallet/info` | GET |
| 获取资金记录 | `/wallet/transactions` | `/api/v1/wallet/transactions` | GET |
| 申请提现 | `/wallet/withdraw` | `/api/v1/wallet/withdraw` | POST |
| 获取提现汇总 | `/wallet/withdraw-summary` | `/api/v1/wallet/withdraw-summary` | GET |
| 充值申请 | `/wallet/recharge-apply` | `/api/v1/wallet/recharge-apply` | POST |
| 充值确认 | `/wallet/recharge-confirm` | `/api/v1/wallet/recharge-confirm` | POST |

### 管理员接口
| 功能 | 旧路径 | 新路径 | 方法 |
|------|--------|--------|------|
| 确认提现 | `/admin/withdraw/confirm` | `/api/v1/admin/withdraw/confirm` | POST |
| 取消提现 | `/admin/withdraw/cancel` | `/api/v1/admin/withdraw/cancel` | POST |

### 健康检查接口
| 功能 | 旧路径 | 新路径 | 方法 |
|------|--------|--------|------|
| 系统健康检查 | `/health/check` | `/api/v1/health/check` | GET |
| 数据库健康检查 | `/health/database` | `/api/v1/health/database` | GET |
| Redis健康检查 | `/health/redis` | `/api/v1/health/redis` | GET |

## 迁移指南

### 前端应用
1. 更新所有API请求的URL，添加 `/api/v1` 前缀
2. 更新API配置文件中的基础URL
3. 测试所有接口功能

### 测试脚本
已更新的测试脚本：
- `test_scripts/init_super_admin.sh`
- `test_scripts/test_auth.sh`
- `test_scripts/test_new_api_paths.sh` (新增)

### Nginx配置
如果使用Nginx反向代理，需要更新配置：

```nginx
# 新的API路径配置
location /api/v1/ {
    proxy_pass http://gin_app;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    
    # 跨域配置
    add_header Access-Control-Allow-Origin *;
    add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS";
    add_header Access-Control-Allow-Headers "Origin, Content-Type, Accept, Authorization, X-Requested-With";
}

# 保持原有的健康检查路径
location /health {
    proxy_pass http://gin_app;
    # ... 其他配置
}
```

## 测试验证

运行以下命令测试新的API路径：

```bash
# 测试新的API路径结构
./test_scripts/test_new_api_paths.sh

# 测试认证功能
./test_scripts/test_auth.sh

# 测试超级管理员创建
./test_scripts/init_super_admin.sh
```

## 版本控制

- 当前版本：v1.0
- API版本：v1
- 路径前缀：`/api/v1`

未来如果需要API版本升级，可以添加 `/api/v2` 等新版本路径，同时保持旧版本兼容。

## 注意事项

1. **向后兼容**：基础的 `/health` 路径保持不变，确保监控系统正常工作
2. **Swagger文档**：已更新为新的路径结构
3. **测试脚本**：主要测试脚本已更新，其他脚本需要手动更新
4. **前端集成**：需要更新前端应用中的API调用路径

## 优势

1. **统一管理**：所有业务接口都有统一的前缀
2. **版本控制**：便于未来API版本管理
3. **安全性**：可以针对 `/api/*` 路径进行统一的安全配置
4. **监控友好**：便于日志分析和监控
5. **文档清晰**：API文档结构更加清晰 