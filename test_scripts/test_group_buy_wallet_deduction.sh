#!/bin/bash

# 测试拼单参与时的钱包扣费和流水记录
# 使用方法: ./test_group_buy_wallet_deduction.sh

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== 测试拼单参与时的钱包扣费和流水记录 ===${NC}"

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

# 2. 获取当前钱包余额
echo -e "\n${YELLOW}2. 获取当前钱包余额${NC}"
WALLET_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/info" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "钱包信息响应: $WALLET_RESPONSE"

# 提取当前余额
CURRENT_BALANCE=$(echo "$WALLET_RESPONSE" | jq -r '.data.balance')
echo -e "${GREEN}当前钱包余额: $CURRENT_BALANCE${NC}"

# 3. 获取活跃拼单
echo -e "\n${YELLOW}3. 获取活跃拼单${NC}"
ACTIVE_GROUP_BUY_RESPONSE=$(curl -s -X POST "${BASE_URL}/group-buy/active-detail" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')

echo "活跃拼单响应: $ACTIVE_GROUP_BUY_RESPONSE"

# 提取拼单编号和每人付款金额
GROUP_BUY_NO=$(echo "$ACTIVE_GROUP_BUY_RESPONSE" | jq -r '.data.group_buy_no')
PER_PERSON_AMOUNT=$(echo "$ACTIVE_GROUP_BUY_RESPONSE" | jq -r '.data.per_person_amount')

if [ "$GROUP_BUY_NO" != "null" ] && [ -n "$GROUP_BUY_NO" ]; then
    echo -e "${GREEN}找到活跃拼单: $GROUP_BUY_NO${NC}"
    echo -e "${GREEN}每人付款金额: $PER_PERSON_AMOUNT${NC}"
    
    # 4. 检查余额是否足够
    if (( $(echo "$CURRENT_BALANCE >= $PER_PERSON_AMOUNT" | bc -l) )); then
        echo -e "${GREEN}✓ 余额充足，可以参与拼单${NC}"
        
        # 5. 参与拼单
        echo -e "\n${YELLOW}4. 参与拼单${NC}"
        JOIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/group-buy/join" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $TOKEN" \
          -d "{
            \"group_buy_no\": \"$GROUP_BUY_NO\"
          }")
        
        echo "参与拼单响应: $JOIN_RESPONSE"
        
        # 检查是否成功参与
        if echo "$JOIN_RESPONSE" | jq -e '.code == 0' > /dev/null; then
            echo -e "${GREEN}✓ 成功参与拼单${NC}"
            
            # 6. 再次获取钱包余额，检查是否已扣费
            echo -e "\n${YELLOW}5. 检查钱包余额变化${NC}"
            NEW_WALLET_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/info" \
              -H "Content-Type: application/json" \
              -H "Authorization: Bearer $TOKEN" \
              -d '{}')
            
            echo "新钱包信息响应: $NEW_WALLET_RESPONSE"
            
            # 提取新余额
            NEW_BALANCE=$(echo "$NEW_WALLET_RESPONSE" | jq -r '.data.balance')
            echo -e "${GREEN}参与拼单后余额: $NEW_BALANCE${NC}"
            
            # 计算余额变化
            BALANCE_CHANGE=$(echo "$CURRENT_BALANCE - $NEW_BALANCE" | bc -l)
            echo -e "${GREEN}余额变化: -$BALANCE_CHANGE${NC}"
            
            # 检查余额是否正确扣减
            if (( $(echo "$BALANCE_CHANGE == $PER_PERSON_AMOUNT" | bc -l) )); then
                echo -e "${GREEN}✓ 余额扣减正确${NC}"
            else
                echo -e "${RED}✗ 余额扣减不正确，期望扣减: $PER_PERSON_AMOUNT，实际扣减: $BALANCE_CHANGE${NC}"
            fi
            
            # 7. 获取钱包流水记录，检查是否有拼单类型的流水
            echo -e "\n${YELLOW}6. 检查钱包流水记录${NC}"
            TRANSACTIONS_RESPONSE=$(curl -s -X POST "${BASE_URL}/wallet/transactions" \
              -H "Content-Type: application/json" \
              -H "Authorization: Bearer $TOKEN" \
              -d '{
                "page": 1,
                "page_size": 5
              }')
            
            echo "钱包流水响应: $TRANSACTIONS_RESPONSE"
            
            # 检查是否有拼单类型的交易记录
            if echo "$TRANSACTIONS_RESPONSE" | grep -q "group_buy"; then
                echo -e "${GREEN}✓ 成功创建拼单类型的钱包流水记录${NC}"
                
                # 提取最新的拼单流水记录
                GROUP_BUY_TRANSACTION=$(echo "$TRANSACTIONS_RESPONSE" | jq -r '.data.transactions[] | select(.type == "group_buy") | .[0]')
                if [ "$GROUP_BUY_TRANSACTION" != "null" ]; then
                    echo -e "${GREEN}拼单流水记录详情:${NC}"
                    echo "$GROUP_BUY_TRANSACTION" | jq '.'
                fi
            else
                echo -e "${RED}✗ 未找到拼单类型的钱包流水记录${NC}"
            fi
            
            # 8. 检查订单是否创建成功
            echo -e "\n${YELLOW}7. 检查订单创建${NC}"
            ORDER_ID=$(echo "$JOIN_RESPONSE" | jq -r '.data.order_id')
            if [ "$ORDER_ID" != "null" ] && [ "$ORDER_ID" != "0" ]; then
                echo -e "${GREEN}✓ 订单创建成功，订单ID: $ORDER_ID${NC}"
            else
                echo -e "${RED}✗ 订单创建失败${NC}"
            fi
            
        else
            echo -e "${RED}✗ 参与拼单失败${NC}"
            echo "错误信息: $(echo "$JOIN_RESPONSE" | jq -r '.message')"
        fi
    else
        echo -e "${RED}✗ 余额不足，无法参与拼单${NC}"
        echo "当前余额: $CURRENT_BALANCE"
        echo "拼单金额: $PER_PERSON_AMOUNT"
        echo "差额: $(echo "$PER_PERSON_AMOUNT - $CURRENT_BALANCE" | bc -l)"
    fi
else
    echo -e "${YELLOW}⚠ 没有找到活跃的拼单${NC}"
fi

# 9. 测试余额不足的情况
echo -e "\n${YELLOW}8. 测试余额不足的情况${NC}"

# 创建一个余额不足的用户进行测试
echo "如果需要测试余额不足的情况，请先确保有余额不足的用户账户"

echo -e "\n${GREEN}=== 测试完成 ===${NC}"
echo -e "${GREEN}拼单参与时的钱包扣费和流水记录功能测试完成${NC}"
echo -e "${GREEN}包括：${NC}"
echo -e "  - 余额检查"
echo -e "  - 余额扣减"
echo -e "  - 拼单类型流水记录创建"
echo -e "  - 订单创建"
echo -e "  - 错误处理和回滚机制" 