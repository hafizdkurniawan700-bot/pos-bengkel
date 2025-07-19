package models

import (
	"time"
)

type TestDrive struct {
	ID            int       `json:"id" db:"id"`
	CustomerID    int       `json:"customer_id" db:"customer_id" validate:"required"`
	VehicleID     int       `json:"vehicle_id" db:"vehicle_id" validate:"required"`
	ScheduledDate time.Time `json:"scheduled_date" db:"scheduled_date" validate:"required"`
	Status        string    `json:"status" db:"status" validate:"oneof=scheduled completed cancelled no_show"`
	Notes         string    `json:"notes" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	
	// Related entities
	Vehicle  *Vehicle  `json:"vehicle,omitempty"`
	Customer *Customer `json:"customer,omitempty"`
}

type TestDriveRequest struct {
	CustomerID    int       `json:"customer_id" validate:"required"`
	VehicleID     int       `json:"vehicle_id" validate:"required"`
	ScheduledDate time.Time `json:"scheduled_date" validate:"required"`
	Status        string    `json:"status" validate:"oneof=scheduled completed cancelled no_show"`
	Notes         string    `json:"notes"`
}

type TestDriveStatusUpdate struct {
	Status string `json:"status" validate:"required,oneof=scheduled completed cancelled no_show"`
	Notes  string `json:"notes"`
}

type TestDriveSearchParams struct {
	CustomerID int    `query:"customer_id"`
	VehicleID  int    `query:"vehicle_id"`
	Status     string `query:"status"`
	DateFrom   string `query:"date_from"`
	DateTo     string `query:"date_to"`
	Limit      int    `query:"limit"`
	Offset     int    `query:"offset"`
}

// ToTestDrive converts TestDriveRequest to TestDrive
func (tdr *TestDriveRequest) ToTestDrive() *TestDrive {
	if tdr.Status == "" {
		tdr.Status = "scheduled"
	}
	
	return &TestDrive{
		CustomerID:    tdr.CustomerID,
		VehicleID:     tdr.VehicleID,
		ScheduledDate: tdr.ScheduledDate,
		Status:        tdr.Status,
		Notes:         tdr.Notes,
	}
}