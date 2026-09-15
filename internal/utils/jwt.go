package utils

import (
	"errors"
	"time"

	"go-vue-admin/internal/global"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint     `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateToken issues a signed JWT carrying user id, username and role keywords.
func GenerateToken(userID uint, username string, roles []string) (string, error) {
	expire, err := time.ParseDuration(global.Config.JWT.Expire)
	if err != nil || expire <= 0 {
		expire = 72 * time.Hour
	}
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-vue-admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.Config.JWT.Secret))
}

func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
