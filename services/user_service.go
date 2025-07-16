package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	userRepo     *database.UserRepository
	loginLogRepo *database.LoginLogRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userRepo:     database.NewUserRepository(),
		loginLogRepo: database.NewLoginLogRepository(),
	}
}

// Register 用户注册
func (s *UserService) Register(req *models.UserRegisterRequest) (*models.UserResponse, error) {
	ctx := context.Background()

	// 判断账号类型
	isEmail := func(account string) bool {
		emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
		return emailRegex.MatchString(account)
	}
	isPhone := func(account string) bool {
		// 国际手机号E.164格式
		phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
		return phoneRegex.MatchString(account)
	}

	if !isEmail(req.Account) && !isPhone(req.Account) {
		return nil, utils.NewAppError(utils.CodeAccountFormatInvalid, "账号格式错误，请输入正确的邮箱或手机号")
	}

	if isEmail(req.Account) {
		emailExists, err := s.userRepo.CheckEmailExists(ctx, req.Account)
		if err != nil {
			return nil, utils.NewAppError(utils.CodeUserQueryFailed, "查询用户失败")
		}
		if emailExists {
			return nil, utils.NewAppError(utils.CodeEmailAlreadyExists, "邮箱已被注册")
		}
	}
	if isPhone(req.Account) {
		phoneExists, err := s.userRepo.CheckPhoneExists(ctx, req.Account)
		if err != nil {
			return nil, utils.NewAppError(utils.CodeUserQueryFailed, "查询用户失败")
		}
		if phoneExists {
			return nil, utils.NewAppError(utils.CodePhoneAlreadyExists, "手机号已被注册")
		}
	}

	if req.Password != req.ConfirmPassword {
		return nil, utils.NewAppError(utils.CodePasswordNotMatch, "两次输入的密码不一致")
	}

	// 验证邀请码是否来自活跃的管理员
	if req.InviteCode != "" {
		adminUserRepo := database.NewAdminUserRepository()
		adminUser, err := adminUserRepo.GetActiveInviteCode(ctx, strings.ToUpper(req.InviteCode))
		if err != nil {

			return nil, utils.NewAppError(utils.CodeInviteCodeAdminDisabled, "邀请码无效或管理员账户已被禁用")
		}

		// 可以在这里添加额外的邀请码验证逻辑
		// 例如：检查管理员是否有权限邀请用户
		if !adminUser.IsActive() {

			return nil, utils.NewAppError(utils.CodeInviteCodeAdminDisabled2, "邀请码对应的管理员账户已被禁用")
		}
	}

	// 自动生成用户名（使用邮箱/手机号前缀 + 随机字符串）
	username := s.generateUsername(req.Account)

	// 使用雪花算法生成八位数用户ID
	userID := utils.GenerateUID()

	// 创建新用户
	user := &models.User{
		Uid:          userID,
		Username:     username,
		Password:     req.Password,
		Status:       1,                                                                                   // 默认待审核
		Experience:   1,                                                                                   // 新注册用户默认等级为1
		InvitedBy:    strings.ToUpper(req.InviteCode),                                                     // 统一存储为大写格式
		BankCardInfo: "{\"card_number\":\"\",\"card_holder\":\"\",\"bank_name\":\"\",\"card_type\":\"\"}", // 无条件赋值
	}
	if isEmail(req.Account) {
		user.Email = req.Account
	}
	if isPhone(req.Account) {
		user.Phone = req.Account
	}
	// 保险：防止意外为空
	if user.BankCardInfo == "" {
		user.BankCardInfo = "{\"card_number\":\"\",\"card_holder\":\"\",\"bank_name\":\"\",\"card_type\":\"\"}"
	}

	// 加密密码
	if err := user.HashPassword(); err != nil {

		return nil, utils.NewAppError(utils.CodePasswordEncryptFailed, "加密密码失败")
	}

	// 保存用户到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {

		return nil, utils.NewAppError(utils.CodeUserCreateFailed, "创建用户失败")
	}

	// 自动为用户创建钱包
	walletService := NewWalletService()
	wallet, err := walletService.CreateWallet(user.Uid)
	if err != nil {
		// 记录钱包创建失败的错误，但不影响用户注册流程
		utils.LogWarn(nil, "用户注册后创建钱包失败 - UID: %s, 错误: %v", user.Uid, err)
	} else {
		utils.LogInfo(nil, "用户注册成功，自动创建钱包 - UID: %s, 钱包ID: %d", user.Uid, wallet.ID)
	}

	response := user.ToResponse()
	return &response, nil
}

// validateRegisterRequest 验证注册请求参数
func (s *UserService) validateRegisterRequest(req *models.UserRegisterRequest) error {
	if req.Account == "" {
		return utils.NewAppError(utils.CodeAccountEmpty, "账号不能为空")
	}
	if req.Password == "" {
		return utils.NewAppError(utils.CodePasswordEmpty, "密码不能为空")
	}
	if len(req.Password) < 6 {
		return utils.NewAppError(utils.CodePasswordTooShort, "密码长度不能少于6位")
	}
	if req.ConfirmPassword == "" {
		return utils.NewAppError(utils.CodePasswordEmpty, "确认密码不能为空")
	}
	return nil
}

// generateUsername 根据邮箱生成用户名
func (s *UserService) generateUsername(account string) string {
	// 提取邮箱前缀
	parts := strings.Split(account, "@")
	prefix := parts[0]

	// 生成随机后缀
	suffix := utils.RandomString(6)

	// 组合用户名
	username := fmt.Sprintf("%s_%s", prefix, suffix)

	return username
}

// Login 用户登录
func (s *UserService) Login(req *models.UserLoginRequest, loginIP, userAgent string) (*models.TokenResponse, error) {
	ctx := context.Background()

	// 判断账号类型
	isEmail := func(account string) bool {
		emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
		return emailRegex.MatchString(account)
	}
	isPhone := func(account string) bool {
		// 国际手机号E.164格式
		phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
		return phoneRegex.MatchString(account)
	}

	var user *models.User
	var err error
	if isEmail(req.Account) {
		user, err = s.userRepo.FindByEmail(ctx, req.Account)
	} else if isPhone(req.Account) {
		user, err = s.userRepo.FindByPhone(ctx, req.Account)
	} else {
		return nil, utils.NewAppError(utils.CodeAccountFormatInvalid, "账号格式错误，请输入正确的邮箱或手机号")
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.recordFailedLogin(ctx, req.Account, loginIP, userAgent, "账号不存在")
			return nil, utils.NewAppError(utils.CodeAccountNotFound, "账号不存在")
		}
		return nil, err
	}

	if !user.CheckPassword(req.Password) {
		s.recordFailedLogin(ctx, req.Account, loginIP, userAgent, "密码错误")
		return nil, utils.NewAppError(utils.CodeLoginCredentialError, "邮箱或手机号或密码错误")
	}

	// 新增：校验用户状态
	if user.DeletedAt != nil {
		s.recordFailedLogin(ctx, req.Account, loginIP, userAgent, "账户已被删除")
		return nil, utils.NewAppError(utils.CodeUserDeletedLogin, "账户已被删除，无法登录")
	}
	if user.Status == 0 {
		s.recordFailedLogin(ctx, req.Account, loginIP, userAgent, "账户已被禁用")
		return nil, utils.NewAppError(utils.CodeUserDisabledLogin, "账户已被禁用，无法登录")
	}
	if user.Status == 2 {
		s.recordFailedLogin(ctx, req.Account, loginIP, userAgent, "账户待审核")
		return nil, utils.NewAppError(utils.CodeUserPendingApproval, "账户待审核，无法登录")
	}

	// 记录成功的登录
	s.recordSuccessfulLogin(ctx, user, loginIP, userAgent)

	// 实现单点登录：检查并撤销旧会话
	tokenService := NewTokenService()
	activeToken, err := tokenService.GetUserActiveToken(ctx, user.Uid)
	if err == nil && activeToken != nil {
		// 将旧token加入黑名单
		tokenService.AddTokenToBlacklist(ctx, activeToken.TokenHash)
	}

	// 生成访问令牌
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Uid, user.Username)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Uid, user.Username)
	if err != nil {
		return nil, err
	}

	// 设置新的活跃token
	deviceInfo := s.extractDeviceInfo(userAgent)
	err = tokenService.SetUserActiveToken(ctx, user.Uid, accessToken, deviceInfo, loginIP, userAgent)
	if err != nil {
		// 记录错误但不影响登录流程
		// 这里不记录日志，因为不是关键错误
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(utils.AccessTokenExpiry.Seconds()),
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *UserService) RefreshToken(req *models.RefreshTokenRequest) (*models.TokenResponse, error) {
	ctx := context.Background()

	// 验证刷新令牌
	if req.RefreshToken == "" {
		return nil, utils.NewAppError(utils.CodeRefreshTokenInvalid, "无效的刷新令牌")
	}

	// 解析刷新令牌
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeRefreshTokenInvalid, "无效的刷新令牌")
	}

	// 检查用户是否存在
	user, err := s.userRepo.FindByUsername(ctx, claims.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(utils.CodeUserNotFound, "用户不存在")
		}
		return nil, utils.NewAppError(utils.CodeUserQueryFailed, "查询用户失败")
	}

	// 检查用户状态
	if user.Status == 0 { // 禁用
		return nil, utils.NewAppError(utils.CodeUserDisabledRefresh, "账户已被禁用，无法刷新令牌")
	}
	if user.Status == 2 { // 待审核
		return nil, utils.NewAppError(utils.CodeUserPendingRefresh, "账户待审核，无法刷新令牌")
	}

	// 生成新的访问令牌
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Uid, user.Username)
	if err != nil {
		return nil, err
	}

	// 生成新的刷新令牌
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, user.Uid, user.Username)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(utils.AccessTokenExpiry.Seconds()),
	}, nil
}

// GetUserByID 根据用户ID获取用户信息
func (s *UserService) GetUserByID(userID uint) (*models.UserResponse, error) {
	ctx := context.Background()

	var user models.User
	if err := s.userRepo.FindByID(ctx, userID, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(utils.CodeUserNotFound, "用户不存在")
		}
		return nil, utils.NewAppError(utils.CodeUserQueryFailed, "查询用户失败")
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		return nil, utils.NewAppError(utils.CodeUserDeleted, "用户已被删除")
	}

	// 从Redis获取用户等级进度和动态计算经验值
	userLevelService := NewUserLevelService()
	rate, err := userLevelService.GetUserLevelRate(ctx, user.Uid)
	if err != nil {
		// 如果获取失败，使用默认值0
		rate = 0
	}

	// 动态计算用户等级（experience）
	level, err := userLevelService.GetUserLevel(ctx, user.Uid)
	if err != nil {
		// 如果获取失败，使用默认值1
		level = 1
	}

	response := user.ToResponse()
	response.Rate = rate        // 设置从Redis获取的等级进度
	response.Experience = level // 动态计算的经验值（等级）

	return &response, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	ctx := context.Background()
	return s.userRepo.FindByEmail(ctx, email)
}

// recordSuccessfulLogin 记录成功登录
func (s *UserService) recordSuccessfulLogin(ctx context.Context, user *models.User, loginIP, userAgent string) {
	log := &models.UserLoginLog{
		Uid:       user.Uid,
		Username:  user.Username,
		Email:     user.Email,
		LoginIP:   loginIP,
		UserAgent: userAgent,
		LoginTime: time.Now().UTC(),
		Status:    1,
	}
	s.loginLogRepo.Create(ctx, log)
}

// recordFailedLogin 记录失败登录
func (s *UserService) recordFailedLogin(ctx context.Context, user interface{}, loginIP, userAgent, reason string) {
	var uid string

	// 安全的类型断言
	switch u := user.(type) {
	case *models.User:
		uid = u.Uid
	case string: // 邮箱字符串
		uid = u
	default:
		// 如果类型不匹配，记录错误但不panic
		// 这里不记录日志，因为不是关键错误
		return
	}

	// 记录失败登录
	logEntry := &models.UserLoginLog{
		Uid:        uid,
		LoginIP:    loginIP,
		UserAgent:  userAgent,
		Status:     0, // 0表示失败
		FailReason: reason,
		LoginTime:  time.Now().UTC(),
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.loginLogRepo.Create(ctx, logEntry); err != nil {
		// 记录错误但不影响主流程
		// 这里不记录日志，因为不是关键错误
	}
}

// BindBankCardRequest 绑定银行卡请求
type BindBankCardRequest struct {
	BankName   string `json:"bank_name" binding:"required"`
	CardHolder string `json:"card_holder" binding:"required"`
	CardNumber string `json:"card_number" binding:"required"`
	CardType   string `json:"card_type" binding:"required"` // 借记卡、信用卡等
}

// BindBankCard 绑定银行卡
func (s *UserService) BindBankCard(req *BindBankCardRequest, uid string) (*models.UserResponse, error) {
	ctx := context.Background()

	// 查找用户
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(utils.CodeUserNotFound, "用户不存在")
		}
		return nil, utils.NewAppError(utils.CodeUserQueryFailed, "查询用户失败")
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		return nil, utils.NewAppError(utils.CodeUserDeletedBindCard, "账户已被删除，无法绑定银行卡")
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		return nil, utils.NewAppError(utils.CodeUserDisabledBindCard, "账户已被禁用，无法绑定银行卡")
	}

	// 验证银行卡信息
	if err := s.validateBankCardInfo(req); err != nil {
		return nil, err
	}

	// 创建银行卡信息结构
	bankCardInfo := &models.BankCardInfo{
		CardNumber: req.CardNumber,
		CardType:   req.CardType,
		BankName:   req.BankName,
		CardHolder: req.CardHolder,
	}

	// 将银行卡信息转换为JSON字符串
	bankCardInfoJSON, err := utils.StructToJSON(bankCardInfo)
	if err != nil {
		return nil, utils.NewAppError(utils.CodeBankCardFormatError, "银行卡信息格式错误")
	}

	// 更新用户的银行卡信息
	user.BankCardInfo = bankCardInfoJSON
	user.UpdatedAt = time.Now().UTC()

	// 保存到数据库
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, utils.NewAppError(utils.CodeBankCardUpdateFailed, "更新银行卡信息失败")
	}

	response := user.ToResponse()
	return &response, nil
}

// validateBankCardInfo 验证银行卡信息
func (s *UserService) validateBankCardInfo(req *BindBankCardRequest) error {
	// 验证银行名称
	if err := utils.ValidateBankName(req.BankName); err != nil {
		return err
	}

	// 验证持卡人姓名
	if err := utils.ValidateCardholderName(req.CardHolder); err != nil {
		return err
	}

	// 验证银行卡号
	if err := utils.ValidateCardNumberFormat(req.CardNumber); err != nil {
		return err
	}

	// 验证卡类型
	if strings.TrimSpace(req.CardType) == "" {
		return utils.NewAppError(utils.CodeCardTypeEmpty, "卡类型不能为空")
	}

	// 验证卡类型是否在允许的范围内
	allowedCardTypes := []string{"借记卡", "信用卡", "储蓄卡"}
	isValidCardType := false
	for _, cardType := range allowedCardTypes {
		if req.CardType == cardType {
			isValidCardType = true
			break
		}
	}
	if !isValidCardType {
		return utils.NewAppError(utils.CodeCardTypeInvalid, "卡类型不正确，支持的类型：借记卡、信用卡、储蓄卡")
	}

	return nil
}

// GetBankCardInfo 获取银行卡信息
func (s *UserService) GetBankCardInfo(uid string) (*models.BankCardInfo, error) {
	ctx := context.Background()

	// 查找用户
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(utils.CodeUserNotFound, "用户不存在")
		}
		return nil, utils.NewAppError(utils.CodeUserQueryFailed, "查询用户失败")
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		return nil, utils.NewAppError(utils.CodeUserDeleted, "用户已被删除")
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		return nil, utils.NewAppError(utils.CodeUserDisabledLogin, "账户已被禁用")
	}

	// 检查是否已绑定银行卡
	if user.BankCardInfo == "" {
		return nil, utils.NewAppError(utils.CodeUserNoBankCard, "用户未绑定银行卡")
	}

	// 解析银行卡信息
	var bankCardInfo models.BankCardInfo
	if err := utils.JSONToStruct(user.BankCardInfo, &bankCardInfo); err != nil {
		return nil, utils.NewAppError(utils.CodeBankCardFormatError, "银行卡信息格式错误")
	}

	return &bankCardInfo, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(req *models.ChangePasswordRequest, uid string) error {
	ctx := context.Background()

	// 检查用户是否存在
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewAppError(utils.CodeUserNotFound, "用户不存在")
		}
		return utils.NewAppError(utils.CodeUserQueryFailed, "查询用户失败")
	}

	// 检查用户状态
	if user.Status == 0 { // 禁用
		return utils.NewAppError(utils.CodeUserDisabledChangePwd, "账户已被禁用，无法修改密码")
	}
	if user.Status == 2 { // 待审核
		return utils.NewAppError(utils.CodeUserPendingApproval, "账户待审核，无法修改密码")
	}

	// 验证当前密码
	if !user.CheckPassword(req.OldPassword) {
		return utils.NewAppError(utils.CodeCurrentPasswordWrong, "当前密码错误")
	}

	// 检查新密码是否与当前密码相同
	if user.CheckPassword(req.NewPassword) {
		return utils.NewAppError(utils.CodeNewPasswordSame, "新密码不能与当前密码相同")
	}

	// 验证密码长度
	if len(req.NewPassword) < 6 {
		return utils.NewAppError(utils.CodePasswordTooShort, "密码长度不能少于6位")
	}
	if len(req.NewPassword) > 50 {
		return utils.NewAppError(utils.CodePasswordTooLong, "密码长度不能超过50位")
	}

	// 更新密码
	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return utils.NewAppError(utils.CodePasswordEncryptFailed2, "密码加密失败")
	}

	user.UpdatedAt = time.Now().UTC()

	// 保存到数据库
	if err := s.userRepo.Update(ctx, user); err != nil {
		return utils.NewAppError(utils.CodePasswordUpdateFailed, "更新密码失败")
	}

	// 可选：使所有现有token失效（这里可以通过Redis实现）
	// 由于项目中没有token黑名单机制，这里暂时跳过
	// 如果需要实现，可以在Redis中维护一个token黑名单

	return nil
}

// extractDeviceInfo 提取设备信息
func (s *UserService) extractDeviceInfo(userAgent string) string {
	// 简单的设备信息提取，可以集成更复杂的解析库
	if len(userAgent) > 200 {
		userAgent = userAgent[:200]
	}
	return userAgent
}
