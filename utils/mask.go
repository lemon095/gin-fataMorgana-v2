package utils

import (
	"strings"
)

// MaskPhone 手机号脱敏 (138****8888)
func MaskPhone(phone string) string {
	if phone == "" {
		return phone
	}

	// 去除空格和特殊字符
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "_", "")

	if len(phone) < 7 {
		return phone
	}

	// 保留前3位和后4位，中间用****代替
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// MaskBankCard 银行卡号脱敏 (6222 **** **** 1234)
func MaskBankCard(cardNumber string) string {
	if cardNumber == "" {
		return cardNumber
	}

	// 去除空格和特殊字符
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	cardNumber = strings.ReplaceAll(cardNumber, "-", "")
	cardNumber = strings.ReplaceAll(cardNumber, "_", "")

	if len(cardNumber) < 8 {
		return cardNumber
	}

	// 保留前4位和后4位，中间用****代替
	return cardNumber[:4] + " **** **** " + cardNumber[len(cardNumber)-4:]
}

// MaskEmail 邮箱脱敏 (a***@example.com)
func MaskEmail(email string) string {
	if email == "" {
		return email
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 1 {
		return email
	}

	// 用户名只显示第一个字符，其他用***代替
	maskedUsername := username[:1] + "***"

	return maskedUsername + "@" + domain
}

// MaskIDCard 身份证号脱敏 (110***********1234)
func MaskIDCard(idCard string) string {
	if idCard == "" {
		return idCard
	}

	// 去除空格和特殊字符
	idCard = strings.ReplaceAll(idCard, " ", "")
	idCard = strings.ReplaceAll(idCard, "-", "")
	idCard = strings.ReplaceAll(idCard, "_", "")

	if len(idCard) < 8 {
		return idCard
	}

	// 保留前3位和后4位，中间用***********代替
	return idCard[:3] + "***********" + idCard[len(idCard)-4:]
}

// MaskName 姓名脱敏 (张**)
func MaskName(name string) string {
	if name == "" {
		return name
	}

	// 去除空格
	name = strings.TrimSpace(name)

	if len(name) <= 1 {
		return name
	}

	// 保留第一个字符，其他用**代替
	return name[:1] + "**"
}

// MaskAddress 地址脱敏 (北京市朝阳区***)
func MaskAddress(address string) string {
	if address == "" {
		return address
	}

	// 去除空格
	address = strings.TrimSpace(address)

	if len(address) <= 6 {
		return address
	}

	// 保留前6个字符，其他用***代替
	return address[:6] + "***"
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
