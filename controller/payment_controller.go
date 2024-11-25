package controller

import (
	"encoding/json"
	"net/http"

	"github.com/devramcc/merchant-bank-go/middleware"
	"github.com/devramcc/merchant-bank-go/model"
	"github.com/devramcc/merchant-bank-go/service"
)

type PaymentController struct {
	paymentService *service.PaymentService
}

func NewPaymentController(paymentService *service.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

func (c *PaymentController) HandlePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	customerID, ok := r.Context().Value(middleware.CustomerIDKey).(int)
	if !ok {
		http.Error(w, "Customer ID not found in context", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		payments, err := c.paymentService.GetCurrentCustomerPayments(customerID)
		if err != nil {
			http.Error(w, "Unable to retrieve payments", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(payments)

	case http.MethodPost:
		var payment model.Payment
		if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := c.paymentService.SavePayments(payment); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Payment created successfully"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
