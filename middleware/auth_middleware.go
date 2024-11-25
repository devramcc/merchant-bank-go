package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/devramcc/merchant-bank-go/repository"
	"github.com/devramcc/merchant-bank-go/service"
	"github.com/golang-jwt/jwt/v5"
)

const CustomerIDKey = "customer_id"

func JWTMiddleware(jwtService *service.JWTService, whitelistAccessTokenRepository *repository.WhitelistAccessTokenRepository, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		customerID, ok := claims["customer_id"].(float64)
		if !ok {
			http.Error(w, "Unauthorized: customer_id not found in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), CustomerIDKey, int(customerID))
		r = r.WithContext(ctx)

		if !whitelistAccessTokenRepository.IsTokenWhitelisted(tokenString) {
			http.Error(w, "Unauthorized: Token not whitelisted", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
