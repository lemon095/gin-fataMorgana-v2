package services

import (
	"context"
	"errors"
	"fmt"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
	"time"
)

// OrderService 订单服务
type OrderService struct {
	orderRepo *database.OrderRepository
}

// NewOrderService 创建订单服务实例
func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo: database.NewOrderRepository(),
	}
}

// GetUserOrders 获取用户订单列表
func (s *OrderService) GetUserOrders(req *models.OrderListRequest, uid string) (*models.OrderListResponse, error) {
	ctx := context.Background()

	// 设置默认分页参数
	page := req.Page
	pageSize := req.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 验证分页参数
	if pageSize > 100 {
		return nil, errors.New("每页数量不能超过100")
	}

	// 获取订单列表
	orders, total, err := s.orderRepo.GetUserOrders(ctx, uid, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取用户订单列表失败: %w", err)
	}

	// 转换为响应格式
	var orderResponses []models.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, order.ToResponse())
	}

	// 计算分页信息
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := models.PaginationInfo{
		CurrentPage: page,
		PageSize:    pageSize,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	return &models.OrderListResponse{
		Orders:     orderResponses,
		Pagination: pagination,
	}, nil
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(uid string, buyAmount float64, description string) (*models.Order, error) {
	ctx := context.Background()

	if buyAmount <= 0 {
		return nil, errors.New("买入金额必须大于0")
	}

	// 生成订单号
	orderNo := utils.GenerateOrderNo()

	// 创建订单
	order := &models.Order{
		OrderNo:      orderNo,
		Uid:          uid,
		BuyAmount:    buyAmount,
		ProfitAmount: 0.00, // 初始利润为0
		Status:       models.OrderStatusPending,
		Description:  description,
		Remark:       "",
	}

	if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	return order, nil
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(orderNo string, status string, profitAmount float64, remark string) error {
	ctx := context.Background()

	// 获取订单
	order, err := s.orderRepo.FindOrderByOrderNo(ctx, orderNo)
	if err != nil {
		return fmt.Errorf("获取订单失败: %w", err)
	}

	// 更新订单状态和利润
	order.Status = status
	if profitAmount > 0 {
		order.ProfitAmount = profitAmount
	}
	if remark != "" {
		order.Remark = remark
	}

	if err := s.orderRepo.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("更新订单失败: %w", err)
	}

	return nil
}

// GetOrderByOrderNo 根据订单号获取订单
func (s *OrderService) GetOrderByOrderNo(orderNo string) (*models.Order, error) {
	ctx := context.Background()

	order, err := s.orderRepo.FindOrderByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	return order, nil
}

// GetUserOrderStats 获取用户订单统计
func (s *OrderService) GetUserOrderStats(uid string) (map[string]interface{}, error) {
	ctx := context.Background()

	stats, err := s.orderRepo.GetUserOrderStats(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("获取用户订单统计失败: %w", err)
	}

	return stats, nil
}

// GetOrdersByStatus 根据状态获取用户订单
func (s *OrderService) GetOrdersByStatus(req *models.OrderListRequest, uid string, status string) (*models.OrderListResponse, error) {
	ctx := context.Background()

	// 设置默认分页参数
	page := req.Page
	pageSize := req.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 验证分页参数
	if pageSize > 100 {
		return nil, errors.New("每页数量不能超过100")
	}

	// 获取订单列表
	orders, total, err := s.orderRepo.GetOrdersByStatus(ctx, uid, status, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取用户订单列表失败: %w", err)
	}

	// 转换为响应格式
	var orderResponses []models.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, order.ToResponse())
	}

	// 计算分页信息
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := models.PaginationInfo{
		CurrentPage: page,
		PageSize:    pageSize,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	return &models.OrderListResponse{
		Orders:     orderResponses,
		Pagination: pagination,
	}, nil
}

// GetOrdersByDateRange 根据日期范围获取用户订单
func (s *OrderService) GetOrdersByDateRange(req *models.OrderListRequest, uid string, startDate, endDate string) (*models.OrderListResponse, error) {
	ctx := context.Background()

	// 设置默认分页参数
	page := req.Page
	pageSize := req.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 验证分页参数
	if pageSize > 100 {
		return nil, errors.New("每页数量不能超过100")
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		return nil, errors.New("开始日期格式错误，应为YYYY-MM-DD")
	}
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		return nil, errors.New("结束日期格式错误，应为YYYY-MM-DD")
	}

	// 获取订单列表
	orders, total, err := s.orderRepo.GetOrdersByDateRange(ctx, uid, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取用户订单列表失败: %w", err)
	}

	// 转换为响应格式
	var orderResponses []models.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, order.ToResponse())
	}

	// 计算分页信息
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := models.PaginationInfo{
		CurrentPage: page,
		PageSize:    pageSize,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	return &models.OrderListResponse{
		Orders:     orderResponses,
		Pagination: pagination,
	}, nil
} 