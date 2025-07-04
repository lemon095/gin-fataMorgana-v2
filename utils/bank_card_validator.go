package utils

import (
	"strings"
)

// 验证银行卡号长度
func ValidateCardNumberLength(cardNumber string) error {
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return NewAppError(CodeBankCardLengthInvalid, "银行卡号长度不正确，应为13-19位")
	}
	return nil
}

// 验证银行卡号格式
func ValidateCardNumberFormat(cardNumber string) error {
	for _, char := range cardNumber {
		if char < '0' || char > '9' {
			return NewAppError(CodeBankCardFormatInvalid, "银行卡号只能包含数字")
		}
	}
	return nil
}

// 验证持卡人姓名
func ValidateCardholderName(cardHolder string) error {
	if strings.TrimSpace(cardHolder) == "" {
		return NewAppError(CodeCardholderNameEmpty, "持卡人姓名不能为空")
	}

	if len(cardHolder) < 2 || len(cardHolder) > 50 {
		return NewAppError(CodeCardholderNameLength, "持卡人姓名长度应在2-50个字符之间")
	}

	// 验证持卡人姓名格式（只允许中文、英文字母和空格）
	for _, char := range cardHolder {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '一' && char <= '龯') || char == ' ') {
			return NewAppError(CodeCardholderNameFormat, "持卡人姓名只能包含中文、英文字母和空格")
		}
	}
	return nil
}

// 验证银行名称
func ValidateBankName(bankName string) error {
	if strings.TrimSpace(bankName) == "" {
		return NewAppError(CodeBankNameEmpty, "银行名称不能为空")
	}

	if len(bankName) < 2 || len(bankName) > 50 {
		return NewAppError(CodeBankNameLength, "银行名称长度应在2-50个字符之间")
	}
	return nil
} 