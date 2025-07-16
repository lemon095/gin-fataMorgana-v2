package controllers

import (
	"context"
	"gin-fataMorgana/database"
	"gin-fataMorgana/middleware"
	"gin-fataMorgana/models"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OrderController 订单控制器
type OrderController struct {
	orderService            *services.OrderService
	operationFailureService *services.OperationFailureService
}

// NewOrderController 创建订单控制器实例
func NewOrderController() *OrderController {
	return &OrderController{
		orderService:            services.NewOrderService(),
		operationFailureService: services.NewOperationFailureService(),
	}
}

// GetOrderList 获取订单列表
func (oc *OrderController) GetOrderList(c *gin.Context) {
	var req models.GetOrderListRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 根据user_id查询用户信息获取uid
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		utils.ErrorWithMessage(c, utils.CodeUserNotFound, "账户已被删除，无法查询订单")
		return
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		utils.ErrorWithMessage(c, utils.CodeAccountLocked, "账户已被禁用，无法查询订单")
		return
	}

	// 获取订单列表
	response, err := oc.orderService.GetOrderList(&req, user.Uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetMyOrderList 获取我的订单列表（只获取当前用户的订单）
func (oc *OrderController) GetMyOrderList(c *gin.Context) {
	var req models.GetOrderListRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 根据user_id查询用户信息获取uid
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		utils.ErrorWithMessage(c, utils.CodeUserNotFound, "账户已被删除，无法查询订单")
		return
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		utils.ErrorWithMessage(c, utils.CodeAccountLocked, "账户已被禁用，无法查询订单")
		return
	}

	// 获取订单列表
	response, err := oc.orderService.GetOrderList(&req, user.Uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	// 创建新的响应结构，包含 order_type 字段
	type OrderWithType struct {
		models.OrderResponse
		OrderType string `json:"order_type"`
	}

	var ordersWithType []OrderWithType
	for _, order := range response.Orders {
		// 根据 period_number 判断订单类型
		orderType := "purchase" // 默认为购买订单
		if !utils.IsAllDigits(order.Number) {
			orderType = "group_buy" // 如果不是全数字，则为拼单
		}

		orderWithType := OrderWithType{
			OrderResponse: order,
			OrderType:     orderType,
		}
		ordersWithType = append(ordersWithType, orderWithType)
	}

	// 创建新的响应
	newResponse := map[string]interface{}{
		"orders":     ordersWithType,
		"pagination": response.Pagination,
	}

	utils.Success(c, newResponse)
}

// CreateOrder 创建订单
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req services.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 根据user_id查询用户信息获取uid
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		utils.ErrorWithMessage(c, utils.CodeUserNotFound, "账户已被删除，无法创建订单")
		return
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		utils.ErrorWithMessage(c, utils.CodeAccountLocked, "账户已被禁用，无法创建订单")
		return
	}

	// 使用从数据库获取的uid，而不是从请求参数中获取
	req.Uid = user.Uid

	// 获取当前用户ID作为操作员
	operatorUid := middleware.GetCurrentUser(c)
	operatorUidStr := "system"
	if operatorUid != 0 {
		operatorUidStr = strconv.FormatUint(uint64(operatorUid), 10)
	}

	// 创建订单
	response, err := oc.orderService.CreateOrder(&req, operatorUidStr)
	if err != nil {
		// 记录操作失败
		oc.operationFailureService.RecordOrderCreateFailure(c.Request.Context(), user.Uid, req, gin.H{
			"error": err.Error(),
			"code":  utils.CodeOperationFailed,
		})

		utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "订单创建成功", response)
}

// GetOrderDetail 获取订单详情
func (oc *OrderController) GetOrderDetail(c *gin.Context) {
	var req models.GetOrderDetailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 根据user_id查询用户信息获取uid
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		utils.ErrorWithMessage(c, utils.CodeUserNotFound, "账户已被删除，无法查询订单")
		return
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		utils.ErrorWithMessage(c, utils.CodeAccountLocked, "账户已被禁用，无法查询订单")
		return
	}

	// 获取订单详情
	response, err := oc.orderService.GetOrderDetail(&req, user.Uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetOrderStats 获取订单统计
func (oc *OrderController) GetOrderStats(c *gin.Context) {
	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 根据user_id查询用户信息获取uid
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		utils.ErrorWithMessage(c, utils.CodeUserNotFound, "账户已被删除，无法查询订单统计")
		return
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		utils.ErrorWithMessage(c, utils.CodeAccountLocked, "账户已被禁用，无法查询订单统计")
		return
	}

	// 获取订单统计
	response, err := oc.orderService.GetOrderStats(user.Uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetPeriodList 获取期数列表
func (oc *OrderController) GetPeriodList(c *gin.Context) {
	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	// 根据user_id查询用户信息获取uid
	userRepo := database.NewUserRepository()
	var user models.User
	err := userRepo.FindByID(context.Background(), userID, &user)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, "获取用户信息失败")
		return
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		utils.ErrorWithMessage(c, utils.CodeUserNotFound, "账户已被删除，无法获取期数信息")
		return
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		utils.ErrorWithMessage(c, utils.CodeAccountLocked, "账户已被禁用，无法获取期数信息")
		return
	}

	// 获取期数列表
	response, err := oc.orderService.GetPeriodList()
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetAllOrderList 获取所有订单列表（只需登录即可查看所有订单）
func (oc *OrderController) GetAllOrderList(c *gin.Context) {
	var req models.GetOrderListRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 只需登录即可
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	response, err := oc.orderService.GetAllOrderList(&req)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}
