package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
)

type OrderCacheService struct {
	// redisRepo *database.RedisRepository // 移除
}

func NewOrderCacheService() *OrderCacheService {
	return &OrderCacheService{}
}

// 生成缓存Key
func (s *OrderCacheService) generateCacheKey(periodNumber string) string {
	return fmt.Sprintf("fataMorgana_%s", periodNumber)
}

// 缓存订单数据
func (s *OrderCacheService) CacheOrder(ctx context.Context, order *models.Order) error {
	if order == nil || order.OrderNo == "" || order.PeriodNumber == "" {
		return utils.NewAppError(utils.CodeInvalidParams, "订单数据无效")
	}

	// 生成缓存Key
	cacheKey := s.generateCacheKey(order.PeriodNumber)
	
	// 将订单数据转换为JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return utils.NewAppError(utils.CodeDatabaseError, "订单数据序列化失败")
	}

	// 使用Hash结构存储，小key为订单号
	err = database.GlobalRedisHelper.HSet(ctx, cacheKey, order.OrderNo, string(orderJSON))
	if err != nil {
		return utils.NewAppError(utils.CodeDatabaseError, "缓存订单数据失败")
	}

	// 设置过期时间（24小时）
	err = database.GlobalRedisHelper.Expire(ctx, cacheKey, 24*time.Hour)
	if err != nil {
		// 过期时间设置失败不影响主流程，只记录日志
		utils.LogWarn(nil, "设置缓存过期时间失败: %v", err)
	}

	return nil
}

// 获取期数下的所有订单
func (s *OrderCacheService) GetOrdersByPeriod(ctx context.Context, periodNumber string) ([]*models.Order, error) {
	if periodNumber == "" {
		return nil, utils.NewAppError(utils.CodeInvalidParams, "期号不能为空")
	}

	// 生成缓存Key
	cacheKey := s.generateCacheKey(periodNumber)
	
	// 获取Hash中的所有数据
	orderMap, err := database.GlobalRedisHelper.HGetAll(ctx, cacheKey)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取缓存订单数据失败")
	}

	var orders []*models.Order
	for _, orderJSON := range orderMap {
		var order models.Order
		err := json.Unmarshal([]byte(orderJSON), &order)
		if err != nil {
			utils.LogWarn(nil, "订单数据反序列化失败: %v", err)
			continue
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

// 获取指定订单
func (s *OrderCacheService) GetOrder(ctx context.Context, periodNumber, orderNo string) (*models.Order, error) {
	if periodNumber == "" || orderNo == "" {
		return nil, utils.NewAppError(utils.CodeInvalidParams, "期号和订单号不能为空")
	}

	// 生成缓存Key
	cacheKey := s.generateCacheKey(periodNumber)
	
	// 获取指定订单数据
	orderJSON, err := database.GlobalRedisHelper.HGet(ctx, cacheKey, orderNo)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取缓存订单数据失败")
	}

	if orderJSON == "" {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "订单不存在")
	}

	// 反序列化订单数据
	var order models.Order
	err = json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "订单数据反序列化失败")
	}

	return &order, nil
}

// 删除订单缓存
func (s *OrderCacheService) DeleteOrder(ctx context.Context, periodNumber, orderNo string) error {
	if periodNumber == "" || orderNo == "" {
		return utils.NewAppError(utils.CodeInvalidParams, "期号和订单号不能为空")
	}

	// 生成缓存Key
	cacheKey := s.generateCacheKey(periodNumber)
	
	// 删除指定订单
	err := database.GlobalRedisHelper.HDel(ctx, cacheKey, orderNo)
	if err != nil {
		return utils.NewAppError(utils.CodeDatabaseError, "删除缓存订单数据失败")
	}

	return nil
}

// 删除期数下的所有订单缓存
func (s *OrderCacheService) DeletePeriodOrders(ctx context.Context, periodNumber string) error {
	if periodNumber == "" {
		return utils.NewAppError(utils.CodeInvalidParams, "期号不能为空")
	}

	// 生成缓存Key
	cacheKey := s.generateCacheKey(periodNumber)
	
	// 删除整个Key
	err := database.GlobalRedisHelper.Del(ctx, cacheKey)
	if err != nil {
		return utils.NewAppError(utils.CodeDatabaseError, "删除期数缓存数据失败")
	}

	return nil
}

// 获取期数下的订单数量
func (s *OrderCacheService) GetOrderCount(ctx context.Context, periodNumber string) (int64, error) {
	if periodNumber == "" {
		return 0, utils.NewAppError(utils.CodeInvalidParams, "期号不能为空")
	}

	// 生成缓存Key
	cacheKey := s.generateCacheKey(periodNumber)
	
	// 获取Hash的长度（订单数量）
	count, err := database.GlobalRedisHelper.HLen(ctx, cacheKey)
	if err != nil {
		return 0, utils.NewAppError(utils.CodeDatabaseError, "获取订单数量失败")
	}

	return count, nil
}

// 检查订单是否在缓存中
func (s *OrderCacheService) IsOrderCached(ctx context.Context, periodNumber, orderNo string) (bool, error) {
	if periodNumber == "" || orderNo == "" {
		return false, utils.NewAppError(utils.CodeInvalidParams, "期号和订单号不能为空")
	}

	// 生成缓存Key
	cacheKey := s.generateCacheKey(periodNumber)
	
	// 检查Hash中是否存在该字段
	exists, err := database.GlobalRedisHelper.HExists(ctx, cacheKey, orderNo)
	if err != nil {
		return false, utils.NewAppError(utils.CodeDatabaseError, "检查订单缓存状态失败")
	}

	return exists, nil
}

// 批量缓存订单
func (s *OrderCacheService) BatchCacheOrders(ctx context.Context, orders []*models.Order) error {
	if len(orders) == 0 {
		return nil
	}

	// 按期数分组
	periodOrders := make(map[string][]*models.Order)
	for _, order := range orders {
		if order.PeriodNumber != "" {
			periodOrders[order.PeriodNumber] = append(periodOrders[order.PeriodNumber], order)
		}
	}

	// 按期数批量缓存
	for periodNumber, periodOrderList := range periodOrders {
		cacheKey := s.generateCacheKey(periodNumber)
		
		// 准备批量数据
		orderMap := make(map[string]string)
		for _, order := range periodOrderList {
			orderJSON, err := json.Marshal(order)
			if err != nil {
				utils.LogWarn(nil, "订单数据序列化失败: %v", err)
				continue
			}
			orderMap[order.OrderNo] = string(orderJSON)
		}

		// 批量设置
		err := database.GlobalRedisHelper.HSet(ctx, cacheKey, orderMap)
		if err != nil {
			return utils.NewAppError(utils.CodeDatabaseError, "批量缓存订单数据失败")
		}

		// 设置过期时间
		err = database.GlobalRedisHelper.Expire(ctx, cacheKey, 24*time.Hour)
		if err != nil {
			return utils.NewAppError(utils.CodeDatabaseError, "批量设置缓存过期时间失败")
		}
	}

	return nil
} 