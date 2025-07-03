package database

import (
	"context"
	"gin-fataMorgana/models"
)

// OrderRepository 订单仓库
type OrderRepository struct {
	*BaseRepository
}

// NewOrderRepository 创建订单仓库实例
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// CreateOrder 创建订单
func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	return r.Create(ctx, order)
}

// FindOrderByOrderNo 根据订单号查找订单
func (r *OrderRepository) FindOrderByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
	var order models.Order
	err := r.FindByCondition(ctx, map[string]interface{}{"order_no": orderNo}, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetUserOrders 获取用户订单列表
func (r *OrderRepository) GetUserOrders(ctx context.Context, uid string, page, pageSize int, statusType int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// 构建查询条件
	conditions := map[string]interface{}{"uid": uid}

	// 根据状态类型添加状态过滤条件
	status := models.GetStatusByType(statusType)
	if status != "" {
		conditions["status"] = status
	}

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.Order{}).Where(conditions).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据
	query := r.db.WithContext(ctx).Where(conditions).Order("created_at DESC").Offset(offset).Limit(pageSize)
	err = query.Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateOrder 更新订单
func (r *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	return r.Update(ctx, order)
}

// GetOrdersByStatus 根据状态获取订单列表
func (r *OrderRepository) GetOrdersByStatus(ctx context.Context, status string, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.Order{}).Where("status = ?", status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据
	err = r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetExpiredOrders 获取已过期的订单
func (r *OrderRepository) GetExpiredOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).
		Where("status = ? AND expire_time < NOW()", models.OrderStatusPending).
		Find(&orders).Error
	return orders, err
}

// GetOrderStats 获取订单统计信息
func (r *OrderRepository) GetOrderStats(ctx context.Context, uid string) (map[string]interface{}, error) {
	var stats struct {
		TotalOrders   int64   `json:"total_orders"`
		PendingOrders int64   `json:"pending_orders"`
		SuccessOrders int64   `json:"success_orders"`
		FailedOrders  int64   `json:"failed_orders"`
		TotalAmount   float64 `json:"total_amount"`
		TotalProfit   float64 `json:"total_profit"`
	}

	// 总订单数
	err := r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ?", uid).Count(&stats.TotalOrders).Error
	if err != nil {
		return nil, err
	}

	// 待处理订单数
	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ? AND status = ?", uid, models.OrderStatusPending).Count(&stats.PendingOrders).Error
	if err != nil {
		return nil, err
	}

	// 成功订单数
	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ? AND status = ?", uid, models.OrderStatusSuccess).Count(&stats.SuccessOrders).Error
	if err != nil {
		return nil, err
	}

	// 失败订单数
	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ? AND status = ?", uid, models.OrderStatusFailed).Count(&stats.FailedOrders).Error
	if err != nil {
		return nil, err
	}

	// 总金额
	err = r.db.WithContext(ctx).Model(&models.Order{}).Select("COALESCE(SUM(amount), 0)").Where("uid = ?", uid).Scan(&stats.TotalAmount).Error
	if err != nil {
		return nil, err
	}

	// 总利润
	err = r.db.WithContext(ctx).Model(&models.Order{}).Select("COALESCE(SUM(profit_amount), 0)").Where("uid = ?", uid).Scan(&stats.TotalProfit).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_orders":   stats.TotalOrders,
		"pending_orders": stats.PendingOrders,
		"success_orders": stats.SuccessOrders,
		"failed_orders":  stats.FailedOrders,
		"total_amount":   stats.TotalAmount,
		"total_profit":   stats.TotalProfit,
	}, nil
}

// GetUserOrderCount 获取用户订单数量
func (r *OrderRepository) GetUserOrderCount(ctx context.Context, uid string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ?", uid).Count(&count).Error
	return count, err
}

// CheckUserPeriodExists 检查用户是否已购买过指定期号
func (r *OrderRepository) CheckUserPeriodExists(ctx context.Context, uid string, periodNumber string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("uid = ? AND period_number = ?", uid, periodNumber).
		Count(&count).Error

	return count > 0, err
}
