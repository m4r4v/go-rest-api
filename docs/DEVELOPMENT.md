# Development Guide

## Overview

This guide provides comprehensive information for developers who want to contribute to, extend, or customize the Go REST API Framework v2.0.

## Development Environment Setup

### Prerequisites

- Go 1.23 or higher
- Git
- Docker and Docker Compose (optional)

### Initial Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/m4r4v/go-rest-api.git
   cd go-rest-api
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   go mod verify
   ```

3. **Setup Development Environment**
   ```bash
   cp .env.example .env
   # Edit .env with development settings
   ```

4. **Verify Installation**
   ```bash
   go build ./cmd/server
   go test ./...
   ```

## Project Structure

### Directory Organization

```
go-rest-api/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go        # Main server application
├── internal/              # Private application code
│   ├── handlers/          # HTTP handlers
│   └── models/           # Data models and business logic
├── pkg/                   # Public library code
│   ├── auth/             # Authentication utilities
│   ├── config/           # Configuration management
│   ├── errors/           # Error handling
│   ├── logger/           # Logging utilities
│   ├── middleware/       # HTTP middleware
│   └── validation/       # Input validation
├── tests/                 # Test files
├── examples/             # Usage examples and demos
├── scripts/              # Build and deployment scripts
├── docs/                 # Documentation
├── .env.example          # Environment template
├── Dockerfile            # Container configuration
├── docker-compose.yml    # Local development setup
├── go.mod               # Go module definition
└── go.sum               # Go module checksums
```

### Package Guidelines

#### Internal Packages (`internal/`)

- **Purpose**: Private application code that cannot be imported by external packages
- **Usage**: Business logic, handlers, models specific to this application
- **Import Rule**: Only importable by code in the same module

#### Public Packages (`pkg/`)

- **Purpose**: Reusable library code that could be imported by other projects
- **Usage**: Utilities, helpers, and components with general applicability
- **Import Rule**: Can be imported by any external package

## Coding Standards

### Go Style Guide

Follow the official Go style guide and these additional conventions:

1. **Naming Conventions**
   ```go
   // Good: Clear, descriptive names
   func CreateUser(user *User) error
   var userRepository UserRepository
   
   // Bad: Unclear abbreviations
   func CrtUsr(u *Usr) error
   var usrRepo UsrRepo
   ```

2. **Error Handling**
   ```go
   // Good: Explicit error handling
   user, err := userService.GetUser(id)
   if err != nil {
       return nil, fmt.Errorf("failed to get user: %w", err)
   }
   
   // Bad: Ignoring errors
   user, _ := userService.GetUser(id)
   ```

3. **Interface Design**
   ```go
   // Good: Small, focused interfaces
   type UserReader interface {
       GetUser(id string) (*User, error)
   }
   
   type UserWriter interface {
       CreateUser(user *User) error
       UpdateUser(id string, user *User) error
   }
   
   // Bad: Large, monolithic interfaces
   type UserService interface {
       GetUser(id string) (*User, error)
       CreateUser(user *User) error
       UpdateUser(id string, user *User) error
       DeleteUser(id string) error
       ListUsers() ([]*User, error)
       // ... many more methods
   }
   ```

### Code Formatting

Use standard Go tools for formatting:

```bash
# Format code
go fmt ./...

# Organize imports
goimports -w .

# Run linter
golangci-lint run
```

### Documentation Standards

1. **Package Documentation**
   ```go
   // Package auth provides JWT-based authentication utilities
   // for the Go REST API Framework.
   //
   // This package includes token generation, validation, and
   // user claim management functionality.
   package auth
   ```

2. **Function Documentation**
   ```go
   // GenerateToken creates a new JWT token for the given user.
   // The token includes user ID, username, and roles as claims.
   // Returns the signed token string or an error if generation fails.
   func GenerateToken(userID, username string, roles []string) (string, error) {
       // Implementation...
   }
   ```

3. **Type Documentation**
   ```go
   // User represents a system user with authentication and authorization data.
   type User struct {
       ID       string    `json:"id"`       // Unique user identifier
       Username string    `json:"username"` // Login username
       Email    string    `json:"email"`    // User email address
       Roles    []string  `json:"roles"`    // Assigned user roles
   }
   ```

## Testing Strategy

### Test Organization

```
tests/
├── unit/              # Unit tests for individual components
├── integration/       # Integration tests for component interaction
├── e2e/              # End-to-end tests for complete workflows
└── fixtures/         # Test data and utilities
```

### Unit Testing

#### Test Structure

```go
func TestAuthService_GenerateToken(t *testing.T) {
    // Arrange
    authService := auth.NewAuthService("test-secret", time.Hour, 4)
    userID := "123"
    username := "testuser"
    roles := []string{"user"}
    
    // Act
    token, err := authService.GenerateToken(userID, username, roles)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    
    // Verify token can be validated
    claims, err := authService.ValidateToken(token)
    assert.NoError(t, err)
    assert.Equal(t, userID, claims.UserID)
    assert.Equal(t, username, claims.Username)
    assert.Equal(t, roles, claims.Roles)
}
```

#### Table-Driven Tests

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name     string
        email    string
        expected bool
    }{
        {"valid email", "user@example.com", true},
        {"invalid email", "invalid-email", false},
        {"empty email", "", false},
        {"email without domain", "user@", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := validateEmail(tt.email)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Integration Testing

```go
func TestUserHandlers_Integration(t *testing.T) {
    // Setup test server
    server := setupTestServer(t)
    defer server.Close()
    
    // Create test user
    user := &User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    // Test user creation
    resp, err := http.Post(server.URL+"/setup", "application/json", 
        strings.NewReader(toJSON(user)))
    assert.NoError(t, err)
    assert.Equal(t, http.StatusAccepted, resp.StatusCode)
    
    // Test login
    loginResp, err := http.Post(server.URL+"/login", "application/json",
        strings.NewReader(toJSON(map[string]string{
            "username": user.Username,
            "password": user.Password,
        })))
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, loginResp.StatusCode)
}
```

### Test Utilities

Create helper functions for common test operations:

```go
// setupTestServer creates a test HTTP server with all routes configured
func setupTestServer(t *testing.T) *httptest.Server {
    authService := auth.NewAuthService("test-secret", time.Hour, 4)
    handlers := NewAPIHandlers(authService)
    
    router := mux.NewRouter()
    setupRoutes(router, handlers)
    
    return httptest.NewServer(router)
}

// toJSON converts a struct to JSON string for test requests
func toJSON(v interface{}) string {
    data, _ := json.Marshal(v)
    return string(data)
}

// assertJSONResponse validates JSON response structure
func assertJSONResponse(t *testing.T, resp *http.Response, expectedStatus int) map[string]interface{} {
    assert.Equal(t, expectedStatus, resp.StatusCode)
    
    var result map[string]interface{}
    err := json.NewDecoder(resp.Body).Decode(&result)
    assert.NoError(t, err)
    
    return result
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestAuthService_GenerateToken ./tests/

# Run tests with race detection
go test -race ./...

# Benchmark tests
go test -bench=. ./...
```

## Adding New Features

### Adding New Endpoints

1. **Define Handler Function**
   ```go
   // internal/handlers/api_handlers.go
   func (h *APIHandlers) CreateResource(w http.ResponseWriter, r *http.Request) {
       // Parse request
       var req CreateResourceRequest
       if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
           writeErrorResponse(w, errors.BadRequest("Invalid request format"))
           return
       }
       
       // Validate input
       if err := validateCreateResourceRequest(&req); err != nil {
           writeErrorResponse(w, err)
           return
       }
       
       // Get user from context
       claims := auth.GetClaimsFromContext(r.Context())
       
       // Business logic
       resource, err := h.createResource(&req, claims.UserID)
       if err != nil {
           writeErrorResponse(w, err)
           return
       }
       
       // Success response
       writeSuccessResponse(w, http.StatusCreated, map[string]interface{}{
           "message": "Resource created successfully",
           "resource": resource,
       })
   }
   ```

2. **Register Route**
   ```go
   // cmd/server/main.go
   router.HandleFunc("/v1/resources", 
       middleware.AuthMiddleware(authService)(
           http.HandlerFunc(handlers.CreateResource))).Methods("POST")
   ```

3. **Add Tests**
   ```go
   // tests/resource_test.go
   func TestCreateResource(t *testing.T) {
       server := setupTestServer(t)
       defer server.Close()
       
       // Test implementation...
   }
   ```

### Adding New Middleware

1. **Create Middleware Function**
   ```go
   // pkg/middleware/middleware.go
   func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
       limiter := rate.NewLimiter(rate.Limit(requestsPerMinute), requestsPerMinute)
       
       return func(next http.Handler) http.Handler {
           return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
               if !limiter.Allow() {
                   writeErrorResponse(w, errors.TooManyRequests("Rate limit exceeded"))
                   return
               }
               next.ServeHTTP(w, r)
           })
       }
   }
   ```

2. **Apply Middleware**
   ```go
   // cmd/server/main.go
   router.Use(middleware.RateLimitMiddleware(100)) // 100 requests per minute
   ```

### Adding New Validation Rules

1. **Create Validation Function**
   ```go
   // pkg/validation/validation.go
   func ValidatePhoneNumber(phone string) error {
       phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
       if !phoneRegex.MatchString(phone) {
           return errors.ValidationError("Invalid phone number format")
       }
       return nil
   }
   ```

2. **Use in Handlers**
   ```go
   if err := validation.ValidatePhoneNumber(req.Phone); err != nil {
       writeErrorResponse(w, err)
       return
   }
   ```

## Database Integration

### Adding External Database Support

1. **Create Database Interface**
   ```go
   // internal/models/database.go
   type DatabaseInterface interface {
       // User operations
       CreateUser(user *User) error
       GetUser(username string) (*User, error)
       UpdateUser(username string, updates *User) error
       DeleteUser(username string) error
       
       // Resource operations
       CreateResource(resource *Resource) error
       GetResource(id string) (*Resource, error)
       UpdateResource(id string, updates *Resource) error
       DeleteResource(id string) error
   }
   ```

2. **Implement PostgreSQL Database**
   ```go
   // pkg/database/postgres.go
   type PostgresDB struct {
       db *sql.DB
   }
   
   func NewPostgresDB(connectionString string) (*PostgresDB, error) {
       db, err := sql.Open("postgres", connectionString)
       if err != nil {
           return nil, err
       }
       
       return &PostgresDB{db: db}, nil
   }
   
   func (p *PostgresDB) CreateUser(user *User) error {
       query := `INSERT INTO users (username, email, password_hash, roles) 
                 VALUES ($1, $2, $3, $4) RETURNING id`
       
       err := p.db.QueryRow(query, user.Username, user.Email, 
           user.PasswordHash, pq.Array(user.Roles)).Scan(&user.ID)
       
       return err
   }
   ```

3. **Update Configuration**
   ```go
   // pkg/config/config.go
   type DatabaseConfig struct {
       Driver   string `env:"DB_DRIVER" envDefault:"memory"`
       Host     string `env:"DB_HOST" envDefault:"localhost"`
       Port     int    `env:"DB_PORT" envDefault:"5432"`
       Username string `env:"DB_USERNAME"`
       Password string `env:"DB_PASSWORD"`
       Database string `env:"DB_DATABASE"`
       SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
   }
   ```

## Performance Optimization

### Profiling

1. **Add Profiling Endpoints**
   ```go
   // cmd/server/main.go
   import _ "net/http/pprof"
   
   // Add pprof routes in development
   if config.Environment == "development" {
       router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
   }
   ```

2. **Profile Application**
   ```bash
   # CPU profiling
   go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
   
   # Memory profiling
   go tool pprof http://localhost:8080/debug/pprof/heap
   
   # Goroutine profiling
   go tool pprof http://localhost:8080/debug/pprof/goroutine
   ```

### Benchmarking

```go
// tests/benchmark_test.go
func BenchmarkAuthService_GenerateToken(b *testing.B) {
    authService := auth.NewAuthService("test-secret", time.Hour, 4)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := authService.GenerateToken("123", "user", []string{"user"})
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkAuthService_ValidateToken(b *testing.B) {
    authService := auth.NewAuthService("test-secret", time.Hour, 4)
    token, _ := authService.GenerateToken("123", "user", []string{"user"})
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := authService.ValidateToken(token)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Debugging

### Logging for Development

```go
// pkg/logger/logger.go
func SetupDevelopmentLogger() {
    logrus.SetLevel(logrus.DebugLevel)
    logrus.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
        ForceColors:   true,
    })
}
```

### Debug Middleware

```go
// pkg/middleware/debug.go
func DebugMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Log request
        logger.Debugf("Request: %s %s", r.Method, r.URL.Path)
        
        // Process request
        next.ServeHTTP(w, r)
        
        // Log response time
        logger.Debugf("Response time: %v", time.Since(start))
    })
}
```


## Contributing Guidelines

### Pull Request Process

1. **Fork and Clone**
   ```bash
   git clone https://github.com/yourusername/go-rest-api.git
   cd go-rest-api
   git remote add upstream https://github.com/m4r4v/go-rest-api.git
   ```

2. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Changes**
   - Follow coding standards
   - Add tests for new functionality
   - Update documentation

4. **Test Changes**
   ```bash
   go test ./...
   go vet ./...
   golangci-lint run
   ```

5. **Commit and Push**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   git push origin feature/your-feature-name
   ```

6. **Create Pull Request**
   - Provide clear description
   - Reference related issues
   - Include test results

### Commit Message Format

Follow conventional commits:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions or modifications
- `chore`: Maintenance tasks

Examples:
```
feat(auth): add JWT token refresh functionality
fix(handlers): resolve user creation validation bug
docs(api): update endpoint documentation
test(auth): add comprehensive JWT validation tests
```

This development guide provides a comprehensive foundation for contributing to and extending the Go REST API Framework while maintaining code quality and consistency.
