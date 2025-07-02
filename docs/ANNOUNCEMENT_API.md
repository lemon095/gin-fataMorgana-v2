# 公告API文档

## 接口概述

公告接口用于获取系统公告信息，支持分页查询和图片展示。**只返回已发布状态(status=1且is_publish=true)的公告**。

**缓存策略**: 公告列表接口支持Redis缓存，缓存时间为1分钟，相同参数的请求会直接返回缓存数据。

## 1. 获取公告列表

### 接口信息
- **接口路径**: `POST /api/v1/announcements/list`
- **请求方法**: POST
- **认证要求**: 无需认证
- **说明**: 获取已发布的公告列表，按创建时间倒序排列，支持缓存（1分钟）

### 请求参数
```json
{
  "page": 1,
  "page_size": 10
}
```

#### 参数说明
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码，从1开始 |
| page_size | int | 否 | 10 | 每页大小，最大100 |

### 返回示例
```json
{
  "code": 0,
  "message": "获取公告列表成功",
  "data": {
    "announcements": [
      {
        "id": 1,
        "title": "系统维护通知",
        "content": "系统将于今晚22:00-24:00进行维护升级，期间可能影响部分功能使用。",
        "tag": "系统通知",
        "status": 1,
        "is_publish": true,
        "link": "https://example.com/maintenance",
        "created_at": "2025-01-01T10:00:00Z",
        "banners": [
          "https://example.com/images/banner1.jpg",
          "https://example.com/images/banner2.jpg"
        ]
      }
    ],
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total": 25,
      "total_pages": 3,
      "has_next": true,
      "has_prev": false
    }
  },
  "timestamp": 1751365370
}
```

#### 数据结构说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint | 公告ID |
| title | string | 公告标题 |
| content | string | 公告内容 |
| tag | string | 公告标签 |
| status | int64 | 公告状态 1:已发布 0:草稿 |
| is_publish | bool | 是否发布 |
| link | string | 文章链接 |
| created_at | string | 创建时间 |
| banners | array | 图片列表 |

#### 图片数据结构
| 字段名 | 类型 | 说明 |
|--------|------|------|
| banners | array | 图片URL数组，按sort字段升序排列 |

#### 分页信息
| 字段名 | 类型 | 说明 |
|--------|------|------|
| current_page | int | 当前页码 |
| page_size | int | 每页大小 |
| total | int64 | 总记录数 |
| total_pages | int | 总页数 |
| has_next | bool | 是否有下一页 |
| has_prev | bool | 是否有上一页 |

### 错误响应
```json
{
  "code": 400,
  "message": "请求参数错误: page must be greater than 0",
  "data": null,
  "timestamp": 1751365370
}
```

## 注意事项

1. **发布状态**: 只返回已发布状态(status=1且is_publish=true)的公告
2. **排序规则**: 返回结果按创建时间倒序排列，最新的公告在最前面
3. **分页限制**: 每页大小最大为100条
4. **图片排序**: 公告图片按sort字段升序排列，直接返回图片URL数组
5. **缓存策略**: 相同参数的请求会返回缓存数据，缓存时间为1分钟
6. **错误处理**: 请根据返回的code字段判断请求是否成功

## 使用示例

### 获取第一页公告
```bash
curl -X POST "http://localhost:9001/api/v1/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 10
  }'
```

### 获取第二页公告
```bash
curl -X POST "http://localhost:9001/api/v1/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 2,
    "page_size": 5
  }'
```

### 使用默认参数
```bash
curl -X POST "http://localhost:9001/api/v1/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{}'
``` 