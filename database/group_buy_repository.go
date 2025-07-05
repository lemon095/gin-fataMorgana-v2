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
// 查询条件：截止时间比当前大，状态为进行中
// 可以按时间最近或随机返回一条数据
func (r *GroupBuyRepository) GetActiveGroupBuyDetail(ctx context.Context, random bool) (*models.GroupBuy, error) {
	var groupBuy models.GroupBuy

	query := r.db.WithContext(ctx).Where("deadline > ? AND status = ?",
		time.Now(), "pending")

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

// GetActiveGroupBuyByUid 根据用户ID获取活跃拼单详情
// 查询条件：用户ID匹配，截止时间比当前大
func (r *GroupBuyRepository) GetActiveGroupBuyByUid(ctx context.Context, uid string) (*models.GroupBuy, error) {
	var groupBuy models.GroupBuy

	err := r.db.WithContext(ctx).Where("uid = ? AND deadline > ?", uid, time.Now()).First(&groupBuy).Error
	if err != nil {
		return nil, err
	}

	return &groupBuy, nil
}

// GetNotStartedGroupBuyByUid 根据用户ID获取未开始的拼单详情
// 查询条件：用户ID匹配，状态为not_started，截止时间比当前大
func (r *GroupBuyRepository) GetNotStartedGroupBuyByUid(ctx context.Context, uid string) (*models.GroupBuy, error) {
	var groupBuy models.GroupBuy

	err := r.db.WithContext(ctx).Where("uid = ? AND status = ? AND deadline > ?", uid, "not_started", time.Now()).First(&groupBuy).Error
	if err != nil {
		return nil, err
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

// GetActiveGroupBuys 获取活跃拼单列表
// 查询条件：用户ID匹配，截止时间比当前大，创建时间不超过当前时间，按创建时间倒序排列
func (r *GroupBuyRepository) GetActiveGroupBuys(ctx context.Context, uid string, page, pageSize int) ([]models.GroupBuy, int64, error) {
	var groupBuys []models.GroupBuy
	var total int64

	// 构建查询条件：用户ID匹配，截止时间比当前大，创建时间不超过当前时间
	query := r.db.WithContext(ctx).Where("uid = ? AND deadline > ? AND created_at <= NOW()", uid, time.Now())

	// 获取总数
	err := query.Model(&models.GroupBuy{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据，按创建时间倒序排列
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&groupBuys).Error
	if err != nil {
		return nil, 0, err
	}

	return groupBuys, total, nil
}
