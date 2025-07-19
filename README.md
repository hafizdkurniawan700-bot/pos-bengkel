# POS Bengkel - Vehicle Sales Management System

A complete MVP Vehicle Sales Management System with Golang backend (Fiber + SQLX + SQLite) and Flutter frontend implementing Material Design 3.

## 🎯 Project Overview

This is a full-stack application that provides a comprehensive vehicle sales management solution with role-based access control, inventory management, customer management, and transaction processing.

## 🛠️ Technology Stack

### Backend (Golang)
- **Framework**: Fiber v2 (High-performance web framework)
- **Database**: SQLite with SQLX (PostgreSQL-compatible schema)
- **Authentication**: JWT with bcrypt password hashing
- **Validation**: go-playground/validator
- **Architecture**: Clean architecture with controllers, models, repositories

### Frontend (Flutter)
- **Framework**: Flutter with Material Design 3
- **State Management**: BLoC pattern with flutter_bloc
- **Architecture**: Clean architecture (data/domain/presentation layers)
- **HTTP Client**: Dio for API communication
- **Storage**: flutter_secure_storage for tokens
- **Dependency Injection**: GetIt

## 🚀 Features Implemented

### Backend Features ✅
- [x] **Authentication System**
  - JWT-based authentication with role-based access control
  - User roles: admin, sales, cashier, customer
  - Secure password hashing with bcrypt
  - Token refresh mechanism

- [x] **Vehicle Management**
  - Complete CRUD operations for vehicles
  - Vehicle search and filtering
  - Status management (available, sold, reserved, maintenance)
  - Detailed vehicle information (brand, model, year, price, specs)

- [x] **Customer Management**
  - Customer registration and profile management
  - Customer search and filtering
  - Link customers to user accounts

- [x] **Transaction Management**
  - Transaction processing with status tracking
  - Automatic vehicle status updates
  - Payment method tracking

- [x] **Database & Infrastructure**
  - SQLite database with PostgreSQL-compatible schema
  - Database migrations and seeding
  - Proper error handling and logging
  - CORS support for Flutter integration

### Frontend Features ✅
- [x] **Beautiful UI with Material Design 3**
  - Custom theme with primary colors (Deep Blue #1565C0, Orange #FF9800)
  - Responsive design for mobile and web
  - Dark/light theme support

- [x] **Authentication Flow**
  - Login screen with role detection
  - Secure token storage
  - Auto-navigation based on user role

- [x] **Role-Based Dashboards**
  - Admin dashboard with comprehensive overview
  - Sales dashboard for sales team
  - Customer dashboard for customers
  - Statistics and activity tracking

- [x] **Core Architecture**
  - Clean architecture implementation
  - BLoC pattern for state management
  - Dependency injection setup
  - Network layer with error handling

## 📁 Project Structure

```
pos-bengkel/
├── backend/                           # Golang backend
│   ├── main.go                       # Application entry point
│   ├── config/                       # Configuration
│   │   ├── database.go              # Database connection & migrations
│   │   └── env.go                   # Environment variables
│   ├── models/                       # Data models
│   │   ├── user.go                  # User model with roles
│   │   ├── vehicle.go               # Vehicle model
│   │   ├── customer.go              # Customer model
│   │   ├── transaction.go           # Transaction model
│   │   └── test_drive.go            # Test drive model
│   ├── controllers/                  # HTTP controllers
│   │   ├── auth_controller.go       # Authentication endpoints
│   │   ├── vehicle_controller.go    # Vehicle CRUD operations
│   │   ├── customer_controller.go   # Customer management
│   │   └── transaction_controller.go# Transaction handling
│   ├── middleware/                   # HTTP middleware
│   │   └── auth.go                  # JWT middleware
│   ├── routes/                       # Route definitions
│   │   └── routes.go                # API routes with middleware
│   ├── utils/                        # Utility functions
│   │   ├── response.go              # Standard API responses
│   │   └── validation.go            # Input validation
│   ├── go.mod                        # Go dependencies
│   └── .env.example                  # Environment variables template
│
└── flutter_app/                      # Flutter frontend
    ├── lib/
    │   ├── main.dart                 # App entry point
    │   ├── core/                     # Core functionality
    │   │   ├── network/              # API client & error handling
    │   │   ├── theme/                # Material Design 3 theme
    │   │   └── constants/            # App constants
    │   ├── data/                     # Data layer
    │   │   ├── models/               # Data models
    │   │   ├── repositories/         # Repository implementations
    │   │   └── datasources/          # Remote data sources
    │   └── presentation/             # Presentation layer
    │       ├── pages/                # Screen widgets
    │       │   ├── auth/             # Authentication screens
    │       │   ├── dashboard/        # Role-based dashboards
    │       │   └── vehicles/         # Vehicle management screens
    │       ├── widgets/              # Reusable UI components
    │       └── bloc/                 # BLoC state management
    └── pubspec.yaml                  # Flutter dependencies
```

## 🔧 API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `GET /api/auth/profile` - Get user profile
- `POST /api/auth/refresh` - Refresh JWT token

### Vehicles
- `GET /api/vehicles` - Get vehicles list (with pagination & filtering)
- `GET /api/vehicles/:id` - Get vehicle by ID
- `POST /api/vehicles` - Create new vehicle (admin/sales only)
- `PUT /api/vehicles/:id` - Update vehicle (admin/sales only)
- `DELETE /api/vehicles/:id` - Delete vehicle (admin only)
- `GET /api/vehicles/search?q=query` - Search vehicles

### Customers
- `GET /api/customers` - Get customers list (authenticated)
- `GET /api/customers/:id` - Get customer by ID
- `POST /api/customers` - Create customer (admin/sales only)
- `PUT /api/customers/:id` - Update customer (admin/sales only)
- `DELETE /api/customers/:id` - Delete customer (admin only)

### Transactions
- `GET /api/transactions` - Get transactions list
- `GET /api/transactions/:id` - Get transaction by ID
- `POST /api/transactions` - Create transaction
- `PUT /api/transactions/:id/status` - Update transaction status

## 🚀 Getting Started

### Backend Setup

1. **Navigate to backend directory**:
   ```bash
   cd backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

4. **The API will be available at**: `http://localhost:8080`

5. **Health check**: `GET http://localhost:8080/health`

6. **Default admin credentials**:
   - Username: `admin`
   - Password: `admin123`

### Frontend Setup

1. **Navigate to Flutter directory**:
   ```bash
   cd flutter_app
   ```

2. **Install dependencies**:
   ```bash
   flutter pub get
   ```

3. **Run code generation**:
   ```bash
   flutter packages pub run build_runner build
   ```

4. **Run the application**:
   ```bash
   flutter run
   ```

## 🔒 User Roles & Permissions

| Feature | Admin | Sales | Cashier | Customer |
|---------|-------|-------|---------|----------|
| View Vehicles | ✅ | ✅ | ✅ | ✅ |
| Create/Edit Vehicles | ✅ | ✅ | ❌ | ❌ |
| Delete Vehicles | ✅ | ❌ | ❌ | ❌ |
| View Customers | ✅ | ✅ | ✅ | ❌ |
| Manage Customers | ✅ | ✅ | ❌ | ❌ |
| Process Transactions | ✅ | ✅ | ✅ | ❌ |
| User Management | ✅ | ❌ | ❌ | ❌ |
| Reports & Analytics | ✅ | ✅ | ❌ | ❌ |

## 🎨 Design System

### Color Palette
- **Primary Blue**: #1565C0 (Deep Blue)
- **Secondary Orange**: #FF9800 (Orange)
- **Success**: #388E3C (Green)
- **Warning**: #F57C00 (Amber)
- **Error**: #D32F2F (Red)

### Typography
- Material Design 3 typography scale
- Custom text styles for different use cases
- Responsive font sizes

## ✅ Success Criteria Met

- [x] **Backend runs without errors** on `go run main.go`
- [x] **API endpoints working** (tested with curl)
- [x] **Database connected** with successful migrations
- [x] **Authentication system** with JWT and role-based access
- [x] **Beautiful UI** with Material Design 3 theme
- [x] **Clean architecture** in both backend and frontend
- [x] **No compilation errors** in Golang backend
- [x] **Professional code structure** following best practices

## 🔄 Future Enhancements

### Phase 4: Advanced Features (Planned)
- [ ] Real-time notifications
- [ ] Advanced search and filtering
- [ ] Image upload for vehicles
- [ ] Test drive scheduling
- [ ] Reporting and analytics
- [ ] Email notifications
- [ ] Multi-language support
- [ ] Advanced user permissions

### Phase 5: Production Ready (Planned)
- [ ] PostgreSQL production database
- [ ] Docker containerization
- [ ] CI/CD pipeline
- [ ] Unit and integration tests
- [ ] API documentation (Swagger)
- [ ] Performance optimization
- [ ] Security audit
- [ ] Deployment guide

## 👨‍💻 Developer

Developed with ❤️ for the Vehicle Sales Management System MVP.

---

**Note**: This is an MVP implementation focusing on core functionality. The Flutter app structure is ready for development, with the backend API fully functional and tested.
