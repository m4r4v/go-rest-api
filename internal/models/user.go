package models

import (
	"errors"
	"time"
)

// Common errors
var (
	ErrUserExists       = errors.New("user already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrResourceNotFound = errors.New("resource not found")
	ErrLogNotFound      = errors.New("log not found")
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`          // Never include password in JSON responses
	Role      string    `json:"role"`       // "super_admin", "admin", or "user"
	CreatedBy string    `json:"created_by"` // ID of user who created this user (empty for super_admin)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserResponse represents a user response without sensitive data
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	Username   string                 `json:"username"`
	Action     string                 `json:"action"`   // "create", "update", "delete", "login"
	Resource   string                 `json:"resource"` // "user", "resource", "auth"
	ResourceID string                 `json:"resource_id,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	IPAddress  string                 `json:"ip_address,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
}

// Resource represents a custom resource in the system
type Resource struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
	CreatedBy   string                 `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedBy: u.CreatedBy,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// IsSuperAdmin checks if the user is a super admin
func (u *User) IsSuperAdmin() bool {
	return u.Role == "super_admin"
}

// IsAdmin checks if the user is an admin (including super admin)
func (u *User) IsAdmin() bool {
	return u.Role == "admin" || u.Role == "super_admin"
}

// CanManageUser checks if this user can manage another user
func (u *User) CanManageUser(targetUser *User) bool {
	if u.IsSuperAdmin() {
		return true
	}
	if u.IsAdmin() && targetUser.CreatedBy == u.ID {
		return true
	}
	return u.ID == targetUser.ID // Users can manage themselves
}

// CanManageResource checks if this user can manage a resource
func (u *User) CanManageResource(resource *Resource) bool {
	if u.IsSuperAdmin() {
		return true
	}
	if u.IsAdmin() && resource.CreatedBy == u.ID {
		return true
	}
	return false
}

// CanViewResource checks if this user can view a resource
func (u *User) CanViewResource(resource *Resource) bool {
	if u.IsSuperAdmin() {
		return true
	}
	if u.IsAdmin() && resource.CreatedBy == u.ID {
		return true
	}
	// Users can view resources created by the same admin who created them
	if u.Role == "user" && resource.CreatedBy == u.CreatedBy {
		return true
	}
	return false
}
