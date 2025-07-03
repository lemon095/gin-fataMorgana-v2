# 用户等级配置 API 文档

## 概述

用户等级配置系统允许根据用户经验值（experience）来确定用户等级，并提供相应的返现比例。新注册用户默认经验值为1，对应等级1。

## 数据库表结构

### member_level 表

| 字段名 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| id | uint | 主键ID | 自增 |
| level | int | 等级 | 必填 |
| name | string | 等级名称 | 必填 |
| min_experience | int | 最小经验值 | 0 |
| max_experience | int | 最大经验值 | 0 |
| cashback_ratio | decimal(5,2) | 返现比例（百分比） | 0.00 |
| status | int | 状态 1:启用 0:禁用 | 1 |
| created_at | time.Time | 创建时间 | 自动 |
| updated_at | time.Time | 更新时间 | 自动 |

## API 接口

### 1. 获取所有等级配置

**接口地址：** `GET /api/member-levels`

**请求参数：** 无

**响应示例：**
```json
{
  "code": 0,
  "message": "获取等级配置成功",
  "data": [
    {
      "id": 1,
      "level": 1,
      "name": "青铜会员",
      "min_experience": 1,
      "max_experience": 99,
      "cashback_ratio": 0.50,
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "timestamp": 1704067200
}
```

### 2. 根据等级获取配置

**接口地址：** `GET /api/member-levels/{level}`

**路径参数：**
- `level`: 等级（整数）

**响应示例：**
```json
{
  "code": 0,
  "message": "获取等级配置成功",
  "data": {
    "id": 1,
    "level": 1,
    "name": "青铜会员",
    "min_experience": 1,
    "max_experience": 99,
    "cashback_ratio": 0.50,
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "timestamp": 1704067200
}
```

### 3. 获取用户等级信息

**接口地址：** `GET /api/member-levels/user-info`

**查询参数：**
- `experience`: 用户经验值（整数）

**响应示例：**
```json
{
  "code": 0,
  "message": "获取用户等级信息成功",
  "data": {
    "current_level": {
      "id": 1,
      "level": 1,
      "name": "青铜会员",
      "min_experience": 1,
      "max_experience": 99,
      "cashback_ratio": 0.50,
      "status": 1
    },
    "next_level": {
      "id": 2,
      "level": 2,
      "name": "白银会员",
      "min_experience": 100,
      "max_experience": 299,
      "cashback_ratio": 1.00,
      "status": 1
    },
    "experience": 1
  },
  "timestamp": 1704067200
}
```

### 4. 计算返现金额

**接口地址：** `GET /api/member-levels/calculate-cashback`

**查询参数：**
- `experience`: 用户经验值（整数）
- `amount`: 金额（浮点数）

**响应示例：**
```json
{
  "code": 0,
  "message": "计算返现金额成功",
  "data": {
    "experience": 1,
    "amount": 100,
    "cashback_amount": 0.5
  },
  "timestamp": 1704067200
}
```

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1003 | 参数错误 |
| 1004 | 操作失败 |
| 404 | 资源不存在 |

## 使用示例

### 初始化等级配置数据

```sql
-- 执行 test_scripts/init_member_levels.sql
INSERT INTO member_level (level, name, min_experience, max_experience, cashback_ratio, status) VALUES
(1, '青铜会员', 1, 99, 0.50, 1),
(2, '白银会员', 100, 299, 1.00, 1),
(3, '黄金会员', 300, 599, 1.50, 1);
```

### 测试API接口

```bash
# 执行测试脚本
chmod +x test_scripts/test_member_level.sh
./test_scripts/test_member_level.sh
```

## 业务逻辑

1. **用户注册**：新用户注册时，`experience` 字段默认设置为 1
2. **等级判定**：根据用户经验值在 `min_experience` 和 `max_experience` 范围内确定等级
3. **返现计算**：返现金额 = 交易金额 × (返现比例 / 100)
4. **等级升级**：当用户经验值达到下一等级的 `min_experience` 时，可以升级

## 注意事项

1. 确保 `member_level` 表已创建并包含相应的等级配置数据
2. 等级配置的 `min_experience` 和 `max_experience` 不能重叠
3. 返现比例以百分比形式存储，计算时需要除以100
4. 只有 `status = 1` 的等级配置才会被查询到 