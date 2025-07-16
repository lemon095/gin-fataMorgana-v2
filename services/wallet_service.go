package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"
)

// WalletService 统一的钱包服务（支持跨进程并发安全）
type WalletService struct {
	walletRepo *database.WalletRepository
	// 缓存服务
	cacheService *WalletCacheService
}

func NewWalletService() *WalletService {
	return &WalletService{
		walletRepo:   database.NewWalletRepository(),
		cacheService: NewWalletCacheService(),
	}
}

// 生成分布式锁的Key
func (s *WalletService) generateLockKey(uid string) string {
	return utils.RedisKeys.GenerateWalletLockKey(uid)
}

// 生成锁的值（用于标识锁的持有者）
func (s *WalletService) generateLockValue() string {
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), utils.RandomString(8))
}

// 尝试获取分布式锁（支持等待和重试）
func (s *WalletService) acquireLock(ctx context.Context, uid string, timeout time.Duration) (string, error) {
	lockKey := s.generateLockKey(uid)
	lockValue := s.generateLockValue()

	// 使用SET NX EX命令实现分布式锁
	// NX: 只在key不存在时设置
	// EX: 设置过期时间（秒）
	success, err := database.GetGlobalRedisHelper().SetNX(ctx, lockKey, lockValue, timeout)
	if err != nil {
		return "", utils.NewAppError(utils.CodeRedisError, "获取分布式锁失败")
	}

	if !success {
		return "", utils.NewAppError(utils.CodeSystemBusy, "系统繁忙，请稍后重试")
	}

	return lockValue, nil
}

// 尝试获取分布式锁（支持等待和重试）
func (s *WalletService) acquireLockWithRetry(ctx context.Context, uid string, timeout time.Duration, maxRetries int, retryDelay time.Duration) (string, error) {
	lockKey := s.generateLockKey(uid)
	lockValue := s.generateLockValue()

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// 使用SET NX EX命令实现分布式锁
		success, err := database.GetGlobalRedisHelper().SetNX(ctx, lockKey, lockValue, timeout)
		if err != nil {
			return "", utils.NewAppError(utils.CodeRedisError, "获取分布式锁失败")
		}

		if success {
			return lockValue, nil
		}

		// 如果不是最后一次尝试，则等待后重试
		if attempt < maxRetries {
			// 检查上下文是否已取消
			select {
			case <-ctx.Done():
				return "", utils.NewAppError(utils.CodeRequestTimeout, "请求超时")
			case <-time.After(retryDelay):
				// 继续重试
			}

			// 指数退避：每次重试延迟时间递增
			retryDelay = time.Duration(float64(retryDelay) * 1.5)
		}
	}

	return "", utils.NewAppError(utils.CodeSystemBusy, "系统繁忙，请稍后重试")
}

// 释放分布式锁
func (s *WalletService) releaseLock(ctx context.Context, uid, lockValue string) error {
	lockKey := s.generateLockKey(uid)

	// 使用Lua脚本确保原子性释放锁
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	// 使用Redis助手的Lock/Unlock方法
	// 注意：这里需要直接使用Redis客户端执行Lua脚本，因为需要原子性检查
	result, err := database.RedisClient.Eval(ctx, luaScript, []string{lockKey}, []string{lockValue}).Result()
	if err != nil {
		return utils.NewAppError(utils.CodeRedisError, "释放分布式锁失败")
	}

	// result为1表示成功释放，0表示锁不存在或不是当前持有者
	if result == 0 {
		return utils.NewAppError(utils.CodeSystemBusy, "锁已过期或被其他进程持有")
	}

	return nil
}

// 原子性余额操作（支持跨进程并发安全，带重试机制）
func (s *WalletService) AtomicBalanceOperation(ctx context.Context, uid string, operation func(*models.Wallet) error) error {
	return s.AtomicBalanceOperationWithRetry(ctx, uid, operation, 3, 100*time.Millisecond)
}

// 原子性余额操作（支持跨进程并发安全，可配置重试）
func (s *WalletService) AtomicBalanceOperationWithRetry(ctx context.Context, uid string, operation func(*models.Wallet) error, maxRetries int, retryDelay time.Duration) error {
	if uid == "" {
		return utils.NewAppError(utils.CodeInvalidParams, "用户ID不能为空")
	}

	// 1. 获取分布式锁（支持重试）
	lockValue, err := s.acquireLockWithRetry(ctx, uid, 30*time.Second, maxRetries, retryDelay)
	if err != nil {
		return err
	}

	// 2. 确保锁会被释放
	defer func() {
		if releaseErr := s.releaseLock(ctx, uid, lockValue); releaseErr != nil {
			utils.LogWarn(nil, "释放分布式锁失败: %v", releaseErr)
		}
	}()

	// 3. 从数据库获取最新钱包数据
	wallet, err := s.walletRepo.FindWalletByUid(ctx, uid)
	if err != nil {
		return utils.NewAppError(utils.CodeWalletGetFailed, "获取钱包数据失败")
	}

	// 4. 记录操作前余额
	balanceBefore := wallet.Balance

	// 5. 执行余额操作
	if err := operation(wallet); err != nil {
		return err
	}

	// 6. 更新钱包
	wallet.UpdatedAt = time.Now().UTC()
	if err := s.walletRepo.UpdateWallet(ctx, wallet); err != nil {
		return utils.NewAppError(utils.CodeWalletUpdateFailed, "更新钱包失败")
	}

	// 7. 立即更新缓存
	if cacheErr := s.cacheService.UpdateWalletBalanceOnEvent(ctx, uid, wallet.Balance); cacheErr != nil {
		// 缓存更新失败不影响主流程，只记录日志
		utils.LogWarn(nil, "更新钱包余额缓存失败: %v", cacheErr)
	}

	// 8. 记录操作日志
	utils.LogInfo(nil, "钱包余额操作成功 - UID: %s, 操作前: %.2f, 操作后: %.2f",
		uid, balanceBefore, wallet.Balance)

	return nil
}

// 扣减余额（跨进程并发安全）
func (s *WalletService) WithdrawBalance(ctx context.Context, uid string, amount float64, description string) error {
	if amount <= 0 {
		return utils.NewAppError(utils.CodeInvalidParams, "扣减金额必须大于0")
	}

	return s.AtomicBalanceOperation(ctx, uid, func(wallet *models.Wallet) error {
		// 检查余额是否足够
		if wallet.Balance < amount {
			return utils.NewAppError(utils.CodeBalanceInsufficient,
				fmt.Sprintf("余额不足，当前余额: %.2f，扣减金额: %.2f", wallet.Balance, amount))
		}

		// 检查钱包是否可以操作
		if !wallet.CanOperate() {
			return utils.NewAppError(utils.CodeWalletFrozenWithdraw, "钱包已被冻结，无法扣减余额")
		}

		// 扣减余额
		return wallet.Withdraw(amount)
	})
}

// 增加余额（跨进程并发安全）
func (s *WalletService) AddBalance(ctx context.Context, uid string, amount float64, description string) error {
	if amount <= 0 {
		return utils.NewAppError(utils.CodeInvalidParams, "增加金额必须大于0")
	}

	return s.AtomicBalanceOperation(ctx, uid, func(wallet *models.Wallet) error {
		// 检查钱包是否被冻结（状态0）
		if wallet.IsFrozen() {
			return utils.NewAppError(utils.CodeWalletFrozenRecharge, "钱包已被冻结，无法增加余额")
		}

		// 状态2（无法提现）不影响充值操作
		// 增加余额
		wallet.Recharge(amount)
		return nil
	})
}

// 转账操作（跨进程并发安全）
func (s *WalletService) TransferBalance(ctx context.Context, fromUid, toUid string, amount float64, description string) error {
	if amount <= 0 {
		return utils.NewAppError(utils.CodeInvalidParams, "转账金额必须大于0")
	}

	if fromUid == toUid {
		return utils.NewAppError(utils.CodeInvalidParams, "不能向自己转账")
	}

	// 获取两个用户的分布式锁，按UID排序避免死锁
	var firstLockValue, secondLockValue string
	var firstUid, secondUid string

	if fromUid < toUid {
		firstUid = fromUid
		secondUid = toUid
	} else {
		firstUid = toUid
		secondUid = fromUid
	}

	// 1. 获取第一个锁
	var err error
	firstLockValue, err = s.acquireLock(ctx, firstUid, 30*time.Second)
	if err != nil {
		return err
	}

	// 2. 获取第二个锁
	secondLockValue, err = s.acquireLock(ctx, secondUid, 30*time.Second)
	if err != nil {
		// 释放第一个锁
		s.releaseLock(ctx, firstUid, firstLockValue)
		return err
	}

	// 3. 确保锁会被释放
	defer func() {
		s.releaseLock(ctx, firstUid, firstLockValue)
		s.releaseLock(ctx, secondUid, secondLockValue)
	}()

	// 4. 获取转出方钱包
	fromWallet, err := s.walletRepo.FindWalletByUid(ctx, fromUid)
	if err != nil {
		return utils.NewAppError(utils.CodeWalletGetFailed, "获取转出方钱包失败")
	}

	// 5. 获取转入方钱包
	toWallet, err := s.walletRepo.FindWalletByUid(ctx, toUid)
	if err != nil {
		return utils.NewAppError(utils.CodeWalletGetFailed, "获取转入方钱包失败")
	}

	// 6. 检查转出方余额
	if fromWallet.Balance < amount {
		return utils.NewAppError(utils.CodeBalanceInsufficient,
			fmt.Sprintf("转出方余额不足，当前余额: %.2f，转账金额: %.2f", fromWallet.Balance, amount))
	}

	// 7. 检查钱包状态
	if fromWallet.IsFrozen() {
		return utils.NewAppError(utils.CodeWalletFrozenWithdraw, "转出方钱包已被冻结")
	}
	if toWallet.IsFrozen() {
		return utils.NewAppError(utils.CodeWalletFrozenRecharge, "转入方钱包已被冻结")
	}
	// 状态2（无法提现）不影响转账操作

	// 8. 执行转账
	fromWallet.Withdraw(amount)
	toWallet.Recharge(amount)

	// 9. 更新钱包
	fromWallet.UpdatedAt = time.Now().UTC()
	toWallet.UpdatedAt = time.Now().UTC()

	if err := s.walletRepo.UpdateWallet(ctx, fromWallet); err != nil {
		return utils.NewAppError(utils.CodeWalletUpdateFailed, "更新转出方钱包失败")
	}

	if err := s.walletRepo.UpdateWallet(ctx, toWallet); err != nil {
		return utils.NewAppError(utils.CodeWalletUpdateFailed, "更新转入方钱包失败")
	}

	// 10. 更新缓存
	if cacheErr := s.cacheService.UpdateWalletBalanceOnEvent(ctx, fromUid, fromWallet.Balance); cacheErr != nil {
		utils.LogWarn(nil, "更新转出方钱包缓存失败: %v", cacheErr)
	}

	if cacheErr := s.cacheService.UpdateWalletBalanceOnEvent(ctx, toUid, toWallet.Balance); cacheErr != nil {
		utils.LogWarn(nil, "更新转入方钱包缓存失败: %v", cacheErr)
	}

	// 11. 记录操作日志
	utils.LogInfo(nil, "转账操作成功 - 从: %s, 到: %s, 金额: %.2f", fromUid, toUid, amount)

	return nil
}

// BalanceOperation 余额操作结构体
type BalanceOperation struct {
	UID    string  `json:"uid"`
	Type   string  `json:"type"` // "withdraw" 或 "add"
	Amount float64 `json:"amount"`
}

// 批量余额操作（跨进程并发安全）
func (s *WalletService) BatchBalanceOperation(ctx context.Context, operations []BalanceOperation) error {
	if len(operations) == 0 {
		return nil
	}

	// 按用户分组操作
	userOperations := make(map[string][]BalanceOperation)
	for _, op := range operations {
		userOperations[op.UID] = append(userOperations[op.UID], op)
	}

	// 并发处理不同用户的操作
	var wg sync.WaitGroup
	errChan := make(chan error, len(userOperations))

	for uid, ops := range userOperations {
		wg.Add(1)
		go func(userUid string, userOps []BalanceOperation) {
			defer wg.Done()
			if err := s.executeUserOperations(ctx, userUid, userOps); err != nil {
				errChan <- err
			}
		}(uid, ops)
	}

	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// 执行单个用户的批量操作
func (s *WalletService) executeUserOperations(ctx context.Context, uid string, operations []BalanceOperation) error {
	return s.AtomicBalanceOperationWithRetry(ctx, uid, func(wallet *models.Wallet) error {
		for _, op := range operations {
			switch op.Type {
			case "withdraw":
				if wallet.Balance < op.Amount {
					return utils.NewAppError(utils.CodeBalanceInsufficient,
						fmt.Sprintf("余额不足，当前余额: %.2f，扣减金额: %.2f", wallet.Balance, op.Amount))
				}
				if wallet.IsFrozen() {
					return utils.NewAppError(utils.CodeWalletFrozenWithdraw, "钱包已被冻结，无法扣减余额")
				}
				// 状态2（无法提现）不影响扣减操作
				wallet.Withdraw(op.Amount)
			case "add":
				if wallet.IsFrozen() {
					return utils.NewAppError(utils.CodeWalletFrozenRecharge, "钱包已被冻结，无法增加余额")
				}
				// 状态2（无法提现）不影响增加操作
				wallet.Recharge(op.Amount)
			default:
				return utils.NewAppError(utils.CodeInvalidParams, "无效的操作类型")
			}
		}
		return nil
	}, 5, 200*time.Millisecond) // 批量操作使用更多重试次数和更长延迟
}

// 批量发奖优化方法（专门用于批量发奖场景）
func (s *WalletService) BatchAddBalanceForRewards(ctx context.Context, rewards []struct {
	UID    string  `json:"uid"`
	Amount float64 `json:"amount"`
	Desc   string  `json:"description"`
}) error {
	if len(rewards) == 0 {
		return nil
	}

	// 按用户分组奖励
	userRewards := make(map[string][]struct {
		UID    string  `json:"uid"`
		Amount float64 `json:"amount"`
		Desc   string  `json:"description"`
	})

	for _, reward := range rewards {
		userRewards[reward.UID] = append(userRewards[reward.UID], reward)
	}

	// 并发处理不同用户的奖励
	var wg sync.WaitGroup
	errChan := make(chan error, len(userRewards))

	for uid, userRewardList := range userRewards {
		wg.Add(1)
		go func(userUid string, rewardList []struct {
			UID    string  `json:"uid"`
			Amount float64 `json:"amount"`
			Desc   string  `json:"description"`
		}) {
			defer wg.Done()
			if err := s.executeUserRewards(ctx, userUid, rewardList); err != nil {
				errChan <- err
			}
		}(uid, userRewardList)
	}

	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// 执行单个用户的批量奖励
func (s *WalletService) executeUserRewards(ctx context.Context, uid string, rewards []struct {
	UID    string  `json:"uid"`
	Amount float64 `json:"amount"`
	Desc   string  `json:"description"`
}) error {
	return s.AtomicBalanceOperationWithRetry(ctx, uid, func(wallet *models.Wallet) error {
		// 检查钱包是否被冻结（状态0）
		if wallet.IsFrozen() {
			return utils.NewAppError(utils.CodeWalletFrozenRecharge, "钱包已被冻结，无法添加奖励")
		}
		// 状态2（无法提现）不影响奖励操作

		// 批量增加余额
		totalAmount := 0.0
		for _, reward := range rewards {
			if reward.Amount <= 0 {
				return utils.NewAppError(utils.CodeInvalidParams, "奖励金额必须大于0")
			}
			totalAmount += reward.Amount
		}

		// 一次性增加总金额
		wallet.Recharge(totalAmount)
		return nil
	}, 10, 500*time.Millisecond) // 批量发奖使用更多重试次数和更长延迟
}

// 获取锁状态信息（用于监控和调试）
func (s *WalletService) GetLockStatus(ctx context.Context, uid string) (map[string]interface{}, error) {
	lockKey := s.generateLockKey(uid)

	// 检查锁是否存在
	exists, err := database.GetGlobalRedisHelper().Exists(ctx, lockKey)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeRedisError, "检查锁状态失败")
	}

	status := map[string]interface{}{
		"uid":         uid,
		"lock_key":    lockKey,
		"lock_exists": exists > 0,
		"timestamp":   time.Now().Unix(),
	}

	if exists > 0 {
		// 获取锁的剩余过期时间
		ttl, err := database.GetGlobalRedisHelper().TTL(ctx, lockKey)
		if err == nil {
			status["ttl_seconds"] = ttl.Seconds()
		}

		// 获取锁的值（持有者标识）
		lockValue, err := database.GetGlobalRedisHelper().Get(ctx, lockKey)
		if err == nil {
			status["lock_value"] = lockValue
		}
	}

	return status, nil
}

// 强制释放锁（仅用于紧急情况）
func (s *WalletService) ForceReleaseLock(ctx context.Context, uid string) error {
	lockKey := s.generateLockKey(uid)

	// 直接删除锁
	err := database.GetGlobalRedisHelper().Del(ctx, lockKey)
	if err != nil {
		return utils.NewAppError(utils.CodeRedisError, "强制释放锁失败")
	}

	utils.LogWarn(nil, "强制释放钱包锁 - UID: %s", uid)
	return nil
}

// 提现请求结构体
type WithdrawRequest struct {
	Uid         string  `json:"uid"` // 移除 binding:"required"，uid 从当前登录用户获取
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// GetUserTransactionsRequest 获取用户交易记录请求
type GetUserTransactionsRequest struct {
	Uid      string `json:"uid" binding:"required"`
	Page     int    `json:"page" binding:"required,min=1"`
	PageSize int    `json:"page_size" binding:"required,min=1,max=100"`
	Type     string `json:"type"` // 可选，交易类型过滤
}

// GetTransactionDetailRequest 获取交易详情请求
type GetTransactionDetailRequest struct {
	TransactionNo string `json:"transaction_no" binding:"required"`
}

// GetUserTransactions 获取用户交易记录
func (s *WalletService) GetUserTransactions(req *GetUserTransactionsRequest) (*models.GetTransactionsResponse, error) {
	ctx := context.Background()

	// 获取交易记录
	transactions, total, err := s.walletRepo.GetTransactionsByUid(ctx, req.Uid, req.Page, req.PageSize, req.Type)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取交易记录失败")
	}

	// 构建响应
	response := &models.GetTransactionsResponse{
		Transactions: transactions,
		Total:        total,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}

	return response, nil
}

// GetWallet 获取钱包信息
func (s *WalletService) GetWallet(uid string) (*models.Wallet, error) {
	ctx := context.Background()

	// 先从缓存获取
	wallet, err := s.cacheService.GetWalletBalanceWithCache(ctx, uid)
	if err == nil && wallet != nil {
		return wallet, nil
	}

	// 缓存未命中，从数据库获取
	wallet, err = s.walletRepo.FindWalletByUid(ctx, uid)
	if err != nil {
		// 如果钱包不存在，自动创建钱包
		if err.Error() == "record not found" {
			utils.LogInfo(nil, "用户钱包不存在，自动创建钱包 - UID: %s", uid)
			wallet, err = s.CreateWallet(uid)
			if err != nil {
				utils.LogError(nil, "获取钱包失败1: %v", err)
				return nil, utils.NewAppError(utils.CodeWalletGetFailed, "获取钱包信息失败")
			}
			return wallet, nil
		}
		utils.LogError(nil, "获取钱包失败2: %v", err)
		return nil, utils.NewAppError(utils.CodeWalletGetFailed, "获取钱包信息失败")
	}

	// 更新缓存
	if cacheErr := s.cacheService.UpdateWalletBalanceOnEvent(ctx, uid, wallet.Balance); cacheErr != nil {
		utils.LogWarn(nil, "更新钱包余额缓存失败: %v", cacheErr)
	}

	return wallet, nil
}

// CreateWallet 创建钱包
func (s *WalletService) CreateWallet(uid string) (*models.Wallet, error) {
	ctx := context.Background()

	// 检查钱包是否已存在
	existingWallet, err := s.walletRepo.FindWalletByUid(ctx, uid)
	if err == nil && existingWallet != nil {
		return existingWallet, nil // 钱包已存在，直接返回
	}

	// 创建新钱包
	wallet := &models.Wallet{
		Uid:       uid,
		Balance:   0,
		Status:    1, // 1表示正常状态
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.walletRepo.CreateWallet(ctx, wallet); err != nil {
		return nil, utils.NewAppError(utils.CodeWalletCreateFailed, "创建钱包失败")
	}

	// 更新缓存
	if cacheErr := s.cacheService.UpdateWalletBalanceOnEvent(ctx, uid, wallet.Balance); cacheErr != nil {
		utils.LogWarn(nil, "更新钱包余额缓存失败: %v", cacheErr)
	}

	return wallet, nil
}

// Recharge 充值申请
func (s *WalletService) Recharge(uid string, amount float64, description string) (string, error) {
	ctx := context.Background()

	// 生成交易号
	transactionNo := utils.GenerateTransactionNo("RECHARGE")

	// 创建充值交易记录
	transaction := &models.WalletTransaction{
		TransactionNo: transactionNo,
		Uid:           uid,
		Type:          models.TransactionTypeRecharge,
		Amount:        amount,
		BalanceBefore: 0, // 充值时余额为0，实际余额在审核通过后更新
		BalanceAfter:  0, // 充值时余额为0，实际余额在审核通过后更新
		Description:   description,
		Status:        models.TransactionStatusPending,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := s.walletRepo.CreateTransaction(ctx, transaction); err != nil {
		return "", utils.NewAppError(utils.CodeTransactionCreateFailed, "创建充值申请失败")
	}

	return transactionNo, nil
}

// AddProfit 添加利润
func (s *WalletService) AddProfit(ctx context.Context, uid string, amount float64, description string) error {
	if uid == "" || amount <= 0 {
		return utils.NewAppError(utils.CodeInvalidParams, "参数无效")
	}

	return s.AtomicBalanceOperation(ctx, uid, func(wallet *models.Wallet) error {
		// 检查钱包状态
		if !wallet.IsActive() {
			return utils.NewAppError(utils.CodeWalletFrozenRecharge, "钱包已被冻结，无法添加利润")
		}

		// 增加余额（利润）
		wallet.Recharge(amount)
		return nil
	})
}

// CreateProfitTransaction 创建利润交易记录
func (s *WalletService) CreateProfitTransaction(ctx context.Context, uid string, amount float64, description string, relatedOrderNo string) (string, error) {
	// 生成交易号
	transactionNo := utils.GenerateTransactionNo("PROFIT")

	// 获取当前钱包余额
	wallet, err := s.GetWallet(uid)
	if err != nil {
		return "", utils.NewAppError(utils.CodeWalletGetFailed, "获取钱包信息失败")
	}

	// 创建利润交易记录
	transaction := &models.WalletTransaction{
		TransactionNo:  transactionNo,
		Uid:            uid,
		Type:           models.TransactionTypeProfit,
		Amount:         amount,
		BalanceBefore:  wallet.Balance,
		BalanceAfter:   wallet.Balance + amount,
		Description:    description,
		RelatedOrderNo: relatedOrderNo,
		Status:         models.TransactionStatusSuccess, // 利润直接成功
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	if err := s.walletRepo.CreateTransaction(ctx, transaction); err != nil {
		return "", utils.NewAppError(utils.CodeTransactionCreateFailed, "创建利润交易记录失败")
	}

	// 更新钱包余额
	if err := s.AddProfit(ctx, uid, amount, description); err != nil {
		return "", err
	}

	return transactionNo, nil
}

// RequestWithdraw 申请提现
func (s *WalletService) RequestWithdraw(req *WithdrawRequest, userUid string) (*models.WithdrawResponse, error) {
	ctx := context.Background()

	// 验证用户权限
	if req.Uid != userUid {
		return nil, utils.NewAppError(utils.CodeForbidden, "只能操作自己的钱包")
	}

	// 检查用户是否已绑定银行卡
	userRepo := database.NewUserRepository()
	user, err := userRepo.FindByUid(ctx, req.Uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeUserNotFound, "用户不存在")
	}

	// 检查银行卡信息是否为空或为默认空值
	if user.BankCardInfo == "" || user.BankCardInfo == "{\"card_number\":\"\",\"card_holder\":\"\",\"bank_name\":\"\",\"card_type\":\"\"}" {
		return nil, utils.NewAppError(utils.CodeBankCardNotBound, "请先绑定银行卡后再进行提现操作")
	}

	// 解析银行卡信息，验证是否包含有效的银行卡号
	var bankCardInfo models.BankCardInfo
	if err := utils.JSONToStruct(user.BankCardInfo, &bankCardInfo); err != nil {
		return nil, utils.NewAppError(utils.CodeBankCardFormatError, "银行卡信息格式错误")
	}

	// 检查银行卡号是否为空
	if bankCardInfo.CardNumber == "" {
		return nil, utils.NewAppError(utils.CodeBankCardNotBound, "请先绑定银行卡后再进行提现操作")
	}

	// 使用原子操作立即扣减余额
	var balanceAfter float64
	err = s.AtomicBalanceOperation(ctx, req.Uid, func(wallet *models.Wallet) error {
		// 检查余额是否足够
		if wallet.Balance < req.Amount {
			return utils.NewAppError(utils.CodeBalanceInsufficient,
				fmt.Sprintf("余额不足，当前余额: %.2f，提现金额: %.2f", wallet.Balance, req.Amount))
		}

		// 检查钱包是否可以提现
		if !wallet.CanWithdraw() {
			if wallet.IsFrozen() {
				return utils.NewAppError(utils.CodeWalletFrozenWithdraw, "钱包已被冻结，无法提现")
			}
			if wallet.IsNoWithdraw() {
				return utils.NewAppError(utils.CodeWalletNoWithdraw, "钱包暂时无法提现")
			}
		}

		// 立即扣减余额
		if err := wallet.Withdraw(req.Amount); err != nil {
			return err
		}

		balanceAfter = wallet.Balance
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 获取扣减前的余额（用于流水记录）
	balanceBefore := balanceAfter + req.Amount

	// 生成交易号
	transactionNo := utils.GenerateTransactionNo("WITHDRAW")

	// 创建提现交易记录（状态为 pending）
	transaction := &models.WalletTransaction{
		TransactionNo: transactionNo,
		Uid:           req.Uid,
		Type:          models.TransactionTypeWithdraw,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter, // 扣减后的余额
		Description:   req.Description,
		Status:        models.TransactionStatusPending, // 状态为待处理
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := s.walletRepo.CreateTransaction(ctx, transaction); err != nil {
		// 如果创建流水失败，需要回滚扣减的余额
		if rollbackErr := s.AddBalance(ctx, req.Uid, req.Amount, "提现申请失败回滚"); rollbackErr != nil {
			// 回滚失败，记录严重错误
			utils.LogError(nil, "提现申请失败且余额回滚失败: %v, 回滚错误: %v", err, rollbackErr)
		}
		return nil, utils.NewAppError(utils.CodeTransactionCreateFailed, "创建提现申请失败")
	}

	response := &models.WithdrawResponse{
		TransactionNo: transactionNo,
		Amount:        req.Amount,
		Balance:       balanceAfter, // 返回扣减后的余额
		Status:        models.TransactionStatusPending,
	}

	return response, nil
}

// GetWithdrawSummary 获取提现汇总信息
func (s *WalletService) GetWithdrawSummary(uid string) (*models.WithdrawSummary, error) {
	ctx := context.Background()

	// 获取提现汇总数据
	summary, err := s.walletRepo.GetWithdrawSummary(ctx, uid)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeDatabaseError, "获取提现汇总失败")
	}

	return summary, nil
}

// GetTransactionDetail 获取交易详情
func (s *WalletService) GetTransactionDetail(req *GetTransactionDetailRequest) (*models.TransactionDetail, error) {
	ctx := context.Background()

	// 获取交易详情
	transaction, err := s.walletRepo.GetTransactionByNo(ctx, req.TransactionNo)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeTransactionDetailGetFailed, "获取交易详情失败")
	}

	// 构建响应
	detail := &models.TransactionDetail{
		TransactionNo: transaction.TransactionNo,
		Uid:           transaction.Uid,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		Balance:       transaction.BalanceAfter, // 使用交易后余额
		Description:   transaction.Description,
		Status:        transaction.Status,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}

	return detail, nil
}

// 用户登录时延长钱包缓存过期时间
func (s *WalletService) ExtendWalletCacheOnLogin(ctx context.Context, uid string) error {
	if uid == "" {
		return utils.NewAppError(utils.CodeInvalidParams, "用户ID不能为空")
	}

	// 调用缓存服务延长过期时间
	return s.cacheService.ExtendWalletCacheOnLogin(ctx, uid)
}

// 清理过期钱包缓存
func (s *WalletService) CleanupExpiredWalletCache(ctx context.Context) error {
	return s.cacheService.CleanupExpiredWalletCache(ctx)
}
