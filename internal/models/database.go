package models

import (
	"sync"
	"time"
)

// Database represents our in-memory NoSQL-like database
type Database struct {
	Users     map[string]*User     `json:"users"`
	Resources map[string]*Resource `json:"resources"`
	Logs      map[string]*LogEntry `json:"logs"`
	Setup     bool                 `json:"setup_completed"`
	mutex     sync.RWMutex         // For thread safety
}

// NewDatabase creates a new database instance
func NewDatabase() *Database {
	return &Database{
		Users:     make(map[string]*User),
		Resources: make(map[string]*Resource),
		Logs:      make(map[string]*LogEntry),
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
		"total_users":     len(db.Users),
		"total_resources": len(db.Resources),
		"total_logs":      len(db.Logs),
		"setup_complete":  db.Setup,
	}
}

// Log Management Methods

// CreateLog adds a new log entry to the database
func (db *Database) CreateLog(logEntry *LogEntry) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	logEntry.Timestamp = time.Now()
	db.Logs[logEntry.ID] = logEntry
	return nil
}

// GetLog retrieves a log entry by ID
func (db *Database) GetLog(id string) (*LogEntry, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	logEntry, exists := db.Logs[id]
	if !exists {
		return nil, ErrLogNotFound
	}
	return logEntry, nil
}

// ListLogs returns logs based on filter criteria
func (db *Database) ListLogs(filter *LogFilter) []*LogEntry {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	var logs []*LogEntry
	for _, logEntry := range db.Logs {
		// Apply filters
		if filter.UserID != "" && logEntry.UserID != filter.UserID {
			continue
		}
		if filter.Username != "" && logEntry.Username != filter.Username {
			continue
		}
		if filter.Level != "" && logEntry.Level != filter.Level {
			continue
		}
		if filter.Action != "" && logEntry.Action != filter.Action {
			continue
		}
		if filter.Resource != "" && logEntry.Resource != filter.Resource {
			continue
		}
		if filter.Method != "" && logEntry.Method != filter.Method {
			continue
		}
		if !filter.StartTime.IsZero() && logEntry.Timestamp.Before(filter.StartTime) {
			continue
		}
		if !filter.EndTime.IsZero() && logEntry.Timestamp.After(filter.EndTime) {
			continue
		}

		logs = append(logs, logEntry)
	}

	// Sort by timestamp (newest first)
	for i := 0; i < len(logs)-1; i++ {
		for j := i + 1; j < len(logs); j++ {
			if logs[i].Timestamp.Before(logs[j].Timestamp) {
				logs[i], logs[j] = logs[j], logs[i]
			}
		}
	}

	// Apply pagination
	if filter.Offset > 0 && filter.Offset < len(logs) {
		logs = logs[filter.Offset:]
	}
	if filter.Limit > 0 && filter.Limit < len(logs) {
		logs = logs[:filter.Limit]
	}

	return logs
}

// GetLogStats returns statistics about logs
func (db *Database) GetLogStats() *LogStats {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	stats := &LogStats{
		TopActions:   make(map[string]int64),
		TopResources: make(map[string]int64),
	}

	uniqueUsers := make(map[string]bool)
	var lastActivity time.Time

	for _, logEntry := range db.Logs {
		stats.TotalLogs++

		switch logEntry.Level {
		case LogLevelInfo:
			stats.InfoLogs++
		case LogLevelWarning:
			stats.WarningLogs++
		case LogLevelError:
			stats.ErrorLogs++
		}

		if logEntry.UserID != "" {
			uniqueUsers[logEntry.UserID] = true
		}

		stats.TopActions[logEntry.Action]++
		stats.TopResources[logEntry.Resource]++

		if logEntry.Timestamp.After(lastActivity) {
			lastActivity = logEntry.Timestamp
		}
	}

	stats.UniqueUsers = int64(len(uniqueUsers))
	stats.LastActivity = lastActivity

	return stats
}

// DeleteOldLogs removes logs older than the specified duration
func (db *Database) DeleteOldLogs(olderThan time.Duration) int {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	cutoff := time.Now().Add(-olderThan)
	deleted := 0

	for id, logEntry := range db.Logs {
		if logEntry.Timestamp.Before(cutoff) {
			delete(db.Logs, id)
			deleted++
		}
	}

	return deleted
}
