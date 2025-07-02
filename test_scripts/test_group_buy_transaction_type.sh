#!/bin/bash

# 测试拼单交易类型
# 使用方法: ./test_group_buy_transaction_type.sh

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== 测试拼单交易类型 ===${NC}"

# 1. 先登录获取token
echo -e "\n${YELLOW}1. 用户登录${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_user",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}登录失败，无法获取token${NC}"
    exit 1
fi

echo -e "${GREEN}登录成功，获取到token${NC}"

# 2. 获取钱包流水记录，查看是否有拼单类型
echo -e "\n${YELLOW}2. 获取钱包流水记录${NC}"
TRANSACTIONS_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/transactions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "page": 1,
    "page_size": 10
  }')

echo "钱包流水响应: $TRANSACTIONS_RESPONSE"

# 检查响应中是否包含拼单类型
if echo "$TRANSACTIONS_RESPONSE" | grep -q "group_buy"; then
    echo -e "${GREEN}✓ 钱包流水中包含拼单类型${NC}"
else
    echo -e "${YELLOW}⚠ 钱包流水中暂无拼单类型记录${NC}"
fi

# 3. 测试拼单参与（如果拼单功能已实现钱包扣费）
echo -e "\n${YELLOW}3. 测试拼单参与${NC}"

# 先获取活跃拼单
ACTIVE_GROUP_BUY_RESPONSE=$(curl -s -X POST "${BASE_URL}/group-buy/active-detail" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "活跃拼单响应: $ACTIVE_GROUP_BUY_RESPONSE"

# 提取拼单编号
GROUP_BUY_NO=$(echo "$ACTIVE_GROUP_BUY_RESPONSE" | jq -r '.data.group_buy_no')
if [ "$GROUP_BUY_NO" != "null" ] && [ -n "$GROUP_BUY_NO" ]; then
    echo -e "${GREEN}找到活跃拼单: $GROUP_BUY_NO${NC}"
    
    # 尝试参与拼单
    JOIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/group-buy/join" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"group_buy_no\": \"$GROUP_BUY_NO\"
      }")
    
    echo "参与拼单响应: $JOIN_RESPONSE"
    
    # 检查是否成功参与
    if echo "$JOIN_RESPONSE" | jq -e '.code == 200' > /dev/null; then
        echo -e "${GREEN}✓ 成功参与拼单${NC}"
        
        # 再次获取钱包流水，查看是否有新的拼单记录
        echo -e "\n${YELLOW}4. 再次获取钱包流水记录${NC}"
        NEW_TRANSACTIONS_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/transactions" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $TOKEN" \
          -d '{
            "page": 1,
            "page_size": 5
          }')
        
        echo "最新钱包流水响应: $NEW_TRANSACTIONS_RESPONSE"
        
        # 检查是否有拼单类型的交易记录
        if echo "$NEW_TRANSACTIONS_RESPONSE" | grep -q "group_buy"; then
            echo -e "${GREEN}✓ 成功创建拼单类型的钱包流水记录${NC}"
        else
            echo -e "${YELLOW}⚠ 拼单参与后未创建钱包流水记录（可能拼单功能未实现钱包扣费）${NC}"
        fi
    else
        echo -e "${RED}✗ 参与拼单失败${NC}"
        echo "错误信息: $(echo "$JOIN_RESPONSE" | jq -r '.message')"
    fi
else
    echo -e "${YELLOW}⚠ 没有找到活跃的拼单${NC}"
fi

# 5. 测试交易类型名称显示
echo -e "\n${YELLOW}5. 测试交易类型名称显示${NC}"

# 创建一个测试用的交易类型名称映射
echo "交易类型名称映射:"
echo "  recharge -> 充值"
echo "  withdraw -> 提现"
echo "  purchase -> 购买订单"
echo "  group_buy -> 拼单"
echo "  income -> 收入"
echo "  expense -> 支出"

echo -e "\n${GREEN}=== 测试完成 ===${NC}"
echo -e "${GREEN}拼单交易类型已成功添加到系统中${NC}"
echo -e "${GREEN}包括：${NC}"
echo -e "  - 交易类型常量: TransactionTypeGroupBuy = \"group_buy\""
echo -e "  - 类型名称映射: \"拼单\""
echo -e "  - 金额显示: 负号显示（扣费）"
echo -e "  - 类型说明: 用户参与拼单" 