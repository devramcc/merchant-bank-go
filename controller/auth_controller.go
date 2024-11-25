package controller

import (
	"encoding/json"
	"net/http"

	"github.com/devramcc/merchant-bank-go/service"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		service: service.NewAuthService(),
	}
}

func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		customers := c.service.GetAll()
		json.NewEncoder(w).Encode(customers)

	case http.MethodPost:
		var customer service.Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
		}
		c.service.Register(customer)
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
		var customer service.Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		customerRet, err := c.service.Login(w, customer)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Login success",
			"customer": customerRet,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
