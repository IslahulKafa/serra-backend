package utils

import (
	"context"
	"errors"
	"net/http"
	"serra/config"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey = contextKey("user_id")

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			WriteError(w, http.StatusUnauthorized, errors.New("missing or invalid Authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(config.Envs.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			WriteError(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			WriteError(w, http.StatusUnauthorized, errors.New("invalid token claims"))
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			WriteError(w, http.StatusUnauthorized, errors.New("user_id not found in token"))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, int64(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
