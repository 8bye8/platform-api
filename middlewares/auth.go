package middlewares

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const AuthUserId = "middleware.auth.userId"

func RequireAuth(publicKey *rsa.PublicKey) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if !strings.HasPrefix(authorization, "Bearer ") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}

			rawToken := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(rawToken, claims, func(tk *jwt.Token) (interface{}, error) {
				if tk.Method != jwt.SigningMethodRS256 {
					return nil, fmt.Errorf("unexpected signing method: %v", tk.Header["alg"])
				}
				return publicKey, nil
			})
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			if !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			userId, ok := claims["sub"].(string)
			if !ok || userId == "" {
				http.Error(w, "token missing subject", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthUserId, userId)
			req := r.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	}
}
