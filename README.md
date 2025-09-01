# Restaurant Management Backend API

A comprehensive REST API backend for restaurant management built with Go, Gin framework, and MongoDB. This system provides complete functionality for managing restaurant operations including user authentication, menu management, order processing, table management, and invoice generation.

## 🏗️ Architecture Overview

This application follows a clean architecture pattern with the following structure:

```
golang-restaurant-management-backend/
├── main.go                    # Application entry point
├── controllers/               # HTTP request handlers
├── models/                   # Database models and structs
├── routes/                   # Route definitions and grouping
├── middleware/               # HTTP middleware (authentication)
├── helpers/                  # Utility functions (JWT tokens)
├── database/                 # Database connection and setup
├── go.mod                    # Go module dependencies
└── go.sum                    # Dependency checksums
```

## 🚀 Features

### Core Functionality

- **User Management**: Registration, authentication, and profile management
- **Menu Management**: Create and organize restaurant menus by category
- **Food Item Management**: CRUD operations for food items with pricing and images
- **Table Management**: Restaurant table tracking and capacity management
- **Order Processing**: Complete order lifecycle management
- **Order Items**: Individual item management within orders
- **Invoice Generation**: Billing and payment tracking
- **JWT Authentication**: Secure token-based authentication system

### Technical Features

- **RESTful API Design**: Clean, predictable API endpoints
- **JWT Token Security**: Access and refresh token system
- **Input Validation**: Comprehensive data validation using struct tags
- **MongoDB Integration**: NoSQL database with efficient queries
- **Middleware Support**: Authentication and logging middleware
- **Error Handling**: Consistent error responses across all endpoints
- **Pagination Support**: Efficient data retrieval for large datasets

## 🛠️ Technology Stack

- **Backend Framework**: [Gin](https://gin-gonic.com/) - High-performance HTTP web framework
- **Database**: [MongoDB](https://www.mongodb.com/) - NoSQL document database
- **Authentication**: JWT (JSON Web Tokens) with refresh token support
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator) for input validation
- **Password Hashing**: bcrypt for secure password storage
- **MongoDB Driver**: Official Go MongoDB driver

## 📋 Prerequisites

Before running this application, ensure you have the following installed:

- **Go**: Version 1.16 or higher
- **MongoDB**: Version 4.0 or higher (running locally on default port 27017)
- **Git**: For cloning the repository

## 🔧 Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd golang-restaurant-management-backend-master
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Environment Setup

Create a `.env` file in the root directory (optional - the app uses defaults):

```env
PORT=8000
SECRET_KEY=your-jwt-secret-key
MONGODB_URI=mongodb://localhost:27017
```

### 4. Start MongoDB

Ensure MongoDB is running locally:

```bash
# On macOS with Homebrew
brew services start mongodb-community

# On Ubuntu/Debian
sudo systemctl start mongod

# On Windows
net start MongoDB
```

### 5. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:8000` (or the port specified in your environment)

## 📚 API Documentation

### Authentication Endpoints (Public)

#### User Registration

```http
POST /users/signup
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "securepassword",
  "phone": "+1234567890"
}
```

#### User Login

```http
POST /users/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword"
}
```

### Protected Endpoints (Require Authentication)

All endpoints below require a valid JWT token in the header:

```http
token: your-jwt-token-here
```

#### User Management

- `GET /users` - Get paginated list of users
- `GET /users/:user_id` - Get specific user details

#### Food Management

- `GET /foods` - Get all food items
- `GET /foods/:food_id` - Get specific food item
- `POST /foods` - Create new food item
- `PATCH /foods/:food_id` - Update food item

#### Menu Management

- `GET /menus` - Get all menus
- `GET /menus/:menu_id` - Get specific menu
- `POST /menus` - Create new menu
- `PATCH /menus/:menu_id` - Update menu

#### Table Management

- `GET /tables` - Get all tables
- `GET /tables/:table_id` - Get specific table
- `POST /tables` - Create new table
- `PATCH /tables/:table_id` - Update table

#### Order Management

- `GET /orders` - Get all orders
- `GET /orders/:order_id` - Get specific order
- `POST /orders` - Create new order
- `PATCH /orders/:order_id` - Update order

#### Order Items Management

- `GET /orderItems` - Get all order items
- `GET /orderItems/:order_item_id` - Get specific order item
- `POST /orderItems` - Create new order item
- `PATCH /orderItems/:order_item_id` - Update order item

#### Invoice Management

- `GET /invoices` - Get all invoices
- `GET /invoices/:invoice_id` - Get specific invoice
- `POST /invoices` - Create new invoice
- `PATCH /invoices/:invoice_id` - Update invoice

## 🗃️ Database Schema

The application uses MongoDB with the following collections:

### Users Collection

```json
{
  "_id": "ObjectId",
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "password": "string (hashed)",
  "avatar": "string (optional)",
  "phone": "string",
  "token": "string",
  "refresh_token": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "user_id": "string"
}
```

### Food Collection

```json
{
  "_id": "ObjectId",
  "name": "string",
  "price": "number",
  "food_image": "string",
  "menu_id": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "food_id": "string"
}
```

### Menu Collection

```json
{
  "_id": "ObjectId",
  "name": "string",
  "category": "string",
  "start_date": "timestamp (optional)",
  "end_date": "timestamp (optional)",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "menu_id": "string"
}
```

### Orders Collection

```json
{
  "_id": "ObjectId",
  "order_date": "timestamp",
  "table_id": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "order_id": "string"
}
```

### Tables Collection

```json
{
  "_id": "ObjectId",
  "number_of_guests": "number",
  "table_number": "number",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "table_id": "string"
}
```

## 🔐 Authentication Flow

1. **Registration**: User creates account with email/password
2. **Login**: User authenticates and receives access + refresh tokens
3. **API Access**: Include access token in `token` header for protected routes
4. **Token Refresh**: Use refresh token to get new access tokens when expired

### Token Lifespans

- **Access Token**: 24 hours
- **Refresh Token**: 7 days (168 hours)

## 🧪 Testing the API

You can test the API using tools like Postman, curl, or any HTTP client:

### Example: Create a user and make authenticated request

```bash
# 1. Register a user
curl -X POST http://localhost:8000/users/signup \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "email": "test@example.com",
    "password": "password123",
    "phone": "1234567890"
  }'

# 2. Login to get tokens
curl -X POST http://localhost:8000/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# 3. Use the token for authenticated requests
curl -X GET http://localhost:8000/foods \
  -H "token: your-jwt-token-here"
```

## 🏃‍♂️ Development Workflow

### Project Structure Explanation

- **main.go**: Application entry point, sets up routes and middleware
- **controllers/**: Contains business logic and HTTP handlers for each entity
- **models/**: Defines data structures and MongoDB document schemas
- **routes/**: Groups related routes and applies middleware
- **middleware/**: Authentication and other HTTP middleware
- **helpers/**: Utility functions, primarily JWT token management
- **database/**: MongoDB connection and collection management

## ⚙️ Configuration

The application can be configured via environment variables:

- `PORT`: Server port (default: 8000)
- `SECRET_KEY`: JWT signing key (recommended for production)
- `MONGODB_URI`: MongoDB connection string (default: localhost:27017)
