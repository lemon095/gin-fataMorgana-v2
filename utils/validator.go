package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators 注册自定义验证器
func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 可以在这里添加更多自定义验证器
		_ = v
	}
}

// ValidationError 验证错误信息
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// 字段名称映射（中文）
var fieldNameMap = map[string]string{
	"Email":           "邮箱",
	"Password":        "密码",
	"ConfirmPassword": "确认密码",
	"Username":        "用户名",
	"InviteCode":      "邀请码",
	"OldPassword":     "当前密码",
	"NewPassword":     "新密码",
	"BankName":        "银行名称",
	"CardHolder":      "持卡人",
	"CardNumber":      "银行卡号",
	"CardType":        "卡类型",
	"Page":            "页码",
	"PageSize":        "每页大小",
	"Type":            "类型",
	"Amount":          "金额",
	"PeriodNumber":    "期号",
	"Uid":             "用户ID",
	"OrderNo":         "订单号",
	"TransactionNo":   "交易流水号",
	"WithdrawAmount":  "提现金额",
	"RechargeAmount":  "充值金额",
	"BankCardInfo":    "银行卡信息",
}

// 验证标签错误信息映射
var tagErrorMap = map[string]string{
	"required":           "不能为空",
	"email":              "格式不正确",
	"min":                "长度不能少于%s位",
	"max":                "长度不能超过%s位",
	"len":                "长度必须为%s位",
	"oneof":              "必须是以下值之一：%s",
	"numeric":            "必须是数字",
	"alpha":              "只能包含字母",
	"alphanum":           "只能包含字母和数字",
	"url":                "必须是有效的URL",
	"uuid":               "必须是有效的UUID",
	"date":               "日期格式不正确",
	"time":               "时间格式不正确",
	"gt":                 "必须大于%s",
	"gte":                "必须大于等于%s",
	"lt":                 "必须小于%s",
	"lte":                "必须小于等于%s",
	"eq":                 "必须等于%s",
	"ne":                 "不能等于%s",
	"unique":             "值已存在",
	"exists":             "值不存在",
	"regexp":             "格式不正确",
	"json":               "JSON格式不正确",
	"base64":             "Base64格式不正确",
	"hexadecimal":        "十六进制格式不正确",
	"hexcolor":           "颜色格式不正确",
	"rgb":                "RGB格式不正确",
	"rgba":               "RGBA格式不正确",
	"hsl":                "HSL格式不正确",
	"hsla":               "HSLA格式不正确",
	"e164":               "电话号码格式不正确",
	"isbn":               "ISBN格式不正确",
	"isbn10":             "ISBN-10格式不正确",
	"isbn13":             "ISBN-13格式不正确",
	"uuid3":              "UUID v3格式不正确",
	"uuid4":              "UUID v4格式不正确",
	"uuid5":              "UUID v5格式不正确",
	"ulid":               "ULID格式不正确",
	"cron":               "Cron表达式格式不正确",
	"mongodb":            "MongoDB ObjectID格式不正确",
	"datetime":           "日期时间格式不正确",
	"image":              "必须是有效的图片文件",
	"file":               "必须是有效的文件",
	"dir":                "必须是有效的目录",
	"path":               "必须是有效的路径",
	"hostname":           "主机名格式不正确",
	"hostname_rfc1123":   "主机名格式不正确",
	"fqdn":               "完全限定域名格式不正确",
	"tcp_addr":           "TCP地址格式不正确",
	"tcp4_addr":          "TCP4地址格式不正确",
	"tcp6_addr":          "TCP6地址格式不正确",
	"udp_addr":           "UDP地址格式不正确",
	"udp4_addr":          "UDP4地址格式不正确",
	"udp6_addr":          "UDP6地址格式不正确",
	"ip_addr":            "IP地址格式不正确",
	"ip4_addr":           "IPv4地址格式不正确",
	"ip6_addr":           "IPv6地址格式不正确",
	"unix_addr":          "Unix地址格式不正确",
	"mac":                "MAC地址格式不正确",
	"hostname_port":      "主机名端口格式不正确",
	"ip_port":            "IP端口格式不正确",
	"ip4_port":           "IPv4端口格式不正确",
	"ip6_port":           "IPv6端口格式不正确",
	"tcp_port":           "TCP端口格式不正确",
	"udp_port":           "UDP端口格式不正确",
	"cidr":               "CIDR格式不正确",
	"cidrv4":             "CIDR v4格式不正确",
	"cidrv6":             "CIDR v6格式不正确",
	"multicast_ip":       "多播IP格式不正确",
	"multicast_ip4":      "多播IPv4格式不正确",
	"multicast_ip6":      "多播IPv6格式不正确",
	"private_ip":         "私有IP格式不正确",
	"private_ip4":        "私有IPv4格式不正确",
	"private_ip6":        "私有IPv6格式不正确",
	"public_ip":          "公网IP格式不正确",
	"public_ip4":         "公网IPv4格式不正确",
	"public_ip6":         "公网IPv6格式不正确",
	"loopback_ip":        "回环IP格式不正确",
	"loopback_ip4":       "回环IPv4格式不正确",
	"loopback_ip6":       "回环IPv6格式不正确",
	"link_local_ip":      "链路本地IP格式不正确",
	"link_local_ip4":     "链路本地IPv4格式不正确",
	"link_local_ip6":     "链路本地IPv6格式不正确",
	"global_unicast_ip":  "全局单播IP格式不正确",
	"global_unicast_ip4": "全局单播IPv4格式不正确",
	"global_unicast_ip6": "全局单播IPv6格式不正确",
	"unspecified_ip":     "未指定IP格式不正确",
	"unspecified_ip4":    "未指定IPv4格式不正确",
	"unspecified_ip6":    "未指定IPv6格式不正确",
}

// FormatValidationErrors 格式化验证错误
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := getFieldName(e.Field())
			tag := e.Tag()
			value := fmt.Sprintf("%v", e.Value())
			message := getErrorMessage(field, tag, e.Param())

			errors = append(errors, ValidationError{
				Field:   field,
				Tag:     tag,
				Value:   value,
				Message: message,
			})
		}
	}

	return errors
}

// getFieldName 获取字段的中文名称
func getFieldName(field string) string {
	if name, exists := fieldNameMap[field]; exists {
		return name
	}
	return field
}

// getErrorMessage 获取错误信息
func getErrorMessage(field, tag, param string) string {
	// 特殊处理一些常见的验证错误
	switch tag {
	case "required":
		return fmt.Sprintf("%s不能为空", field)
	case "email":
		return fmt.Sprintf("%s格式不正确", field)
	case "min":
		return fmt.Sprintf("%s长度不能少于%s位", field, param)
	case "max":
		return fmt.Sprintf("%s长度不能超过%s位", field, param)
	case "len":
		return fmt.Sprintf("%s长度必须为%s位", field, param)
	case "oneof":
		return fmt.Sprintf("%s必须是以下值之一：%s", field, param)
	case "gt", "gte", "lt", "lte", "eq", "ne":
		return fmt.Sprintf("%s%s", field, getTagErrorMessage(tag, param))
	default:
		if errorMsg, exists := tagErrorMap[tag]; exists {
			if param != "" {
				return fmt.Sprintf("%s%s", field, fmt.Sprintf(errorMsg, param))
			}
			return fmt.Sprintf("%s%s", field, errorMsg)
		}
		return fmt.Sprintf("%s格式不正确", field)
	}
}

// getTagErrorMessage 获取标签错误信息
func getTagErrorMessage(tag, param string) string {
	switch tag {
	case "gt":
		return fmt.Sprintf("必须大于%s", param)
	case "gte":
		return fmt.Sprintf("必须大于等于%s", param)
	case "lt":
		return fmt.Sprintf("必须小于%s", param)
	case "lte":
		return fmt.Sprintf("必须小于等于%s", param)
	case "eq":
		return fmt.Sprintf("必须等于%s", param)
	case "ne":
		return fmt.Sprintf("不能等于%s", param)
	default:
		return "格式不正确"
	}
}

// CreateValidationErrorResponse 创建验证错误响应
func CreateValidationErrorResponse(err error) Response {
	errors := FormatValidationErrors(err)
	var messages []string
	for _, e := range errors {
		messages = append(messages, e.Message)
	}
	mainMessage := "请求参数验证失败"
	if len(messages) > 0 {
		mainMessage = strings.Join(messages, "；")
	}
	return Response{
		Code:      CodeInvalidParams, // 统一用参数错误码
		Message:   mainMessage,
		Data:      nil, // 错误时data为nil
		Timestamp: time.Now().UnixMilli(),
	}
}

// HandleValidationError 处理验证错误并返回响应
func HandleValidationError(c *gin.Context, err error) {

	response := CreateValidationErrorResponse(err)
	c.JSON(getHTTPStatus(CodeInvalidParams), response)
}

// ValidateStruct 验证结构体并返回友好错误
func ValidateStruct(obj interface{}) error {
	validate := validator.New()

	// 注册自定义验证器（如果需要）
	// validate.RegisterValidation("custom", customValidator)

	err := validate.Struct(obj)
	if err != nil {
		return err
	}

	return nil
}

// ValidateAndHandleError 验证结构体并处理错误
func ValidateAndHandleError(c *gin.Context, obj interface{}) bool {
	err := ValidateStruct(obj)
	if err != nil {
		HandleValidationError(c, err)
		return false
	}
	return true
}

// IsAllDigits 检查字符串是否全为数字
func IsAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
