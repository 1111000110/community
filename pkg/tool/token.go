package tool

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
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

func GetUserId(ctx context.Context) (int64, error) {
	userIdVal := ctx.Value("user_id")
	if userIdVal == nil {
		return 0, errors.New("user_id not found in context")
	}
	switch v := userIdVal.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case json.Number:
		id, err := v.Int64()
		if err != nil {
			return 0, errors.Wrap(err, "invalid user_id (json.Number)")
		}
		return id, nil
	default:
		return 0, errors.Errorf("unsupported user_id type: %T", v)
	}
}
