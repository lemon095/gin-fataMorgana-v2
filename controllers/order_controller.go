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
	var req models.GetOrderListRequest

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
	response, err := oc.orderService.GetOrderList(&req, uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}

// CreateOrder 创建订单
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req services.CreateOrderRequest

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
	req.Uid = uid

	// 获取当前用户ID作为操作员
	operatorUid := middleware.GetCurrentUser(c)
	operatorUidStr := "system"
	if operatorUid != 0 {
		operatorUidStr = strconv.FormatUint(uint64(operatorUid), 10)
	}

	// 创建订单
	response, err := oc.orderService.CreateOrder(&req, operatorUidStr)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeOperationFailed, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "订单创建成功", response)
}

// GetOrderDetail 获取订单详情
func (oc *OrderController) GetOrderDetail(c *gin.Context) {
	var req models.GetOrderDetailRequest

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
	response, err := oc.orderService.GetOrderDetail(&req, uid)
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

	uid := strconv.FormatUint(uint64(userID), 10)

	// 获取订单统计
	response, err := oc.orderService.GetOrderStats(uid)
	if err != nil {
		utils.ErrorWithMessage(c, utils.CodeDatabaseError, err.Error())
		return
	}

	utils.Success(c, response)
}

 