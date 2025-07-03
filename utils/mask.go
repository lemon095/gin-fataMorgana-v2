package utils

import (
	"strings"
)

// MaskPhone 手机号脱敏 (138****8888)
func MaskPhone(phone string) string {
	if len(phone) < 7 {
		return phone // 如果手机号太短，直接返回
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// MaskBankCard 银行卡号脱敏 (6222 **** **** 1234)
func MaskBankCard(cardNumber string) string {
	if len(cardNumber) < 8 {
		return cardNumber // 如果卡号太短，直接返回
	}
	return cardNumber[:4] + " **** **** " + cardNumber[len(cardNumber)-4:]
}

// MaskEmail 邮箱脱敏 (a***@example.com)
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // 如果邮箱格式不正确，直接返回
	}
	
	username := parts[0]
	domain := parts[1]
	
	if len(username) <= 1 {
		return email // 如果用户名太短，直接返回
	}
	
	maskedUsername := username[:1] + "***" + username[len(username)-1:]
	return maskedUsername + "@" + domain
}

// MaskIDCard 身份证号脱敏 (110***********1234)
func MaskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard // 如果身份证号太短，直接返回
	}
	return idCard[:3] + "***********" + idCard[len(idCard)-4:]
}

// MaskName 姓名脱敏 (张**)
func MaskName(name string) string {
	if len(name) <= 1 {
		return name // 如果姓名太短，直接返回
	}
	
	if len(name) == 2 {
		return name[:1] + "*"
	}
	
	return name[:1] + "*" + name[len(name)-1:]
}

// MaskAddress 地址脱敏 (北京市朝阳区***)
func MaskAddress(address string) string {
	if len(address) <= 6 {
		return address // 如果地址太短，直接返回
	}
	
	// 保留前3位和后3位，中间用*替换
	visibleLength := 6
	middleLength := len(address) - visibleLength
	
	if middleLength <= 0 {
		return address
	}
	
	middle := strings.Repeat("*", middleLength)
	return address[:3] + middle + address[len(address)-3:]
}

// IsSensitiveField 判断是否为敏感字段
func IsSensitiveField(fieldName string) bool {
	sensitiveFields := []string{
		"phone", "mobile", "tel", "telephone",
		"bank_card", "card_number", "card_no",
		"email", "mail",
		"id_card", "identity", "id_no",
		"name", "real_name", "full_name",
		"address", "addr",
		"password", "pwd", "passwd",
	}

	fieldName = strings.ToLower(fieldName)
	for _, field := range sensitiveFields {
		if strings.Contains(fieldName, field) {
			return true
		}
	}

	return false
}

// MaskSensitiveData 根据字段名自动脱敏
func MaskSensitiveData(fieldName string, value string) string {
	if value == "" {
		return value
	}

	fieldName = strings.ToLower(fieldName)

	switch {
	case strings.Contains(fieldName, "phone") || strings.Contains(fieldName, "mobile") || strings.Contains(fieldName, "tel"):
		return MaskPhone(value)
	case strings.Contains(fieldName, "bank_card") || strings.Contains(fieldName, "card_number") || strings.Contains(fieldName, "card_no"):
		return MaskBankCard(value)
	case strings.Contains(fieldName, "email") || strings.Contains(fieldName, "mail"):
		return MaskEmail(value)
	case strings.Contains(fieldName, "id_card") || strings.Contains(fieldName, "identity") || strings.Contains(fieldName, "id_no"):
		return MaskIDCard(value)
	case strings.Contains(fieldName, "name") || strings.Contains(fieldName, "real_name") || strings.Contains(fieldName, "full_name"):
		return MaskName(value)
	case strings.Contains(fieldName, "address") || strings.Contains(fieldName, "addr"):
		return MaskAddress(value)
	default:
		return value
	}
}
