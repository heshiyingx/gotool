package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CustomClaims[T any] struct {
	JwtPayLoad T
	jwt.RegisteredClaims
}

// GenToken 创建 Token
func GenToken[T any](payLoad T, accessSecret string, during time.Duration) (string, error) {
	claim := CustomClaims[T]{
		JwtPayLoad: payLoad,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(during)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(accessSecret))
}

// ParseToken 解析 token
func ParseToken[T any](tokenStr string, accessSecret string) (*CustomClaims[T], error) {

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims[T]{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims[T]); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")

}

// 判断是否过期
func (c *CustomClaims[T]) IsExpired() bool {
	return c.ExpiresAt.Time.Before(time.Now())
}
