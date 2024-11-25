package repository

import (
	"encoding/json"
	"os"

	"github.com/devramcc/merchant-bank-go/model"
)

type PaymentRepository struct {
	filePath string
}

func NewPaymentRepository(filePath string) *PaymentRepository {
	return &PaymentRepository{
		filePath: filePath,
	}
}

func (r *PaymentRepository) LoadPayments() ([]model.Payment, error) {
	var payments []model.Payment
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return payments, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) SavePayments(payments []model.Payment) error {
	data, err := json.MarshalIndent(payments, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return err
	}
	return nil
}
