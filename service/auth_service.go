package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/devramcc/merchant-bank-go/model"
	"github.com/devramcc/merchant-bank-go/repository"
)

type AuthService struct {
	customers                      []model.Customer
	whitelistAccessTokens          []string
	nextID                         int
	customerRepository             *repository.CustomerRepository
	whitelistAccessTokenRepository *repository.WhitelistAccessTokenRepository
	hashService                    *HashService
	jwtService                     *JWTService
}

func NewAuthService(customerRepository *repository.CustomerRepository, whitelistAccessTokenRepository *repository.WhitelistAccessTokenRepository, hashService *HashService, jwtService *JWTService) *AuthService {
	service := &AuthService{
		customers:                      []model.Customer{},
		whitelistAccessTokens:          []string{},
		nextID:                         1,
		customerRepository:             customerRepository,
		whitelistAccessTokenRepository: whitelistAccessTokenRepository,
		hashService:                    hashService,
		jwtService:                     jwtService,
	}
	service.loadCustomers()
	service.loadWhitelistTokens()
	return service
}

func (s *AuthService) loadCustomers() {
	customers, err := s.customerRepository.LoadCustomers()
	if err != nil {
		log.Fatalf("Failed to load customers: %v", err)
	}
	s.customers = customers
	for _, customer := range s.customers {
		if customer.ID >= s.nextID {
			s.nextID = customer.ID + 1
		}
	}
}

func (s *AuthService) saveCustomers() {
	if err := s.customerRepository.SaveCustomers(s.customers); err != nil {
		log.Fatalf("Failed to save customers: %v", err)
	}
}

func (s *AuthService) loadWhitelistTokens() {
	data, err := os.ReadFile("./database/whitelistAccessToken.json")
	if err != nil {
		if os.IsNotExist(err) {
			s.whitelistAccessTokens = []string{}
			return
		}
		log.Fatalf("Failed to read whitelist file: %v", err)
	}
	if err := json.Unmarshal(data, &s.whitelistAccessTokens); err != nil {
		log.Fatalf("Failed to parse whitelist file: %v", err)
	}
}

func (s *AuthService) saveWhitelistAccessToken(token string) {
	if err := s.whitelistAccessTokenRepository.SaveWhitelistToken(token); err != nil {
		log.Printf("Failed to save whitelist token: %v", err)
	}
}

func (s *AuthService) removeWhitelistAccessToken(token string) {
	if err := s.whitelistAccessTokenRepository.RemoveWhitelistToken(token); err != nil {
		log.Printf("Failed to remove whitelist token: %v", err)
	}
}

func (s *AuthService) Register(customer model.Customer) error {
	hashedPassword, err := s.hashService.HashPassword(customer.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	customer.ID = s.nextID
	customer.Password = hashedPassword
	s.nextID++
	s.customers = append(s.customers, customer)

	s.saveCustomers()

	return nil
}

func (s *AuthService) Login(w http.ResponseWriter, customerLog model.Customer) (string, error) {
	for _, customer := range s.customers {
		if customer.Username == customerLog.Username {
			if err := s.hashService.CheckPassword(customer.Password, customerLog.Password); err == nil {
				token, err := s.jwtService.GenerateToken(customer.ID, customer.Username)
				if err != nil {
					http.Error(w, "Failed to generate token", http.StatusInternalServerError)
					return "", err
				}
				s.saveWhitelistAccessToken(token)
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

func (s *AuthService) Logout(w http.ResponseWriter, token string) {
	s.removeWhitelistAccessToken(token)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "logout success")
}

func (s *AuthService) GetAll() []model.Customer {
	return s.customers
}
