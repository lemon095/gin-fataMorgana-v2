#!/bin/bash

# 测试新公告表结构的接口
echo "=== 测试新公告表结构的接口 ==="

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 测试获取公告列表
echo "1. 测试获取公告列表"
curl -X POST "${BASE_URL}/announcements/list" \
  -H "Content-Type: application/json" \
  -d '{
    "page": 1,
    "page_size": 10
  }' \
  | jq '.'

echo ""
echo "=== 测试完成 ==="

# 预期返回结构示例：
# {
#   "code": 200,
#   "message": "success",
#   "data": {
#     "announcements": [
#       {
#         "id": 1,
#         "title": "系统公告",
#         "content": "这是纯文本内容",
#         "rich_content": "<p>这是<strong>富文本</strong>内容</p>",
#         "tag": "系统",
#         "status": 1,
#         "is_publish": true,
#         "created_at": "2024-12-01T10:00:00Z",
#         "banners": [
#           "https://example.com/banner1.jpg",
#           "https://example.com/banner2.jpg"
#         ]
#       }
#     ],
#     "pagination": {
#       "current_page": 1,
#       "page_size": 10,
#       "total": 1,
#       "total_pages": 1,
#       "has_next": false,
#       "has_prev": false
#     }
#   }
# } 