package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// BankCardValidator 银行卡校验器
type BankCardValidator struct{}

// NewBankCardValidator 创建银行卡校验器
func NewBankCardValidator() *BankCardValidator {
	return &BankCardValidator{}
}

// ValidateCardNumber 验证银行卡号
func (v *BankCardValidator) ValidateCardNumber(cardNumber string) error {
	// 1. 基本格式验证
	if err := v.validateBasicFormat(cardNumber); err != nil {
		return err
	}

	// 2. Luhn算法校验（开发阶段暂时注释）
	// TODO: 生产环境需要启用Luhn算法校验
	// if !v.luhnCheck(cardNumber) {
	// 	return errors.New("银行卡号校验失败，请检查卡号是否正确")
	// }

	return nil
}

// validateBasicFormat 基本格式验证
func (v *BankCardValidator) validateBasicFormat(cardNumber string) error {
	// 去除空格
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	
	// 检查长度
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return errors.New("银行卡号长度不正确，应为13-19位")
	}

	// 检查是否全为数字
	if !regexp.MustCompile(`^\d+$`).MatchString(cardNumber) {
		return errors.New("银行卡号只能包含数字")
	}

	return nil
}

// luhnCheck Luhn算法校验
func (v *BankCardValidator) luhnCheck(cardNumber string) bool {
	sum := 0
	alternate := false

	// 从右到左遍历
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(cardNumber[i]))

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// ValidateCardHolder 验证持卡人姓名
func (v *BankCardValidator) ValidateCardHolder(cardHolder string) error {
	if cardHolder == "" {
		return errors.New("持卡人姓名不能为空")
	}
	
	if len(cardHolder) < 2 || len(cardHolder) > 50 {
		return errors.New("持卡人姓名长度应在2-50个字符之间")
	}
	
	// 检查是否包含特殊字符
	if !regexp.MustCompile(`^[\p{Han}a-zA-Z\s]+$`).MatchString(cardHolder) {
		return errors.New("持卡人姓名只能包含中文、英文字母和空格")
	}
	
	return nil
}

// ValidateBankName 验证银行名称
func (v *BankCardValidator) ValidateBankName(bankName string) error {
	if bankName == "" {
		return errors.New("银行名称不能为空")
	}
	
	if len(bankName) < 2 || len(bankName) > 50 {
		return errors.New("银行名称长度应在2-50个字符之间")
	}
	
	return nil
} 