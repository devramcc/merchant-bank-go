package controller

import (
	"encoding/json"
	"net/http"

	"github.com/devramcc/merchant-bank-go/model"
	"github.com/devramcc/merchant-bank-go/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		customers := c.authService.GetAll()
		json.NewEncoder(w).Encode(customers)

	case http.MethodPost:
		var customer model.Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := c.authService.Register(customer); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Customer created successfully"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var customer model.Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		token, err := c.authService.Login(w, customer)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Login success",
			"token":   token,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Authorization header is missing or invalid", http.StatusBadRequest)
			return
		}

		token := authHeader[7:]

		c.authService.Logout(w, token)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
