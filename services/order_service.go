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
	orderRepo  *database.OrderRepository
	walletRepo *database.WalletRepository
}

// NewOrderService 创建订单服务实例
func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:  database.NewOrderRepository(),
		walletRepo: database.NewWalletRepository(),
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Uid          string  `json:"uid" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	ProfitAmount float64 `json:"profit_amount" binding:"required,gte=0"`
	LikeCount    int     `json:"like_count" binding:"gte=0"`
	ShareCount   int     `json:"share_count" binding:"gte=0"`
	FollowCount  int     `json:"follow_count" binding:"gte=0"`
	FavoriteCount int    `json:"favorite_count" binding:"gte=0"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	OrderNo string  `json:"order_no"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
	Message string  `json:"message"`
}

// GetOrderListResponse 获取订单列表响应
type GetOrderListResponse struct {
	Orders     []models.OrderResponse `json:"orders"`
	Pagination PaginationInfo         `json:"pagination"`
}

// GetOrderStatsResponse 获取订单统计响应
type GetOrderStatsResponse struct {
	Stats map[string]interface{} `json:"stats"`
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(req *CreateOrderRequest, operatorUid string) (*CreateOrderResponse, error) {
	ctx := context.Background()

	// 创建订单对象
	order := &models.Order{
		Uid:          req.Uid,
		Amount:       req.Amount,
		ProfitAmount: req.ProfitAmount,
		LikeCount:    req.LikeCount,
		ShareCount:   req.ShareCount,
		FollowCount:  req.FollowCount,
		FavoriteCount: req.FavoriteCount,
		Status:       models.OrderStatusPending,
		ExpireTime:   time.Now().Add(5 * time.Minute), // 创建时间+5分钟
		AuditorUid:   operatorUid,
	}

	// 验证订单数据
	if err := order.ValidateOrderData(); err != nil {
		return nil, fmt.Errorf("订单数据验证失败: %w", err)
	}

	// 初始化任务状态
	order.InitializeTaskStatuses()

	// 生成订单号
	order.OrderNo = s.generateOrderNo()

	// 获取用户钱包
	wallet, err := s.walletRepo.FindWalletByUid(ctx, req.Uid)
	if err != nil {
		return nil, fmt.Errorf("获取钱包失败: %w", err)
	}

	// 检查钱包状态
	if !wallet.IsActive() {
		return nil, fmt.Errorf("钱包已被冻结，无法创建订单")
	}

	// 检查余额是否足够
	if wallet.Balance < req.Amount {
		return nil, fmt.Errorf("余额不足，当前余额: %.2f，订单金额: %.2f", wallet.Balance, req.Amount)
	}

	// 记录交易前余额
	balanceBefore := wallet.Balance

	// 扣减余额
	if err := wallet.Withdraw(req.Amount); err != nil {
		return nil, fmt.Errorf("扣减余额失败: %w", err)
	}

	// 更新钱包
	if err := s.walletRepo.UpdateWallet(ctx, wallet); err != nil {
		return nil, fmt.Errorf("更新钱包失败: %w", err)
	}

	// 创建订单
	if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
		// 如果创建订单失败，需要回滚扣减的余额
		wallet.Recharge(req.Amount)
		s.walletRepo.UpdateWallet(ctx, wallet)
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	// 创建交易流水记录
	transaction := &models.WalletTransaction{
		TransactionNo:  s.generateTransactionNo(),
		Uid:            req.Uid,
		Type:           models.TransactionTypeOrderBuy,
		Amount:         req.Amount,
		BalanceBefore:  balanceBefore,
		BalanceAfter:   wallet.Balance,
		Status:         models.TransactionStatusSuccess,
		Description:    fmt.Sprintf("购买订单 %s", order.OrderNo),
		RelatedOrderNo: order.OrderNo,
		OperatorUid:    operatorUid,
	}

	if err := s.walletRepo.CreateTransaction(ctx, transaction); err != nil {
		// 如果创建交易记录失败，记录日志但不影响订单创建
		fmt.Printf("创建交易记录失败: %v\n", err)
	}

	return &CreateOrderResponse{
		OrderNo: order.OrderNo,
		Amount:  req.Amount,
		Status:  models.OrderStatusPending,
		Message: "订单创建成功",
	}, nil
}

// GetOrderList 获取订单列表
func (s *OrderService) GetOrderList(req *models.GetOrderListRequest, uid string) (*GetOrderListResponse, error) {
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

	// 验证状态类型参数
	if req.Status < 1 || req.Status > 3 {
		return nil, errors.New("状态类型参数无效，必须是1(进行中)、2(已完成)或3(全部)")
	}

	// 获取订单列表
	orders, total, err := s.orderRepo.GetUserOrders(ctx, uid, page, pageSize, req.Status)
	if err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
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

	pagination := PaginationInfo{
		CurrentPage: page,
		PageSize:    pageSize,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	return &GetOrderListResponse{
		Orders:     orderResponses,
		Pagination: pagination,
	}, nil
}

// GetOrderDetail 获取订单详情
func (s *OrderService) GetOrderDetail(req *models.GetOrderDetailRequest, uid string) (*models.OrderResponse, error) {
	ctx := context.Background()

	// 获取订单
	order, err := s.orderRepo.FindOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("获取订单详情失败: %w", err)
	}

	// 检查订单是否属于当前用户
	if order.Uid != uid {
		return nil, errors.New("无权访问此订单")
	}

	// 转换为响应格式
	response := order.ToResponse()

	return &response, nil
}

// GetOrderStats 获取订单统计
func (s *OrderService) GetOrderStats(uid string) (*GetOrderStatsResponse, error) {
	ctx := context.Background()

	// 获取订单统计
	stats, err := s.orderRepo.GetOrderStats(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("获取订单统计失败: %w", err)
	}

	return &GetOrderStatsResponse{
		Stats: stats,
	}, nil
}

// generateOrderNo 生成订单号
func (s *OrderService) generateOrderNo() string {
	// 格式：ORD + 年月日 + 时分秒 + 4位随机数
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := utils.RandomString(4)
	return fmt.Sprintf("ORD%s%s", timestamp, random)
}

// generateTransactionNo 生成交易流水号
func (s *OrderService) generateTransactionNo() string {
	// 格式：TX + 年月日 + 时分秒 + 4位随机数
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := utils.RandomString(4)
	return fmt.Sprintf("TX%s%s", timestamp, random)
} 