package main

import (
	"log"
	"net/http"

	"github.com/devramcc/merchant-bank-go/controller"
)

func main() {
	mux := http.NewServeMux()

	authController := controller.NewAuthController()

	mux.HandleFunc("/auth", authController.HandleRegister)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
