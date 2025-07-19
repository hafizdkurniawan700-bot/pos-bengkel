package models

import (
	"time"
)

type Customer struct {
	ID        int       `json:"id" db:"id"`
	UserID    *int      `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name" validate:"required,max=100"`
	Phone     string    `json:"phone" db:"phone" validate:"max=20"`
	Address   string    `json:"address" db:"address"`
	NIK       string    `json:"nik" db:"nik" validate:"max=20"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	
	// Related user information if exists
	User *UserResponse `json:"user,omitempty"`
}

type CustomerRequest struct {
	UserID  *int   `json:"user_id"`
	Name    string `json:"name" validate:"required,max=100"`
	Phone   string `json:"phone" validate:"max=20"`
	Address string `json:"address"`
	NIK     string `json:"nik" validate:"max=20"`
}

type CustomerSearchParams struct {
	Query  string `query:"q"`
	Name   string `query:"name"`
	Phone  string `query:"phone"`
	NIK    string `query:"nik"`
	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
}

// ToCustomer converts CustomerRequest to Customer
func (cr *CustomerRequest) ToCustomer() *Customer {
	return &Customer{
		UserID:  cr.UserID,
		Name:    cr.Name,
		Phone:   cr.Phone,
		Address: cr.Address,
		NIK:     cr.NIK,
	}
}