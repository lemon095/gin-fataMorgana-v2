package services

import (
	"context"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
)

// MemberLevelService 用户等级配置表Service
type MemberLevelService struct {
	memberLevelRepo *database.MemberLevelRepository
}

// NewMemberLevelService 创建用户等级配置表Service实例
func NewMemberLevelService(memberLevelRepo *database.MemberLevelRepository) *MemberLevelService {
	return &MemberLevelService{
		memberLevelRepo: memberLevelRepo,
	}
}

// NewMemberLevelRepository 创建用户等级配置表Repository实例
func NewMemberLevelRepository() *database.MemberLevelRepository {
	return database.NewMemberLevelRepository(database.DB)
}

// GetUserLevel 获取用户当前等级配置
func (s *MemberLevelService) GetUserLevel(ctx context.Context, experience int) (*models.MemberLevel, error) {
	return s.memberLevelRepo.GetByExperience(ctx, experience)
}

// GetLevelByLevel 根据等级获取配置
func (s *MemberLevelService) GetLevelByLevel(ctx context.Context, level int) (*models.MemberLevel, error) {
	return s.memberLevelRepo.GetByLevel(ctx, level)
}

// GetAllLevels 获取所有等级配置
func (s *MemberLevelService) GetAllLevels(ctx context.Context) ([]models.MemberLevel, error) {
	return s.memberLevelRepo.GetAllActive(ctx)
}

// CalculateCashback 计算返现金额
func (s *MemberLevelService) CalculateCashback(ctx context.Context, experience int, amount float64) (float64, error) {
	level, err := s.memberLevelRepo.GetByExperience(ctx, experience)
	if err != nil {
		return 0, err
	}

	cashbackAmount := amount * (level.CashbackRatio / 100.0)
	return cashbackAmount, nil
}
