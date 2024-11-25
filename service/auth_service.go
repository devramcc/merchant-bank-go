package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Customer struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthService struct {
	customers             []Customer
	whitelistAccessTokens []string
	nextID                int
	filePath              string
	hashService           *HashService
	jwtService            *JWTService
}

func NewAuthService(filePath string, whitelistAccessTokenFilePath string, hashService *HashService, jwtService *JWTService) *AuthService {
	service := &AuthService{
		customers:             []Customer{},
		whitelistAccessTokens: []string{},
		nextID:                1,
		filePath:              filePath,
		hashService:           hashService,
		jwtService:            jwtService,
	}
	service.loadCustomers()
	service.loadWhitelistTokens()
	return service
}

func (s *AuthService) loadCustomers() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatalf("Failed to read customers file: %v", err)
	}
	if err := json.Unmarshal(data, &s.customers); err != nil {
		log.Fatalf("Failed to parse customers file: %v", err)
	}

	for _, customer := range s.customers {
		if customer.ID >= s.nextID {
			s.nextID = customer.ID + 1
		}
	}
}

func (s *AuthService) saveCustomers() {
	data, err := json.MarshalIndent(s.customers, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal customers: %v", err)
	}
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		log.Fatalf("Failed to write customers file: %v", err)
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
	s.whitelistAccessTokens = append(s.whitelistAccessTokens, token)
	data, err := json.MarshalIndent(s.whitelistAccessTokens, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal whitelist tokens: %v", err)
		return
	}
	if err := os.WriteFile("./database/whitelistAccessToken.json", data, 0644); err != nil {
		log.Printf("Failed to write whitelist tokens file: %v", err)
	}
}

func (s *AuthService) removeWhitelistAccessToken(token string) {
	newTokens := []string{}
	for _, t := range s.whitelistAccessTokens {
		if t != token {
			newTokens = append(newTokens, t)
		}
	}
	s.whitelistAccessTokens = newTokens

	data, err := json.MarshalIndent(newTokens, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal whitelist tokens: %v", err)
		return
	}
	if err := os.WriteFile("./database/whitelistAccessToken.json", data, 0644); err != nil {
		log.Printf("Failed to write whitelist tokens file: %v", err)
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

	s.saveCustomers()

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

func (s *AuthService) GetAll() []Customer {
	return s.customers
}
