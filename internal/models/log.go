package models

import (
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
)

// LogEntry represents a single log entry for user interactions
type LogEntry struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id,omitempty"`
	Username   string                 `json:"username,omitempty"`
	Level      LogLevel               `json:"level"`
	Message    string                 `json:"message"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	Method     string                 `json:"method"`
	StatusCode int                    `json:"status_code"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	RequestID  string                 `json:"request_id,omitempty"`
	Duration   time.Duration          `json:"duration"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Error      string                 `json:"error,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
}

// LogFilter represents filters for querying logs
type LogFilter struct {
	UserID    string    `json:"user_id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Level     LogLevel  `json:"level,omitempty"`
	Action    string    `json:"action,omitempty"`
	Resource  string    `json:"resource,omitempty"`
	Method    string    `json:"method,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}

// LogStats represents statistics about logs
type LogStats struct {
	TotalLogs    int64            `json:"total_logs"`
	InfoLogs     int64            `json:"info_logs"`
	WarningLogs  int64            `json:"warning_logs"`
	ErrorLogs    int64            `json:"error_logs"`
	UniqueUsers  int64            `json:"unique_users"`
	TopActions   map[string]int64 `json:"top_actions"`
	TopResources map[string]int64 `json:"top_resources"`
	LastActivity time.Time        `json:"last_activity"`
}

// ToResponse converts LogEntry to a response format
func (l *LogEntry) ToResponse() map[string]interface{} {
	response := map[string]interface{}{
		"id":          l.ID,
		"level":       l.Level,
		"message":     l.Message,
		"action":      l.Action,
		"resource":    l.Resource,
		"method":      l.Method,
		"status_code": l.StatusCode,
		"ip_address":  l.IPAddress,
		"user_agent":  l.UserAgent,
		"duration":    l.Duration.String(),
		"timestamp":   l.Timestamp.Format(time.RFC3339),
	}

	if l.UserID != "" {
		response["user_id"] = l.UserID
	}
	if l.Username != "" {
		response["username"] = l.Username
	}
	if l.RequestID != "" {
		response["request_id"] = l.RequestID
	}
	if l.Error != "" {
		response["error"] = l.Error
	}
	if l.Metadata != nil && len(l.Metadata) > 0 {
		response["metadata"] = l.Metadata
	}

	return response
}
