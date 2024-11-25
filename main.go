package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/devramcc/merchant-bank-go/controller"
	"github.com/devramcc/merchant-bank-go/middleware"
	"github.com/devramcc/merchant-bank-go/repository"
	"github.com/devramcc/merchant-bank-go/service"
)

func main() {
	mux := http.NewServeMux()

	// JSON Database
	customerFilePath := "./database/customers.json"
	whitelistAccessTokenFilePath := "./database/whitelistAccessTokens.json"
	paymentFilePath := "./database/payments.json"

	// Repository
	customerRepository := repository.NewCustomerRepository(customerFilePath)
	whitelistAccessTokenRepository := repository.NewWhitelistAccessTokenRepository(whitelistAccessTokenFilePath)
	paymentRepository := repository.NewPaymentRepository(paymentFilePath)

	// Service
	hashService := service.NewHashService()
	jwtService := service.NewJWTService("mysecretkey", time.Hour)
	authService := service.NewAuthService(customerRepository, whitelistAccessTokenRepository, hashService, jwtService)
	paymentService := service.NewPaymentService(paymentRepository)

	// Controller
	authController := controller.NewAuthController(authService)
	paymentController := controller.NewPaymentController(paymentService)

	// Routes
	mux.HandleFunc("/auth", authController.HandleRegister)
	mux.HandleFunc("/auth/login", authController.HandleLogin)
	mux.HandleFunc("/auth/logout", authController.HandleLogout)

	// Protected Route
	mux.HandleFunc("/protected", middleware.JWTMiddleware(jwtService, whitelistAccessTokenRepository, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "protected route")
	}))
	mux.HandleFunc("/payment", middleware.JWTMiddleware(jwtService, whitelistAccessTokenRepository, paymentController.HandlePayment))

	// Start Server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
