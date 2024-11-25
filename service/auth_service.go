package service

import (
	"fmt"
	"net/http"
)

type Customer struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthService struct {
	customers []Customer
	nextID    int
}

func NewAuthService() AuthService {
	return AuthService{
		customers: []Customer{},
		nextID:    1,
	}
}

func (s *AuthService) Register(customer Customer) {
	customer.ID = s.nextID
	s.nextID++
	s.customers = append(s.customers, customer)
}

func (s *AuthService) Login(w http.ResponseWriter, customerLog Customer) (Customer, error) {
	for _, customer := range s.customers {
		if customer.Username == customerLog.Username {
			if customer.Password == customerLog.Password {
				return customer, nil
			} else {
				http.Error(w, "invalid credential.", http.StatusBadRequest)
				return Customer{}, fmt.Errorf("invalid credential")
			}
		}
	}
	http.Error(w, "customer not found.", http.StatusBadRequest)
	return Customer{}, fmt.Errorf("customer not found")
}

func (s *AuthService) GetAll() []Customer {
	return s.customers
}
