#!/bin/bash

# 测试期数列表接口
echo "=== 测试期数列表接口 ==="

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 测试获取期数列表
echo "1. 测试获取期数列表"
curl -X POST "${BASE_URL}/order/period/list" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{}' \
  | jq '.'

echo ""
echo "=== 测试完成 ===" 