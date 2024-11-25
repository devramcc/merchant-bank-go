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
	customers   []Customer
	nextID      int
	hashService *HashService
}

func NewAuthService(hashService *HashService) *AuthService {
	return &AuthService{
		customers:   []Customer{},
		nextID:      1,
		hashService: hashService,
	}
}

func (s *AuthService) Register(customer Customer) error {
	hashedPassword, err := s.hashService.HashPassword(customer.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	customer.ID = s.nextID
	customer.Password = hashedPassword
	s.nextID++
	s.customers = append(s.customers, customer)
	return nil
}

func (s *AuthService) Login(w http.ResponseWriter, customerLog Customer) (Customer, error) {
	for _, customer := range s.customers {
		if customer.Username == customerLog.Username {
			err := s.hashService.CheckPassword(customer.Password, customerLog.Password)
			if err != nil {
				http.Error(w, "invalid credentials", http.StatusBadRequest)
				return Customer{}, fmt.Errorf("invalid credentials")
			}
			return customer, nil
		}
	}
	http.Error(w, "customer not found.", http.StatusBadRequest)
	return Customer{}, fmt.Errorf("customer not found")
}

func (s *AuthService) GetAll() []Customer {
	return s.customers
}
