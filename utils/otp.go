package utils

import (
	"serra/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateOTPToken(email, otp string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"otp":   otp,
		"exp":   time.Now().Add(ttl).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Envs.JWTSecret))
}

func VerifyOTPToken(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(config.Envs.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", jwt.ErrTokenMalformed
	}

	email, _ := claims["email"].(string)
	otp, _ := claims["otp"].(string)

	return email, otp, nil
}
