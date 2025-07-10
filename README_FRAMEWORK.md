# Go REST API Framework v2.0

## Overview

This Go REST API Framework is designed to help developers test dynamic backends by creating resources and JSON responses to be consumed by others. It effectively acts as a mock API for frontend development, providing a way to create and consume custom JSON resources without requiring a persistent database.

## Core Purpose

The primary purpose of this framework is to help developers test dynamic backends by providing:

- **Mock API Creation**: Create custom JSON endpoints for frontend testing
- **Dynamic Resource Management**: Add, modify, and delete API resources on-the-fly
- **User Management**: Role-based access control with Super Admin, Admin, and User roles
- **Volatile Data**: All data is stored in-memory and destroyed when the container restarts
- **Production Ready**: Can be used in production with database connectivity via environment variables

## Architecture

### Key Features

- **Stateless Design**: No persistent local state, enabling horizontal scaling
- **Role-Based Access Control**: Three distinct user roles with different permissions
- **JWT Authentication**: Secure token-based authentication with 24-hour expiry
- **Standardized JSON Responses**: Consistent response format across all endpoints
- **Input Validation**: Comprehensive request validation and sanitization
- **Cloud-Native**: Optimized for containerized deployment (Google Cloud Run)

### User Roles

1. **Super Admin**: 
   - GOD role, created during initial setup
   - One-time only setup (container must be destroyed for another Super Admin)
   - Full access to all resources and users

2. **Admin**: 
   - Can Add, Edit, Delete, List users and resources
   - Manage other users and their resources

3. **User**: 
   - Can only read users and resources
   - Can create and manage their own resources
   - Can invite others to consume their resource data

## API Endpoints

### Public Endpoints

#### `/setup` - Initial Super Admin Setup
- **Method**: POST only
- **Purpose**: One-time setup of the Super Admin account
- **Payload**: `username`, `email` (valid), `password` (strong)

**Example Request:**
```bash
curl -X POST http://localhost:8080/setup \
  -H "Content-Type: application/json" \
  -d '{"username": "johndoe", "email": "john@doe.com", "password": "str0ngP4ssw0rd"}'
```

**Example Response:**
```json
{
  "http_status_code": "202",
  "http_status_message": "Accepted",
  "resource": "/setup",
  "app": "Go REST API Framework",
  "timestamp": "2025-07-10T03:37:42Z",
  "response": {
    "message": "Now please login in order to get you the authorization token",
    "login_endpoint": "/login"
  }
}
```

#### `/login` - User Authentication
- **Method**: POST only
- **Purpose**: Authenticate user and receive JWT token
- **Payload**: `username`, `password`

**Example Request:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "johndoe", "password": "str0ngP4ssw0rd"}'
```

**Example Response:**
```json
{
  "http_status_code": "201",
  "http_status_message": "Created",
  "resource": "/login",
  "app": "Go REST API Framework",
  "timestamp": "2025-07-10T03:37:42Z",
  "response": {
    "token": "JWT Token that lasts 24 hours"
  }
}
```

#### `/status` - Server Status Check
- **Method**: GET only
- **Purpose**: Check server health and status
- **Payload**: None

**Example Request:**
```bash
curl -X GET http://localhost:8080/status \
  -H "Content-Type: application/json"
```

**Example Response:**
```json
{
  "http_status_code": "200",
  "http_status_message": "OK",
  "resource": "/status",
  "app": "Go REST API Framework",
  "timestamp": "2025-07-10T03:37:42Z",
  "response": {
    "status": "healthy",
    "version": "2.0.0"
  }
}
```

#### `/v1/ping` - Sample Resource
- **Method**: GET only
- **Purpose**: Sample endpoint demonstrating resource structure
- **Payload**: None

**Example Request:**
```bash
curl -X GET http://localhost:8080/v1/ping \
  -H "Content-Type: application/json"
```

**Example Response:**
```json
{
  "http_status_code": "200",
  "http_status_message": "OK",
  "resource": "/v1/ping",
  "app": "Go REST API Framework",
  "timestamp": "2025-07-10T03:37:42Z",
  "response": {
    "message": "pong"
  }
}
```

## Response Format

All API responses follow a standardized JSON structure:

```json
{
  "http_status_code": "200",
  "http_status_message": "OK",
  "resource": "/endpoint",
  "app": "Go REST API Framework",
  "timestamp": "2025-07-10T03:37:42Z",
  "response": {
    // Any kind of data here, up to the user
  }
}
```

### Error Response Format

```json
{
  "http_status_code": "400",
  "http_status_message": "Bad Request",
  "resource": "/endpoint",
  "app": "Go REST API Framework",
  "timestamp": "2025-07-10T03:37:42Z",
  "response": {
    "error": {
      "message": "Error description"
    }
  }
}
```

## Authentication & Authorization

### JWT Token Structure
- **Expiry**: 24 hours
- **Claims**: User ID, Username, Roles
- **Algorithm**: HS256

### Authorization Header
```
Authorization: Bearer <JWT_TOKEN>
```

## Getting Started

### Prerequisites
- Go 1.23+
- Docker (optional)

### Local Development

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-rest-api
```

2. **Install dependencies**
```bash
go mod download
```

3. **Run the server**
```bash
go run cmd/server/main.go
```

4. **Test the setup**
```bash
# Run the demo script
./examples/framework_demo.sh
```

### Docker Deployment

1. **Build the image**
```bash
docker build -t go-rest-api .
```

2. **Run the container**
```bash
docker run -p 8080:8080 go-rest-api
```

### Cloud Run Deployment

```bash
# Deploy to Google Cloud Run
./deploy.sh
```

## Environment Variables

### Required for Production
- `PORT`: Server port (default: 8080)
- `JWT_SECRET`: Secret key for JWT signing
- `DB_*`: Database connection parameters (when using persistent storage)

### Optional
- `LOG_LEVEL`: Logging level (default: info)
- `BCRYPT_COST`: Password hashing cost (default: 12)

## Use Cases

### Frontend Development
- Create mock API endpoints for testing
- Simulate different response scenarios
- Test authentication flows
- Prototype API integrations

### API Testing
- Create test data sets
- Simulate backend responses
- Test error handling
- Validate request/response formats

### Development Teams
- Share mock APIs across team members
- Create consistent test environments
- Prototype new features
- Test integration scenarios

## Framework Benefits

1. **Rapid Prototyping**: Quickly create mock APIs for testing
2. **No Database Required**: Volatile in-memory storage for development
3. **Production Ready**: Can connect to real databases via environment variables
4. **Standardized Responses**: Consistent JSON format across all endpoints
5. **Security**: JWT-based authentication with role-based access control
6. **Cloud Native**: Optimized for containerized deployment
7. **Developer Friendly**: Comprehensive logging and error handling

## Example Workflow

1. **Deploy the framework** to your preferred platform
2. **Setup Super Admin** using the `/setup` endpoint
3. **Login** to receive JWT token
4. **Create resources** with custom JSON responses
5. **Share endpoints** with frontend developers
6. **Test integrations** using the mock API
7. **Iterate quickly** by modifying resources as needed

## Security Considerations

- All passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- Input validation on all endpoints
- CORS support for cross-origin requests
- No sensitive data in logs
- Role-based access control

## Limitations

- **Volatile Storage**: Data is lost when container restarts
- **Single Instance**: No built-in clustering support
- **Memory Bound**: Limited by available system memory
- **No Persistence**: Requires external database for production persistence

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

---

**Go REST API Framework v2.0** - Empowering developers to test dynamic backends with ease.
