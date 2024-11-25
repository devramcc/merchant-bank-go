package service

import (
	"log"

	"github.com/devramcc/merchant-bank-go/model"
	"github.com/devramcc/merchant-bank-go/repository"
)

type PaymentService struct {
	payments          []model.Payment
	nextID            int
	paymentRepository *repository.PaymentRepository
}

func NewPaymentService(paymentRepository *repository.PaymentRepository) *PaymentService {
	service := &PaymentService{
		payments:          []model.Payment{},
		nextID:            1,
		paymentRepository: paymentRepository,
	}
	service.loadPayments()
	return service
}

func (s *PaymentService) loadPayments() {
	payments, err := s.paymentRepository.LoadPayments()
	if err != nil {
		log.Fatalf("Failed to load payments: %v", err)
	}
	s.payments = payments
	for _, payment := range s.payments {
		if payment.ID >= s.nextID {
			s.nextID = payment.ID + 1
		}
	}
}

func (s *PaymentService) SavePayments(payment model.Payment) error {
	payment.ID = s.nextID
	s.nextID++
	s.payments = append(s.payments, payment)
	if err := s.paymentRepository.SavePayments(s.payments); err != nil {
		log.Fatalf("Failed to save payments: %v", err)
		return err
	}
	return nil
}
