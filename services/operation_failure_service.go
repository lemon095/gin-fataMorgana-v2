package services

import (
	"context"
	"encoding/json"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
)

// OperationFailureService 操作失败记录服务
type OperationFailureService struct {
	failureRepo *database.OperationFailureRepository
}

func NewOperationFailureService() *OperationFailureService {
	return &OperationFailureService{
		failureRepo: database.NewOperationFailureRepository(),
	}
}

// RecordFailure 记录操作失败
func (s *OperationFailureService) RecordFailure(ctx context.Context, uid *string, operationType string, requestData, responseData interface{}) error {
	failure := &models.OperationFailure{
		Uid:           uid,
		OperationType: operationType,
	}

	// 处理请求数据
	if requestData != nil {
		requestJSON, err := json.Marshal(requestData)
		if err == nil {
			failure.RequestData = &models.JSONString{Data: requestJSON}
		}
	}

	// 处理响应数据
	if responseData != nil {
		responseJSON, err := json.Marshal(responseData)
		if err == nil {
			failure.ResponseData = &models.JSONString{Data: responseJSON}
		}
	}

	return s.failureRepo.Create(ctx, failure)
}

// RecordLoginFailure 记录登录失败
func (s *OperationFailureService) RecordLoginFailure(ctx context.Context, uid *string, requestData, responseData interface{}) error {
	return s.RecordFailure(ctx, uid, models.OperationTypeLogin, requestData, responseData)
}

// RecordRegisterFailure 记录注册失败
func (s *OperationFailureService) RecordRegisterFailure(ctx context.Context, requestData, responseData interface{}) error {
	// 注册时没有UID，所以传nil
	return s.RecordFailure(ctx, nil, models.OperationTypeRegister, requestData, responseData)
}

// RecordOrderCreateFailure 记录订单创建失败
func (s *OperationFailureService) RecordOrderCreateFailure(ctx context.Context, uid string, requestData, responseData interface{}) error {
	return s.RecordFailure(ctx, &uid, models.OperationTypeOrderCreate, requestData, responseData)
}

// RecordWalletWithdrawFailure 记录钱包提现失败
func (s *OperationFailureService) RecordWalletWithdrawFailure(ctx context.Context, uid string, requestData, responseData interface{}) error {
	return s.RecordFailure(ctx, &uid, models.OperationTypeWalletWithdraw, requestData, responseData)
}

// RecordBankCardBindFailure 记录银行卡绑定失败
func (s *OperationFailureService) RecordBankCardBindFailure(ctx context.Context, uid string, requestData, responseData interface{}) error {
	return s.RecordFailure(ctx, &uid, models.OperationTypeBankCardBind, requestData, responseData)
}

// GetUserFailures 获取用户失败记录
func (s *OperationFailureService) GetUserFailures(uid string, page, pageSize int) ([]*models.OperationFailure, int64, error) {
	ctx := context.Background()
	return s.failureRepo.GetByUID(ctx, uid, page, pageSize)
}

// GetFailuresByType 根据操作类型获取失败记录
func (s *OperationFailureService) GetFailuresByType(operationType string, page, pageSize int) ([]*models.OperationFailure, int64, error) {
	ctx := context.Background()
	return s.failureRepo.GetByOperationType(ctx, operationType, page, pageSize)
}

// GetRecentFailures 获取最近的失败记录
func (s *OperationFailureService) GetRecentFailures(limit int) ([]*models.OperationFailure, error) {
	ctx := context.Background()
	return s.failureRepo.GetRecentFailures(ctx, limit)
}

// GetFailureByID 根据ID获取失败记录
func (s *OperationFailureService) GetFailureByID(id uint) (*models.OperationFailure, error) {
	ctx := context.Background()
	return s.failureRepo.GetByID(ctx, id)
}
