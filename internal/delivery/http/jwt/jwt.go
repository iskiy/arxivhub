package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const minSecretSize = 32

type Payload struct {
	ID        int64     `json:"id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type JWTManager struct {
	secret string
}

func NewJWTManager(secret string) (*JWTManager, error) {
	if len(secret) < minSecretSize {
		return nil, fmt.Errorf("secret length is less than %d", minSecretSize)
	}
	return &JWTManager{
		secret: secret,
	}, nil
}

func (manager *JWTManager) CreateToken(userID int64, duration time.Duration) (string, *Claims, error) {
	claims := NewClaims(userID, duration)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(manager.secret))

	return token, claims, err
}

func (manager *JWTManager) ParseToken(signedToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(manager.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
