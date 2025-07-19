package controllers

import (
	"database/sql"
	"fmt"
	"pos-bengkel-backend/config"
	"pos-bengkel-backend/models"
	"pos-bengkel-backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TransactionController struct{}

func NewTransactionController() *TransactionController {
	return &TransactionController{}
}

// GetTransactions returns a list of transactions with optional filtering
func (tc *TransactionController) GetTransactions(c *fiber.Ctx) error {
	var params models.TransactionSearchParams
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
		FROM transactions t
		LEFT JOIN vehicles v ON t.vehicle_id = v.id
		LEFT JOIN customers c ON t.customer_id = c.id
		WHERE 1=1
	`
	var args []interface{}
	argCounter := 0

	// Add filters
	if params.CustomerID > 0 {
		argCounter++
		baseQuery += fmt.Sprintf(" AND t.customer_id = $%d", argCounter)
		args = append(args, params.CustomerID)
	}
	if params.VehicleID > 0 {
		argCounter++
		baseQuery += fmt.Sprintf(" AND t.vehicle_id = $%d", argCounter)
		args = append(args, params.VehicleID)
	}
	if params.Status != "" {
		argCounter++
		baseQuery += fmt.Sprintf(" AND t.status = $%d", argCounter)
		args = append(args, params.Status)
	}
	if params.DateFrom != "" {
		argCounter++
		baseQuery += fmt.Sprintf(" AND t.created_at >= $%d", argCounter)
		args = append(args, params.DateFrom)
	}
	if params.DateTo != "" {
		argCounter++
		baseQuery += fmt.Sprintf(" AND t.created_at <= $%d", argCounter)
		args = append(args, params.DateTo)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := config.DB.Get(&total, countQuery, args...)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Get transactions with related data
	selectQuery := `
		SELECT 
			t.id, t.vehicle_id, t.customer_id, t.amount, t.status, t.payment_method, t.notes, t.created_at, t.updated_at,
			v.id as "vehicle.id", v.brand as "vehicle.brand", v.model as "vehicle.model", v.year as "vehicle.year", v.price as "vehicle.price",
			c.id as "customer.id", c.name as "customer.name", c.phone as "customer.phone"
		` + baseQuery + " ORDER BY t.created_at DESC LIMIT $" + strconv.Itoa(argCounter+1) + " OFFSET $" + strconv.Itoa(argCounter+2)
	args = append(args, params.Limit, params.Offset)

	rows, err := config.DB.Query(selectQuery, args...)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		var vehicle models.Vehicle
		var customer models.Customer
		var vehicleID, vehicleBrand, vehicleModel sql.NullString
		var vehicleYear sql.NullInt64
		var vehiclePrice sql.NullFloat64
		var customerID, customerName, customerPhone sql.NullString

		err := rows.Scan(
			&transaction.ID, &transaction.VehicleID, &transaction.CustomerID, &transaction.Amount, 
			&transaction.Status, &transaction.PaymentMethod, &transaction.Notes, &transaction.CreatedAt, &transaction.UpdatedAt,
			&vehicleID, &vehicleBrand, &vehicleModel, &vehicleYear, &vehiclePrice,
			&customerID, &customerName, &customerPhone,
		)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error scanning transaction data")
		}

		// Add vehicle info if exists
		if vehicleID.Valid {
			var vehicleIDInt int
			fmt.Sscanf(vehicleID.String, "%d", &vehicleIDInt)
			vehicle.ID = vehicleIDInt
			vehicle.Brand = vehicleBrand.String
			vehicle.Model = vehicleModel.String
			vehicle.Year = int(vehicleYear.Int64)
			vehicle.Price = vehiclePrice.Float64
			transaction.Vehicle = &vehicle
		}

		// Add customer info if exists
		if customerID.Valid {
			var customerIDInt int
			fmt.Sscanf(customerID.String, "%d", &customerIDInt)
			customer.ID = customerIDInt
			customer.Name = customerName.String
			customer.Phone = customerPhone.String
			transaction.Customer = &customer
		}

		transactions = append(transactions, transaction)
	}

	page := (params.Offset / params.Limit) + 1
	return utils.PaginatedSuccessResponse(c, "Transactions retrieved successfully", transactions, page, params.Limit, total)
}

// GetTransaction returns a single transaction by ID
func (tc *TransactionController) GetTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	transactionID, err := strconv.Atoi(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid transaction ID")
	}

	query := `
		SELECT 
			t.id, t.vehicle_id, t.customer_id, t.amount, t.status, t.payment_method, t.notes, t.created_at, t.updated_at,
			v.id as "vehicle.id", v.brand as "vehicle.brand", v.model as "vehicle.model", v.year as "vehicle.year", 
			v.price as "vehicle.price", v.status as "vehicle.status",
			c.id as "customer.id", c.name as "customer.name", c.phone as "customer.phone", c.address as "customer.address"
		FROM transactions t
		LEFT JOIN vehicles v ON t.vehicle_id = v.id
		LEFT JOIN customers c ON t.customer_id = c.id
		WHERE t.id = $1
	`

	row := config.DB.QueryRow(query, transactionID)
	var transaction models.Transaction
	var vehicle models.Vehicle
	var customer models.Customer
	var vehicleID, vehicleBrand, vehicleModel, vehicleStatus sql.NullString
	var vehicleYear sql.NullInt64
	var vehiclePrice sql.NullFloat64
	var customerID, customerName, customerPhone, customerAddress sql.NullString

	err = row.Scan(
		&transaction.ID, &transaction.VehicleID, &transaction.CustomerID, &transaction.Amount, 
		&transaction.Status, &transaction.PaymentMethod, &transaction.Notes, &transaction.CreatedAt, &transaction.UpdatedAt,
		&vehicleID, &vehicleBrand, &vehicleModel, &vehicleYear, &vehiclePrice, &vehicleStatus,
		&customerID, &customerName, &customerPhone, &customerAddress,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "Transaction")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Add vehicle info if exists
	if vehicleID.Valid {
		var vehicleIDInt int
		fmt.Sscanf(vehicleID.String, "%d", &vehicleIDInt)
		vehicle.ID = vehicleIDInt
		vehicle.Brand = vehicleBrand.String
		vehicle.Model = vehicleModel.String
		vehicle.Year = int(vehicleYear.Int64)
		vehicle.Price = vehiclePrice.Float64
		vehicle.Status = vehicleStatus.String
		transaction.Vehicle = &vehicle
	}

	// Add customer info if exists
	if customerID.Valid {
		var customerIDInt int
		fmt.Sscanf(customerID.String, "%d", &customerIDInt)
		customer.ID = customerIDInt
		customer.Name = customerName.String
		customer.Phone = customerPhone.String
		customer.Address = customerAddress.String
		transaction.Customer = &customer
	}

	return utils.SuccessResponse(c, "Transaction retrieved successfully", transaction)
}

// CreateTransaction creates a new transaction
func (tc *TransactionController) CreateTransaction(c *fiber.Ctx) error {
	var req models.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
	}

	// Check if vehicle exists and is available
	var vehicle models.Vehicle
	vehicleQuery := "SELECT id, status, price FROM vehicles WHERE id = $1"
	err := config.DB.Get(&vehicle, vehicleQuery, req.VehicleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "Vehicle")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	if vehicle.Status != "available" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Vehicle is not available for sale")
	}

	// Check if customer exists
	var customer models.Customer
	customerQuery := "SELECT id FROM customers WHERE id = $1"
	err = config.DB.Get(&customer, customerQuery, req.CustomerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "Customer")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Convert to transaction model
	transaction := req.ToTransaction()

	// Start database transaction
	tx, err := config.DB.Beginx()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	// Insert transaction
	var transactionID int
	insertQuery := `
		INSERT INTO transactions (vehicle_id, customer_id, amount, status, payment_method, notes) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id
	`
	err = tx.Get(&transactionID, insertQuery, 
		transaction.VehicleID, transaction.CustomerID, transaction.Amount, 
		transaction.Status, transaction.PaymentMethod, transaction.Notes)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create transaction")
	}

	// If transaction is completed, update vehicle status
	if transaction.Status == "completed" {
		updateVehicleQuery := "UPDATE vehicles SET status = 'sold', updated_at = CURRENT_TIMESTAMP WHERE id = $1"
		_, err = tx.Exec(updateVehicleQuery, transaction.VehicleID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update vehicle status")
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	// Get the created transaction
	getQuery := "SELECT * FROM transactions WHERE id = $1"
	err = config.DB.Get(transaction, getQuery, transactionID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve created transaction")
	}

	return utils.CreatedResponse(c, "Transaction created successfully", transaction)
}

// UpdateTransactionStatus updates a transaction's status
func (tc *TransactionController) UpdateTransactionStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	transactionID, err := strconv.Atoi(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid transaction ID")
	}

	var req models.TransactionStatusUpdate
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
	}

	// Check if transaction exists
	var existingTransaction models.Transaction
	checkQuery := "SELECT id, vehicle_id, status FROM transactions WHERE id = $1"
	err = config.DB.Get(&existingTransaction, checkQuery, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "Transaction")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Start database transaction
	tx, err := config.DB.Beginx()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	// Update transaction status
	updateQuery := `
		UPDATE transactions SET 
			status = $1, notes = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	_, err = tx.Exec(updateQuery, req.Status, req.Notes, transactionID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update transaction")
	}

	// Update vehicle status based on transaction status
	var vehicleStatus string
	switch req.Status {
	case "completed":
		vehicleStatus = "sold"
	case "cancelled", "refunded":
		vehicleStatus = "available"
	}

	if vehicleStatus != "" {
		updateVehicleQuery := "UPDATE vehicles SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2"
		_, err = tx.Exec(updateVehicleQuery, vehicleStatus, existingTransaction.VehicleID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update vehicle status")
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	// Get updated transaction
	var updatedTransaction models.Transaction
	getQuery := "SELECT * FROM transactions WHERE id = $1"
	err = config.DB.Get(&updatedTransaction, getQuery, transactionID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve updated transaction")
	}

	return utils.SuccessResponse(c, "Transaction status updated successfully", updatedTransaction)
}