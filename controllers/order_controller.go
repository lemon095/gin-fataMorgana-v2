package controllers

import (
	"gin-fataMorgana/middleware"
	"gin-fataMorgana/models"
	"gin-fataMorgana/services"
	"gin-fataMorgana/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OrderController 订单控制器
type OrderController struct {
	orderService *services.OrderService
}

// NewOrderController 创建订单控制器实例
func NewOrderController() *OrderController {
	return &OrderController{
		orderService: services.NewOrderService(),
	}
}

// GetOrderList 获取订单列表
func (oc *OrderController) GetOrderList(c *gin.Context) {
	var req models.OrderListRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	uid := strconv.FormatUint(uint64(userID), 10)

	// 获取订单列表
	response, err := oc.orderService.GetUserOrders(&req, uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}

// CreateOrder 创建订单
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req struct {
		BuyAmount   float64 `json:"buy_amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	uid := strconv.FormatUint(uint64(userID), 10)

	// 创建订单
	order, err := oc.orderService.CreateOrder(uid, req.BuyAmount, req.Description)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "订单创建成功", order.ToResponse())
}

// GetOrderDetail 获取订单详情
func (oc *OrderController) GetOrderDetail(c *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	uid := strconv.FormatUint(uint64(userID), 10)

	// 获取订单详情
	order, err := oc.orderService.GetOrderByOrderNo(req.OrderNo)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	// 验证订单是否属于当前用户
	if order.Uid != uid {
		utils.Forbidden(c)
		return
	}

	utils.Success(c, order.ToResponse())
}

// GetOrderStats 获取订单统计
func (oc *OrderController) GetOrderStats(c *gin.Context) {
	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	uid := strconv.FormatUint(uint64(userID), 10)

	// 获取订单统计
	stats, err := oc.orderService.GetUserOrderStats(uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, stats)
}

// GetOrdersByStatus 根据状态获取订单列表
func (oc *OrderController) GetOrdersByStatus(c *gin.Context) {
	var req struct {
		models.OrderListRequest
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	uid := strconv.FormatUint(uint64(userID), 10)

	// 获取订单列表
	response, err := oc.orderService.GetOrdersByStatus(&req.OrderListRequest, uid, req.Status)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetOrdersByDateRange 根据日期范围获取订单列表
func (oc *OrderController) GetOrdersByDateRange(c *gin.Context) {
	var req struct {
		models.OrderListRequest
		StartDate string `json:"start_date" binding:"required"`
		EndDate   string `json:"end_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParamsWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	uid := strconv.FormatUint(uint64(userID), 10)

	// 获取订单列表
	response, err := oc.orderService.GetOrdersByDateRange(&req.OrderListRequest, uid, req.StartDate, req.EndDate)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
} 