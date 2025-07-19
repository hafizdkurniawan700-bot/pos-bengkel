package controllers

import (
	"database/sql"
	"fmt"
	"pos-bengkel-backend/config"
	"pos-bengkel-backend/middleware"
	"pos-bengkel-backend/models"
	"pos-bengkel-backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CustomerController struct{}

func NewCustomerController() *CustomerController {
	return &CustomerController{}
}

// GetCustomers returns a list of customers with optional filtering
func (cc *CustomerController) GetCustomers(c *fiber.Ctx) error {
	var params models.CustomerSearchParams
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

	// Build query
	baseQuery := `
		FROM customers c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE 1=1
	`
	var args []interface{}
	argCounter := 0

	// Add filters
	if params.Query != "" {
		argCounter++
		baseQuery += fmt.Sprintf(" AND (c.name ILIKE $%d OR c.phone ILIKE $%d OR c.nik ILIKE $%d)", argCounter, argCounter, argCounter)
		args = append(args, "%"+params.Query+"%")
	}
	if params.Name != "" {
		argCounter++
		baseQuery += fmt.Sprintf(" AND c.name ILIKE $%d", argCounter)
		args = append(args, "%"+params.Name+"%")
	}
	if params.Phone != "" {
		argCounter++
		baseQuery += fmt.Sprintf(" AND c.phone ILIKE $%d", argCounter)
		args = append(args, "%"+params.Phone+"%")
	}
	if params.NIK != "" {
		argCounter++
		baseQuery += fmt.Sprintf(" AND c.nik = $%d", argCounter)
		args = append(args, params.NIK)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := config.DB.Get(&total, countQuery, args...)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Get customers with user info
	selectQuery := `
		SELECT 
			c.id, c.user_id, c.name, c.phone, c.address, c.nik, c.created_at, c.updated_at,
			u.id as "user.id", u.username as "user.username", u.email as "user.email", 
			u.role as "user.role", u.created_at as "user.created_at", u.updated_at as "user.updated_at"
		` + baseQuery + " ORDER BY c.created_at DESC LIMIT $" + strconv.Itoa(argCounter+1) + " OFFSET $" + strconv.Itoa(argCounter+2)
	args = append(args, params.Limit, params.Offset)

	rows, err := config.DB.Query(selectQuery, args...)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		var user models.UserResponse
		var userID, userCreatedAt, userUpdatedAt sql.NullString
		var username, email, role sql.NullString

		err := rows.Scan(
			&customer.ID, &customer.UserID, &customer.Name, &customer.Phone, &customer.Address, &customer.NIK, &customer.CreatedAt, &customer.UpdatedAt,
			&userID, &username, &email, &role, &userCreatedAt, &userUpdatedAt,
		)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error scanning customer data")
		}

		// Add user info if exists
		if userID.Valid {
			var userIDInt int
			fmt.Sscanf(userID.String, "%d", &userIDInt)
			user.ID = userIDInt
			user.Username = username.String
			user.Email = email.String
			user.Role = role.String
			customer.User = &user
		}

		customers = append(customers, customer)
	}

	page := (params.Offset / params.Limit) + 1
	return utils.PaginatedSuccessResponse(c, "Customers retrieved successfully", customers, page, params.Limit, total)
}

// GetCustomer returns a single customer by ID
func (cc *CustomerController) GetCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid customer ID")
	}

	query := `
		SELECT 
			c.id, c.user_id, c.name, c.phone, c.address, c.nik, c.created_at, c.updated_at,
			u.id as "user.id", u.username as "user.username", u.email as "user.email", 
			u.role as "user.role", u.created_at as "user.created_at", u.updated_at as "user.updated_at"
		FROM customers c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`

	row := config.DB.QueryRow(query, customerID)
	var customer models.Customer
	var user models.UserResponse
	var userID, userCreatedAt, userUpdatedAt sql.NullString
	var username, email, role sql.NullString

	err = row.Scan(
		&customer.ID, &customer.UserID, &customer.Name, &customer.Phone, &customer.Address, &customer.NIK, &customer.CreatedAt, &customer.UpdatedAt,
		&userID, &username, &email, &role, &userCreatedAt, &userUpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "Customer")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Add user info if exists
	if userID.Valid {
		var userIDInt int
		fmt.Sscanf(userID.String, "%d", &userIDInt)
		user.ID = userIDInt
		user.Username = username.String
		user.Email = email.String
		user.Role = role.String
		customer.User = &user
	}

	return utils.SuccessResponse(c, "Customer retrieved successfully", customer)
}

// CreateCustomer creates a new customer
func (cc *CustomerController) CreateCustomer(c *fiber.Ctx) error {
	var req models.CustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
	}

	// Sanitize input
	req.Name = utils.SanitizeString(req.Name)
	req.Phone = utils.SanitizeString(req.Phone)
	req.Address = utils.SanitizeString(req.Address)
	req.NIK = utils.SanitizeString(req.NIK)

	// Check if NIK already exists (if provided)
	if req.NIK != "" {
		var count int
		checkQuery := "SELECT COUNT(*) FROM customers WHERE nik = $1"
		err := config.DB.Get(&count, checkQuery, req.NIK)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
		}
		if count > 0 {
			return utils.ErrorResponse(c, fiber.StatusConflict, "NIK already exists")
		}
	}

	// Convert to customer model
	customer := req.ToCustomer()

	// Insert customer
	var customerID int
	query := `
		INSERT INTO customers (user_id, name, phone, address, nik) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`
	err := config.DB.Get(&customerID, query, customer.UserID, customer.Name, customer.Phone, customer.Address, customer.NIK)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create customer")
	}

	// Get the created customer
	getQuery := "SELECT * FROM customers WHERE id = $1"
	err = config.DB.Get(customer, getQuery, customerID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve created customer")
	}

	return utils.CreatedResponse(c, "Customer created successfully", customer)
}

// UpdateCustomer updates an existing customer
func (cc *CustomerController) UpdateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid customer ID")
	}

	var req models.CustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
	}

	// Check if customer exists
	var existingCustomer models.Customer
	checkQuery := "SELECT id, nik FROM customers WHERE id = $1"
	err = config.DB.Get(&existingCustomer, checkQuery, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "Customer")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Check if NIK already exists for another customer (if provided and different)
	if req.NIK != "" && req.NIK != existingCustomer.NIK {
		var count int
		checkNIKQuery := "SELECT COUNT(*) FROM customers WHERE nik = $1 AND id != $2"
		err := config.DB.Get(&count, checkNIKQuery, req.NIK, customerID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
		}
		if count > 0 {
			return utils.ErrorResponse(c, fiber.StatusConflict, "NIK already exists")
		}
	}

	// Update customer
	updateQuery := `
		UPDATE customers SET 
			user_id = $1, name = $2, phone = $3, address = $4, nik = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`
	_, err = config.DB.Exec(updateQuery, req.UserID, req.Name, req.Phone, req.Address, req.NIK, customerID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update customer")
	}

	// Get updated customer
	var updatedCustomer models.Customer
	getQuery := "SELECT * FROM customers WHERE id = $1"
	err = config.DB.Get(&updatedCustomer, getQuery, customerID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve updated customer")
	}

	return utils.SuccessResponse(c, "Customer updated successfully", updatedCustomer)
}

// DeleteCustomer deletes a customer
func (cc *CustomerController) DeleteCustomer(c *fiber.Ctx) error {
	// Check if user has admin role
	user := middleware.GetUserFromContext(c)
	if user.Role != "admin" {
		return utils.ForbiddenResponse(c, "Only admins can delete customers")
	}

	id := c.Params("id")
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid customer ID")
	}

	// Check if customer exists
	var customer models.Customer
	checkQuery := "SELECT id FROM customers WHERE id = $1"
	err = config.DB.Get(&customer, checkQuery, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "Customer")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Delete customer
	deleteQuery := "DELETE FROM customers WHERE id = $1"
	_, err = config.DB.Exec(deleteQuery, customerID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete customer")
	}

	return utils.SuccessResponse(c, "Customer deleted successfully", nil)
}