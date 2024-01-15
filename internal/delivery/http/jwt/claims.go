package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}

func NewClaims(userID int64, duration time.Duration) *Claims {
	return &Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}
}
