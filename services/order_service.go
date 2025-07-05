package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
	"log"
	"time"
)

// OrderService 订单服务
type OrderService struct {
	orderRepo       *database.OrderRepository
	walletRepo      *database.WalletRepository
	memberLevelRepo *database.MemberLevelRepository
}

// NewOrderService 创建订单服务实例
func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:       database.NewOrderRepository(),
		walletRepo:      database.NewWalletRepository(),
		memberLevelRepo: database.NewMemberLevelRepository(database.DB),
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Uid           string  `json:"uid"`                              // 从token中获取，不需要在请求中传递
	PeriodNumber  string  `json:"period_number" binding:"required"` // 期数编号
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	LikeCount     int     `json:"like_count" binding:"gte=0"`
	ShareCount    int     `json:"share_count" binding:"gte=0"`
	FollowCount   int     `json:"follow_count" binding:"gte=0"`
	FavoriteCount int     `json:"favorite_count" binding:"gte=0"`
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

	// 校验期数
	if err := s.validatePeriod(ctx, req.PeriodNumber); err != nil {
		return nil, err
	}

	// 检查用户是否已购买过该期号
	if err := s.checkUserPeriodDuplicate(ctx, req.Uid, req.PeriodNumber); err != nil {
		return nil, err
	}

	// 获取价格配置并验证金额
	if err := s.validateOrderAmount(ctx, req); err != nil {
		return nil, err
	}

	// 获取用户信息以获取经验值
	userRepo := database.NewUserRepository()
	user, err := userRepo.FindByUid(ctx, req.Uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeUserInfoGetFailed, "获取用户信息失败")
	}

	// 根据用户经验值获取等级配置并计算利润金额
	profitAmount := s.calculateProfitAmount(ctx, user.Experience, req.Amount)

	// 创建订单对象
	order := &models.Order{
		Uid:           req.Uid,
		Amount:        req.Amount,
		ProfitAmount:  profitAmount,
		LikeCount:     req.LikeCount,
		ShareCount:    req.ShareCount,
		FollowCount:   req.FollowCount,
		FavoriteCount: req.FavoriteCount,
		Status:        models.OrderStatusPending,
		ExpireTime:    time.Now().UTC().Add(5 * time.Minute), // 创建时间+5分钟
		AuditorUid:    operatorUid,
	}

	// 验证订单数据
	if err := order.ValidateOrderData(); err != nil {
		return nil, utils.NewAppError(utils.CodeOrderDataValidateFailed, "订单数据验证失败")
	}

	// 初始化任务状态
	order.InitializeTaskStatuses()

	// 生成订单号
	order.OrderNo = utils.GenerateOrderNo()

	// 设置期号
	order.PeriodNumber = req.PeriodNumber

	// 获取用户钱包
	wallet, err := s.walletRepo.FindWalletByUid(ctx, req.Uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeWalletGetFailed, "获取钱包失败")
	}

	// 检查钱包状态
	if !wallet.IsActive() {
		return nil, utils.NewAppError(utils.CodeWalletFrozenOrder, "钱包已被冻结，无法创建订单")
	}

	// 检查余额是否足够
	if wallet.Balance < req.Amount {
		return nil, utils.NewAppError(utils.CodeBalanceInsufficient, "余额不足")
	}

	// 记录交易前余额
	balanceBefore := wallet.Balance

	// 扣减余额
	if err := wallet.Withdraw(req.Amount); err != nil {
		return nil, utils.NewAppError(utils.CodeBalanceDeductFailed, "扣减余额失败")
	}

	// 更新钱包
	if err := s.walletRepo.UpdateWallet(ctx, wallet); err != nil {
		return nil, utils.NewAppError(utils.CodeWalletUpdateFailed, "更新钱包失败")
	}

	// 创建订单
	if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
		// 如果创建订单失败，需要回滚扣减的余额
		wallet.Recharge(req.Amount)
		s.walletRepo.UpdateWallet(ctx, wallet)
		return nil, utils.NewAppError(utils.CodeOrderCreateFailed, "创建订单失败")
	}

	// 创建交易流水记录
	transaction := &models.WalletTransaction{
		TransactionNo:  utils.GenerateTransactionNo("ORDER"),
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

// validatePeriod 校验期数
func (s *OrderService) validatePeriod(ctx context.Context, periodNumber string) error {
	// 创建期数Repository
	periodRepo := database.NewLotteryPeriodRepository()

	// 根据期数编号获取期数
	period, err := periodRepo.GetPeriodByNumber(ctx, periodNumber)
	if err != nil {
		return utils.NewAppError(utils.CodePeriodNotFound, "期数不存在")
	}

	// 检查当前时间是否在期数的订单开始时间和订单结束时间范围内
	now := time.Now().UTC()
	if now.Before(period.OrderStartTime) {
		return utils.NewAppError(utils.CodePeriodNotStarted, "期数还未开始")
	}

	if now.After(period.OrderEndTime) || now.Equal(period.OrderEndTime) {
		return utils.NewAppError(utils.CodePeriodEnded, "期数已结束")
	}

	return nil
}

// validateOrderAmount 验证订单金额
func (s *OrderService) validateOrderAmount(ctx context.Context, req *CreateOrderRequest) error {
	// 获取价格配置
	purchaseConfig, err := s.getPurchaseConfig(ctx)
	if err != nil {
		return utils.NewAppError(utils.CodePriceConfigGetFailed, "获取价格配置失败")
	}

	// 计算订单金额
	calculatedAmount := float64(req.LikeCount)*purchaseConfig.LikeAmount +
		float64(req.ShareCount)*purchaseConfig.ShareAmount +
		float64(req.FollowCount)*purchaseConfig.ForwardAmount +
		float64(req.FavoriteCount)*purchaseConfig.FavoriteAmount

	// 比较计算金额和请求金额
	if calculatedAmount != req.Amount {
		return utils.NewAppError(utils.CodeOrderAmountMismatch, "订单金额不匹配")
	}

	return nil
}

// GetOrderList 获取订单列表
func (s *OrderService) GetOrderList(req *models.GetOrderListRequest, uid string) (*GetOrderListResponse, error) {
	ctx := context.Background()

	// 限制page_size最大值，超出时设置为默认值20
	if req.PageSize > 20 {
		req.PageSize = 20
	}

	// 验证状态类型参数
	if req.Status < 1 || req.Status > 3 {
		return nil, utils.NewAppError(utils.CodeOrderStatusInvalid, "状态类型参数无效，必须是1(进行中)、2(已完成)或3(拼单数据)")
	}

	// 如果status为3，从拼单表获取数据
	if req.Status == 3 {
		return s.getGroupBuyList(ctx, uid, req.Page, req.PageSize)
	}

	// 获取订单列表
	orders, total, err := s.orderRepo.GetUserOrders(ctx, uid, req.Page, req.PageSize, req.Status)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeOrderListGetFailed, "获取订单列表失败")
	}

	// 转换为响应格式
	var orderResponses []models.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, order.ToResponse())
	}

	// 计算分页信息
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	hasNext := req.Page < totalPages
	hasPrev := req.Page > 1

	return &GetOrderListResponse{
		Orders: orderResponses,
		Pagination: PaginationInfo{
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
			Total:       total,
			TotalPages:  totalPages,
			HasNext:     hasNext,
			HasPrev:     hasPrev,
		},
	}, nil
}

// getGroupBuyList 获取拼单列表
func (s *OrderService) getGroupBuyList(ctx context.Context, uid string, page, pageSize int) (*GetOrderListResponse, error) {
	// 创建拼单Repository
	groupBuyRepo := database.NewGroupBuyRepository()

	// 获取拼单列表，查询符合开始时间和截至时间在当前范围内的拼单
	groupBuys, total, err := groupBuyRepo.GetActiveGroupBuys(ctx, uid, page, pageSize)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeGroupBuyListGetFailed, "获取拼单列表失败")
	}

	// 转换为订单响应格式
	var orderResponses []models.OrderResponse
	for _, groupBuy := range groupBuys {
		// 将拼单数据转换为订单响应格式
		orderResponse := models.OrderResponse{
			ID:           groupBuy.ID,
			OrderNo:      groupBuy.GroupBuyNo,
			Uid:          groupBuy.Uid,
			Amount:       groupBuy.PerPersonAmount,
			ProfitAmount: 0, // 拼单没有利润金额
			Status:       groupBuy.Status,
			StatusName:   s.getGroupBuyStatusName(groupBuy.Status),
			ExpireTime:   groupBuy.Deadline,
			CreatedAt:    groupBuy.CreatedAt,
			UpdatedAt:    groupBuy.UpdatedAt,
			IsExpired:    time.Now().UTC().After(groupBuy.Deadline),
			RemainingTime: func() int64 {
				if time.Now().UTC().After(groupBuy.Deadline) {
					return 0
				}
				return int64(time.Until(groupBuy.Deadline).Seconds())
			}(),
		}
		orderResponses = append(orderResponses, orderResponse)
	}

	// 计算分页信息
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	return &GetOrderListResponse{
		Orders: orderResponses,
		Pagination: PaginationInfo{
			CurrentPage: page,
			PageSize:    pageSize,
			Total:       total,
			TotalPages:  totalPages,
			HasNext:     hasNext,
			HasPrev:     hasPrev,
		},
	}, nil
}

// getGroupBuyStatusName 获取拼单状态名称
func (s *OrderService) getGroupBuyStatusName(status string) string {
	statusNames := map[string]string{
		models.GroupBuyStatusNotStarted: "未开启",
		models.GroupBuyStatusPending:    "进行中",
		models.GroupBuyStatusSuccess:    "已完成",
	}
	return statusNames[status]
}

// GetOrderDetail 获取订单详情
func (s *OrderService) GetOrderDetail(req *models.GetOrderDetailRequest, uid string) (*models.OrderResponse, error) {
	ctx := context.Background()

	// 获取订单
	order, err := s.orderRepo.FindOrderByOrderNo(ctx, req.OrderNo)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeOrderDetailGetFailed, "获取订单详情失败")
	}

	// 检查订单是否属于当前用户
	if order.Uid != uid {
		return nil, utils.NewAppError(utils.CodeOrderAccessDenied, "无权访问此订单")
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
		return nil, utils.NewAppError(utils.CodeOrderStatsGetFailed, "获取订单统计失败")
	}

	return &GetOrderStatsResponse{
		Stats: stats,
	}, nil
}

// calculateProfitAmount 根据用户经验值和订单金额计算利润金额
func (s *OrderService) calculateProfitAmount(ctx context.Context, experience int, amount float64) float64 {
	// 根据经验值获取等级配置
	level, err := s.memberLevelRepo.GetByExperience(ctx, experience)
	if err != nil {
		// 如果获取等级配置失败，返回0利润
		return 0.0
	}

	// 计算利润金额：订单金额 × (返现比例 / 100)
	profitAmount := amount * (level.CashbackRatio / 100.0)
	return profitAmount
}

// GetPeriodList 获取期数列表
func (s *OrderService) GetPeriodList() (*models.PeriodListResponse, error) {
	ctx := context.Background()

	// 添加时间调试信息
	now := time.Now().UTC()
	log.Printf("当前时间(UTC): %s", now.Format("2006-01-02 15:04:05 UTC"))

	// 创建期数Repository
	periodRepo := database.NewLotteryPeriodRepository()

	// 获取当前活跃期数
	period, err := periodRepo.GetCurrentPeriod(ctx)
	if err != nil {
		return nil, utils.NewAppError(utils.CodePeriodInfoGetFailed, "获取期数信息失败")
	}

	// 添加期数时间调试信息
	log.Printf("期数信息: ID=%d, 期号=%s, 开始时间=%s, 结束时间=%s, 状态=%s", 
		period.ID, period.PeriodNumber, 
		period.OrderStartTime.Format("2006-01-02 15:04:05 UTC"),
		period.OrderEndTime.Format("2006-01-02 15:04:05 UTC"),
		period.GetStatus())

	// 获取价格配置
	purchaseConfig, err := s.getPurchaseConfig(ctx)
	if err != nil {
		return nil, utils.NewAppError(utils.CodePriceConfigGetFailed, "获取价格配置失败")
	}

	// 构建响应
	response := &models.PeriodListResponse{
		ID:           period.ID,
		PeriodNumber: period.PeriodNumber,
		StartTime:    period.OrderStartTime.Format("2006-01-02T15:04:05Z"),
		EndTime:      period.OrderEndTime.Format("2006-01-02T15:04:05Z"),
		Status:       period.GetStatus(),
		IsExpired:    period.IsExpired(),
		RemainingTime: func() int64 {
			if period.IsExpired() {
				return 0
			}
			return int64(time.Until(period.OrderEndTime).Seconds())
		}(),
		LikeAmount:     purchaseConfig.LikeAmount,
		ShareAmount:    purchaseConfig.ShareAmount,
		ForwardAmount:  purchaseConfig.ForwardAmount,
		FavoriteAmount: purchaseConfig.FavoriteAmount,
	}

	return response, nil
}

// getPurchaseConfig 从Redis获取价格配置
func (s *OrderService) getPurchaseConfig(ctx context.Context) (*models.PurchaseConfig, error) {
	// 从Redis获取价格配置
	configJSON, err := database.RedisClient.Get(ctx, "purchase_config").Result()
	if err != nil {
		// 如果Redis中没有数据，返回默认配置
		return &models.PurchaseConfig{
			LikeAmount:     0,
			ShareAmount:    0,
			ForwardAmount:  0,
			FavoriteAmount: 0,
		}, nil
	}

	// 解析JSON
	var config models.PurchaseConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, utils.NewAppError(utils.CodePriceConfigParseFailed, "解析价格配置失败")
	}

	return &config, nil
}

// generateUserOrderNumber 生成用户订单编号
func (s *OrderService) generateUserOrderNumber(ctx context.Context, uid string) string {
	// 实现生成用户订单编号的逻辑
	// 这里可以根据需要生成唯一的用户订单编号
	return utils.RandomString(10) // 临时返回一个随机字符串
}

// checkUserPeriodDuplicate 检查用户是否已购买过该期号
func (s *OrderService) checkUserPeriodDuplicate(ctx context.Context, uid string, periodNumber string) error {
	// 查询用户是否已经用该期号创建过订单
	exists, err := s.orderRepo.CheckUserPeriodExists(ctx, uid, periodNumber)
	if err != nil {
		return utils.NewAppError(utils.CodePeriodDuplicateCheckFailed, "检查期号重复失败")
	}

	if exists {
		return utils.NewAppError(utils.CodePeriodAlreadyBought, "您已经购买过期号的订单")
	}

	return nil
}

// GetAllOrderList 获取所有订单列表（只需登录即可查看所有订单）
func (s *OrderService) GetAllOrderList(req *models.GetOrderListRequest) (*GetOrderListResponse, error) {
	ctx := context.Background()

	if req.PageSize > 20 {
		req.PageSize = 20
	}

	if req.Status < 1 || req.Status > 3 {
		return nil, utils.NewAppError(utils.CodeOrderStatusInvalid, "状态类型参数无效，必须是1(进行中)、2(已完成)或3(拼单数据)")
	}

	if req.Status == 3 {
		// 拼单数据不支持全量查询，直接返回空
		return &GetOrderListResponse{Orders: []models.OrderResponse{}, Pagination: PaginationInfo{CurrentPage: req.Page, PageSize: req.PageSize, Total: 0, TotalPages: 0, HasNext: false, HasPrev: false}}, nil
	}

	// 获取所有订单
	status := models.GetStatusByType(req.Status)
	orders, total, err := s.orderRepo.GetOrdersByStatus(ctx, status, req.Page, req.PageSize)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeOrderListGetFailed, "获取订单列表失败")
	}

	var orderResponses []models.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, order.ToResponse())
	}

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	hasNext := req.Page < totalPages
	hasPrev := req.Page > 1

	return &GetOrderListResponse{
		Orders: orderResponses,
		Pagination: PaginationInfo{
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
			Total:       total,
			TotalPages:  totalPages,
			HasNext:     hasNext,
			HasPrev:     hasPrev,
		},
	}, nil
}
