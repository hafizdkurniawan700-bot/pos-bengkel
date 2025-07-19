package models

import (
	"time"
)

type Transaction struct {
	ID            int       `json:"id" db:"id"`
	VehicleID     int       `json:"vehicle_id" db:"vehicle_id" validate:"required"`
	CustomerID    int       `json:"customer_id" db:"customer_id" validate:"required"`
	Amount        float64   `json:"amount" db:"amount" validate:"required,min=0"`
	Status        string    `json:"status" db:"status" validate:"oneof=pending completed cancelled refunded"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	Notes         string    `json:"notes" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	
	// Related entities
	Vehicle  *Vehicle  `json:"vehicle,omitempty"`
	Customer *Customer `json:"customer,omitempty"`
}

type TransactionRequest struct {
	VehicleID     int     `json:"vehicle_id" validate:"required"`
	CustomerID    int     `json:"customer_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,min=0"`
	Status        string  `json:"status" validate:"oneof=pending completed cancelled refunded"`
	PaymentMethod string  `json:"payment_method"`
	Notes         string  `json:"notes"`
}

type TransactionStatusUpdate struct {
	Status string `json:"status" validate:"required,oneof=pending completed cancelled refunded"`
	Notes  string `json:"notes"`
}

type TransactionSearchParams struct {
	CustomerID int    `query:"customer_id"`
	VehicleID  int    `query:"vehicle_id"`
	Status     string `query:"status"`
	DateFrom   string `query:"date_from"`
	DateTo     string `query:"date_to"`
	Limit      int    `query:"limit"`
	Offset     int    `query:"offset"`
}

// ToTransaction converts TransactionRequest to Transaction
func (tr *TransactionRequest) ToTransaction() *Transaction {
	if tr.Status == "" {
		tr.Status = "pending"
	}
	
	return &Transaction{
		VehicleID:     tr.VehicleID,
		CustomerID:    tr.CustomerID,
		Amount:        tr.Amount,
		Status:        tr.Status,
		PaymentMethod: tr.PaymentMethod,
		Notes:         tr.Notes,
	}
}