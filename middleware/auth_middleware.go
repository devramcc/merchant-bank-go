package middleware

import (
	"net/http"
	"strings"

	"github.com/devramcc/merchant-bank-go/repository"
	"github.com/devramcc/merchant-bank-go/service"
)

func JWTMiddleware(jwtService *service.JWTService, whitelistAccessTokenRepository *repository.WhitelistAccessTokenRepository, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !whitelistAccessTokenRepository.IsTokenWhitelisted(tokenString) {
			http.Error(w, "Unauthorized: Token not whitelisted", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
