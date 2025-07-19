package controllers

import (
	"database/sql"
	"log"
	"pos-bengkel-backend/config"
	"pos-bengkel-backend/middleware"
	"pos-bengkel-backend/models"
	"pos-bengkel-backend/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type VehicleController struct{}

func NewVehicleController() *VehicleController {
	return &VehicleController{}
}

// GetVehicles returns a list of vehicles with optional filtering
func (vc *VehicleController) GetVehicles(c *fiber.Ctx) error {
	var params models.VehicleSearchParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	// Set default pagination
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	// For simplicity, let's start with a basic query without complex filtering
	var vehicles []models.Vehicle
	query := "SELECT * FROM vehicles ORDER BY created_at DESC LIMIT ? OFFSET ?"
	err := config.DB.Select(&vehicles, query, params.Limit, params.Offset)
	if err != nil {
		log.Printf("Error querying vehicles: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM vehicles"
	err = config.DB.Get(&total, countQuery)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	page := (params.Offset / params.Limit) + 1
	return utils.PaginatedSuccessResponse(c, "Vehicles retrieved successfully", vehicles, page, params.Limit, total)
}

// GetVehicle returns a single vehicle by ID
func (vc *VehicleController) GetVehicle(c *fiber.Ctx) error {
	id := c.Params("id")
	vehicleID, err := strconv.Atoi(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid vehicle ID")
	}

	var vehicle models.Vehicle
	query := "SELECT * FROM vehicles WHERE id = $1"
	err = config.DB.Get(&vehicle, query, vehicleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "Vehicle")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	return utils.SuccessResponse(c, "Vehicle retrieved successfully", vehicle)
}

// CreateVehicle creates a new vehicle
func (vc *VehicleController) CreateVehicle(c *fiber.Ctx) error {
	var req models.VehicleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
	}

	// Convert to vehicle model
	vehicle := req.ToVehicle()

	// Insert vehicle
	var vehicleID int
	query := `
		INSERT INTO vehicles (brand, model, year, price, status, images, description, engine_type, transmission, fuel_type, mileage, color) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) 
		RETURNING id
	`
	err := config.DB.Get(&vehicleID, query,
		vehicle.Brand, vehicle.Model, vehicle.Year, vehicle.Price, vehicle.Status,
		vehicle.Images, vehicle.Description, vehicle.EngineType, vehicle.Transmission,
		vehicle.FuelType, vehicle.Mileage, vehicle.Color)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create vehicle")
	}

	// Get the created vehicle
	getQuery := "SELECT * FROM vehicles WHERE id = $1"
	err = config.DB.Get(vehicle, getQuery, vehicleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve created vehicle")
	}

	return utils.CreatedResponse(c, "Vehicle created successfully", vehicle)
}

// UpdateVehicle updates an existing vehicle
func (vc *VehicleController) UpdateVehicle(c *fiber.Ctx) error {
	id := c.Params("id")
	vehicleID, err := strconv.Atoi(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid vehicle ID")
	}

	var req models.VehicleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
	}

	// Check if vehicle exists
	var existingVehicle models.Vehicle
	checkQuery := "SELECT id FROM vehicles WHERE id = $1"
	err = config.DB.Get(&existingVehicle, checkQuery, vehicleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "Vehicle")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Update vehicle
	updateQuery := `
		UPDATE vehicles SET 
			brand = $1, model = $2, year = $3, price = $4, status = $5, 
			images = $6, description = $7, engine_type = $8, transmission = $9, 
			fuel_type = $10, mileage = $11, color = $12, updated_at = CURRENT_TIMESTAMP
		WHERE id = $13
	`
	_, err = config.DB.Exec(updateQuery,
		req.Brand, req.Model, req.Year, req.Price, req.Status,
		req.Images, req.Description, req.EngineType, req.Transmission,
		req.FuelType, req.Mileage, req.Color, vehicleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update vehicle")
	}

	// Get updated vehicle
	var updatedVehicle models.Vehicle
	getQuery := "SELECT * FROM vehicles WHERE id = $1"
	err = config.DB.Get(&updatedVehicle, getQuery, vehicleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve updated vehicle")
	}

	return utils.SuccessResponse(c, "Vehicle updated successfully", updatedVehicle)
}

// DeleteVehicle deletes a vehicle
func (vc *VehicleController) DeleteVehicle(c *fiber.Ctx) error {
	// Check if user has admin role
	user := middleware.GetUserFromContext(c)
	if user.Role != "admin" {
		return utils.ForbiddenResponse(c, "Only admins can delete vehicles")
	}

	id := c.Params("id")
	vehicleID, err := strconv.Atoi(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid vehicle ID")
	}

	// Check if vehicle exists
	var vehicle models.Vehicle
	checkQuery := "SELECT id FROM vehicles WHERE id = $1"
	err = config.DB.Get(&vehicle, checkQuery, vehicleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "Vehicle")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Delete vehicle
	deleteQuery := "DELETE FROM vehicles WHERE id = $1"
	_, err = config.DB.Exec(deleteQuery, vehicleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete vehicle")
	}

	return utils.SuccessResponse(c, "Vehicle deleted successfully", nil)
}

// SearchVehicles searches vehicles by query
func (vc *VehicleController) SearchVehicles(c *fiber.Ctx) error {
	query := c.Query("q", "")
	if strings.TrimSpace(query) == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Search query is required")
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	var vehicles []models.Vehicle
	searchQuery := `
		SELECT * FROM vehicles 
		WHERE brand ILIKE $1 OR model ILIKE $1 OR description ILIKE $1 
		ORDER BY 
			CASE 
				WHEN brand ILIKE $1 THEN 1
				WHEN model ILIKE $1 THEN 2
				ELSE 3
			END,
			created_at DESC
		LIMIT $2
	`
	
	err := config.DB.Select(&vehicles, searchQuery, "%"+query+"%", limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	return utils.SuccessResponse(c, "Search completed successfully", vehicles)
}