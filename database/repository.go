package database

import (
	"context"
	"errors"
	"fmt"

	"gin-fataMorgana/models"

	"gorm.io/gorm"
)

// Repository 基础仓库接口
type Repository interface {
	Create(ctx context.Context, model interface{}) error
	FindByID(ctx context.Context, id uint, model interface{}) error
	FindByCondition(ctx context.Context, condition interface{}, model interface{}) error
	Update(ctx context.Context, model interface{}) error
	Delete(ctx context.Context, id uint) error
	SoftDelete(ctx context.Context, id uint) error
	Exists(ctx context.Context, condition interface{}) (bool, error)
	Count(ctx context.Context, condition interface{}) (int64, error)
}

// BaseRepository 基础仓库实现
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓库
func NewBaseRepository() *BaseRepository {
	return &BaseRepository{db: DB}
}

// Create 创建记录
func (r *BaseRepository) Create(ctx context.Context, model interface{}) error {
	return r.db.WithContext(ctx).Create(model).Error
}

// FindByID 根据ID查找记录
func (r *BaseRepository) FindByID(ctx context.Context, id uint, model interface{}) error {
	result := r.db.WithContext(ctx).First(model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("记录不存在")
		}
		return result.Error
	}
	return nil
}

// FindByCondition 根据条件查找记录
func (r *BaseRepository) FindByCondition(ctx context.Context, condition interface{}, model interface{}) error {
	result := r.db.WithContext(ctx).Where(condition).First(model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("记录不存在")
		}
		return result.Error
	}
	return nil
}

// Update 更新记录
func (r *BaseRepository) Update(ctx context.Context, model interface{}) error {
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete 硬删除记录
func (r *BaseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// SoftDelete 软删除记录
func (r *BaseRepository) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// Exists 检查记录是否存在
func (r *BaseRepository) Exists(ctx context.Context, condition interface{}) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where(condition).Count(&count).Error
	return count > 0, err
}

// Count 统计记录数量
func (r *BaseRepository) Count(ctx context.Context, condition interface{}) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where(condition).Count(&count).Error
	return count, err
}

// UserRepository 用户仓库
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository 创建用户仓库
func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// FindByEmail 根据邮箱查找用户（包括软删除检查）
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmailIncludeDeleted 根据邮箱查找用户（包括已删除的）
func (r *UserRepository) FindByEmailIncludeDeleted(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Unscoped().Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// EmailExists 检查邮箱是否存在（不包括已删除的）
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// EmailExistsIncludeDeleted 检查邮箱是否存在（包括已删除的）
func (r *UserRepository) EmailExistsIncludeDeleted(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// IsUserDeleted 检查用户是否已删除
func (r *UserRepository) IsUserDeleted(ctx context.Context, email string) (bool, error) {
	var user models.User
	err := r.db.WithContext(ctx).Unscoped().Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return user.DeletedAt != nil, nil
}

// IsUserDisabled 检查用户是否被禁用
func (r *UserRepository) IsUserDisabled(ctx context.Context, email string) (bool, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		return false, err
	}
	return user.Status == 0, nil
}

// FindByUid 根据UID查找用户
func (r *UserRepository) FindByUid(ctx context.Context, uid string) (*models.User, error) {
	var user models.User
	err := r.FindByCondition(ctx, map[string]interface{}{"uid": uid}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.FindByCondition(ctx, map[string]interface{}{"username": username}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UidExists 检查uid是否存在
func (r *UserRepository) UidExists(ctx context.Context, uid string) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"uid": uid})
}

// CheckEmailExists 检查邮箱是否存在
func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	// 先尝试从缓存获取邮箱检查结果
	if cachedExists, hit, err := GetCachedEmailExists(ctx, email); err == nil && hit {
		return cachedExists, nil
	} else {
		// 缓存未命中，查询数据库
		var count int64
		err = r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error
		if err != nil {
			return false, err
		}
		emailExists := count > 0
		// 缓存结果
		CacheEmailExists(ctx, email, emailExists)
		return emailExists, nil
	}
}

// GetActiveUsers 获取活跃用户列表
func (r *UserRepository) GetActiveUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("status = ?", 1).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, err
}

// GetUsersByExperience 根据经验值范围获取用户
func (r *UserRepository) GetUsersByExperience(ctx context.Context, minExp, maxExp int) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("experience BETWEEN ? AND ?", minExp, maxExp).
		Order("experience DESC").
		Find(&users).Error
	return users, err
}

// GetUsersByCreditScore 根据信用分范围获取用户
func (r *UserRepository) GetUsersByCreditScore(ctx context.Context, minScore, maxScore int) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("credit_score BETWEEN ? AND ?", minScore, maxScore).
		Order("credit_score DESC").
		Find(&users).Error
	return users, err
}

// SearchUsers 搜索用户
func (r *UserRepository) SearchUsers(ctx context.Context, keyword string, limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, err
}

// GetUserStats 获取用户统计信息
func (r *UserRepository) GetUserStats(ctx context.Context) (map[string]interface{}, error) {
	var stats struct {
		TotalUsers     int64   `json:"total_users"`
		ActiveUsers    int64   `json:"active_users"`
		InactiveUsers  int64   `json:"inactive_users"`
		AvgExperience  float64 `json:"avg_experience"`
		AvgCreditScore float64 `json:"avg_credit_score"`
	}

	// 总用户数
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// 活跃用户数
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("status = ?", 1).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}

	// 非活跃用户数
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("status = ?", 0).Count(&stats.InactiveUsers).Error; err != nil {
		return nil, err
	}

	// 平均经验值
	if err := r.db.WithContext(ctx).Model(&models.User{}).Select("AVG(experience)").Scan(&stats.AvgExperience).Error; err != nil {
		return nil, err
	}

	// 平均信用分
	if err := r.db.WithContext(ctx).Model(&models.User{}).Select("AVG(credit_score)").Scan(&stats.AvgCreditScore).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_users":      stats.TotalUsers,
		"active_users":     stats.ActiveUsers,
		"inactive_users":   stats.InactiveUsers,
		"avg_experience":   stats.AvgExperience,
		"avg_credit_score": stats.AvgCreditScore,
	}, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// FindByPhone 根据手机号查找用户
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckPhoneExists 检查手机号是否存在
func (r *UserRepository) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
