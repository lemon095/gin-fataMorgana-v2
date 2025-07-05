package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 完整的错误码定义
const (
	CodeSuccess           = 0    // 成功
	CodeError             = 1    // 一般错误
	CodeAuth              = 401  // 认证错误
	CodeForbidden         = 403  // 禁止访问
	CodeNotFound          = 404  // 资源不存在
	CodeValidation        = 422  // 数据验证错误
	CodeServer            = 500  // 服务器错误
	CodeDatabaseError     = 1001 // 数据库错误
	CodeRedisError        = 1002 // Redis错误
	CodeInvalidParams     = 1003 // 参数错误
	CodeOperationFailed   = 1004 // 操作失败
	CodeUserNotFound      = 1005 // 用户不存在
	CodeUserAlreadyExists = 1006 // 用户已存在
	CodeValidationFailed  = 1007 // 验证失败
	CodeAccountLocked     = 1008 // 账户锁定
	CodeRegisterClosed    = 1009 // 注册关闭
	CodeGroupBuyNotFound  = 1010 // 拼单不存在
	CodeGroupBuyExpired   = 3003 // 拼单已过期
	CodeGroupBuyOccupied  = 1012 // 拼单已被占用
	CodeEmailAlreadyExists      = 2001 // 邮箱已被注册
	CodePhoneAlreadyExists      = 2002 // 手机号已被注册
	CodeEmailFormatInvalid      = 2003 // 邮箱格式不正确
	CodePhoneFormatInvalid      = 2004 // 手机号格式不正确
	CodeAccountEmpty            = 2005 // 账号不能为空
	CodePasswordEmpty           = 2006 // 密码不能为空
	CodePasswordNotMatch        = 2007 // 两次输入的密码不一致
	CodeInviteCodeInvalid       = 2008 // 邀请码无效
	CodeAccountNotFound         = 2009 // 账号不存在
	CodePasswordIncorrect       = 2010 // 密码错误
	CodeBalanceInsufficient     = 3001 // 余额不足
	CodeGroupBuyFull            = 3002 // 拼单已满员
	CodeOrderAmountMismatch     = 3004 // 订单金额不匹配
	CodeWithdrawAmountExceeded  = 3005 // 提现金额超限
	CodeUserPendingApproval      = 1011 // 账户待审核，无法登录

	// 新增错误码
	CodeAccountFormatInvalid    = 2011 // 账号格式错误，请输入正确的邮箱或手机号
	CodeLoginCredentialError    = 2012 // 邮箱或手机号或密码错误
	CodeUserDeletedLogin        = 2013 // 账户已被删除，无法登录
	CodeUserDisabledLogin       = 2014 // 账户已被禁用，无法登录
	CodeRefreshTokenInvalid     = 2015 // 无效的刷新令牌
	CodeUserDeletedRefresh      = 2016 // 账户已被删除，无法刷新令牌
	CodeUserDisabledRefresh     = 2017 // 账户已被禁用，无法刷新令牌
	CodeUserPendingRefresh      = 2018 // 账户待审核，无法刷新令牌
	CodeUserDeleted             = 2019 // 用户已被删除
	CodeUserDeletedBindCard     = 2020 // 账户已被删除，无法绑定银行卡
	CodeUserDisabledBindCard    = 2021 // 账户已被禁用，无法绑定银行卡
	CodeCardTypeEmpty           = 2022 // 卡类型不能为空
	CodeCardTypeInvalid         = 2023 // 卡类型不正确，支持的类型：借记卡、信用卡、储蓄卡
	CodeUserNoBankCard          = 2024 // 用户未绑定银行卡
	CodeUserDeletedChangePwd    = 2025 // 用户已被删除，无法修改密码
	CodeUserDisabledChangePwd   = 2026 // 账户已被禁用，无法修改密码
	CodeCurrentPasswordWrong    = 2027 // 当前密码错误
	CodeNewPasswordSame         = 2028 // 新密码不能与当前密码相同
	CodePasswordTooShort        = 2029 // 密码长度不能少于6位
	CodePasswordTooLong         = 2030 // 密码长度不能超过50位
	CodeInviteCodeAdminDisabled = 2031 // 邀请码无效或管理员账户已被禁用
	CodeInviteCodeAdminDisabled2 = 2032 // 邀请码对应的管理员账户已被禁用

	// 业务逻辑错误码
	CodeOrderStatusInvalid      = 3006 // 状态类型参数无效，必须是1(进行中)、2(已完成)或3(拼单数据)
	CodeOrderAccessDenied       = 3007 // 无权访问此订单
	CodeOrderAmountInvalid      = 3008 // 订单金额必须大于0
	CodeProfitAmountInvalid     = 3009 // 利润金额不能为负数
	CodeTaskCountInvalid        = 3010 // 至少需要有一个任务数量大于0
	CodeTaskCountNegative       = 3011 // 任务数量不能为负数
	CodeWalletFrozenOrder       = 3012 // 钱包已被冻结，无法创建订单
	CodeWalletFrozenRecharge    = 3013 // 钱包已被冻结，无法充值
	CodeRechargeAmountInvalid   = 3014 // 充值金额必须大于0
	CodeRechargeAmountExceeded  = 3015 // 单笔充值金额不能超过100万元
	CodeWalletFrozenWithdraw    = 3016 // 钱包已被冻结，无法提现
	CodeWithdrawAmountInvalid   = 3017 // 提现金额必须大于0
	// CodeWithdrawAmountExceeded2 = 3018 // 单笔提现金额不能超过100万元（已移除）
	CodeBankCardNotBound        = 3019 // 请先绑定银行卡后再进行提现操作
	CodeWithdrawPasswordWrong   = 3020 // 登录密码错误
	// CodeDailyWithdrawExceeded   = 3021 // 超过每日提现限额（已移除）
	CodePeriodNotFound          = 3022 // 期数不存在
	CodePeriodNotStarted        = 3023 // 期数还未开始
	CodePeriodEnded             = 3024 // 期数已结束
	CodePeriodAlreadyBought     = 3025 // 您已经购买过期号的订单

	// 银行卡相关错误码
	CodeBankCardLengthInvalid   = 4001 // 银行卡号长度不正确，应为13-19位
	CodeBankCardFormatInvalid   = 4002 // 银行卡号只能包含数字
	CodeCardholderNameEmpty     = 4003 // 持卡人姓名不能为空
	CodeCardholderNameLength    = 4004 // 持卡人姓名长度应在2-50个字符之间
	CodeCardholderNameFormat    = 4005 // 持卡人姓名只能包含中文、英文字母和空格
	CodeBankNameEmpty           = 4006 // 银行名称不能为空
	CodeBankNameLength          = 4007 // 银行名称长度应在2-50个字符之间

	// Token相关错误码
	CodeTokenValidationFailed   = 5001 // token验证失败
	CodeTokenSingleLogin        = 5002 // 您的账号已在其他设备登录，请重新登录
	CodeTokenInvalid            = 5003 // 无效的令牌
	CodeTokenExpired            = 5004 // 令牌已过期

	// 系统配置错误码
	CodeConfigReadFailed        = 6001 // 读取配置文件失败
	CodeConfigParseFailed       = 6002 // 解析配置文件失败
	CodeConfigNotLoaded         = 6003 // 配置未加载
	CodeDBHostEmpty             = 6004 // 数据库主机地址不能为空
	CodeDBUserEmpty             = 6005 // 数据库用户名不能为空
	CodeDBNameEmpty             = 6006 // 数据库名称不能为空
	CodeRedisHostEmpty          = 6007 // Redis主机地址不能为空
	CodeJWTSecretEmpty          = 6008 // JWT密钥不能为空
	CodeDBConnectFailed         = 6009 // 连接数据库失败
	CodeDBInstanceFailed        = 6010 // 获取数据库实例失败
	CodeDBNotInitialized        = 6011 // 数据库未初始化
	CodeDBMigrationFailed       = 6012 // 数据库迁移失败
	CodeRedisConnectFailed      = 6013 // Redis连接测试失败

	// 工具类错误码
	CodeInviteCodeGenFailed     = 7001 // 无法生成唯一邀请码，请稍后重试
	CodeIdempotencyKeyExists    = 7002 // 重复请求，幂等键已存在
	CodeRecordNotFound          = 7003 // 记录不存在

	// 中间件错误码
	CodeRateLimitExceeded       = 8001 // 请求过于频繁，请稍后再试
	CodeSystemBusy              = 8002 // 系统繁忙，请稍后再试
	CodeRegisterNotAllowed      = 8003 // 当前系统不允许注册

	// 服务层错误码
	CodeLeaderboardGetFailed    = 9001 // 获取热榜数据失败
	CodePasswordEncryptFailed   = 9002 // 加密密码失败
	CodeUserCreateFailed        = 9003 // 创建用户失败
	CodeUserQueryFailed         = 9004 // 查询用户失败
	CodeBankCardFormatError     = 9005 // 银行卡信息格式错误
	CodeBankCardUpdateFailed    = 9006 // 更新银行卡信息失败
	CodePasswordEncryptFailed2  = 9007 // 密码加密失败
	CodePasswordUpdateFailed    = 9008 // 更新密码失败
	CodeUserInfoGetFailed       = 9009 // 获取用户信息失败
	CodeOrderDataValidateFailed = 9010 // 订单数据验证失败
	CodeWalletGetFailed         = 9011 // 获取钱包失败
	CodeBalanceDeductFailed     = 9012 // 扣减余额失败
	CodeWalletUpdateFailed      = 9013 // 更新钱包失败
	CodeOrderCreateFailed       = 9014 // 创建订单失败
	CodePriceConfigGetFailed    = 9015 // 获取价格配置失败
	CodeOrderListGetFailed      = 9016 // 获取订单列表失败
	CodeGroupBuyListGetFailed   = 9017 // 获取拼单列表失败
	CodeOrderDetailGetFailed    = 9018 // 获取订单详情失败
	CodeOrderStatsGetFailed     = 9019 // 获取订单统计失败
	CodePeriodInfoGetFailed     = 9020 // 获取期数信息失败
	CodePriceConfigParseFailed  = 9021 // 解析价格配置失败
	CodePeriodDuplicateCheckFailed = 9022 // 检查期号重复失败
	CodeUserLevelGetFailed      = 9023 // 获取用户等级信息失败
	CodeUserLevelParseFailed    = 9024 // 解析用户等级信息失败
	CodeUserLevelSerializeFailed = 9025 // 序列化用户等级信息失败
	CodeUserLevelStoreFailed    = 9026 // 存储用户等级信息失败
	CodeUserFundRecordGetFailed = 9027 // 获取用户资金记录失败
	CodeUserNotExists           = 9028 // 用户不存在
	CodeUserDisabledCreateWallet = 9029 // 用户账户已被禁用，无法创建钱包
	CodeWalletCreateFailed      = 9030 // 创建钱包失败
	CodeTransactionCreateFailed = 9031 // 创建交易记录失败
	CodeRechargeCreateFailed    = 9032 // 创建充值申请失败
	CodeWithdrawCreateFailed    = 9033 // 创建提现申请失败
	CodeTodayWithdrawQueryFailed = 9034 // 查询今日提现记录失败
	CodePendingWithdrawQueryFailed = 9035 // 查询待处理提现记录失败
	CodeTransactionDetailGetFailed = 9036 // 获取交易详情失败
)

// ResponseMessage 完整的响应消息映射
var ResponseMessage = map[int]string{
	CodeSuccess:           "操作成功",
	CodeError:             "操作失败",
	CodeAuth:              "认证失败",
	CodeForbidden:         "禁止访问",
	CodeNotFound:          "资源不存在",
	CodeValidation:        "数据验证失败",
	CodeServer:            "服务器内部错误",
	CodeDatabaseError:     "数据库操作失败",
	CodeRedisError:        "Redis操作失败",
	CodeInvalidParams:     "参数错误",
	CodeOperationFailed:   "操作失败",
	CodeUserNotFound:      "用户不存在",
	CodeUserAlreadyExists: "用户已存在",
	CodeValidationFailed:  "验证失败",
	CodeAccountLocked:     "账户已被锁定",
	CodeRegisterClosed:    "当前系统不允许注册",
	CodeGroupBuyNotFound:  "拼单不存在或已被删除",
	CodeGroupBuyExpired:   "拼单已过期",
	CodeGroupBuyOccupied:  "拼单已被其他用户参与",
	CodeEmailAlreadyExists:     "邮箱已被注册",
	CodePhoneAlreadyExists:     "手机号已被注册",
	CodeEmailFormatInvalid:     "邮箱格式不正确",
	CodePhoneFormatInvalid:     "手机号格式不正确",
	CodeAccountEmpty:           "账号不能为空",
	CodePasswordEmpty:          "密码不能为空",
	CodePasswordNotMatch:       "两次输入的密码不一致",
	CodeInviteCodeInvalid:      "邀请码无效",
	CodeAccountNotFound:        "账号不存在",
	CodePasswordIncorrect:      "密码错误",
	CodeBalanceInsufficient:    "余额不足",
	CodeGroupBuyFull:           "拼单已满员",
	CodeOrderAmountMismatch:    "订单金额不匹配",
	CodeWithdrawAmountExceeded: "提现金额超限",
	CodeUserPendingApproval:     "账户待审核，无法登录",

	// 新增错误消息
	CodeAccountFormatInvalid:    "账号格式错误，请输入正确的邮箱或手机号",
	CodeLoginCredentialError:    "邮箱或手机号或密码错误",
	CodeUserDeletedLogin:        "账户已被删除，无法登录",
	CodeUserDisabledLogin:       "账户已被禁用，无法登录",
	CodeRefreshTokenInvalid:     "无效的刷新令牌",
	CodeUserDeletedRefresh:      "账户已被删除，无法刷新令牌",
	CodeUserDisabledRefresh:     "账户已被禁用，无法刷新令牌",
	CodeUserPendingRefresh:      "账户待审核，无法刷新令牌",
	CodeUserDeleted:             "用户已被删除",
	CodeUserDeletedBindCard:     "账户已被删除，无法绑定银行卡",
	CodeUserDisabledBindCard:    "账户已被禁用，无法绑定银行卡",
	CodeCardTypeEmpty:           "卡类型不能为空",
	CodeCardTypeInvalid:         "卡类型不正确，支持的类型：借记卡、信用卡、储蓄卡",
	CodeUserNoBankCard:          "用户未绑定银行卡",
	CodeUserDeletedChangePwd:    "用户已被删除，无法修改密码",
	CodeUserDisabledChangePwd:   "账户已被禁用，无法修改密码",
	CodeCurrentPasswordWrong:    "当前密码错误",
	CodeNewPasswordSame:         "新密码不能与当前密码相同",
	CodePasswordTooShort:        "密码长度不能少于6位",
	CodePasswordTooLong:         "密码长度不能超过50位",
	CodeInviteCodeAdminDisabled: "邀请码无效或管理员账户已被禁用",
	CodeInviteCodeAdminDisabled2: "邀请码对应的管理员账户已被禁用",

	// 业务逻辑错误消息
	CodeOrderStatusInvalid:      "状态类型参数无效，必须是1(进行中)、2(已完成)或3(拼单数据)",
	CodeOrderAccessDenied:       "无权访问此订单",
	CodeOrderAmountInvalid:      "订单金额必须大于0",
	CodeProfitAmountInvalid:     "利润金额不能为负数",
	CodeTaskCountInvalid:        "至少需要有一个任务数量大于0",
	CodeTaskCountNegative:       "任务数量不能为负数",
	CodeWalletFrozenOrder:       "钱包已被冻结，无法创建订单",
	CodeWalletFrozenRecharge:    "钱包已被冻结，无法充值",
	CodeRechargeAmountInvalid:   "充值金额必须大于0",
	CodeRechargeAmountExceeded:  "单笔充值金额不能超过100万元",
	CodeWalletFrozenWithdraw:    "钱包已被冻结，无法提现",
	CodeWithdrawAmountInvalid:   "提现金额必须大于0",
	// CodeWithdrawAmountExceeded2: "单笔提现金额不能超过100万元", // 已移除
	CodeBankCardNotBound:        "请先绑定银行卡后再进行提现操作",
	CodeWithdrawPasswordWrong:   "登录密码错误",
	// CodeDailyWithdrawExceeded:   "超过每日提现限额", // 已移除
	CodePeriodNotFound:          "期数不存在",
	CodePeriodNotStarted:        "期数还未开始",
	CodePeriodEnded:             "期数已结束",
	CodePeriodAlreadyBought:     "您已经购买过期号的订单",

	// 银行卡相关错误消息
	CodeBankCardLengthInvalid:   "银行卡号长度不正确，应为13-19位",
	CodeBankCardFormatInvalid:   "银行卡号只能包含数字",
	CodeCardholderNameEmpty:     "持卡人姓名不能为空",
	CodeCardholderNameLength:    "持卡人姓名长度应在2-50个字符之间",
	CodeCardholderNameFormat:    "持卡人姓名只能包含中文、英文字母和空格",
	CodeBankNameEmpty:           "银行名称不能为空",
	CodeBankNameLength:          "银行名称长度应在2-50个字符之间",

	// Token相关错误消息
	CodeTokenValidationFailed:   "token验证失败",
	CodeTokenSingleLogin:        "您的账号已在其他设备登录，请重新登录",
	CodeTokenInvalid:            "无效的令牌",
	CodeTokenExpired:            "令牌已过期",

	// 系统配置错误消息
	CodeConfigReadFailed:        "读取配置文件失败",
	CodeConfigParseFailed:       "解析配置文件失败",
	CodeConfigNotLoaded:         "配置未加载",
	CodeDBHostEmpty:             "数据库主机地址不能为空",
	CodeDBUserEmpty:             "数据库用户名不能为空",
	CodeDBNameEmpty:             "数据库名称不能为空",
	CodeRedisHostEmpty:          "Redis主机地址不能为空",
	CodeJWTSecretEmpty:          "JWT密钥不能为空",
	CodeDBConnectFailed:         "连接数据库失败",
	CodeDBInstanceFailed:        "获取数据库实例失败",
	CodeDBNotInitialized:        "数据库未初始化",
	CodeDBMigrationFailed:       "数据库迁移失败",
	CodeRedisConnectFailed:      "Redis连接测试失败",

	// 工具类错误消息
	CodeInviteCodeGenFailed:     "无法生成唯一邀请码，请稍后重试",
	CodeIdempotencyKeyExists:    "重复请求，幂等键已存在",
	CodeRecordNotFound:          "记录不存在",

	// 中间件错误消息
	CodeRateLimitExceeded:       "请求过于频繁，请稍后再试",
	CodeSystemBusy:              "系统繁忙，请稍后再试",
	CodeRegisterNotAllowed:      "当前系统不允许注册",

	// 服务层错误消息
	CodeLeaderboardGetFailed:    "获取热榜数据失败",
	CodePasswordEncryptFailed:   "加密密码失败",
	CodeUserCreateFailed:        "创建用户失败",
	CodeUserQueryFailed:         "查询用户失败",
	CodeBankCardFormatError:     "银行卡信息格式错误",
	CodeBankCardUpdateFailed:    "更新银行卡信息失败",
	CodePasswordEncryptFailed2:  "密码加密失败",
	CodePasswordUpdateFailed:    "更新密码失败",
	CodeUserInfoGetFailed:       "获取用户信息失败",
	CodeOrderDataValidateFailed: "订单数据验证失败",
	CodeWalletGetFailed:         "获取钱包失败",
	CodeBalanceDeductFailed:     "扣减余额失败",
	CodeWalletUpdateFailed:      "更新钱包失败",
	CodeOrderCreateFailed:       "创建订单失败",
	CodePriceConfigGetFailed:    "获取价格配置失败",
	CodeOrderListGetFailed:      "获取订单列表失败",
	CodeGroupBuyListGetFailed:   "获取拼单列表失败",
	CodeOrderDetailGetFailed:    "获取订单详情失败",
	CodeOrderStatsGetFailed:     "获取订单统计失败",
	CodePeriodInfoGetFailed:     "获取期数信息失败",
	CodePriceConfigParseFailed:  "解析价格配置失败",
	CodePeriodDuplicateCheckFailed: "检查期号重复失败",
	CodeUserLevelGetFailed:      "获取用户等级信息失败",
	CodeUserLevelParseFailed:    "解析用户等级信息失败",
	CodeUserLevelSerializeFailed: "序列化用户等级信息失败",
	CodeUserLevelStoreFailed:    "存储用户等级信息失败",
	CodeUserFundRecordGetFailed: "获取用户资金记录失败",
	CodeUserNotExists:           "用户不存在",
	CodeUserDisabledCreateWallet: "用户账户已被禁用，无法创建钱包",
	CodeWalletCreateFailed:      "创建钱包失败",
	CodeTransactionCreateFailed: "创建交易记录失败",
	CodeRechargeCreateFailed:    "创建充值申请失败",
	CodeWithdrawCreateFailed:    "创建提现申请失败",
	CodeTodayWithdrawQueryFailed: "查询今日提现记录失败",
	CodePendingWithdrawQueryFailed: "查询待处理提现记录失败",
	CodeTransactionDetailGetFailed: "获取交易详情失败",
}

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Message:   ResponseMessage[CodeSuccess],
		Data:      data,
		Timestamp: time.Now().UTC().UnixMilli(),
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().UnixMilli(),
	})
}

// Error 错误响应
func Error(c *gin.Context, code int) {
	message := ResponseMessage[code]
	if message == "" {
		message = "未知错误"
	}

	c.JSON(getHTTPStatus(code), Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UTC().UnixMilli(),
	})
}

// ErrorWithMessage 带自定义消息的错误响应
func ErrorWithMessage(c *gin.Context, code int, message string) {
	c.JSON(getHTTPStatus(code), Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UTC().UnixMilli(),
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code int, data interface{}) {
	message := ResponseMessage[code]
	if message == "" {
		message = "未知错误"
	}

	c.JSON(getHTTPStatus(code), Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().UnixMilli(),
	})
}

// getHTTPStatus 根据错误码获取HTTP状态码
func getHTTPStatus(code int) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeAuth:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeValidation:
		return http.StatusUnprocessableEntity
	case CodeServer:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

// 常用响应函数

// InvalidParams 参数错误
func InvalidParams(c *gin.Context) {
	Error(c, CodeValidation)
}

// InvalidParamsWithMessage 带消息的参数错误
func InvalidParamsWithMessage(c *gin.Context, message string) {
	ErrorWithMessage(c, CodeValidation, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context) {
	Error(c, CodeAuth)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context) {
	Error(c, CodeForbidden)
}

// NotFound 资源不存在
func NotFound(c *gin.Context) {
	Error(c, CodeNotFound)
}

// InternalError 内部服务器错误
func InternalError(c *gin.Context) {
	Error(c, CodeServer)
}

// DatabaseError 数据库错误
func DatabaseError(c *gin.Context) {
	Error(c, CodeDatabaseError)
}

// RedisError Redis错误
func RedisError(c *gin.Context) {
	Error(c, CodeRedisError)
}

// LoginFailed 登录失败
func LoginFailed(c *gin.Context) {
	ErrorWithMessage(c, CodeAuth, "邮箱或密码错误")
}

// UserNotFound 用户不存在
func UserNotFound(c *gin.Context) {
	ErrorWithMessage(c, CodeNotFound, "用户不存在")
}

// UserAlreadyExists 用户已存在
func UserAlreadyExists(c *gin.Context) {
	ErrorWithMessage(c, CodeValidation, "用户已存在")
}

// EmailAlreadyExists 邮箱已存在
func EmailAlreadyExists(c *gin.Context) {
	ErrorWithMessage(c, CodeValidation, "邮箱已被注册")
}

// TokenExpired Token过期
func TokenExpired(c *gin.Context) {
	ErrorWithMessage(c, CodeAuth, "Token已过期")
}

// TokenInvalid Token无效
func TokenInvalid(c *gin.Context) {
	ErrorWithMessage(c, CodeAuth, "Token无效")
}

// InviteCodeInvalid 邀请码无效
func InviteCodeInvalid(c *gin.Context) {
	ErrorWithMessage(c, CodeValidation, "邀请码无效")
}

// AccountLocked 账户锁定
func AccountLocked(c *gin.Context) {
	ErrorWithMessage(c, CodeAccountLocked, "账户已被锁定")
}

// 账户待审核
func UserPendingApproval(c *gin.Context) {
	ErrorWithMessage(c, CodeUserPendingApproval, ResponseMessage[CodeUserPendingApproval])
}
