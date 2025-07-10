package tool

import (
	"golang.org/x/crypto/bcrypt"
)

// 哈希密码
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}
	// 生成哈希值
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 比较密码和哈希值
func ComparePassword(hashedPassword, password string) bool {
	// 比较密码和哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
