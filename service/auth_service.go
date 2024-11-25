package service

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

func (s *AuthService) GetAll() []Customer {
	return s.customers
}
