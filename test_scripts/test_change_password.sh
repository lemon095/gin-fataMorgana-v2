#!/bin/bash

# 测试修改密码接口
echo "=== 测试修改密码接口 ==="

# 设置API基础URL
API_BASE="http://localhost:8080/api/v1"

# 测试用户登录
echo "1. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser1@example.com",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token: $TOKEN"

if [ -n "$TOKEN" ]; then
    # 测试修改密码 - 成功情况
    echo -e "\n2. 测试修改密码（成功）..."
    CHANGE_PASSWORD_RESPONSE=$(curl -s -X POST "$API_BASE/auth/change-password" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "old_password": "123456",
        "new_password": "newpassword123"
      }')

    echo "修改密码响应: $CHANGE_PASSWORD_RESPONSE"

    # 测试用新密码登录
    echo -e "\n3. 测试用新密码登录..."
    NEW_LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
      -H "Content-Type: application/json" \
      -d '{
        "email": "testuser1@example.com",
        "password": "newpassword123"
      }')

    echo "新密码登录响应: $NEW_LOGIN_RESPONSE"

    # 提取新token
    NEW_TOKEN=$(echo $NEW_LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo "新Token: $NEW_TOKEN"

    if [ -n "$NEW_TOKEN" ]; then
        # 测试修改回原密码
        echo -e "\n4. 测试修改回原密码..."
        REVERT_PASSWORD_RESPONSE=$(curl -s -X POST "$API_BASE/auth/change-password" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $NEW_TOKEN" \
          -d '{
            "old_password": "newpassword123",
            "new_password": "123456"
          }')

        echo "修改回原密码响应: $REVERT_PASSWORD_RESPONSE"
    fi

    # 测试错误情况
    echo -e "\n5. 测试当前密码错误..."
    WRONG_PASSWORD_RESPONSE=$(curl -s -X POST "$API_BASE/auth/change-password" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "old_password": "wrongpassword",
        "new_password": "newpassword123"
      }')

    echo "当前密码错误响应: $WRONG_PASSWORD_RESPONSE"

    # 测试新密码格式错误
    echo -e "\n6. 测试新密码格式错误（太短）..."
    SHORT_PASSWORD_RESPONSE=$(curl -s -X POST "$API_BASE/auth/change-password" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "old_password": "123456",
        "new_password": "123"
      }')

    echo "新密码太短响应: $SHORT_PASSWORD_RESPONSE"

    # 测试新密码与旧密码相同
    echo -e "\n7. 测试新密码与旧密码相同..."
    SAME_PASSWORD_RESPONSE=$(curl -s -X POST "$API_BASE/auth/change-password" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "old_password": "123456",
        "new_password": "123456"
      }')

    echo "新密码与旧密码相同响应: $SAME_PASSWORD_RESPONSE"

    # 测试参数错误
    echo -e "\n8. 测试参数错误（缺少新密码）..."
    MISSING_PARAM_RESPONSE=$(curl -s -X POST "$API_BASE/auth/change-password" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "old_password": "123456"
      }')

    echo "缺少参数响应: $MISSING_PARAM_RESPONSE"

else
    echo "跳过需要认证的测试（登录失败）"
fi

# 测试未认证访问
echo -e "\n9. 测试未认证访问..."
UNAUTH_RESPONSE=$(curl -s -X POST "$API_BASE/auth/change-password" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "123456",
    "new_password": "newpassword123"
  }')

echo "未认证访问响应: $UNAUTH_RESPONSE"

echo -e "\n=== 测试完成 ==="
echo "接口说明："
echo "POST /api/v1/auth/change-password - 修改密码"
echo "  - 需要认证：Authorization: Bearer {token}"
echo "  - 请求参数：old_password (当前密码), new_password (新密码)"
echo "  - 业务逻辑：验证当前密码，检查新密码格式，更新密码哈希"
echo "  - 返回：成功消息或错误信息" 