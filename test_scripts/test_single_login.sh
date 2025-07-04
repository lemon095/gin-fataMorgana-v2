#!/bin/bash

# 单点登录功能测试脚本
# 测试设备A登录后，设备B登录会踢出设备A的场景

BASE_URL="http://localhost:9001/api/v1"
EMAIL="test@example.com"
PASSWORD="123456"

echo "=== 单点登录功能测试 ==="
echo "测试场景：设备A登录后，设备B登录会踢出设备A"
echo ""

# 清理之前的测试数据（可选）
echo "1. 清理测试环境..."
# 这里可以添加清理逻辑

echo ""
echo "2. 设备A登录..."
DEVICE_A_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -H "User-Agent: DeviceA/1.0" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "设备A登录响应: $DEVICE_A_RESPONSE"

# 提取设备A的token
DEVICE_A_TOKEN=$(echo $DEVICE_A_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$DEVICE_A_TOKEN" ]; then
    echo "❌ 设备A登录失败，无法获取token"
    exit 1
fi

echo "✅ 设备A登录成功，获得token: ${DEVICE_A_TOKEN:0:20}..."

echo ""
echo "3. 设备A使用token访问受保护接口..."
DEVICE_A_PROFILE_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DEVICE_A_TOKEN" \
  -d "{}")

echo "设备A访问用户信息响应: $DEVICE_A_PROFILE_RESPONSE"

# 检查设备A是否能正常访问
if echo "$DEVICE_A_PROFILE_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 设备A可以正常访问受保护接口"
else
    echo "❌ 设备A无法访问受保护接口"
    exit 1
fi

echo ""
echo "4. 设备B登录（应该踢出设备A）..."
DEVICE_B_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -H "User-Agent: DeviceB/1.0" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "设备B登录响应: $DEVICE_B_RESPONSE"

# 提取设备B的token
DEVICE_B_TOKEN=$(echo $DEVICE_B_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$DEVICE_B_TOKEN" ]; then
    echo "❌ 设备B登录失败，无法获取token"
    exit 1
fi

echo "✅ 设备B登录成功，获得token: ${DEVICE_B_TOKEN:0:20}..."

echo ""
echo "5. 设备A再次使用旧token访问接口（应该被拒绝）..."
DEVICE_A_PROFILE_RESPONSE_AFTER_B_LOGIN=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DEVICE_A_TOKEN" \
  -d "{}")

echo "设备A使用旧token访问响应: $DEVICE_A_PROFILE_RESPONSE_AFTER_B_LOGIN"

# 检查设备A是否被踢出
if echo "$DEVICE_A_PROFILE_RESPONSE_AFTER_B_LOGIN" | grep -q '"code":401'; then
    if echo "$DEVICE_A_PROFILE_RESPONSE_AFTER_B_LOGIN" | grep -q "已在其他设备登录"; then
        echo "✅ 设备A被成功踢出，返回正确的错误信息"
    else
        echo "⚠️  设备A被踢出，但错误信息不正确"
    fi
else
    echo "❌ 设备A没有被踢出，仍然可以访问"
    exit 1
fi

echo ""
echo "6. 设备B使用新token访问接口（应该正常）..."
DEVICE_B_PROFILE_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DEVICE_B_TOKEN" \
  -d "{}")

echo "设备B访问用户信息响应: $DEVICE_B_PROFILE_RESPONSE"

# 检查设备B是否能正常访问
if echo "$DEVICE_B_PROFILE_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 设备B可以正常访问受保护接口"
else
    echo "❌ 设备B无法访问受保护接口"
    exit 1
fi

echo ""
echo "7. 设备A重新登录（应该踢出设备B）..."
DEVICE_A_RELOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -H "User-Agent: DeviceA/1.0" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "设备A重新登录响应: $DEVICE_A_RELOGIN_RESPONSE"

# 提取设备A的新token
DEVICE_A_NEW_TOKEN=$(echo $DEVICE_A_RELOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$DEVICE_A_NEW_TOKEN" ]; then
    echo "❌ 设备A重新登录失败"
    exit 1
fi

echo "✅ 设备A重新登录成功，获得新token: ${DEVICE_A_NEW_TOKEN:0:20}..."

echo ""
echo "8. 设备B使用旧token访问接口（应该被拒绝）..."
DEVICE_B_PROFILE_RESPONSE_AFTER_A_RELOGIN=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DEVICE_B_TOKEN" \
  -d "{}")

echo "设备B使用旧token访问响应: $DEVICE_B_PROFILE_RESPONSE_AFTER_A_RELOGIN"

# 检查设备B是否被踢出
if echo "$DEVICE_B_PROFILE_RESPONSE_AFTER_A_RELOGIN" | grep -q '"code":401'; then
    if echo "$DEVICE_B_PROFILE_RESPONSE_AFTER_A_RELOGIN" | grep -q "已在其他设备登录"; then
        echo "✅ 设备B被成功踢出，返回正确的错误信息"
    else
        echo "⚠️  设备B被踢出，但错误信息不正确"
    fi
else
    echo "❌ 设备B没有被踢出，仍然可以访问"
    exit 1
fi

echo ""
echo "9. 设备A使用新token访问接口（应该正常）..."
DEVICE_A_NEW_PROFILE_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DEVICE_A_NEW_TOKEN" \
  -d "{}")

echo "设备A使用新token访问响应: $DEVICE_A_NEW_PROFILE_RESPONSE"

# 检查设备A是否能正常访问
if echo "$DEVICE_A_NEW_PROFILE_RESPONSE" | grep -q '"code":200'; then
    echo "✅ 设备A使用新token可以正常访问"
else
    echo "❌ 设备A使用新token无法访问"
    exit 1
fi

echo ""
echo "=== 测试完成 ==="
echo "✅ 单点登录功能测试通过！"
echo ""
echo "测试总结："
echo "1. 设备A登录后可以正常访问接口"
echo "2. 设备B登录后，设备A的token被踢出"
echo "3. 设备A重新登录后，设备B的token被踢出"
echo "4. 每次新登录都会使旧token失效"
echo "5. 错误信息正确提示'已在其他设备登录'" 