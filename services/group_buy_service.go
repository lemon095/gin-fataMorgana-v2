package services

import (
	"context"
	"math/rand"
	"time"

	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
)

// GroupBuyService 拼单服务
type GroupBuyService struct {
	groupBuyRepo *database.GroupBuyRepository
}

// NewGroupBuyService 创建拼单服务实例
func NewGroupBuyService() *GroupBuyService {
	return &GroupBuyService{
		groupBuyRepo: database.NewGroupBuyRepository(),
	}
}

// GetActiveGroupBuyDetail 获取活跃拼单详情
func (s *GroupBuyService) GetActiveGroupBuyDetail(ctx context.Context, random bool) (*models.GetGroupBuyDetailResponse, error) {
	// 获取拼单数据
	groupBuy, err := s.groupBuyRepo.GetActiveGroupBuyDetail(ctx, random)
	if err != nil {
		// 如果没有找到数据，返回空数据而不是错误
		if err.Error() == "record not found" {
			return &models.GetGroupBuyDetailResponse{}, nil
		}
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取拼单详情失败，请稍后重试", err)
	}

	// 转换为响应格式
	response := groupBuy.ToDetailResponse()
	return &response, nil
}

// JoinGroupBuy 确认参与拼单
func (s *GroupBuyService) JoinGroupBuy(ctx context.Context, groupBuyNo, uid string) (*models.JoinGroupBuyResponse, error) {
	// 1. 根据拼单编号查询拼单信息
	groupBuy, err := s.groupBuyRepo.GetGroupBuyByNo(ctx, groupBuyNo)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, utils.NewAppError(utils.CodeGroupBuyNotFound, "拼单不存在或已被删除", err)
		}
		return nil, utils.NewAppError(utils.CodeDatabaseError, "查询拼单信息失败", err)
	}

	// 2. 检查拼单是否可以参与
	// 检查OrderNo是否有值
	if groupBuy.OrderNo != nil && *groupBuy.OrderNo != "" {
		return nil, utils.NewAppError(utils.CodeGroupBuyOccupied, "该拼单已被其他用户参与", nil)
	}

	// 检查ParticipantUid是否有值
	if groupBuy.ParticipantUid != "" {
		return nil, utils.NewAppError(utils.CodeGroupBuyOccupied, "该拼单已被其他用户参与", nil)
	}

	// 检查截止时间是否已经过了
	if time.Now().After(groupBuy.Deadline) {
		return nil, utils.NewAppError(utils.CodeGroupBuyExpired, "该拼单已超过截止时间", nil)
	}

	// 3. 生成订单编号
	orderNo := utils.GenerateOrderNo()

	// 4. 创建订单数据
	order := &models.Order{
		OrderNo:        orderNo,
		Uid:            uid,
		Amount:         groupBuy.PerPersonAmount,
		ProfitAmount:   0, // 利润金额暂时为0
		LikeCount:      rand.Intn(8000) + 1, // 随机从1-8000生成
		ShareCount:     rand.Intn(8000) + 1,
		FollowCount:    rand.Intn(8000) + 1,
		FavoriteCount:  rand.Intn(8000) + 1,
		LikeStatus:     "pending",
		ShareStatus:    "pending",
		FollowStatus:   "pending",
		FavoriteStatus: "pending",
		Status:         "pending",
		ExpireTime:     time.Now().Add(24 * time.Hour), // 设置24小时后过期
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 5. 保存订单
	err = s.groupBuyRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "创建订单失败，请稍后重试", err)
	}

	// 6. 更新拼单信息
	groupBuy.OrderNo = &orderNo
	groupBuy.Status = "success"
	groupBuy.Deadline = time.Now()
	groupBuy.ParticipantUid = uid
	groupBuy.UpdatedAt = time.Now()

	err = s.groupBuyRepo.UpdateGroupBuy(ctx, groupBuy)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "更新拼单状态失败，请稍后重试", err)
	}

	// 7. 返回订单ID
	response := &models.JoinGroupBuyResponse{
		OrderID: order.ID,
	}

	return response, nil
} 