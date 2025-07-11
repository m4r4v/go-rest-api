# Security Documentation

## Security Overview

The Go REST API Framework v2.0 implements comprehensive security measures designed for production environments. This document details the security architecture, implementation, and best practices.

## Authentication System

### JWT (JSON Web Tokens)

**Implementation Details**:
- **Algorithm**: HMAC-SHA256 (HS256) for token signing
- **Token Structure**: Header + Payload + Signature
- **Expiration**: Configurable (default: 24 hours)
- **Claims**: User ID, username, roles, standard JWT claims

**Token Lifecycle**:
1. User provides credentials via `/login` endpoint
2. Server validates credentials against stored hash
3. JWT token generated with user claims
4. Token returned to client for subsequent requests
5. Token validated on each protected endpoint access
6. Token expires after configured duration

**Security Features**:
- Cryptographically signed tokens prevent tampering
- Stateless design eliminates server-side session storage
- Configurable expiration reduces exposure window
- Role-based claims enable fine-grained authorization

### Password Security

**Hashing Algorithm**: bcrypt with configurable cost factor

**Implementation**:
```go
// Password hashing with configurable cost
func (a *AuthService) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), a.bcryptCost)
    return string(bytes), err
}

// Password verification
func (a *AuthService) CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

**Security Properties**:
- **Salt Generation**: Automatic unique salt per password
- **Adaptive Cost**: Configurable work factor (default: 12)
- **Time Resistance**: Computational cost increases over time
- **Rainbow Table Resistance**: Salted hashes prevent precomputed attacks

## Authorization System

### Role-Based Access Control (RBAC)

**Role Hierarchy**:
```
super_admin (highest privilege)
    ├── Can create admin and user accounts
    ├── Can manage all resources and users
    └── Complete system access

admin (administrative privilege)
    ├── Can create user accounts
    ├── Can manage all resources
    ├── Can manage users (except super_admin)
    └── Cannot create admin accounts

user (basic privilege)
    ├── Can manage own resources
    ├── Can update own profile
    └── Cannot access admin endpoints
```

**Permission Matrix**:

| Operation | user | admin | super_admin |
|-----------|------|-------|-------------|
| View own profile | ✓ | ✓ | ✓ |
| Update own profile | ✓ | ✓ | ✓ |
| Create resources | ✓ | ✓ | ✓ |
| View own resources | ✓ | ✓ | ✓ |
| Update own resources | ✓ | ✓ | ✓ |
| Delete own resources | ✓ | ✓ | ✓ |
| View all resources | ✗ | ✓ | ✓ |
| Manage all resources | ✗ | ✓ | ✓ |
| View all users | ✗ | ✓ | ✓ |
| Create user accounts | ✗ | ✓ | ✓ |
| Update user accounts | ✗ | ✓ | ✓ |
| Delete user accounts | ✗ | ✓ | ✓ |
| Create admin accounts | ✗ | ✗ | ✓ |
| Manage admin accounts | ✗ | ✗ | ✓ |

### Middleware Implementation

**Authentication Middleware**:
```go
func AuthMiddleware(authService *auth.AuthService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract and validate JWT token
            // Set user context for downstream handlers
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

**Authorization Middleware**:
```go
func RequireRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Check user roles against required roles
            // Allow or deny access based on permissions
            next.ServeHTTP(w, r)
        })
    }
}
```

## Input Validation and Sanitization

### Request Validation

**JSON Schema Validation**:
- All request payloads validated against predefined schemas
- Type checking for all input fields
- Required field validation
- Format validation (email, password strength, etc.)

**Email Validation**:
```go
func validateEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}
```

**Password Strength Requirements**:
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character

### SQL Injection Prevention

**Parameterized Queries**: When using external databases, all queries use parameterized statements:

```go
// Safe query example
stmt, err := db.Prepare("SELECT * FROM users WHERE username = ? AND active = ?")
rows, err := stmt.Query(username, true)
```

**Input Sanitization**: All user inputs are sanitized before database operations.

## Cross-Origin Resource Sharing (CORS)

### Production-Safe CORS Configuration

**Environment-Based Origins**:
```go
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        allowedOrigins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
        origin := r.Header.Get("Origin")
        
        if origin != "" && isAllowedOrigin(origin, allowedOrigins) {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }
        // Additional CORS headers...
    })
}
```

**Security Features**:
- **No Wildcard Origins**: Prevents unauthorized cross-origin access
- **Credential Support**: Secure handling of authenticated requests
- **Method Restrictions**: Only allowed HTTP methods permitted
- **Header Validation**: Restricted allowed headers

**Configuration Examples**:

Development:
```bash
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

Production:
```bash
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
```

## HTTP Security Headers

### Implemented Security Headers

**Content Security**:
```go
w.Header().Set("Content-Type", "application/json; charset=utf-8")
w.Header().Set("X-Content-Type-Options", "nosniff")
```

**Cache Control**:
```go
w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
w.Header().Set("Pragma", "no-cache")
w.Header().Set("Expires", "0")
```

**Framework Identification**:
```go
w.Header().Set("X-API-Framework", "Go-REST-API-v2.0")
```

### Additional Recommended Headers

For production deployment, consider adding:

```go
// HTTPS enforcement
w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

// Clickjacking protection
w.Header().Set("X-Frame-Options", "DENY")

// XSS protection
w.Header().Set("X-XSS-Protection", "1; mode=block")

// Content Security Policy
w.Header().Set("Content-Security-Policy", "default-src 'self'")
```

## Error Handling and Information Disclosure

### Secure Error Responses

**Standardized Error Format**:
```json
{
  "success": false,
  "status_code": 401,
  "status": "Unauthorized",
  "timestamp": "2025-01-11T16:30:00Z",
  "endpoint": "/v1/users/me",
  "method": "GET",
  "user": null,
  "user_id": null,
  "response": {
    "error": {
      "code": "UNAUTHORIZED",
      "message": "Authentication required"
    }
  }
}
```

**Information Disclosure Prevention**:
- Generic error messages for authentication failures
- No stack traces in production responses
- No database error details exposed to clients
- Consistent response timing to prevent enumeration attacks

### Panic Recovery

**Recovery Middleware**:
```go
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                logger.Errorf("Panic recovered: %v", err)
                writeErrorResponse(w, errors.InternalServerError("An internal server error occurred"))
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

## Container Security

### Docker Security Measures

**Non-Root User**:
```dockerfile
# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Switch to non-root user
USER appuser
```

**Minimal Base Image**:
```dockerfile
# Use Alpine Linux for reduced attack surface
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
```

**Security Scanning**:
- Regular vulnerability scanning of base images
- Dependency vulnerability monitoring
- Container image signing and verification

### Runtime Security

**Resource Limits**:
```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "250m"
  limits:
    memory: "128Mi"
    cpu: "500m"
```

**Security Context**:
```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1001
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
```

## Logging and Monitoring

### Security Event Logging

**Authentication Events**:
- Successful logins with user identification
- Failed login attempts with source IP
- Token validation failures
- Account lockout events (when implemented)

**Authorization Events**:
- Unauthorized access attempts
- Permission escalation attempts
- Admin operation logging
- Resource access patterns

**Security Monitoring**:
```go
logger.Infof("Authentication successful: user=%s, ip=%s", username, clientIP)
logger.Warnf("Authentication failed: username=%s, ip=%s", username, clientIP)
logger.Errorf("Unauthorized access attempt: endpoint=%s, user=%s", endpoint, username)
```

### Audit Trail

**Request Logging**:
- All API requests logged with timestamps
- User identification for authenticated requests
- Request/response correlation IDs
- Performance metrics for anomaly detection

**Data Access Logging**:
- Resource creation, modification, deletion
- User account changes
- Permission modifications
- Administrative actions

## Threat Model

### Identified Threats and Mitigations

**1. Authentication Bypass**
- **Threat**: Unauthorized access to protected resources
- **Mitigation**: Strong JWT validation, secure token storage
- **Detection**: Failed authentication monitoring

**2. Authorization Escalation**
- **Threat**: Users accessing resources beyond their permissions
- **Mitigation**: Role-based middleware, permission validation
- **Detection**: Unauthorized access attempt logging

**3. Token Theft**
- **Threat**: Stolen JWT tokens used for unauthorized access
- **Mitigation**: HTTPS enforcement, token expiration, secure storage
- **Detection**: Unusual access patterns, multiple concurrent sessions

**4. Injection Attacks**
- **Threat**: SQL injection, NoSQL injection, command injection
- **Mitigation**: Parameterized queries, input validation, sanitization
- **Detection**: Malformed request monitoring

**5. Cross-Site Scripting (XSS)**
- **Threat**: Malicious script execution in client browsers
- **Mitigation**: Content-Type headers, output encoding, CSP headers
- **Detection**: Suspicious payload analysis

**6. Cross-Site Request Forgery (CSRF)**
- **Threat**: Unauthorized actions performed on behalf of authenticated users
- **Mitigation**: Stateless JWT tokens, SameSite cookies
- **Detection**: Unusual request patterns

**7. Denial of Service (DoS)**
- **Threat**: Service unavailability through resource exhaustion
- **Mitigation**: Rate limiting, resource limits, graceful degradation
- **Detection**: Traffic pattern analysis, resource monitoring

## Security Configuration

### Environment Variables

**Required Security Configuration**:
```bash
# Strong JWT secret (minimum 32 characters)
JWT_SECRET=your-256-bit-secret-key-here

# Restricted CORS origins
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com

# Strong password hashing
BCRYPT_COST=14
```

**Optional Security Configuration**:
```bash
# Token expiration
JWT_EXPIRATION=24h

# Logging level
LOG_LEVEL=warn

# Server timeouts
SERVER_READ_TIMEOUT=60s
SERVER_WRITE_TIMEOUT=60s
SERVER_IDLE_TIMEOUT=60s
```

### Production Security Checklist

**Pre-Deployment**:
- [ ] Strong JWT secret configured (256+ bits)
- [ ] CORS origins restricted to production domains
- [ ] HTTPS enforced for all communications
- [ ] Security headers implemented
- [ ] Input validation comprehensive
- [ ] Error messages sanitized
- [ ] Logging configured for security events

**Runtime Security**:
- [ ] Container running as non-root user
- [ ] Resource limits configured
- [ ] Network policies implemented
- [ ] Regular security updates applied
- [ ] Vulnerability scanning automated
- [ ] Incident response plan documented

**Monitoring and Response**:
- [ ] Security event monitoring active
- [ ] Alerting configured for suspicious activity
- [ ] Log aggregation and analysis implemented
- [ ] Backup and recovery procedures tested
- [ ] Security incident response team identified

## Security Best Practices

### Development Practices

1. **Secure by Default**: All endpoints require authentication unless explicitly public
2. **Principle of Least Privilege**: Users granted minimum necessary permissions
3. **Defense in Depth**: Multiple security layers implemented
4. **Input Validation**: All inputs validated and sanitized
5. **Error Handling**: Secure error responses without information disclosure

### Operational Practices

1. **Regular Updates**: Keep dependencies and base images updated
2. **Security Scanning**: Automated vulnerability scanning in CI/CD
3. **Access Control**: Restrict administrative access to production systems
4. **Monitoring**: Continuous security monitoring and alerting
5. **Incident Response**: Documented procedures for security incidents

### Compliance Considerations

**Data Protection**:
- Personal data handling according to GDPR/CCPA requirements
- Data retention and deletion policies
- Encryption of sensitive data at rest and in transit

**Industry Standards**:
- OWASP Top 10 compliance
- NIST Cybersecurity Framework alignment
- ISO 27001 security controls implementation

This security documentation provides comprehensive coverage of the framework's security implementation and should be regularly reviewed and updated as threats evolve and new security measures are implemented.
