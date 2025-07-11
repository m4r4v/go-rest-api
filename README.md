# Go REST API Framework v2.0

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![API Version](https://img.shields.io/badge/API-v2.0-orange.svg)]()
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)]()
[![Cloud Run](https://img.shields.io/badge/Google%20Cloud%20Run-Ready-blue.svg)]()

A production-ready, enterprise-grade REST API framework built with Go, featuring advanced authentication, role-based access control, comprehensive security measures, and standardized response formats. This framework demonstrates modern software engineering practices and architectural patterns suitable for academic research and enterprise deployment.

## 🎯 Project Overview

This project represents a comprehensive implementation of a modern REST API framework, showcasing:

- **Advanced Authentication & Authorization**: JWT-based authentication with role-based access control (RBAC)
- **Enterprise Security**: Comprehensive security headers, CORS handling, and input validation
- **Standardized Architecture**: Clean code principles with separation of concerns
- **Production Readiness**: Docker containerization, Cloud Run deployment, and graceful shutdown
- **Academic Documentation**: Extensive documentation covering architecture, security, and development practices

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Quick Start](#-quick-start)
- [API Documentation](#-api-documentation)
- [Authentication](#-authentication)
- [Security](#-security)
- [Development](#-development)
- [Deployment](#-deployment)
- [Testing](#-testing)
- [Contributing](#-contributing)
- [Academic References](#-academic-references)

## ✨ Features

### Core Functionality
- **RESTful API Design**: Follows REST architectural principles and HTTP standards
- **JWT Authentication**: Secure token-based authentication with configurable expiration
- **Role-Based Access Control**: Multi-tier permission system (user, admin, super_admin)
- **Resource Management**: CRUD operations with ownership-based access control
- **User Management**: Complete user lifecycle management with admin controls

### Technical Excellence
- **Standardized Responses**: Consistent API response format across all endpoints
- **Comprehensive Logging**: Structured logging with configurable levels
- **Input Validation**: Robust request validation with detailed error messages
- **Error Handling**: Centralized error handling with appropriate HTTP status codes
- **Middleware Architecture**: Modular middleware for cross-cutting concerns

### Security Features
- **Password Hashing**: bcrypt-based password hashing with configurable cost
- **CORS Protection**: Configurable Cross-Origin Resource Sharing
- **Security Headers**: Comprehensive security headers (OWASP recommendations)
- **Request Validation**: Input sanitization and validation
- **Rate Limiting Ready**: Architecture supports rate limiting implementation

### Production Features
- **Docker Support**: Multi-stage Docker builds for optimized containers
- **Cloud Deployment**: Google Cloud Run ready with health checks
- **Graceful Shutdown**: Proper resource cleanup and connection handling
- **Environment Configuration**: Flexible configuration management
- **Health Monitoring**: Built-in health check endpoints

## 🏗️ Architecture

The framework follows a clean architecture pattern with clear separation of concerns:

```
go-rest-api/
├── cmd/server/           # Application entry point
├── internal/             # Private application code
│   ├── handlers/         # HTTP request handlers
│   └── models/           # Data models and business logic
├── pkg/                  # Public packages
│   ├── auth/            # Authentication service
│   ├── config/          # Configuration management
│   ├── errors/          # Error handling
│   ├── logger/          # Logging utilities
│   ├── middleware/      # HTTP middleware
│   └── validation/      # Input validation
├── docs/                # Documentation
├── examples/            # Usage examples
├── scripts/             # Utility scripts
└── tests/               # Test files
```

### Key Architectural Decisions

1. **Dependency Injection**: Services are injected into handlers for testability
2. **Interface Segregation**: Small, focused interfaces for better modularity
3. **Single Responsibility**: Each package has a single, well-defined purpose
4. **Configuration Management**: Environment-based configuration with sensible defaults
5. **Error Handling**: Centralized error types with consistent formatting

For detailed architectural documentation, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)
- Git

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/m4r4v/go-rest-api.git
   cd go-rest-api
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Build and run**:
   ```bash
   go build -o server cmd/server/main.go
   ./server
   ```

### Initial Setup

1. **Create super admin account**:
   ```bash
   curl -X POST http://localhost:8080/setup \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","email":"admin@example.com","password":"admin123"}'
   ```

2. **Login to get authentication token**:
   ```bash
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"admin123"}'
   ```

3. **Use the token for authenticated requests**:
   ```bash
   curl -X GET http://localhost:8080/v1/users/me \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

## 📚 API Documentation

### Response Format

All API responses follow a standardized format for consistency:

#### Success Response (New Standardized Format)
```json
{
  "success": true,
  "status_code": 200,
  "status": "OK",
  "timestamp": "2025-07-11T17:15:05-04:00",
  "endpoint": "/v1/users/me",
  "method": "GET",
  "user": "admin",
  "user_id": "1",
  "response": {
    "message": "Operation completed successfully",
    "data": { /* actual response data */ }
  }
}
```

#### Error Response
```json
{
  "success": false,
  "status_code": 400,
  "status": "Bad Request",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data"
  },
  "timestamp": "2025-07-11T17:15:05-04:00"
}
```

### Detailed Endpoint Documentation

| Endpoint | Method | Authentication | Role Required | Parameters | Request Body | Description |
|----------|--------|----------------|---------------|------------|--------------|-------------|
| `/setup` | POST | ❌ None | None | None | `{"username": "string", "email": "string", "password": "string"}` | Initial super admin setup (one-time only) |
| `/login` | POST | ❌ None | None | None | `{"username": "string", "password": "string"}` | User authentication - returns JWT token |
| `/status` | GET | ❌ None | None | None | None | Server status and health information |
| `/health` | GET | ❌ None | None | None | None | Health check endpoint for monitoring |
| `/v1/ping` | GET | ❌ None | None | None | None | Simple ping endpoint for connectivity test |
| `/v1/users/me` | GET | ✅ JWT Token | user | None | None | Get current user profile information |
| `/v1/users/me` | PUT | ✅ JWT Token | user | None | `{"email": "string", "password": "string"}` | Update current user profile |
| `/v1/resources` | GET | ✅ JWT Token | user | None | None | List all resources (user sees all, filtered by ownership) |
| `/v1/resources` | POST | ✅ JWT Token | user | None | `{"name": "string", "description": "string", "data": {}}` | Create new resource |
| `/v1/resources/{id}` | GET | ✅ JWT Token | user | `id` (UUID) | None | Get specific resource by ID |
| `/v1/resources/{id}` | PUT | ✅ JWT Token | user/admin | `id` (UUID) | `{"name": "string", "description": "string", "data": {}}` | Update resource (owner or admin only) |
| `/v1/resources/{id}` | DELETE | ✅ JWT Token | user/admin | `id` (UUID) | None | Delete resource (owner or admin only) |
| `/v1/admin/users` | GET | ✅ JWT Token | admin | None | None | List all users in the system |
| `/v1/admin/users` | POST | ✅ JWT Token | admin | None | `{"username": "string", "email": "string", "password": "string", "role": "user\|admin\|super_admin"}` | Create new user (super_admin can create admin users) |
| `/v1/admin/users/{id}` | GET | ✅ JWT Token | admin | `id` (UUID) | None | Get specific user by ID |
| `/v1/admin/users/{id}` | PUT | ✅ JWT Token | admin | `id` (UUID) | `{"email": "string", "password": "string", "role": "string"}` | Update user information |
| `/v1/admin/users/{id}` | DELETE | ✅ JWT Token | admin | `id` (UUID) | None | Delete user from system |

### Request/Response Examples

#### Setup Admin Account
**Request:**
```bash
POST /setup
Content-Type: application/json

{
  "username": "admin",
  "email": "admin@example.com", 
  "password": "admin123"
}
```

**Response:**
```json
{
  "http_status_code": "202",
  "http_status_message": "Accepted",
  "resource": "/setup",
  "app": "Go REST API Framework",
  "timestamp": "2025-07-11T17:11:26-04:00",
  "response": {
    "message": "Now please login in order to get you the authorization token",
    "login_endpoint": "/login"
  }
}
```

#### User Login
**Request:**
```bash
POST /login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "http_status_code": "201",
  "http_status_message": "Created",
  "resource": "/login",
  "app": "Go REST API Framework",
  "timestamp": "2025-07-11T17:11:47-04:00",
  "response": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### Update User Profile (Standardized Format)
**Request:**
```bash
PUT /v1/users/me
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "email": "newemail@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "status_code": 201,
  "status": "Created",
  "timestamp": "2025-07-11T17:15:05-04:00",
  "endpoint": "/v1/users/me",
  "method": "PUT",
  "user": "admin",
  "user_id": "1",
  "response": {
    "message": "Profile updated successfully",
    "user": {
      "id": "1",
      "username": "admin",
      "email": "newemail@example.com",
      "role": "super_admin",
      "created_by": "",
      "created_at": "2025-07-11T17:11:26.663103137-04:00",
      "updated_at": "2025-07-11T17:15:05.550642489-04:00"
    }
  }
}
```

#### Create Resource (Standardized Format)
**Request:**
```bash
POST /v1/resources
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "name": "Test Resource",
  "description": "A test resource for demonstration",
  "data": {
    "type": "demo",
    "value": 123
  }
}
```

**Response:**
```json
{
  "success": true,
  "status_code": 201,
  "status": "Created",
  "timestamp": "2025-07-11T17:15:31-04:00",
  "endpoint": "/v1/resources",
  "method": "POST",
  "user": "admin",
  "user_id": "1",
  "response": {
    "message": "Resource created successfully",
    "resource": {
      "id": "55ce3ded-6b7a-4ada-9d93-c08bf5a85b6e",
      "name": "Test Resource",
      "description": "A test resource for demonstration",
      "data": {
        "type": "demo",
        "value": 123
      },
      "created_by": "1",
      "created_at": "2025-07-11T17:15:31.258439775-04:00",
      "updated_at": "2025-07-11T17:15:31.258439984-04:00"
    }
  }
}
```

#### Error Response Example
**Request:**
```bash
POST /v1/resources
Authorization: Bearer <INVALID_TOKEN>
```

**Response:**
```json
{
  "success": false,
  "status_code": 401,
  "status": "Unauthorized",
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  },
  "timestamp": "2025-07-11T17:15:05-04:00"
}
```

For complete API documentation with working examples, see the [examples/](examples/) directory.

## 🔐 Authentication

### JWT Token Structure

The framework uses JSON Web Tokens (JWT) for authentication with the following claims:

```json
{
  "user_id": "unique-user-identifier",
  "username": "user-login-name",
  "roles": ["user", "admin", "super_admin"],
  "iss": "go-rest-api",
  "sub": "user-id",
  "exp": 1752354707,
  "nbf": 1752268307,
  "iat": 1752268307
}
```

### Role-Based Access Control

The framework implements a hierarchical role system:

1. **user**: Basic authenticated user
   - Access to own profile
   - Resource creation and management (own resources)

2. **admin**: Administrative user
   - All user permissions
   - User management capabilities
   - Access to all resources

3. **super_admin**: Super administrator
   - All admin permissions
   - Can create other admin users
   - System-wide access

### Security Implementation

- **Password Hashing**: bcrypt with configurable cost factor
- **Token Expiration**: Configurable JWT expiration times
- **Secure Headers**: OWASP-recommended security headers
- **Input Validation**: Comprehensive request validation
- **CORS Protection**: Configurable cross-origin policies

For detailed security documentation, see [docs/SECURITY.md](docs/SECURITY.md).

## 🛡️ Security

### Security Headers

The framework automatically applies security headers to all responses:

```http
X-Content-Type-Options: nosniff
Cache-Control: no-cache, no-store, must-revalidate
Pragma: no-cache
Expires: 0
X-API-Framework: Go-REST-API-v2.0
```

### Input Validation

All endpoints implement comprehensive input validation:

- **JSON Schema Validation**: Structured validation rules
- **Data Type Checking**: Strict type enforcement
- **Length Constraints**: Minimum and maximum length validation
- **Format Validation**: Email, password strength, etc.
- **Sanitization**: Input cleaning and normalization

### Authentication Security

- **Secure Password Storage**: bcrypt hashing with salt
- **Token Security**: Signed JWT tokens with expiration
- **Role Validation**: Strict role-based access control
- **Session Management**: Stateless authentication design

## 🔧 Development

### Project Structure

```
internal/
├── handlers/
│   └── api_handlers.go      # HTTP request handlers
└── models/
    ├── database.go          # In-memory database implementation
    ├── user.go             # User model and operations
    └── log.go              # Logging model

pkg/
├── auth/
│   └── auth.go             # Authentication service
├── config/
│   └── config.go           # Configuration management
├── errors/
│   └── errors.go           # Custom error types
├── logger/
│   └── logger.go           # Logging utilities
├── middleware/
│   └── middleware.go       # HTTP middleware
└── validation/
    └── validation.go       # Input validation
```

### Development Workflow

1. **Setup Development Environment**:
   ```bash
   # Install dependencies
   go mod download
   
   # Run tests
   go test ./...
   
   # Run with hot reload (using air)
   air
   ```

2. **Code Quality**:
   ```bash
   # Format code
   go fmt ./...
   
   # Lint code
   golangci-lint run
   
   # Vet code
   go vet ./...
   ```

3. **Testing**:
   ```bash
   # Run all tests
   go test ./...
   
   # Run tests with coverage
   go test -cover ./...
   
   # Run specific test
   go test ./tests -run TestAuth
   ```

For detailed development guidelines, see [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md).

## 🚢 Deployment

### Docker Deployment

1. **Build Docker image**:
   ```bash
   docker build -t go-rest-api .
   ```

2. **Run container**:
   ```bash
   docker run -p 8080:8080 \
     -e JWT_SECRET=your-secret-key \
     -e PORT=8080 \
     go-rest-api
   ```

### Google Cloud Run

1. **Deploy to Cloud Run**:
   ```bash
   # Build and deploy
   gcloud run deploy go-rest-api \
     --source . \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated
   ```

2. **Configure environment variables**:
   ```bash
   gcloud run services update go-rest-api \
     --set-env-vars JWT_SECRET=your-secret-key
   ```

### Environment Configuration

Key environment variables:

```bash
# Server Configuration
PORT=8080
HOST=0.0.0.0
LOG_LEVEL=info

# Authentication
JWT_SECRET=your-super-secret-key
JWT_EXPIRATION=24h
BCRYPT_COST=12

# CORS
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=*
```

For complete deployment documentation, see [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md).

## 🧪 Testing

### Test Coverage

The framework includes comprehensive tests covering:

- **Authentication Tests**: JWT generation, validation, and role checking
- **Handler Tests**: HTTP endpoint testing with various scenarios
- **Model Tests**: Data model validation and operations
- **Middleware Tests**: Authentication and authorization middleware
- **Integration Tests**: End-to-end API testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file
go test ./tests/auth_test.go
```

### Example Test Scenarios

```bash
# Test authentication flow
./examples/api_demo.sh

# Test framework features
./examples/framework_demo.sh
```

## 📖 Academic References

This project demonstrates several important software engineering concepts and patterns:

### Design Patterns
- **Repository Pattern**: Data access abstraction
- **Dependency Injection**: Loose coupling and testability
- **Middleware Pattern**: Cross-cutting concerns
- **Factory Pattern**: Service creation and configuration

### Software Engineering Principles
- **SOLID Principles**: Single responsibility, open/closed, etc.
- **Clean Architecture**: Separation of concerns and dependencies
- **RESTful Design**: Resource-oriented API design
- **Security by Design**: Built-in security considerations

### Technologies and Standards
- **JWT (RFC 7519)**: JSON Web Token standard
- **OAuth 2.0 Concepts**: Authorization framework principles
- **HTTP/1.1 (RFC 7231)**: HTTP protocol compliance
- **JSON API**: Consistent API response formatting

### Academic Applications
- **Software Architecture Research**: Clean architecture implementation
- **Security Research**: Authentication and authorization patterns
- **API Design Studies**: RESTful service design principles
- **DevOps Practices**: Containerization and cloud deployment

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit changes**: `git commit -m 'Add amazing feature'`
4. **Push to branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request**

### Development Standards
- Follow Go conventions and best practices
- Write comprehensive tests for new features
- Update documentation for API changes
- Ensure all tests pass before submitting

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

### Why MIT License?

The MIT License was chosen for this project for several important reasons:

1. **Academic Freedom**: Allows unrestricted use in academic research, teaching, and educational projects
2. **Enterprise Adoption**: Enables companies to use, modify, and integrate the framework without licensing concerns
3. **Open Source Compatibility**: Compatible with most other open source licenses, promoting collaboration
4. **Minimal Restrictions**: Only requires attribution, allowing maximum flexibility for users
5. **Industry Standard**: Widely recognized and trusted in the software development community
6. **Research Friendly**: Supports academic publications, thesis work, and research projects without legal barriers

The MIT License aligns with our goal of creating an educational and research-friendly framework that can benefit both academic institutions and commercial organizations while maintaining the open source spirit of knowledge sharing.

## 🙏 Acknowledgments

- **Go Community**: For excellent tooling and libraries
- **Gorilla Toolkit**: For robust HTTP routing and middleware
- **JWT.io**: For JWT implementation guidance
- **OWASP**: For security best practices and guidelines

## 📞 Support

For questions, issues, or contributions:

- **GitHub Issues**: [Create an issue](https://github.com/m4r4v/go-rest-api/issues)
- **Documentation**: Check the [docs/](docs/) directory
- **Examples**: See [examples/](examples/) for usage patterns

---

Built with ❤️ using **Go** by **[m4r4v](https://github.com/m4r4v)**
