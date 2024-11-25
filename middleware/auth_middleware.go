package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/devramcc/merchant-bank-go/service"
)

func JWTMiddleware(jwtService *service.JWTService, next http.HandlerFunc) http.HandlerFunc {
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

		whitelistData, err := os.ReadFile("./database/whitelistAccessToken.json")
		if err != nil {
			http.Error(w, "Failed to read whitelist", http.StatusInternalServerError)
			return
		}

		var whitelistTokens []string
		if err := json.Unmarshal(whitelistData, &whitelistTokens); err != nil {
			http.Error(w, "Failed to parse whitelist", http.StatusInternalServerError)
			return
		}

		tokenExists := false
		for _, t := range whitelistTokens {
			if t == tokenString {
				tokenExists = true
				break
			}
		}

		if !tokenExists {
			http.Error(w, "Unauthorized: Token not whitelisted", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
