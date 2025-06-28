package database

import (
	"context"
	"gin-fataMorgana/models"
)

// WalletRepository 钱包仓库
type WalletRepository struct {
	*BaseRepository
}

// NewWalletRepository 创建钱包仓库实例
func NewWalletRepository() *WalletRepository {
	return &WalletRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// CreateWallet 创建钱包
func (r *WalletRepository) CreateWallet(ctx context.Context, wallet *models.Wallet) error {
	return r.Create(ctx, wallet)
}

// FindWalletByUid 根据UID查找钱包
func (r *WalletRepository) FindWalletByUid(ctx context.Context, uid string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.FindByCondition(ctx, map[string]interface{}{"uid": uid}, &wallet)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// UpdateWallet 更新钱包
func (r *WalletRepository) UpdateWallet(ctx context.Context, wallet *models.Wallet) error {
	return r.Update(ctx, wallet)
}

// CreateTransaction 创建交易记录
func (r *WalletRepository) CreateTransaction(ctx context.Context, transaction *models.WalletTransaction) error {
	return r.Create(ctx, transaction)
}

// GetUserTransactions 获取用户资金记录
func (r *WalletRepository) GetUserTransactions(ctx context.Context, uid string, page, pageSize int) ([]models.WalletTransaction, int64, error) {
	var transactions []models.WalletTransaction
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.WalletTransaction{}).Where("uid = ?", uid).Count(&total).Error
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
		Find(&transactions).Error

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetTransactionByNo 根据交易流水号获取交易记录
func (r *WalletRepository) GetTransactionByNo(ctx context.Context, transactionNo string) (*models.WalletTransaction, error) {
	var transaction models.WalletTransaction
	err := r.FindByCondition(ctx, map[string]interface{}{"transaction_no": transactionNo}, &transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionsByType 根据类型获取用户交易记录
func (r *WalletRepository) GetTransactionsByType(ctx context.Context, uid, transactionType string, page, pageSize int) ([]models.WalletTransaction, int64, error) {
	var transactions []models.WalletTransaction
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.WalletTransaction{}).Where("uid = ? AND type = ?", uid, transactionType).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据
	err = r.db.WithContext(ctx).
		Where("uid = ? AND type = ?", uid, transactionType).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&transactions).Error

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetTransactionsByDateRange 根据日期范围获取用户交易记录
func (r *WalletRepository) GetTransactionsByDateRange(ctx context.Context, uid string, startDate, endDate string, page, pageSize int) ([]models.WalletTransaction, int64, error) {
	var transactions []models.WalletTransaction
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.WalletTransaction{}).
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
		Find(&transactions).Error

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// UpdateTransaction 更新交易记录
func (r *WalletRepository) UpdateTransaction(ctx context.Context, transaction *models.WalletTransaction) error {
	return r.Update(ctx, transaction)
}
