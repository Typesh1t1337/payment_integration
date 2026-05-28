package middleware

import (
	"context"
	"net/http"
	"payment_integration/internal/a_user/service"
	"payment_integration/internal/transport/http_transport"
	"strings"
	"time"
)

func AuthMiddleware(jwt *service.JwtService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Header.Get("Authorization")
			if accessToken == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			accessToken = strings.TrimPrefix(accessToken, "Bearer ")
			if accessToken == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			claims, err := jwt.ParseToken(accessToken)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if claims.Exp < time.Now().Unix() {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if claims.Type != "access" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), http_transport.UserIDContextKey, claims.Sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}