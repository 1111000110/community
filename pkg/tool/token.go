package tool

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 定义自定义声明
type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// 通用Token生成函数
func CreateToken(userID int64, duration time.Duration, SecretKey string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey)) // 将 SecretKey 转换为 []byte
}
