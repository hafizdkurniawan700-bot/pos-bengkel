package controllers

import (
	"database/sql"
	"pos-bengkel-backend/config"
	"pos-bengkel-backend/middleware"
	"pos-bengkel-backend/models"
	"pos-bengkel-backend/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// Login authenticates a user and returns a JWT token
func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
	}

	// Sanitize input
	req.Username = utils.SanitizeString(req.Username)

	// Find user by username or email
	var user models.User
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at 
		FROM users 
		WHERE username = $1 OR email = $1
	`
	err := config.DB.Get(&user, query, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.UnauthorizedResponse(c, "Invalid credentials")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return utils.UnauthorizedResponse(c, "Invalid credentials")
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	// Return response
	response := models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}

	return utils.SuccessResponse(c, "Login successful", response)
}

// Register creates a new user account
func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ValidationErrorResponse(c, utils.FormatValidationError(err))
	}

	// Sanitize input
	req.Username = utils.SanitizeString(req.Username)
	req.Email = utils.SanitizeEmail(req.Email)

	// Check if username or email already exists
	var count int
	checkQuery := "SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2"
	err := config.DB.Get(&count, checkQuery, req.Username, req.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	if count > 0 {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Username or email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	// Insert new user
	var userID int
	insertQuery := `
		INSERT INTO users (username, email, password_hash, role) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`
	err = config.DB.Get(&userID, insertQuery, req.Username, req.Email, string(hashedPassword), req.Role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	// Get the created user
	var user models.User
	getUserQuery := `
		SELECT id, username, email, role, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`
	err = config.DB.Get(&user, getUserQuery, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve user")
	}

	return utils.CreatedResponse(c, "User registered successfully", user.ToResponse())
}

// GetProfile returns the current user's profile
func (ac *AuthController) GetProfile(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return utils.UnauthorizedResponse(c, "User not found in context")
	}

	// Get full user details from database
	var fullUser models.User
	query := `
		SELECT id, username, email, role, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`
	err := config.DB.Get(&fullUser, query, user.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundResponse(c, "User")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	return utils.SuccessResponse(c, "Profile retrieved successfully", fullUser.ToResponse())
}

// RefreshToken generates a new JWT token
func (ac *AuthController) RefreshToken(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return utils.UnauthorizedResponse(c, "User not found in context")
	}

	// Generate new token
	token, err := middleware.GenerateJWT(user.UserID, user.Username, user.Email, user.Role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	return utils.SuccessResponse(c, "Token refreshed successfully", fiber.Map{"token": token})
}