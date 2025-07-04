package utils

import (
	"math/rand"
	"time"
)

const (
	// 字符集
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

	// 邀请码字符集（字母和数字混合，大小写随机）
	inviteCodeBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var src = rand.NewSource(time.Now().UnixNano())

// RandomString 生成指定长度的随机字符串
func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// GenerateInviteCode 生成邀请码（6位字母数字混合，大小写随机，不重复）
func GenerateInviteCode() string {
	// 使用map来确保字符不重复
	usedChars := make(map[byte]bool)
	b := make([]byte, 6)

	for i := 0; i < 6; i++ {
		for {
			char := inviteCodeBytes[rand.Intn(len(inviteCodeBytes))]
			if !usedChars[char] {
				b[i] = char
				usedChars[char] = true
				break
			}
		}
	}

	return string(b)
}

// GenerateUniqueInviteCode 生成唯一的邀请码
func GenerateUniqueInviteCode(checkExists func(string) (bool, error)) (string, error) {
	maxAttempts := 100 // 最大尝试次数，避免无限循环

	for i := 0; i < maxAttempts; i++ {
		inviteCode := GenerateInviteCode()
		exists, err := checkExists(inviteCode)
		if err != nil {
			return "", err
		}
		if !exists {
			return inviteCode, nil
		}
	}

	return "", NewAppError(CodeInviteCodeGenFailed, "无法生成唯一邀请码，请稍后重试")
}

// GenerateUniqueInviteCodeBatch 批量生成唯一邀请码（减少数据库查询）
func GenerateUniqueInviteCodeBatch(checkExistsBatch func([]string) (map[string]bool, error)) (string, error) {
	maxAttempts := 20 // 减少尝试次数，因为每次批量检查多个
	batchSize := 10   // 每次批量生成10个邀请码

	for i := 0; i < maxAttempts; i++ {
		// 批量生成邀请码
		codes := make([]string, batchSize)
		for j := 0; j < batchSize; j++ {
			codes[j] = GenerateInviteCode()
		}

		// 批量检查哪些已存在
		existsMap, err := checkExistsBatch(codes)
		if err != nil {
			return "", err
		}

		// 找到第一个不存在的邀请码
		for _, code := range codes {
			if !existsMap[code] {
				return code, nil
			}
		}
	}

	return "", NewAppError(CodeInviteCodeGenFailed, "无法生成唯一邀请码，请稍后重试")
}
