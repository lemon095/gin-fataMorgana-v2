package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
)

// GroupBuyService 拼单服务
type GroupBuyService struct {
	groupBuyRepo *database.GroupBuyRepository
	walletRepo   *database.WalletRepository
}

// NewGroupBuyService 创建拼单服务实例
func NewGroupBuyService() *GroupBuyService {
	return &GroupBuyService{
		groupBuyRepo: database.NewGroupBuyRepository(),
		walletRepo:   database.NewWalletRepository(),
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

	// 3. 检查用户钱包余额
	wallet, err := s.walletRepo.FindWalletByUid(ctx, uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取钱包失败，请稍后重试", err)
	}

	// 检查钱包状态
	if !wallet.IsActive() {
		return nil, utils.NewAppError(utils.CodeOperationFailed, "钱包已被冻结，无法参与拼单", nil)
	}

	// 检查余额是否足够
	if wallet.Balance < groupBuy.PerPersonAmount {
		return nil, utils.NewAppError(utils.CodeOperationFailed, 
			fmt.Sprintf("余额不足，当前余额: %.2f，拼单金额: %.2f", wallet.Balance, groupBuy.PerPersonAmount), nil)
	}

	// 4. 记录交易前余额
	balanceBefore := wallet.Balance

	// 5. 扣减余额
	if err := wallet.Withdraw(groupBuy.PerPersonAmount); err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "扣减余额失败，请稍后重试", err)
	}

	// 6. 更新钱包
	if err := s.walletRepo.UpdateWallet(ctx, wallet); err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "更新钱包失败，请稍后重试", err)
	}

	// 7. 生成订单编号
	orderNo := utils.GenerateOrderNo()

	// 8. 创建订单数据
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

	// 9. 保存订单
	err = s.groupBuyRepo.CreateOrder(ctx, order)
	if err != nil {
		// 如果创建订单失败，需要回滚扣减的余额
		wallet.Recharge(groupBuy.PerPersonAmount)
		s.walletRepo.UpdateWallet(ctx, wallet)
		return nil, utils.NewAppError(utils.CodeDatabaseError, "创建订单失败，请稍后重试", err)
	}

	// 10. 创建钱包流水记录
	transaction := &models.WalletTransaction{
		TransactionNo:  s.generateTransactionNo(),
		Uid:            uid,
		Type:           models.TransactionTypeGroupBuy,
		Amount:         groupBuy.PerPersonAmount,
		BalanceBefore:  balanceBefore,
		BalanceAfter:   wallet.Balance,
		Status:         models.TransactionStatusSuccess,
		Description:    fmt.Sprintf("参与拼单 %s", groupBuy.GroupBuyNo),
		RelatedOrderNo: orderNo,
		OperatorUid:    "system",
	}

	if err := s.walletRepo.CreateTransaction(ctx, transaction); err != nil {
		// 如果创建交易记录失败，记录日志但不影响拼单参与
		fmt.Printf("创建拼单交易记录失败: %v\n", err)
	}

	// 11. 更新拼单信息
	groupBuy.OrderNo = &orderNo
	groupBuy.Status = "success"
	groupBuy.Deadline = time.Now()
	groupBuy.ParticipantUid = uid
	groupBuy.UpdatedAt = time.Now()
	groupBuy.Complete = "pending" // 进行中

	err = s.groupBuyRepo.UpdateGroupBuy(ctx, groupBuy)
	if err != nil {
		// 如果更新拼单失败，需要回滚扣减的余额和创建的订单
		wallet.Recharge(groupBuy.PerPersonAmount)
		s.walletRepo.UpdateWallet(ctx, wallet)
		// 注意：这里可能需要删除已创建的订单，但为了简化，我们只回滚余额
		return nil, utils.NewAppError(utils.CodeDatabaseError, "更新拼单状态失败，请稍后重试", err)
	}

	// 12. 返回订单ID
	response := &models.JoinGroupBuyResponse{
		OrderID: order.ID,
	}

	return response, nil
}

// generateTransactionNo 生成交易流水号
func (s *GroupBuyService) generateTransactionNo() string {
	// 格式：TX + 年月日 + 时分秒 + 4位随机数
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := utils.RandomString(4)
	return fmt.Sprintf("TX%s%s", timestamp, random)
} 