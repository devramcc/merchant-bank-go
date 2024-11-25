package model

type Payment struct {
	ID         int `json:"id"`
	CustomerID int `json:"customer_id"`
	Amount     int `json:"amount"`
}
