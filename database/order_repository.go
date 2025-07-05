package database

import (
	"context"
	"gin-fataMorgana/models"
)

type OrderRepository struct {
	*BaseRepository
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *OrderRepository) GetOrdersByStatus(ctx context.Context, status string, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&models.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *OrderRepository) GetUserOrders(ctx context.Context, uid string, page, pageSize int, statusType int) ([]models.Order, int64, error) {
	status := models.GetStatusByType(statusType)
	return r.GetUserOrdersByStatus(ctx, uid, status, page, pageSize)
}

func (r *OrderRepository) GetUserOrdersByStatus(ctx context.Context, uid string, status string, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	query := r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ?", uid)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	return r.Create(ctx, order)
}

func (r *OrderRepository) FindOrderByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
	var order models.Order
	err := r.FindByCondition(ctx, map[string]interface{}{"order_no": orderNo}, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	return r.Update(ctx, order)
}

func (r *OrderRepository) GetOrderStats(ctx context.Context, uid string) (map[string]interface{}, error) {
	var stats struct {
		TotalOrders   int64   `json:"total_orders"`
		PendingOrders int64   `json:"pending_orders"`
		SuccessOrders int64   `json:"success_orders"`
		FailedOrders  int64   `json:"failed_orders"`
		TotalAmount   float64 `json:"total_amount"`
		TotalProfit   float64 `json:"total_profit"`
	}

	err := r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ?", uid).Count(&stats.TotalOrders).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ? AND status = ?", uid, models.OrderStatusPending).Count(&stats.PendingOrders).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ? AND status = ?", uid, models.OrderStatusSuccess).Count(&stats.SuccessOrders).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&models.Order{}).Where("uid = ? AND status = ?", uid, models.OrderStatusFailed).Count(&stats.FailedOrders).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&models.Order{}).Select("COALESCE(SUM(amount), 0)").Where("uid = ?", uid).Scan(&stats.TotalAmount).Error
	if err != nil {
		return nil, err
	}

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

func (r *OrderRepository) CheckUserPeriodExists(ctx context.Context, uid string, periodNumber string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("uid = ? AND period_number = ?", uid, periodNumber).
		Count(&count).Error
	return count > 0, err
}