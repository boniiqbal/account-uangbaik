package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyCustomClaims struct {
	UserID    uint
	LoginType string
	jwt.StandardClaims
}

// GetExpiryTime for jwt
func GetExpiryTime(exp time.Duration) int64 {
	return time.Now().Add(time.Hour * exp).Unix()
}

// CreateClaims .
func CreateClaims(loginType string, userID uint, userAgent string, exp time.Duration) MyCustomClaims {
	return MyCustomClaims{
		userID,
		loginType,
		jwt.StandardClaims{
			ExpiresAt: GetExpiryTime(exp),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "api.yabaik.id",
		},
	}
}
