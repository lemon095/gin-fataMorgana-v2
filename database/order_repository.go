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

// UpdateOrder 更新订单
func (r *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	return r.Update(ctx, order)
}

// GetUserOrders 获取用户订单列表
func (r *OrderRepository) GetUserOrders(ctx context.Context, uid string, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ?", uid).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据
	err = r.db.WithContext(ctx).
		Where("uid = ?", uid).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetOrdersByStatus 根据状态获取用户订单
func (r *OrderRepository) GetOrdersByStatus(ctx context.Context, uid, status string, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ? AND status = ?", uid, status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据
	err = r.db.WithContext(ctx).
		Where("uid = ? AND status = ?", uid, status).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetOrdersByDateRange 根据日期范围获取用户订单
func (r *OrderRepository) GetOrdersByDateRange(ctx context.Context, uid string, startDate, endDate string, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("uid = ? AND DATE(created_at) BETWEEN ? AND ?", uid, startDate, endDate).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据
	err = r.db.WithContext(ctx).
		Where("uid = ? AND DATE(created_at) BETWEEN ? AND ?", uid, startDate, endDate).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetUserOrderStats 获取用户订单统计
func (r *OrderRepository) GetUserOrderStats(ctx context.Context, uid string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总订单数
	var totalOrders int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ?", uid).Count(&totalOrders).Error
	if err != nil {
		return nil, err
	}
	stats["total_orders"] = totalOrders

	// 成功订单数
	var successOrders int64
	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ? AND status = ?", uid, models.OrderStatusSuccess).Count(&successOrders).Error
	if err != nil {
		return nil, err
	}
	stats["success_orders"] = successOrders

	// 待处理订单数
	var pendingOrders int64
	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ? AND status = ?", uid, models.OrderStatusPending).Count(&pendingOrders).Error
	if err != nil {
		return nil, err
	}
	stats["pending_orders"] = pendingOrders

	// 失败订单数
	var failedOrders int64
	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ? AND status = ?", uid, models.OrderStatusFailed).Count(&failedOrders).Error
	if err != nil {
		return nil, err
	}
	stats["failed_orders"] = failedOrders

	// 总买入金额
	var totalBuyAmount float64
	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ?", uid).Select("COALESCE(SUM(buy_amount), 0)").Scan(&totalBuyAmount).Error
	if err != nil {
		return nil, err
	}
	stats["total_buy_amount"] = totalBuyAmount

	// 总利润金额
	var totalProfitAmount float64
	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ?", uid).Select("COALESCE(SUM(profit_amount), 0)").Scan(&totalProfitAmount).Error
	if err != nil {
		return nil, err
	}
	stats["total_profit_amount"] = totalProfitAmount

	return stats, nil
} 