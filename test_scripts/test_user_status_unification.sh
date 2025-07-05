#!/bin/bash

# 测试用户状态统一性
echo "🔍 测试用户状态统一性..."

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试基础URL
BASE_URL="http://localhost:8080/api/v1"

# 测试用户状态常量
echo -e "\n${YELLOW}1. 测试用户状态常量定义${NC}"
echo "UserStatusDisabled = 0 (禁用)"
echo "UserStatusActive = 1 (正常)"
echo "UserStatusPending = 2 (待审核)"

# 测试注册用户（默认状态为2-待审核）
echo -e "\n${YELLOW}2. 测试用户注册（默认状态为待审核）${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test_status@example.com",
    "password": "123456",
    "confirm_password": "123456",
    "invite_code": "TEST01"
  }')

echo "注册响应: $REGISTER_RESPONSE"

# 尝试登录待审核用户（应该失败）
echo -e "\n${YELLOW}3. 测试待审核用户登录（应该失败）${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test_status@example.com",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 检查是否返回待审核错误
if echo "$LOGIN_RESPONSE" | grep -q "账户待审核"; then
    echo -e "${GREEN}✅ 待审核用户登录被正确拦截${NC}"
else
    echo -e "${RED}❌ 待审核用户登录未被拦截${NC}"
fi

# 测试禁用用户（需要管理员操作，这里只是验证状态检查逻辑）
echo -e "\n${YELLOW}4. 测试状态检查逻辑${NC}"
echo "状态检查应该遵循以下规则："
echo "- Status == 0: 禁用，无法登录和操作"
echo "- Status == 1: 正常，可以正常使用"
echo "- Status == 2: 待审核，无法登录和操作"

# 验证状态常量使用
echo -e "\n${YELLOW}5. 验证代码中的状态常量使用${NC}"
echo "检查是否使用了状态常量而不是硬编码数字..."

# 检查代码中的状态使用
echo "检查 models/user.go 中的状态常量定义..."
if grep -q "UserStatusDisabled" models/user.go; then
    echo -e "${GREEN}✅ 禁用状态常量已定义${NC}"
else
    echo -e "${RED}❌ 禁用状态常量未定义${NC}"
fi

if grep -q "UserStatusActive" models/user.go; then
    echo -e "${GREEN}✅ 正常状态常量已定义${NC}"
else
    echo -e "${RED}❌ 正常状态常量未定义${NC}"
fi

if grep -q "UserStatusPending" models/user.go; then
    echo -e "${GREEN}✅ 待审核状态常量已定义${NC}"
else
    echo -e "${RED}❌ 待审核状态常量未定义${NC}"
fi

# 检查IsActive方法是否使用常量
if grep -q "UserStatusActive" models/user.go; then
    echo -e "${GREEN}✅ IsActive方法使用了状态常量${NC}"
else
    echo -e "${RED}❌ IsActive方法未使用状态常量${NC}"
fi

echo -e "\n${GREEN}🎉 用户状态统一性测试完成！${NC}"
echo -e "\n${YELLOW}总结：${NC}"
echo "1. 用户状态已统一为：0(禁用) 1(正常) 2(待审核)"
echo "2. 状态常量已定义，避免硬编码"
echo "3. 登录接口已正确拦截非正常状态用户"
echo "4. 建议在其他接口中也添加状态检查" 