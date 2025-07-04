# Go REST API Framework v2.0 - Complete Tutorial

A powerful, flexible Go REST API framework that allows users to create and manage their own APIs with role-based access control and dynamic endpoint creation.

## 📚 Table of Contents

1. [Quick Start](#-quick-start)
2. [Initial Setup Tutorial](#-initial-setup-tutorial)
3. [Authentication Tutorial](#-authentication-tutorial)
4. [User Management Tutorial](#-user-management-tutorial)
5. [Resource Management Tutorial](#-resource-management-tutorial)
6. [Dynamic Endpoints Tutorial](#-dynamic-endpoints-tutorial)
7. [Admin Operations Tutorial](#-admin-operations-tutorial)
8. [Security & Best Practices](#-security--best-practices)
9. [Troubleshooting Guide](#-troubleshooting-guide)
10. [Advanced Usage](#-advanced-usage)

---

## 🚀 Quick Start

### Prerequisites
- Go 1.19 or higher
- curl (for testing)
- jq (optional, for JSON formatting)

### Installation

1. **Clone and Setup**
   ```bash
   git clone https://github.com/m4r4v/go-rest-api.git
   cd go-rest-api
   go mod download
   ```

2. **Start the Server**
   ```bash
   go run cmd/server/main.go cmd/server/router.go
   ```
   
   You should see:
   ```
   {"level":"info","msg":"Starting Go REST API Framework v2.0","time":"2025-07-03T22:00:00"}
   {"level":"info","msg":"Server starting on localhost:8080","time":"2025-07-03T22:00:00"}
   ```

---

## 🔧 Initial Setup Tutorial

### Step 1: First Time Setup

The framework requires an initial admin user to be created. This is a **one-time operation**.

```bash
curl -X POST http://localhost:8080/v1/setup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "password123"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "admin_id": "1",
    "message": "Setup completed successfully",
    "username": "admin"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 2: Verify Setup

Try to run setup again to confirm it's protected:

```bash
curl -X POST http://localhost:8080/v1/setup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "another_admin",
    "email": "another@example.com",
    "password": "password123"
  }' | jq
```

**Expected Response:**
```json
{
  "success": false,
  "status_code": 400,
  "status": "Bad Request",
  "error": {
    "code": "BAD_REQUEST",
    "message": "Setup already completed"
  },
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

✅ **Setup Complete!** Your admin user is ready.

---

## 🔐 Authentication Tutorial

### Step 1: Admin Login

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "admin",
    "role": "admin"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 2: Save Your Token

**Important:** Save the token for subsequent requests:

```bash
# Extract and save the token
TOKEN=$(curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}' | \
  jq -r '.data.token')

echo "Your token: $TOKEN"
```

### Step 3: Test Authentication

Get your user information:

```bash
curl -X GET http://localhost:8080/v1/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "1",
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:00:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 4: Test Unauthorized Access

Try accessing a protected endpoint without a token:

```bash
curl -X GET http://localhost:8080/v1/auth/me | jq
```

**Expected Response:**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authorization header is required",
    "status": 401
  }
}
```

✅ **Authentication Working!** You can now access protected endpoints.

---

## 👥 User Management Tutorial

### Step 1: Create a Regular User (Admin Only)

```bash
curl -X POST http://localhost:8080/v1/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "userpass123",
    "role": "user"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "uuid-here",
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:00:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 2: Create Another Admin User

```bash
curl -X POST http://localhost:8080/v1/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "jane_admin",
    "email": "jane@example.com",
    "password": "adminpass123",
    "role": "admin"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "uuid-jane-here",
    "username": "jane_admin",
    "email": "jane@example.com",
    "role": "admin",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:00:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 3: List All Users

```bash
curl -X GET http://localhost:8080/v1/admin/users \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "data": [
    {
      "id": "1",
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin",
      "created_at": "2025-07-03T22:00:00-04:00",
      "updated_at": "2025-07-03T22:00:00-04:00"
    },
    {
      "id": "uuid-here",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2025-07-03T22:00:00-04:00",
      "updated_at": "2025-07-03T22:00:00-04:00"
    }
  ],
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 4: Login as Regular User

```bash
# Login as the regular user
USER_TOKEN=$(curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "john_doe", "password": "userpass123"}' | \
  jq -r '.data.token')

echo "User token: $USER_TOKEN"
```

### Step 5: Test User Permissions

Try to access admin-only endpoint as regular user:

```bash
curl -X GET http://localhost:8080/v1/admin/users \
  -H "Authorization: Bearer $USER_TOKEN" | jq
```

**Expected Response:**
```json
{
  "success": false,
  "status_code": 403,
  "status": "Forbidden",
  "error": {
    "code": "FORBIDDEN",
    "message": "Insufficient permissions"
  },
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 6: Update User Profile (Self)

Regular users can update their own profile:

```bash
curl -X PUT http://localhost:8080/v1/users/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d '{
    "email": "john.doe.updated@example.com"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "uuid-here",
    "username": "john_doe",
    "email": "john.doe.updated@example.com",
    "role": "user",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:01:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:01:00-04:00"
}
```

### Step 7: Update User by Admin

Admins can update any user:

```bash
# First, get the user ID from the list
USER_ID="uuid-from-previous-list"

curl -X PUT http://localhost:8080/v1/admin/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "email": "john.admin.updated@example.com",
    "role": "admin"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "uuid-here",
    "username": "john_doe",
    "email": "john.admin.updated@example.com",
    "role": "admin",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:01:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:01:00-04:00"
}
```

✅ **User Management Complete!** You can create, list, and manage users.

---

## 📦 Resource Management Tutorial

### Step 1: Create Your First Resource

```bash
curl -X POST http://localhost:8080/v1/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "My First API",
    "description": "A simple greeting API",
    "data": {
      "endpoint": "/api/hello",
      "method": "GET",
      "response": {
        "message": "Hello, World!",
        "version": "1.0",
        "status": "active"
      }
    }
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "resource-uuid-here",
    "name": "My First API",
    "description": "A simple greeting API",
    "data": {
      "endpoint": "/api/hello",
      "method": "GET",
      "response": {
        "message": "Hello, World!",
        "version": "1.0",
        "status": "active"
      }
    },
    "created_by": "1",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:00:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 2: Create a More Complex Resource

```bash
curl -X POST http://localhost:8080/v1/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "User Profile API",
    "description": "Returns user profile information",
    "data": {
      "endpoint": "/api/profile",
      "method": "GET",
      "response": {
        "user": {
          "id": 123,
          "name": "John Doe",
          "email": "john@example.com",
          "preferences": {
            "theme": "dark",
            "notifications": true
          }
        },
        "last_login": "2025-07-03T22:00:00Z",
        "account_status": "active"
      }
    }
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "profile-uuid-here",
    "name": "User Profile API",
    "description": "Returns user profile information",
    "data": {
      "endpoint": "/api/profile",
      "method": "GET",
      "response": {
        "user": {
          "id": 123,
          "name": "John Doe",
          "email": "john@example.com",
          "preferences": {
            "theme": "dark",
            "notifications": true
          }
        },
        "last_login": "2025-07-03T22:00:00Z",
        "account_status": "active"
      }
    },
    "created_by": "1",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:00:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 3: List All Resources

```bash
curl -X GET http://localhost:8080/v1/resources \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "data": [
    {
      "id": "resource-uuid-here",
      "name": "My First API",
      "description": "A simple greeting API",
      "data": {
        "endpoint": "/api/hello",
        "method": "GET",
        "response": {
          "message": "Hello, World!",
          "version": "1.0",
          "status": "active"
        }
      },
      "created_by": "1",
      "created_at": "2025-07-03T22:00:00-04:00",
      "updated_at": "2025-07-03T22:00:00-04:00"
    },
    {
      "id": "profile-uuid-here",
      "name": "User Profile API",
      "description": "Returns user profile information",
      "data": {
        "endpoint": "/api/profile",
        "method": "GET",
        "response": {
          "user": {
            "id": 123,
            "name": "John Doe",
            "email": "john@example.com",
            "preferences": {
              "theme": "dark",
              "notifications": true
            }
          },
          "last_login": "2025-07-03T22:00:00Z",
          "account_status": "active"
        }
      },
      "created_by": "1",
      "created_at": "2025-07-03T22:00:00-04:00",
      "updated_at": "2025-07-03T22:00:00-04:00"
    }
  ],
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 4: Get Specific Resource

```bash
# Use the resource ID from the creation response
RESOURCE_ID="resource-uuid-from-step-1"

curl -X GET http://localhost:8080/v1/resources/$RESOURCE_ID \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "resource-uuid-here",
    "name": "My First API",
    "description": "A simple greeting API",
    "data": {
      "endpoint": "/api/hello",
      "method": "GET",
      "response": {
        "message": "Hello, World!",
        "version": "1.0",
        "status": "active"
      }
    },
    "created_by": "1",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:00:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 5: Update a Resource

```bash
curl -X PUT http://localhost:8080/v1/resources/$RESOURCE_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "My Updated API",
    "description": "An updated greeting API with more features",
    "data": {
      "endpoint": "/api/hello",
      "method": "GET",
      "response": {
        "message": "Hello, Updated World!",
        "version": "2.0",
        "status": "active",
        "features": ["greeting", "versioning", "status"]
      }
    }
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "resource-uuid-here",
    "name": "My Updated API",
    "description": "An updated greeting API with more features",
    "data": {
      "endpoint": "/api/hello",
      "method": "GET",
      "response": {
        "message": "Hello, Updated World!",
        "version": "2.0",
        "status": "active",
        "features": ["greeting", "versioning", "status"]
      }
    },
    "created_by": "1",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:01:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:01:00-04:00"
}
```

### Step 6: Create Resource as Regular User

```bash
curl -X POST http://localhost:8080/v1/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d '{
    "name": "User Created API",
    "description": "API created by regular user",
    "data": {
      "endpoint": "/api/user-data",
      "method": "GET",
      "response": {
        "message": "This was created by a regular user",
        "creator": "john_doe"
      }
    }
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "user-resource-uuid-here",
    "name": "User Created API",
    "description": "API created by regular user",
    "data": {
      "endpoint": "/api/user-data",
      "method": "GET",
      "response": {
        "message": "This was created by a regular user",
        "creator": "john_doe"
      }
    },
    "created_by": "uuid-here",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:00:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 7: Test Resource Ownership

Try to update another user's resource (should fail):

```bash
# Try to update admin's resource with user token
curl -X PUT http://localhost:8080/v1/resources/$RESOURCE_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d '{
    "name": "Unauthorized Update"
  }' | jq
```

**Expected Response:**
```json
{
  "success": false,
  "status_code": 403,
  "status": "Forbidden",
  "error": {
    "code": "FORBIDDEN",
    "message": "You can only update your own resources"
  },
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

✅ **Resource Management Complete!** You can create, list, update, and manage resources with proper permissions.

---

## 🔗 Dynamic Endpoints Tutorial

### Step 1: Understanding Dynamic Endpoints

When you create a resource with `endpoint`, `method`, and `response` data, the framework automatically creates a live API endpoint that requires authentication.

### Step 2: Test Your Dynamic Endpoint

From the previous tutorial, you created an endpoint at `/api/hello`. Let's test it:

**Without Authentication (Should Fail):**
```bash
curl -X GET http://localhost:8080/v1/api/hello | jq
```

**Expected Response:**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authorization header is required",
    "status": 401
  }
}
```

**With Authentication (Should Work):**
```bash
curl -X GET http://localhost:8080/v1/api/hello \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "success": true,
  "status_code": 200,
  "status": "OK",
  "response": {
    "message": "Hello, Updated World!",
    "version": "2.0",
    "status": "active",
    "features": ["greeting", "versioning", "status"]
  },
  "endpoint": "/v1/api/hello",
  "method": "GET",
  "user": "admin",
  "user_id": "1",
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 3: Create Different HTTP Methods

**Create a POST endpoint:**
```bash
curl -X POST http://localhost:8080/v1/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Data Submission API",
    "description": "Simulates data submission",
    "data": {
      "endpoint": "/api/submit",
      "method": "POST",
      "response": {
        "status": "success",
        "message": "Data submitted successfully",
        "submission_id": "12345",
        "timestamp": "2025-07-03T22:00:00Z"
      }
    }
  }' | jq
```

**Test the POST endpoint:**
```bash
curl -X POST http://localhost:8080/v1/api/submit \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Step 4: Create Complex API Responses

**Create an API that returns complex data:**
```bash
curl -X POST http://localhost:8080/v1/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "E-commerce Product API",
    "description": "Returns product catalog data",
    "data": {
      "endpoint": "/api/products",
      "method": "GET",
      "response": {
        "products": [
          {
            "id": 1,
            "name": "Laptop",
            "price": 999.99,
            "category": "Electronics",
            "in_stock": true,
            "ratings": {
              "average": 4.5,
              "count": 128
            }
          },
          {
            "id": 2,
            "name": "Coffee Mug",
            "price": 12.99,
            "category": "Home",
            "in_stock": true,
            "ratings": {
              "average": 4.8,
              "count": 45
            }
          }
        ],
        "total_count": 2,
        "page": 1,
        "per_page": 10
      }
    }
  }' | jq
```

**Test the complex endpoint:**
```bash
curl -X GET http://localhost:8080/v1/api/products \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Step 5: Update Dynamic Endpoints

When you update a resource, the dynamic endpoint is automatically updated:

```bash
curl -X PUT http://localhost:8080/v1/resources/$RESOURCE_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "data": {
      "endpoint": "/api/hello",
      "method": "GET",
      "response": {
        "message": "Hello, This is the updated version!",
        "version": "3.0",
        "status": "updated",
        "last_modified": "2025-07-03T22:00:00Z"
      }
    }
  }' | jq
```

**Test the updated endpoint:**
```bash
curl -X GET http://localhost:8080/v1/api/hello \
  -H "Authorization: Bearer $TOKEN" | jq
```

✅ **Dynamic Endpoints Complete!** You can create live, authenticated API endpoints instantly.

---

## 👑 Admin Operations Tutorial

### Step 1: System Status Monitoring

Check server status and database statistics:

```bash
curl -X GET http://localhost:8080/v1/status \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "data": {
    "status": "healthy",
    "timestamp": "2025-07-03T22:00:00-04:00",
    "version": "2.0.0",
    "database": {
      "users_count": 3,
      "resources_count": 5,
      "setup_completed": true
    }
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 2: Health Check

```bash
curl -X GET http://localhost:8080/v1/health \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "data": {
    "status": "healthy",
    "timestamp": "2025-07-03T22:00:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:00:00-04:00"
}
```

### Step 3: User Administration

**Promote a user to admin:**
```bash
curl -X PUT http://localhost:8080/v1/admin/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "role": "admin"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "uuid-here",
    "username": "john_doe",
    "email": "john@example.com",
    "role": "admin",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:02:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:02:00-04:00"
}
```

**Reset a user's password:**
```bash
curl -X PUT http://localhost:8080/v1/admin/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "password": "newpassword123"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "uuid-here",
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:02:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:02:00-04:00"
}
```

**Delete a user:**
```bash
curl -X DELETE http://localhost:8080/v1/admin/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "data": {
    "message": "User deleted successfully",
    "user_id": "uuid-here"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:02:00-04:00"
}
```

### Step 4: Resource Administration

**View all resources (admin can see all):**
```bash
curl -X GET http://localhost:8080/v1/resources \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "data": [
    {
      "id": "resource-1-uuid",
      "name": "My First API",
      "description": "A simple greeting API",
      "data": {
        "endpoint": "/api/hello",
        "method": "GET",
        "response": {
          "message": "Hello, World!",
          "version": "1.0",
          "status": "active"
        }
      },
      "created_by": "1",
      "created_at": "2025-07-03T22:00:00-04:00",
      "updated_at": "2025-07-03T22:00:00-04:00"
    },
    {
      "id": "resource-2-uuid",
      "name": "User Created API",
      "description": "API created by regular user",
      "data": {
        "endpoint": "/api/user-data",
        "method": "GET",
        "response": {
          "message": "This was created by a regular user",
          "creator": "john_doe"
        }
      },
      "created_by": "user-uuid",
      "created_at": "2025-07-03T22:00:00-04:00",
      "updated_at": "2025-07-03T22:00:00-04:00"
    }
  ],
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:02:00-04:00"
}
```

**Update any resource (admin override):**
```bash
curl -X PUT http://localhost:8080/v1/resources/$ANY_RESOURCE_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Admin Updated Resource",
    "description": "Updated by admin"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "resource-uuid-here",
    "name": "Admin Updated Resource",
    "description": "Updated by admin",
    "data": {
      "endpoint": "/api/user-data",
      "method": "GET",
      "response": {
        "message": "This was created by a regular user",
        "creator": "john_doe"
      }
    },
    "created_by": "user-uuid",
    "created_at": "2025-07-03T22:00:00-04:00",
    "updated_at": "2025-07-03T22:02:00-04:00"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:02:00-04:00"
}
```

**Delete any resource (admin override):**
```bash
curl -X DELETE http://localhost:8080/v1/resources/$ANY_RESOURCE_ID \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**
```json
{
  "data": {
    "message": "Resource deleted successfully",
    "resource_id": "resource-uuid-here"
  },
  "status": "OK",
  "status_code": 200,
  "success": true,
  "timestamp": "2025-07-03T22:02:00-04:00"
}
```

✅ **Admin Operations Complete!** You have full system control.

---

## 🔒 Security & Best Practices

### Step 1: Token Security

**Always use HTTPS in production:**
```bash
# In production, use HTTPS
curl -X POST https://your-domain.com/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

**Store tokens securely:**
```bash
# Don't expose tokens in logs or URLs
# Use environment variables or secure storage
export API_TOKEN="your-jwt-token-here"
curl -H "Authorization: Bearer $API_TOKEN" https://your-api.com/v1/status
```

### Step 2: Password Security

**Use strong passwords:**
```bash
curl -X POST http://localhost:8080/v1/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "secure_user",
    "email": "secure@example.com",
    "password": "MyStr0ng!P@ssw0rd#2025",
    "role": "user"
  }'
```

### Step 3: Input Validation

**The framework validates all inputs:**
```bash
# This will fail validation
curl -X POST http://localhost:8080/v1/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "ab",
    "email": "invalid-email",
    "password": "123",
    "role": "invalid_role"
  }' | jq
```

### Step 4: Rate Limiting

The framework includes built-in rate limiting. Excessive requests will be throttled.

### Step 5: CORS Configuration

The framework includes CORS headers for web applications:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

✅ **Security Best Practices Applied!**

---

## 🔧 Troubleshooting Guide

### Common Issues and Solutions

#### 1. "Connection Refused" Error
```bash
curl: (7) Failed to connect to localhost port 8080: Connection refused
```
**Solution:** Make sure the server is running:
```bash
go run cmd/server/main.go cmd/server/router.go
```

#### 2. "Invalid JSON" Error
```bash
{"error":{"code":"BAD_REQUEST","message":"Invalid JSON: EOF"}}
```
**Solution:** Make sure you're using the `-d` flag with curl:
```bash
# Wrong:
curl -X POST http://localhost:8080/v1/resources '{"name": "test"}'

# Correct:
curl -X POST http://localhost:8080/v1/resources -d '{"name": "test"}'
```

#### 3. "Authorization Header Required" Error
```bash
{"error":{"code":"UNAUTHORIZED","message":"Authorization header is required"}}
```
**Solution:** Include the Authorization header:
```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/v1/resources
```

#### 4. "Setup Already Completed" Error
**Solution:** This is normal. Setup can only be run once. Use login instead:
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

#### 5. "Forbidden" Error
```bash
{"error":{"code":"FORBIDDEN","message":"Insufficient permissions"}}
```
**Solution:** You're trying to access admin-only endpoints with a user token. Login as admin or use admin token.

#### 6. Token Expired
**Solution:** Login again to get a new token:
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}' | \
  jq -r '.data.token')
```

### Debug Mode

Enable debug logging by setting environment variable:
```bash
LOG_LEVEL=debug go run cmd/server/main.go cmd/server/router.go
```

✅ **Troubleshooting Complete!**

---

## 🚀 Advanced Usage

### Step 1: Environment Configuration

Create a `.env` file:
```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_EXPIRATION=24h

# Logger Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

### Step 2: Production Deployment

**Build for production:**
```bash
go build -o api-server cmd/server/main.go cmd/server/router.go
./api-server
```

**Docker deployment:**
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-server cmd/server/main.go cmd/server/router.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-server .
EXPOSE 8080
CMD ["./api-server"]
```

### Step 3: API Testing Script

Create a comprehensive test script:
```bash
#!/bin/bash
# test_api.sh

BASE_URL="http://localhost:8080"

echo "🚀 Testing Go REST API Framework v2.0"

# 1. Setup
echo "📝 Setting up admin user..."
curl -s -X POST $BASE_URL/v1/setup \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "email": "admin@test.com", "password": "test123"}' | jq

# 2. Login
echo "🔐 Logging in..."
TOKEN=$(curl -s -X POST $BASE_URL/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "test123"}' | jq -r '.data.token')

# 3. Create user
echo "👤 Creating user..."
curl -s -X POST $BASE_URL/v1/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"username": "testuser", "email": "test@test.com", "password": "test123", "role": "user"}' | jq

# 4. Create resource
echo "📦 Creating resource..."
curl -s -X POST $BASE_URL/v1/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Test API", "description": "Test", "data": {"endpoint": "/api/test", "method": "GET", "response": {"message": "Test successful!"}}}' | jq

# 5. Test dynamic endpoint
echo "🔗 Testing dynamic endpoint..."
curl -s -X GET $BASE_URL/v1/api/test \
  -H "Authorization: Bearer $TOKEN" | jq

echo "✅ All tests completed!"
```

### Step 4: Integration with Frontend

**JavaScript example:**
```javascript
class APIClient {
  constructor(baseURL) {
    this.baseURL = baseURL;
    this.token = localStorage.getItem('api_token');
  }

  async login(username, password) {
    const response = await fetch(`${this.baseURL}/v1/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    });
    const data = await response.json();
    if (data.success) {
      this.token = data.data.token;
      localStorage.setItem('api_token', this.token);
    }
    return data;
  }

  async createResource(name, description, endpoint, method, response) {
    return await fetch(`${this.baseURL}/v1/resources`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.token}`
      },
      body: JSON.stringify({
        name, description,
        data: { endpoint, method, response }
      })
    }).then(r => r.json());
  }

  async callDynamicEndpoint(endpoint) {
    return await fetch(`${this.baseURL}/v1${endpoint}`, {
      headers: { 'Authorization': `Bearer ${this.token}` }
    }).then(r => r.json());
  }
}

// Usage example
const api = new APIClient('http://localhost:8080');

// Login
await api.login('admin', 'password123');

// Create a resource
await api.createResource(
  'Weather API',
  'Returns weather data',
  '/api/weather',
  'GET',
  { temperature: 22, condition: 'sunny', humidity: 65 }
);

// Call the dynamic endpoint
const weather = await api.callDynamicEndpoint('/api/weather');
console.log(weather);
```

### Step 5: Monitoring and Logging

**Monitor server logs:**
```bash
# Run with structured logging
go run cmd/server/main.go cmd/server/router.go 2>&1 | jq

# Filter for specific log levels
go run cmd/server/main.go cmd/server/router.go 2>&1 | jq 'select(.level == "error")'

# Monitor HTTP requests
go run cmd/server/main.go cmd/server/router.go 2>&1 | jq 'select(.msg == "HTTP Request")'
```

### Step 6: Performance Testing

**Load testing with curl:**
```bash
#!/bin/bash
# load_test.sh

# Get token
TOKEN=$(curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}' | jq -r '.data.token')

# Test endpoint performance
for i in {1..100}; do
  curl -s -X GET http://localhost:8080/v1/api/hello \
    -H "Authorization: Bearer $TOKEN" \
    -w "Time: %{time_total}s\n" \
    -o /dev/null &
done
wait
```

### Step 7: Backup and Recovery

**Export all resources:**
```bash
# Export resources to JSON
curl -X GET http://localhost:8080/v1/resources \
  -H "Authorization: Bearer $TOKEN" | jq '.data' > resources_backup.json

# Export users (admin only)
curl -X GET http://localhost:8080/v1/admin/users \
  -H "Authorization: Bearer $TOKEN" | jq '.data' > users_backup.json
```

✅ **Advanced Usage Complete!** You're now ready for production deployment.

---

## 🎯 Use Cases & Examples

### E-commerce API
```bash
# Create product catalog endpoint
curl -X POST http://localhost:8080/v1/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Product Catalog",
    "description": "E-commerce product listing",
    "data": {
      "endpoint": "/api/products",
      "method": "GET",
      "response": {
        "products": [
          {"id": 1, "name": "Laptop", "price": 999.99, "stock": 50},
          {"id": 2, "name": "Mouse", "price": 29.99, "stock": 200}
        ],
        "total": 2,
        "page": 1
      }
    }
  }'
```

### User Authentication API
```bash
# Create user profile endpoint
curl -X POST http://localhost:8080/v1/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "User Profile",
    "description": "User profile information",
    "data": {
      "endpoint": "/api/user/profile",
      "method": "GET",
      "response": {
        "user": {
          "id": 123,
          "username": "john_doe",
          "email": "john@example.com",
          "profile": {
            "firstName": "John",
            "lastName": "Doe",
            "avatar": "https://example.com/avatar.jpg"
          }
        }
      }
    }
  }'
```

### IoT Data API
```bash
# Create sensor data endpoint
curl -X POST http://localhost:8080/v1/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Sensor Data",
    "description": "IoT sensor readings",
    "data": {
      "endpoint": "/api/sensors",
      "method": "GET",
      "response": {
        "sensors": [
          {"id": "temp_01", "type": "temperature", "value": 23.5, "unit": "°C"},
          {"id": "hum_01", "type": "humidity", "value": 65, "unit": "%"},
          {"id": "press_01", "type": "pressure", "value": 1013.25, "unit": "hPa"}
        ],
        "timestamp": "2025-07-03T22:00:00Z",
        "location": "Office Building A"
      }
    }
  }'
```

---

## 📊 API Reference Summary

### Authentication Endpoints
- `POST /v1/setup` - Initial admin setup (one-time)
- `POST /v1/auth/login` - User login
- `GET /v1/auth/me` - Get current user info

### User Management (Admin)
- `POST /v1/admin/users` - Create user
- `GET /v1/admin/users` - List all users
- `PUT /v1/admin/users/{id}` - Update user
- `DELETE /v1/admin/users/{id}` - Delete user

### User Self-Management
- `PUT /v1/users/me` - Update own profile

### Resource Management
- `POST /v1/resources` - Create resource
- `GET /v1/resources` - List resources
- `GET /v1/resources/{id}` - Get resource
- `PUT /v1/resources/{id}` - Update resource
- `DELETE /v1/resources/{id}` - Delete resource

### System Endpoints
- `GET /v1/status` - Server status
- `GET /v1/health` - Health check

### Dynamic Endpoints
- `{METHOD} /v1{endpoint}` - User-created endpoints

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/m4r4v/go-rest-api/issues)
- **Documentation**: This README
- **Examples**: See the tutorial sections above

---

## 🎉 Conclusion

Congratulations! You now have a complete understanding of the Go REST API Framework v2.0. This framework provides:

✅ **Secure Authentication** with JWT tokens  
✅ **Role-Based Access Control** (Admin/User)  
✅ **Dynamic Endpoint Creation** with instant API deployment  
✅ **Resource Management** with ownership controls  
✅ **Production-Ready** security and logging  
✅ **Developer-Friendly** with comprehensive error handling  

Start building your APIs today! 🚀
