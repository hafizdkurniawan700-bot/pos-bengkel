package models

import (
	"time"
)

type Vehicle struct {
	ID           int     `json:"id" db:"id"`
	Brand        string  `json:"brand" db:"brand" validate:"required,max=50"`
	Model        string  `json:"model" db:"model" validate:"required,max=50"`
	Year         int     `json:"year" db:"year" validate:"required,min=1900,max=2030"`
	Price        float64 `json:"price" db:"price" validate:"required,min=0"`
	Status       string  `json:"status" db:"status" validate:"oneof=available sold reserved maintenance"`
	Images       *string `json:"images" db:"images"` // JSON array of image URLs - can be NULL
	Description  *string `json:"description" db:"description"`
	EngineType   *string `json:"engine_type" db:"engine_type"`
	Transmission *string `json:"transmission" db:"transmission"`
	FuelType     *string `json:"fuel_type" db:"fuel_type"`
	Mileage      int     `json:"mileage" db:"mileage" validate:"min=0"`
	Color        *string `json:"color" db:"color"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type VehicleRequest struct {
	Brand        string  `json:"brand" validate:"required,max=50"`
	Model        string  `json:"model" validate:"required,max=50"`
	Year         int     `json:"year" validate:"required,min=1900,max=2030"`
	Price        float64 `json:"price" validate:"required,min=0"`
	Status       string  `json:"status" validate:"oneof=available sold reserved maintenance"`
	Images       *string `json:"images"`
	Description  *string `json:"description"`
	EngineType   *string `json:"engine_type"`
	Transmission *string `json:"transmission"`
	FuelType     *string `json:"fuel_type"`
	Mileage      int     `json:"mileage" validate:"min=0"`
	Color        *string `json:"color"`
}

type VehicleSearchParams struct {
	Query        string  `query:"q"`
	Brand        string  `query:"brand"`
	Model        string  `query:"model"`
	YearMin      int     `query:"year_min"`
	YearMax      int     `query:"year_max"`
	PriceMin     float64 `query:"price_min"`
	PriceMax     float64 `query:"price_max"`
	Status       string  `query:"status"`
	FuelType     string  `query:"fuel_type"`
	Transmission string  `query:"transmission"`
	Limit        int     `query:"limit"`
	Offset       int     `query:"offset"`
}

// ToVehicle converts VehicleRequest to Vehicle
func (vr *VehicleRequest) ToVehicle() *Vehicle {
	if vr.Status == "" {
		vr.Status = "available"
	}
	
	return &Vehicle{
		Brand:        vr.Brand,
		Model:        vr.Model,
		Year:         vr.Year,
		Price:        vr.Price,
		Status:       vr.Status,
		Images:       vr.Images,
		Description:  vr.Description,
		EngineType:   vr.EngineType,
		Transmission: vr.Transmission,
		FuelType:     vr.FuelType,
		Mileage:      vr.Mileage,
		Color:        vr.Color,
	}
}