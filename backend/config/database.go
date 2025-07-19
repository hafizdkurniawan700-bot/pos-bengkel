package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	
	// Build connection string
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	
	DB, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Printf("Failed to connect to PostgreSQL: %v", err)
		log.Println("Using SQLite as fallback for development")
		
		// Fallback to SQLite for development
		DB, err = sqlx.Connect("sqlite3", "./pos_bengkel.db")
		if err != nil {
			log.Fatal("Failed to connect to SQLite database:", err)
		}
		IsPostgreSQL = false
		log.Println("✅ Connected to SQLite database")
	} else {
		IsPostgreSQL = true
		log.Println("✅ Connected to PostgreSQL database")
	}
	
	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	
	log.Println("✅ Database connection verified")
	
	// Run migrations
	runMigrations()
}

func runMigrations() {
	log.Println("🔄 Running database migrations...")
	
	// Create tables (SQLite compatible)
	createTables := `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'sales', 'cashier', 'customer')),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Customers table
	CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(100) NOT NULL,
		phone VARCHAR(20),
		address TEXT,
		nik VARCHAR(20) UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Vehicles table
	CREATE TABLE IF NOT EXISTS vehicles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		brand VARCHAR(50) NOT NULL,
		model VARCHAR(50) NOT NULL,
		year INTEGER NOT NULL,
		price DECIMAL(15,2) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'sold', 'reserved', 'maintenance')),
		images TEXT, -- JSON array of image URLs
		description TEXT,
		engine_type VARCHAR(50),
		transmission VARCHAR(20),
		fuel_type VARCHAR(20),
		mileage INTEGER DEFAULT 0,
		color VARCHAR(30),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Transactions table
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		vehicle_id INTEGER REFERENCES vehicles(id) ON DELETE CASCADE,
		customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
		amount DECIMAL(15,2) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'cancelled', 'refunded')),
		payment_method VARCHAR(30),
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Test drives table
	CREATE TABLE IF NOT EXISTS test_drives (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
		vehicle_id INTEGER REFERENCES vehicles(id) ON DELETE CASCADE,
		scheduled_date DATETIME NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'completed', 'cancelled', 'no_show')),
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_customers_user_id ON customers(user_id);
	CREATE INDEX IF NOT EXISTS idx_vehicles_status ON vehicles(status);
	CREATE INDEX IF NOT EXISTS idx_transactions_customer_id ON transactions(customer_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_vehicle_id ON transactions(vehicle_id);
	CREATE INDEX IF NOT EXISTS idx_test_drives_customer_id ON test_drives(customer_id);
	CREATE INDEX IF NOT EXISTS idx_test_drives_vehicle_id ON test_drives(vehicle_id);
	`
	
	_, err := DB.Exec(createTables)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	
	// Insert default admin user if not exists
	seedData()
	
	log.Println("✅ Database migrations completed")
}

func seedData() {
	// Check if admin user exists
	var count int
	err := DB.Get(&count, "SELECT COUNT(*) FROM users WHERE role = 'admin'")
	if err != nil {
		log.Printf("Error checking admin user: %v", err)
		return
	}
	
	if count == 0 {
		// Create default admin user (password: admin123)
		// In production, this should be changed immediately
		adminQuery := `
		INSERT INTO users (username, email, password_hash, role) 
		VALUES ($1, $2, $3, $4)
		`
		
		// Generate bcrypt hash for "admin123"
		passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error generating password hash: %v", err)
			return
		}
		
		_, err = DB.Exec(adminQuery, "admin", "admin@posbengkel.com", string(passwordHash), "admin")
		if err != nil {
			log.Printf("Error creating admin user: %v", err)
			return
		}
		
		log.Println("✅ Default admin user created (username: admin, password: admin123)")
	}
	
	// Add some sample vehicles if none exist
	var vehicleCount int
	err = DB.Get(&vehicleCount, "SELECT COUNT(*) FROM vehicles")
	if err != nil {
		log.Printf("Error checking vehicles: %v", err)
		return
	}
	
	if vehicleCount == 0 {
		sampleVehicles := []map[string]interface{}{
			{
				"brand": "Toyota", "model": "Camry", "year": 2023, "price": 35000.00,
				"engine_type": "2.5L 4-Cylinder", "transmission": "Automatic", 
				"fuel_type": "Gasoline", "color": "Silver", "mileage": 5000,
				"description": "Brand new Toyota Camry with excellent fuel efficiency",
			},
			{
				"brand": "Honda", "model": "Civic", "year": 2022, "price": 28000.00,
				"engine_type": "1.5L Turbo", "transmission": "CVT", 
				"fuel_type": "Gasoline", "color": "Blue", "mileage": 15000,
				"description": "Reliable Honda Civic with modern features",
			},
			{
				"brand": "BMW", "model": "X3", "year": 2023, "price": 55000.00,
				"engine_type": "2.0L Turbo", "transmission": "Automatic", 
				"fuel_type": "Gasoline", "color": "Black", "mileage": 2000,
				"description": "Luxury BMW X3 SUV with premium features",
			},
		}
		
		for _, vehicle := range sampleVehicles {
			vehicleQuery := `
			INSERT INTO vehicles (brand, model, year, price, engine_type, transmission, fuel_type, color, mileage, description) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			`
			_, err = DB.Exec(vehicleQuery, 
				vehicle["brand"], vehicle["model"], vehicle["year"], vehicle["price"],
				vehicle["engine_type"], vehicle["transmission"], vehicle["fuel_type"], 
				vehicle["color"], vehicle["mileage"], vehicle["description"])
			if err != nil {
				log.Printf("Error inserting sample vehicle: %v", err)
			}
		}
		
		log.Println("✅ Sample vehicles added")
	}
}