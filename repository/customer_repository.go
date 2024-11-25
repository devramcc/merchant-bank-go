package repository

import (
	"encoding/json"
	"os"

	"github.com/devramcc/merchant-bank-go/model"
)

type CustomerRepository struct {
	filePath string
}

func NewCustomerRepository(filePath string) *CustomerRepository {
	return &CustomerRepository{
		filePath: filePath,
	}
}

func (r *CustomerRepository) LoadCustomers() ([]model.Customer, error) {
	var customers []model.Customer
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return customers, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &customers); err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *CustomerRepository) SaveCustomers(customers []model.Customer) error {
	data, err := json.MarshalIndent(customers, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return err
	}
	return nil
}
