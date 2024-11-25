package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/devramcc/merchant-bank-go/controller"
	"github.com/devramcc/merchant-bank-go/middleware"
	"github.com/devramcc/merchant-bank-go/service"
)

func main() {
	mux := http.NewServeMux()

	// JSON Database
	customerFilePath := "./database/customers.json"
	whitelistAccessTokenFilePath := "./database/whitelistAccessToken.json"

	// Service
	hashService := &service.HashService{}
	jwtService := service.NewJWTService("mysecretkey", time.Hour)
	authService := service.NewAuthService(customerFilePath, whitelistAccessTokenFilePath, hashService, jwtService)

	// Controller
	authController := controller.NewAuthController(authService)

	// Routes
	mux.HandleFunc("/auth", authController.HandleRegister)
	mux.HandleFunc("/auth/login", authController.HandleLogin)
	mux.HandleFunc("/auth/logout", authController.HandleLogout)

	// Protected Route
	mux.HandleFunc("/protected", middleware.JWTMiddleware(jwtService, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "protected route")
	}))

	// Start Server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
