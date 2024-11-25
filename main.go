package main

import (
	"log"
	"net/http"

	"github.com/devramcc/merchant-bank-go/controller"
	"github.com/devramcc/merchant-bank-go/service"
)

func main() {
	mux := http.NewServeMux()

	// Service
	hashService := &service.HashService{}
	authService := service.NewAuthService(hashService)

	// Controller
	authController := controller.NewAuthController(authService)

	// Route
	mux.HandleFunc("/auth", authController.HandleRegister)
	mux.HandleFunc("/auth/login", authController.HandleLogin)

	// Start Server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
