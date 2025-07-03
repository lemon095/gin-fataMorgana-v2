package services

import (
	"context"
	"errors"
	"fmt"
	"log"
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

	// 验证请求参数
	if err := s.validateRegisterRequest(req); err != nil {
		log.Printf("用户注册参数验证失败: %v", err)
		return nil, fmt.Errorf("请求参数验证失败: %w", err)
	}

	// 先验证密码确认
	if req.Password != req.ConfirmPassword {
		log.Printf("用户注册失败：密码确认不匹配，邮箱: %s", req.Email)
		return nil, errors.New("两次输入的密码不一致")
	}

	// 检查邮箱是否存在
	emailExists, err := s.userRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		log.Printf("检查邮箱是否存在失败，邮箱: %s, 错误: %v", req.Email, err)
		return nil, fmt.Errorf("验证邮箱失败: %w", err)
	}

	if emailExists {
		log.Printf("用户注册失败：邮箱已存在，邮箱: %s", req.Email)
		return nil, errors.New("邮箱已被注册")
	}

	// 验证邀请码是否来自活跃的管理员
	if req.InviteCode != "" {
		adminUserRepo := database.NewAdminUserRepository()
		adminUser, err := adminUserRepo.GetActiveInviteCode(ctx, req.InviteCode)
		if err != nil {
			log.Printf("邀请码验证失败，邀请码: %s, 错误: %v", req.InviteCode, err)
			return nil, errors.New("邀请码无效或管理员账户已被禁用")
		}

		// 可以在这里添加额外的邀请码验证逻辑
		// 例如：检查管理员是否有权限邀请用户
		if !adminUser.IsActive() {
			log.Printf("邀请码对应的管理员账户已被禁用，邀请码: %s", req.InviteCode)
			return nil, errors.New("邀请码对应的管理员账户已被禁用")
		}
	}

	// 自动生成用户名（使用邮箱前缀 + 随机字符串）
	username := s.generateUsername(req.Email)

	// 使用雪花算法生成八位数用户ID
	userID := utils.GenerateUID()

	// 创建新用户
	user := &models.User{
		Uid:          userID,
		Username:     username,
		Email:        req.Email,
		Password:     req.Password,
		Status:       1, // 默认启用
		Experience:   1, // 新注册用户默认等级为1
		InvitedBy:    req.InviteCode,
		BankCardInfo: "{\"card_number\":\"\",\"card_holder\":\"\",\"bank_name\":\"\",\"card_type\":\"\"}", // 无条件赋值
	}
	// 保险：防止意外为空
	if user.BankCardInfo == "" {
		user.BankCardInfo = "{\"card_number\":\"\",\"card_holder\":\"\",\"bank_name\":\"\",\"card_type\":\"\"}"
	}

	// 加密密码
	if err := user.HashPassword(); err != nil {
		log.Printf("加密用户密码失败，邮箱: %s, 错误: %v", req.Email, err)
		return nil, fmt.Errorf("加密密码失败: %w", err)
	}

	// 保存用户到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Printf("创建用户失败，邮箱: %s, 错误: %v", req.Email, err)
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	log.Printf("用户注册成功，UID: %s, 邮箱: %s", userID, req.Email)

	response := user.ToResponse()
	return &response, nil
}

// validateRegisterRequest 验证注册请求参数
func (s *UserService) validateRegisterRequest(req *models.UserRegisterRequest) error {
	if req.Email == "" {
		return errors.New("邮箱不能为空")
	}
	if req.Password == "" {
		return errors.New("密码不能为空")
	}
	if len(req.Password) < 6 {
		return errors.New("密码长度不能少于6位")
	}
	if req.ConfirmPassword == "" {
		return errors.New("确认密码不能为空")
	}
	return nil
}

// generateUsername 根据邮箱生成用户名
func (s *UserService) generateUsername(email string) string {
	// 提取邮箱前缀
	parts := strings.Split(email, "@")
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

	// 首先检查用户是否存在（包括已删除的）
	user, err := s.userRepo.FindByEmailIncludeDeleted(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录失败的登录尝试
			s.recordFailedLogin(ctx, req.Email, loginIP, userAgent, "邮箱或密码错误")
			return nil, errors.New("邮箱或密码错误")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		s.recordFailedLogin(ctx, user, loginIP, userAgent, "账户已被删除")
		return nil, errors.New("账户已被删除，无法登录")
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		s.recordFailedLogin(ctx, user, loginIP, userAgent, "账户已被禁用")
		return nil, errors.New("账户已被禁用，无法登录")
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		s.recordFailedLogin(ctx, user, loginIP, userAgent, "密码错误")
		return nil, errors.New("邮箱或密码错误")
	}

	// 记录成功的登录
	s.recordSuccessfulLogin(ctx, user, loginIP, userAgent)

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

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(utils.AccessTokenExpiry.Seconds()),
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *UserService) RefreshToken(refreshToken string) (*models.TokenResponse, error) {
	ctx := context.Background()

	// 验证刷新令牌
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 查找用户
	user, err := s.userRepo.FindByUsername(ctx, claims.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		return nil, errors.New("账户已被删除，无法刷新令牌")
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		return nil, errors.New("账户已被禁用，无法刷新令牌")
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
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		return nil, errors.New("用户已被删除")
	}

	// 从Redis获取用户等级进度
	userLevelService := NewUserLevelService()
	rate, err := userLevelService.GetUserLevelRate(ctx, user.Uid)
	if err != nil {
		// 如果获取失败，使用默认值0
		rate = 0
	}

	response := user.ToResponse()
	response.Rate = rate // 设置从Redis获取的等级进度

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
		LoginTime: time.Now(),
		Status:    1, // 成功
	}
	s.loginLogRepo.Create(ctx, log)
}

// recordFailedLogin 记录失败登录
func (s *UserService) recordFailedLogin(ctx context.Context, user interface{}, loginIP, userAgent, reason string) {
	var log *models.UserLoginLog

	switch u := user.(type) {
	case *models.User:
		log = &models.UserLoginLog{
			Uid:        u.Uid,
			Username:   u.Username,
			Email:      u.Email,
			LoginIP:    loginIP,
			UserAgent:  userAgent,
			LoginTime:  time.Now(),
			Status:     0, // 失败
			FailReason: reason,
		}
	case string: // 邮箱字符串
		log = &models.UserLoginLog{
			Uid:        "",
			Username:   "",
			Email:      u,
			LoginIP:    loginIP,
			UserAgent:  userAgent,
			LoginTime:  time.Now(),
			Status:     0, // 失败
			FailReason: reason,
		}
	}

	if log != nil {
		s.loginLogRepo.Create(ctx, log)
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
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		return nil, errors.New("账户已被删除，无法绑定银行卡")
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		return nil, errors.New("账户已被禁用，无法绑定银行卡")
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
		return nil, fmt.Errorf("银行卡信息格式错误: %w", err)
	}

	// 更新用户的银行卡信息
	user.BankCardInfo = bankCardInfoJSON
	user.UpdatedAt = time.Now()

	// 保存到数据库
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("更新银行卡信息失败: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// validateBankCardInfo 验证银行卡信息
func (s *UserService) validateBankCardInfo(req *BindBankCardRequest) error {
	validator := utils.NewBankCardValidator()

	// 验证银行名称
	if err := validator.ValidateBankName(req.BankName); err != nil {
		return err
	}

	// 验证持卡人姓名
	if err := validator.ValidateCardHolder(req.CardHolder); err != nil {
		return err
	}

	// 验证银行卡号
	if err := validator.ValidateCardNumber(req.CardNumber); err != nil {
		return err
	}

	// 验证卡类型
	if strings.TrimSpace(req.CardType) == "" {
		return errors.New("卡类型不能为空")
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
		return errors.New("卡类型不正确，支持的类型：借记卡、信用卡、储蓄卡")
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
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户是否已被删除
	if user.DeletedAt != nil {
		return nil, errors.New("账户已被删除")
	}

	// 检查用户是否被禁用
	if user.Status == 0 {
		return nil, errors.New("账户已被禁用")
	}

	// 检查是否已绑定银行卡
	if user.BankCardInfo == "" {
		return nil, errors.New("用户未绑定银行卡")
	}

	// 解析银行卡信息
	var bankCardInfo models.BankCardInfo
	if err := utils.JSONToStruct(user.BankCardInfo, &bankCardInfo); err != nil {
		return nil, fmt.Errorf("银行卡信息格式错误: %w", err)
	}

	return &bankCardInfo, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	ctx := context.Background()

	// 1. 根据用户ID获取用户信息
	var user models.User
	if err := s.userRepo.FindByID(ctx, userID, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 2. 检查用户是否已被删除
	if user.DeletedAt != nil {
		return errors.New("用户已被删除，无法修改密码")
	}

	// 3. 检查用户是否被禁用
	if user.Status == 0 {
		return errors.New("账户已被禁用，无法修改密码")
	}

	// 4. 验证旧密码是否正确
	if !user.CheckPassword(oldPassword) {
		return errors.New("当前密码错误")
	}

	// 5. 验证新密码格式
	if err := s.validateNewPassword(newPassword); err != nil {
		return err
	}

	// 6. 检查新密码是否与旧密码相同
	if oldPassword == newPassword {
		return errors.New("新密码不能与当前密码相同")
	}

	// 7. 更新密码
	user.Password = newPassword
	if err := user.HashPassword(); err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	user.UpdatedAt = time.Now()

	// 8. 保存到数据库
	if err := s.userRepo.Update(ctx, &user); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 9. 可选：使所有现有token失效（这里可以通过Redis实现）
	// 由于项目中没有token黑名单机制，这里暂时跳过
	// 如果需要实现，可以在Redis中维护一个token黑名单

	return nil
}

// validateNewPassword 验证新密码格式
func (s *UserService) validateNewPassword(password string) error {
	if password == "" {
		return errors.New("新密码不能为空")
	}

	if len(password) < 6 {
		return errors.New("新密码长度不能少于6位")
	}

	if len(password) > 50 {
		return errors.New("新密码长度不能超过50位")
	}

	// 可以添加更多密码强度验证
	// 例如：必须包含数字、字母、特殊字符等
	// 这里暂时使用简单的长度验证

	return nil
}
