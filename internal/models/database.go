package models

import (
	"sync"
	"time"
)

// Database represents our in-memory NoSQL-like database
type Database struct {
	Users     map[string]*User     `json:"users"`
	Resources map[string]*Resource `json:"resources"`
	AuditLogs map[string]*AuditLog `json:"audit_logs"`
	Setup     bool                 `json:"setup_completed"`
	mutex     sync.RWMutex         // For thread safety
}

// NewDatabase creates a new database instance
func NewDatabase() *Database {
	return &Database{
		Users:     make(map[string]*User),
		Resources: make(map[string]*Resource),
		AuditLogs: make(map[string]*AuditLog),
		Setup:     false,
	}
}

// User Management Methods

// CreateUser adds a new user to the database
func (db *Database) CreateUser(user *User) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check if username already exists
	if _, exists := db.Users[user.Username]; exists {
		return ErrUserExists
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	db.Users[user.Username] = user
	return nil
}

// GetUser retrieves a user by username
func (db *Database) GetUser(username string) (*User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	user, exists := db.Users[username]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func (db *Database) GetUserByID(id string) (*User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, user := range db.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// UpdateUser updates an existing user
func (db *Database) UpdateUser(username string, updates *User) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	user, exists := db.Users[username]
	if !exists {
		return ErrUserNotFound
	}

	// Update fields
	if updates.Email != "" {
		user.Email = updates.Email
	}
	if updates.Password != "" {
		user.Password = updates.Password
	}
	if updates.Role != "" {
		user.Role = updates.Role
	}

	user.UpdatedAt = time.Now()
	return nil
}

// DeleteUser removes a user from the database
func (db *Database) DeleteUser(username string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.Users[username]; !exists {
		return ErrUserNotFound
	}

	delete(db.Users, username)
	return nil
}

// ListUsers returns all users
func (db *Database) ListUsers() []*User {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	users := make([]*User, 0, len(db.Users))
	for _, user := range db.Users {
		users = append(users, user)
	}
	return users
}

// Resource Management Methods

// CreateResource adds a new resource to the database
func (db *Database) CreateResource(resource *Resource) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()
	db.Resources[resource.ID] = resource
	return nil
}

// GetResource retrieves a resource by ID
func (db *Database) GetResource(id string) (*Resource, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	resource, exists := db.Resources[id]
	if !exists {
		return nil, ErrResourceNotFound
	}
	return resource, nil
}

// UpdateResource updates an existing resource
func (db *Database) UpdateResource(id string, updates *Resource) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	resource, exists := db.Resources[id]
	if !exists {
		return ErrResourceNotFound
	}

	// Update fields
	if updates.Name != "" {
		resource.Name = updates.Name
	}
	if updates.Description != "" {
		resource.Description = updates.Description
	}
	if updates.Data != nil {
		resource.Data = updates.Data
	}

	resource.UpdatedAt = time.Now()
	return nil
}

// DeleteResource removes a resource from the database
func (db *Database) DeleteResource(id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.Resources[id]; !exists {
		return ErrResourceNotFound
	}

	delete(db.Resources, id)
	return nil
}

// ListResources returns all resources
func (db *Database) ListResources() []*Resource {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	resources := make([]*Resource, 0, len(db.Resources))
	for _, resource := range db.Resources {
		resources = append(resources, resource)
	}
	return resources
}

// Setup Methods

// IsSetupComplete checks if initial setup is done
func (db *Database) IsSetupComplete() bool {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return db.Setup
}

// CompleteSetup marks the setup as complete
func (db *Database) CompleteSetup() {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.Setup = true
}

// GetStats returns database statistics
func (db *Database) GetStats() map[string]interface{} {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	return map[string]interface{}{
		"total_users":      len(db.Users),
		"total_resources":  len(db.Resources),
		"total_audit_logs": len(db.AuditLogs),
		"setup_complete":   db.Setup,
	}
}

// Multi-tenant User Management Methods

// ListUsersByCreator returns users filtered by creator (for multi-tenancy)
func (db *Database) ListUsersByCreator(creatorID string, isSuperAdmin bool) []*User {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	users := make([]*User, 0)
	for _, user := range db.Users {
		if isSuperAdmin {
			// Super admin can see all users
			users = append(users, user)
		} else if user.CreatedBy == creatorID {
			// Admin can only see users they created
			users = append(users, user)
		}
	}
	return users
}

// GetUserByIDWithOwnership retrieves a user by ID with ownership check
func (db *Database) GetUserByIDWithOwnership(userID, requesterID string, isSuperAdmin bool) (*User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, user := range db.Users {
		if user.ID == userID {
			if isSuperAdmin {
				return user, nil
			}
			if user.CreatedBy == requesterID || user.ID == requesterID {
				return user, nil
			}
			return nil, ErrUserNotFound // Hide existence for security
		}
	}
	return nil, ErrUserNotFound
}

// DeleteUserWithCascade deletes a user and all their created users/resources
func (db *Database) DeleteUserWithCascade(username string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	user, exists := db.Users[username]
	if !exists {
		return ErrUserNotFound
	}

	// If this is an admin, delete all users they created
	if user.IsAdmin() {
		for uname, u := range db.Users {
			if u.CreatedBy == user.ID {
				delete(db.Users, uname)
			}
		}

		// Delete all resources they created
		for rid, r := range db.Resources {
			if r.CreatedBy == user.ID {
				delete(db.Resources, rid)
			}
		}
	}

	// Delete the user itself
	delete(db.Users, username)
	return nil
}

// Multi-tenant Resource Management Methods

// ListResourcesByCreator returns resources filtered by creator (for multi-tenancy)
func (db *Database) ListResourcesByCreator(creatorID string, isSuperAdmin bool) []*Resource {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	resources := make([]*Resource, 0)
	for _, resource := range db.Resources {
		if isSuperAdmin {
			// Super admin can see all resources
			resources = append(resources, resource)
		} else if resource.CreatedBy == creatorID {
			// Admin/User can only see resources created by their creator
			resources = append(resources, resource)
		}
	}
	return resources
}

// GetResourceWithOwnership retrieves a resource by ID with ownership check
func (db *Database) GetResourceWithOwnership(resourceID, requesterID string, isSuperAdmin bool, requesterCreatedBy string) (*Resource, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	resource, exists := db.Resources[resourceID]
	if !exists {
		return nil, ErrResourceNotFound
	}

	if isSuperAdmin {
		return resource, nil
	}
	if resource.CreatedBy == requesterID {
		return resource, nil
	}
	// Users can view resources created by the same admin who created them
	if resource.CreatedBy == requesterCreatedBy {
		return resource, nil
	}

	return nil, ErrResourceNotFound // Hide existence for security
}

// Audit Log Management Methods

// CreateAuditLog adds a new audit log entry
func (db *Database) CreateAuditLog(auditLog *AuditLog) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	auditLog.Timestamp = time.Now()
	db.AuditLogs[auditLog.ID] = auditLog
	return nil
}

// ListAuditLogs returns audit logs (super admin only)
func (db *Database) ListAuditLogs(limit int) []*AuditLog {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	logs := make([]*AuditLog, 0)
	for _, log := range db.AuditLogs {
		logs = append(logs, log)
	}

	// Sort by timestamp (newest first)
	for i := 0; i < len(logs)-1; i++ {
		for j := i + 1; j < len(logs); j++ {
			if logs[i].Timestamp.Before(logs[j].Timestamp) {
				logs[i], logs[j] = logs[j], logs[i]
			}
		}
	}

	// Apply limit
	if limit > 0 && limit < len(logs) {
		logs = logs[:limit]
	}

	return logs
}

// Dynamic Endpoint Validation

// GetExistingRoutes returns a list of existing API routes for conflict detection
func (db *Database) GetExistingRoutes() []string {
	return []string{
		"/setup",
		"/login",
		"/status",
		"/health",
		"/v1/ping",
		"/v1/users/me",
		"/v1/resources",
		"/v1/admin/users",
	}
}

// ValidateEndpointConflict checks if a dynamic endpoint conflicts with existing routes
func (db *Database) ValidateEndpointConflict(endpoint string) bool {
	existingRoutes := db.GetExistingRoutes()

	// Check against existing static routes
	for _, route := range existingRoutes {
		if endpoint == route || endpoint == "/v1"+route {
			return true // Conflict detected
		}
	}

	// Check against existing dynamic endpoints
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, resource := range db.Resources {
		if resource.Data != nil {
			if existingEndpoint, ok := resource.Data["endpoint"].(string); ok {
				if endpoint == existingEndpoint || endpoint == "/v1"+existingEndpoint {
					return true // Conflict detected
				}
			}
		}
	}

	return false // No conflict
}
