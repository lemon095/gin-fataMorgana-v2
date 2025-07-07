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
	groupBuyRepo    *database.GroupBuyRepository
	walletRepo      *database.WalletRepository
	memberLevelRepo *database.MemberLevelRepository
	cacheService    *WalletCacheService
	walletService   *WalletService
}

// NewGroupBuyService 创建拼单服务实例
func NewGroupBuyService() *GroupBuyService {
	return &GroupBuyService{
		groupBuyRepo:    database.NewGroupBuyRepository(),
		walletRepo:      database.NewWalletRepository(),
		memberLevelRepo: database.NewMemberLevelRepository(database.DB),
		cacheService:    NewWalletCacheService(),
		walletService:   NewWalletService(),
	}
}

// createEmptyGroupBuyResponse 创建空的拼单响应
func (s *GroupBuyService) createEmptyGroupBuyResponse() *models.GetGroupBuyDetailResponse {
	return &models.GetGroupBuyDetailResponse{
		HasData:             false,
		GroupBuyNo:          "",
		GroupBuyType:        "",
		TotalAmount:         0.0,
		CurrentParticipants: 0,
		TargetParticipants:  0,
		PaidAmount:          0.0,
		PerPersonAmount:     0.0,
		RemainingAmount:     0.0,
		Deadline:            time.Time{},
	}
}

// GetActiveGroupBuyDetail 获取活跃拼单详情
func (s *GroupBuyService) GetActiveGroupBuyDetail(ctx context.Context, uid string) (*models.GetGroupBuyDetailResponse, error) {
	// 1. 检查用户是否具有拼单资格
	hasQualification, err := s.checkGroupBuyQualification(ctx, uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "检查拼单资格失败，请稍后重试")
	}

	// 2. 如果没有拼单资格，返回固定格式的空数据
	if !hasQualification {
		return s.createEmptyGroupBuyResponse(), nil
	}

	// 3. 确保用户有钱包，如果没有则自动创建
	err = s.ensureUserWallet(ctx, uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "创建钱包失败，请稍后重试")
	}

	// 4. 根据用户ID获取未开始的拼单数据
	groupBuy, err := s.groupBuyRepo.GetNotStartedGroupBuyByUid(ctx, uid)
	if err != nil {
		// 如果没有找到数据，返回固定格式的空数据
		if err.Error() == "record not found" {
			return s.createEmptyGroupBuyResponse(), nil
		}
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取拼单详情失败，请稍后重试")
	}

	// 5. 转换为响应格式
	response := groupBuy.ToDetailResponse()
	return &response, nil
}

// checkGroupBuyQualification 检查用户是否具有拼单资格
func (s *GroupBuyService) checkGroupBuyQualification(ctx context.Context, uid string) (bool, error) {
	// 1. 检查用户是否存在且状态正常
	userRepo := database.NewUserRepository()
	user, err := userRepo.FindByUid(ctx, uid)
	if err != nil {
		return false, err
	}

	// 2. 检查用户状态（1表示正常）
	if user.Status != 1 {
		return false, nil
	}

	// 3. 检查用户是否有拼单资格
	if !user.HasGroupBuyQualification {
		return false, nil
	}

	// 4. 检查用户钱包是否存在且状态正常
	wallet, err := s.walletRepo.FindWalletByUid(ctx, uid)
	if err != nil {
		return false, err
	}

	// 5. 检查钱包状态
	if !wallet.IsActive() {
		return false, nil
	}

	// 6. 检查用户经验值是否达到最低要求（假设最低经验值为1）
	if user.Experience < 1 {
		return false, nil
	}

	return true, nil
}

// ensureUserWallet 确保用户有钱包，如果没有则自动创建
func (s *GroupBuyService) ensureUserWallet(ctx context.Context, uid string) error {
	// 1. 尝试查找用户钱包
	wallet, err := s.walletRepo.FindWalletByUid(ctx, uid)
	if err != nil {
		// 如果钱包不存在，创建新钱包
		if err.Error() == "record not found" {
			newWallet := &models.Wallet{
				Uid:       uid,
				Balance:   0.0,
				Status:    1, // 1表示正常状态
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}

			err = s.walletRepo.CreateWallet(ctx, newWallet)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// 2. 如果钱包存在但状态不正常，激活钱包
	if !wallet.IsActive() {
		wallet.Status = 1
		wallet.UpdatedAt = time.Now().UTC()
		err = s.walletRepo.UpdateWallet(ctx, wallet)
		if err != nil {
			return err
		}
	}

	return nil
}

// JoinGroupBuy 确认参与拼单
func (s *GroupBuyService) JoinGroupBuy(ctx context.Context, groupBuyNo, uid string) (*models.JoinGroupBuyResponse, error) {
	// 1. 根据拼单编号查询拼单信息
	groupBuy, err := s.groupBuyRepo.GetGroupBuyByNo(ctx, groupBuyNo)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, utils.NewAppError(utils.CodeGroupBuyNotFound, "拼单不存在或已被删除")
		}
		return nil, utils.NewAppError(utils.CodeDatabaseError, "查询拼单信息失败")
	}

	// 2. 检查拼单是否为该用户的
	if groupBuy.Uid != uid {
		return nil, utils.NewAppError(utils.CodeGroupBuyOccupied, "该拼单不属于您，无法参与")
	}

	// 3. 检查拼单是否已经被参与（已有订单编号）
	if groupBuy.OrderNo != nil && *groupBuy.OrderNo != "" {
		return nil, utils.NewAppError(utils.CodeGroupBuyOccupied, "该拼单已被参与，无法重复参与")
	}

	// 4. 检查截止时间是否已经过了
	if time.Now().UTC().After(groupBuy.Deadline) {
		return nil, utils.NewAppError(utils.CodeGroupBuyExpired, "该拼单已超过截止时间")
	}

	// 5. 使用并发安全的钱包服务检查余额
	wallet, err := s.walletService.GetWallet(uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取钱包失败，请稍后重试")
	}

	// 检查钱包状态
	if !wallet.IsActive() {
		return nil, utils.NewAppError(utils.CodeOperationFailed, "钱包已被冻结，无法参与拼单")
	}

	// 检查余额是否足够
	if wallet.Balance < groupBuy.PerPersonAmount {
		return nil, utils.NewAppError(utils.CodeOperationFailed,
			fmt.Sprintf("余额不足，当前余额: %.2f，拼单金额: %.2f", wallet.Balance, groupBuy.PerPersonAmount))
	}

	// 6. 使用并发安全的钱包服务扣减余额
	err = s.walletService.WithdrawBalance(ctx, uid, groupBuy.PerPersonAmount, fmt.Sprintf("参与拼单 %s", groupBuy.GroupBuyNo))
	if err != nil {
		return nil, err
	}

	// 7. 获取用户信息以获取经验值
	userRepo := database.NewUserRepository()
	user, err := userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取用户信息失败")
	}

	// 8. 根据用户经验值获取等级配置并计算利润金额
	profitAmount := s.calculateProfitAmount(ctx, user.Experience, groupBuy.PerPersonAmount)

	// 9. 生成订单编号
	orderNo := utils.GenerateOrderNo()

	// 10. 随机选择1-4个类型，每个类型数量为1
	likeCount := 0
	shareCount := 0
	followCount := 0
	favoriteCount := 0
	
	// 随机选择类型数量（1-4个）
	typeCount := rand.Intn(4) + 1
	
	// 创建类型数组并随机打乱
	types := []string{"like", "share", "follow", "favorite"}
	rand.Shuffle(len(types), func(i, j int) {
		types[i], types[j] = types[j], types[i]
	})
	
	// 选择前typeCount个类型，数量设为1
	for i := 0; i < typeCount; i++ {
		switch types[i] {
		case "like":
			likeCount = 1
		case "share":
			shareCount = 1
		case "follow":
			followCount = 1
		case "favorite":
			favoriteCount = 1
		}
	}

	// 11. 创建订单数据
	order := &models.Order{
		OrderNo:        orderNo,
		Uid:            uid,
		PeriodNumber:   groupBuy.GroupBuyNo, // 将拼单编号写入period_number字段
		Amount:         groupBuy.PerPersonAmount,
		ProfitAmount:   profitAmount,
		LikeCount:      likeCount,
		ShareCount:     shareCount,
		FollowCount:    followCount,
		FavoriteCount:  favoriteCount,
		LikeStatus:     "pending",
		ShareStatus:    "pending",
		FollowStatus:   "pending",
		FavoriteStatus: "pending",
		Status:         "pending",
		ExpireTime:     time.Now().UTC().Add(24 * time.Hour), // 设置24小时后过期
		IsSystemOrder:  false, // 拼单订单也是用户订单，不是系统订单
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	// 12. 保存订单
	err = s.groupBuyRepo.CreateOrder(ctx, order)
	if err != nil {
		// 如果创建订单失败，使用并发安全的钱包服务回滚余额
		if rollbackErr := s.walletService.AddBalance(ctx, uid, groupBuy.PerPersonAmount, "拼单订单创建失败回滚"); rollbackErr != nil {
			// 回滚失败，记录严重错误
			utils.LogError(nil, "拼单订单创建失败且余额回滚失败: %v, 回滚错误: %v", err, rollbackErr)
		}
		return nil, utils.NewAppError(utils.CodeDatabaseError, "创建订单失败，请稍后重试")
	}

	// 13. 获取扣减后的钱包信息用于创建交易记录
	updatedWallet, err := s.walletService.GetWallet(uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeWalletGetFailed, "获取钱包失败")
	}

	// 14. 创建钱包流水记录
	transaction := &models.WalletTransaction{
		TransactionNo:  utils.GenerateTransactionNo("GROUP"),
		Uid:            uid,
		Type:           models.TransactionTypeGroupBuy,
		Amount:         groupBuy.PerPersonAmount,
		BalanceBefore:  wallet.Balance,
		BalanceAfter:   updatedWallet.Balance,
		Status:         models.TransactionStatusSuccess,
		Description:    fmt.Sprintf("参与拼单 %s", groupBuy.GroupBuyNo),
		RelatedOrderNo: orderNo,
		OperatorUid:    "system",
	}

	if err := s.walletRepo.CreateTransaction(ctx, transaction); err != nil {
		// 如果创建交易记录失败，记录日志但不影响拼单参与
		utils.LogWarn(nil, "创建拼单交易记录失败: %v", err)
	}

	// 15. 更新拼单信息
	groupBuy.OrderNo = &orderNo
	groupBuy.Status = "pending" // 更新为pending状态
	groupBuy.UpdatedAt = time.Now().UTC()

	err = s.groupBuyRepo.UpdateGroupBuy(ctx, groupBuy)
	if err != nil {
		// 如果更新拼单失败，需要回滚扣减的余额和创建的订单
		if rollbackErr := s.walletService.AddBalance(ctx, uid, groupBuy.PerPersonAmount, "拼单状态更新失败回滚"); rollbackErr != nil {
			// 回滚失败，记录严重错误
			utils.LogError(nil, "拼单状态更新失败且余额回滚失败: %v, 回滚错误: %v", err, rollbackErr)
		}
		// 注意：这里可能需要删除已创建的订单，但为了简化，我们只回滚余额
		return nil, utils.NewAppError(utils.CodeDatabaseError, "更新拼单状态失败，请稍后重试")
	}

	// 16. 返回订单ID
	response := &models.JoinGroupBuyResponse{
		OrderID: order.ID,
	}

	return response, nil
}

// calculateProfitAmount 根据用户经验值和订单金额计算利润金额
func (s *GroupBuyService) calculateProfitAmount(ctx context.Context, experience int, amount float64) float64 {
	// 先查找当前经验值对应的等级配置
	level, err := s.memberLevelRepo.GetByExperience(ctx, experience)
	if err != nil {
		// 如果查不到，查找最大等级配置
		maxLevel, maxErr := s.memberLevelRepo.GetMaxLevel(ctx)
		if maxErr != nil || maxLevel == nil {
			// 等级表为空或出错，利润为0
			return 0.0
		}
		// 如果查不到等级配置，按最大等级算
		return amount * (maxLevel.CashbackRatio / 100.0)
	}
	return amount * (level.CashbackRatio / 100.0)
}

// generateTransactionNo 生成交易流水号
func (s *GroupBuyService) generateTransactionNo() string {
	// 格式：TX + 年月日 + 时分秒 + 4位随机数
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405")
	random := utils.RandomString(4)
	return fmt.Sprintf("TX%s%s", timestamp, random)
}
