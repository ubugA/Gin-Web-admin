package utils

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// stringReplaceAll 替换所有匹配的目标字符串
func stringReplaceAll(originalStr string, targets, replacements []string) string {
	for i := 0; i < len(targets); i++ {
		originalStr = strings.ReplaceAll(originalStr, targets[i], replacements[i])
	}

	return originalStr
}

// GenerateHashedPassword 生成密码
func GenerateHashedPassword(password string) (string, error) {
	salt := "4Q&6p4euU=Yx"
	passwordWithSalt := password + salt

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// VerifyHashedPassword 验证密码
func VerifyHashedPassword(hashedPassword, password string) bool {
	salt := "4Q&6p4euU=Yx"
	passwordWithSalt := password + salt

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword+salt), []byte(passwordWithSalt))
	return err == nil
}
