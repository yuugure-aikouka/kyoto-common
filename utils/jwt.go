package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateAuthToken(userId int32, expiryDuration time.Duration, key string) (string, error) {
	now := time.Now()
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     userId,
		"exp":     now.Add(expiryDuration).Unix(),
		"iat":     now.Unix(),
		"jti":     uuid.New().String(),
		"user_id": userId,
	})

	return authToken.SignedString([]byte(key))
}

func ParseAuthToken(tokenString, key string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
}
