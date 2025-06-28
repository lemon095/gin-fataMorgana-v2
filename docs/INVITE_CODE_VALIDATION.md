# 邀请码校验功能说明

## 概述

本项目使用管理员用户表（`admin_users`）来管理邀请码，用户注册时需要提供有效的管理员邀请码才能成功注册。

## 功能特点

### 1. 邀请码校验
- 用户注册时必须提供有效的邀请码
- 邀请码必须来自活跃的管理员账户
- 支持邀请码的唯一性校验

### 2. 管理员用户表结构
```sql
CREATE TABLE admin_users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    admin_id VARCHAR(8) UNIQUE NOT NULL COMMENT '管理员唯一ID',
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码哈希',
    remark VARCHAR(500) COMMENT '备注',
    status INT DEFAULT 1 COMMENT '账户状态 1:正常 0:禁用',
    avatar VARCHAR(255) COMMENT '头像URL',
    role VARCHAR(20) NOT NULL DEFAULT '业务员' COMMENT '身份角色',
    my_invite_code VARCHAR(6) UNIQUE COMMENT '我的邀请码',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL COMMENT '软删除时间'
);
```

### 3. 支持的角色类型
- 超级管理员
- 经理
- 主管
- 业务员

## 使用流程

### 1. 初始化管理员用户
```bash
# 运行初始化脚本创建管理员用户
./init_admin.sh
```

### 2. 获取邀请码
```sql
-- 查询活跃管理员的邀请码
SELECT username, my_invite_code, role, status 
FROM admin_users 
WHERE status = 1;
```

### 3. 用户注册
用户注册时需要提供有效的邀请码：

```json
POST /auth/register
{
    "email": "user@example.com",
    "password": "password123",
    "confirm_password": "password123",
    "invite_code": "ABC123"
}
```

## 校验逻辑

### 1. 邀请码校验步骤
1. 检查邀请码是否为空
2. 查询管理员用户表，查找对应的邀请码
3. 验证管理员账户状态是否为活跃（status = 1）
4. 如果验证通过，允许用户注册

### 2. 错误处理
- 邀请码无效：`邀请码无效或管理员账户已被禁用`
- 管理员账户被禁用：`邀请码对应的管理员账户已被禁用`
- 邀请码为空：`邀请码不能为空`

## 数据库操作

### 1. 创建管理员用户
```sql
INSERT INTO admin_users (
    admin_id, 
    username, 
    password, 
    remark, 
    status, 
    role, 
    my_invite_code
) VALUES (
    '12345678',
    'admin_user',
    'hashed_password',
    '系统管理员',
    1,
    '超级管理员',
    'ABC123'
);
```

### 2. 禁用管理员账户
```sql
UPDATE admin_users 
SET status = 0 
WHERE username = 'admin_user';
```

### 3. 查询邀请码使用情况
```sql
-- 查询某个邀请码的使用次数
SELECT COUNT(*) as usage_count 
FROM users 
WHERE invited_by = 'ABC123';
```

## 测试

### 1. 运行邀请码校验测试
```bash
./test_invite_validation.sh
```

### 2. 手动测试步骤
1. 启动服务：`go run main.go`
2. 创建管理员用户：`./init_admin.sh`
3. 获取邀请码：查询数据库
4. 使用邀请码注册用户
5. 验证注册结果

## 注意事项

### 1. 安全性
- 邀请码应该定期更换
- 管理员账户密码应该使用强密码
- 定期检查管理员账户状态

### 2. 性能优化
- 邀请码查询已添加缓存
- 使用索引优化查询性能
- 支持批量邀请码校验

### 3. 扩展性
- 可以添加邀请码使用次数限制
- 可以添加邀请码有效期
- 可以添加邀请码权限等级

## 相关文件

- `models/admin_user.go` - 管理员用户模型
- `database/admin_user_repository.go` - 管理员用户仓库
- `services/user_service.go` - 用户服务（包含邀请码校验）
- `init_admin.sh` - 管理员用户初始化脚本
- `test_invite_validation.sh` - 邀请码校验测试脚本 