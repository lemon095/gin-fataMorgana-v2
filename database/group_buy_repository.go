package database

import (
	"context"
	"math/rand"
	"time"

	"gin-fataMorgana/models"

	"gorm.io/gorm"
)

// GroupBuyRepository 拼单仓库
type GroupBuyRepository struct {
	*BaseRepository
}

// NewGroupBuyRepository 创建拼单仓库实例
func NewGroupBuyRepository() *GroupBuyRepository {
	return &GroupBuyRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// GetActiveGroupBuyDetail 获取活跃拼单详情
// 查询条件：AutoStart为true，截止时间比当前大，完成状态为cancelled
// 可以按时间最近或随机返回一条数据
func (r *GroupBuyRepository) GetActiveGroupBuyDetail(ctx context.Context, random bool) (*models.GroupBuy, error) {
	var groupBuy models.GroupBuy
	
	query := r.db.WithContext(ctx).Where("auto_start = ? AND deadline > ? AND complete = ?", 
		true, time.Now(), "cancelled")
	
	if random {
		// 随机返回一条数据
		// 先获取符合条件的总数
		var count int64
		if err := query.Model(&models.GroupBuy{}).Count(&count).Error; err != nil {
			return nil, err
		}
		
		if count == 0 {
			return nil, gorm.ErrRecordNotFound
		}
		
		// 随机选择一条记录
		offset := rand.Intn(int(count))
		err := query.Offset(offset).Limit(1).First(&groupBuy).Error
		if err != nil {
			return nil, err
		}
	} else {
		// 按时间最近返回一条数据
		err := query.Order("deadline ASC").First(&groupBuy).Error
		if err != nil {
			return nil, err
		}
	}
	
	return &groupBuy, nil
}

// GetGroupBuyByNo 根据拼单编号获取拼单详情
func (r *GroupBuyRepository) GetGroupBuyByNo(ctx context.Context, groupBuyNo string) (*models.GroupBuy, error) {
	var groupBuy models.GroupBuy
	err := r.db.WithContext(ctx).Where("group_buy_no = ?", groupBuyNo).First(&groupBuy).Error
	if err != nil {
		return nil, err
	}
	return &groupBuy, nil
}

// UpdateGroupBuy 更新拼单信息
func (r *GroupBuyRepository) UpdateGroupBuy(ctx context.Context, groupBuy *models.GroupBuy) error {
	return r.db.WithContext(ctx).Save(groupBuy).Error
}

// CreateOrder 创建订单
func (r *GroupBuyRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
} 