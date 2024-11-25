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
	jwtService  *JWTService
}

func NewAuthService(hashService *HashService, jwtService *JWTService) *AuthService {
	return &AuthService{
		customers:   []Customer{},
		nextID:      1,
		hashService: hashService,
		jwtService:  jwtService,
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

func (s *AuthService) Login(w http.ResponseWriter, customerLog Customer) (string, error) {
	for _, customer := range s.customers {
		if customer.Username == customerLog.Username {
			if err := s.hashService.CheckPassword(customer.Password, customerLog.Password); err == nil {
				token, err := s.jwtService.GenerateToken(customer.ID, customer.Username)
				if err != nil {
					http.Error(w, "Failed to generate token", http.StatusInternalServerError)
					return "", err
				}
				return token, nil
			} else {
				http.Error(w, "invalid credentials", http.StatusBadRequest)
				return "", fmt.Errorf("invalid credentials")
			}
		}
	}
	http.Error(w, "customer not found.", http.StatusBadRequest)
	return "", fmt.Errorf("customer not found")
}

func (s *AuthService) GetAll() []Customer {
	return s.customers
}
