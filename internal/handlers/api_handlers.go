package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/m4r4v/go-rest-api/internal/models"
	"github.com/m4r4v/go-rest-api/pkg/auth"
	"github.com/m4r4v/go-rest-api/pkg/errors"
	"github.com/m4r4v/go-rest-api/pkg/logger"
	"github.com/m4r4v/go-rest-api/pkg/validation"
)

// DynamicRouter interface for adding dynamic endpoints
type DynamicRouter interface {
	AddDynamicEndpoint(endpoint, method string, response interface{})
	RemoveDynamicEndpoint(endpoint, method string)
}

// APIHandlers contains all HTTP handlers for the new architecture
type APIHandlers struct {
	authService *auth.AuthService
	db          *models.Database
}

// NewAPIHandlers creates a new API handlers instance
func NewAPIHandlers(authService *auth.AuthService) *APIHandlers {
	return &APIHandlers{
		authService: authService,
		db:          models.NewDatabase(),
	}
}

// Setup Endpoints

// Setup handles the initial admin setup
func (h *APIHandlers) Setup(w http.ResponseWriter, r *http.Request) {
	// Check if setup is already complete
	if h.db.IsSetupComplete() {
		h.writeStandardError(w, http.StatusBadRequest, "/setup", "Setup already completed")
		return
	}

	var req validation.RegisterRequest
	if err := validation.ValidateJSON(r, &req); err != nil {
		h.writeStandardError(w, http.StatusBadRequest, "/setup", err.Error())
		return
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		h.writeStandardError(w, http.StatusInternalServerError, "/setup", "Failed to process setup")
		return
	}

	// Create admin user
	adminUser := &models.User{
		ID:       "1",
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "admin",
	}

	// Save admin user
	if err := h.db.CreateUser(adminUser); err != nil {
		logger.Errorf("Failed to create admin user: %v", err)
		h.writeStandardError(w, http.StatusInternalServerError, "/setup", "Failed to create admin user")
		return
	}

	// Mark setup as complete
	h.db.CompleteSetup()

	logger.Infof("Initial setup completed. Admin user created: %s", adminUser.Username)

	response := map[string]interface{}{
		"message":        "Now please login in order to get you the authorization token",
		"login_endpoint": "/login",
	}

	h.writeStandardResponse(w, http.StatusAccepted, "/setup", response)
}

// Authentication Endpoints

// Login handles user login
func (h *APIHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req validation.LoginRequest
	if err := validation.ValidateJSON(r, &req); err != nil {
		h.writeStandardError(w, http.StatusBadRequest, "/login", err.Error())
		return
	}

	// Find user in database
	user, err := h.db.GetUser(req.Username)
	if err != nil {
		h.writeStandardError(w, http.StatusUnauthorized, "/login", "Invalid credentials")
		return
	}

	// Validate password
	if !h.authService.CheckPassword(req.Password, user.Password) {
		h.writeStandardError(w, http.StatusUnauthorized, "/login", "Invalid credentials")
		return
	}

	// Generate JWT token
	roles := []string{"user"}
	if user.Role == "admin" {
		roles = append(roles, "admin")
	}

	token, err := h.authService.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		logger.Errorf("Failed to generate token: %v", err)
		h.writeStandardError(w, http.StatusInternalServerError, "/login", "Failed to generate token")
		return
	}

	response := map[string]interface{}{
		"token": token,
	}

	h.writeStandardResponse(w, http.StatusCreated, "/login", response)
}

// GetMe returns current user information
func (h *APIHandlers) GetMe(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())

	user, err := h.db.GetUserByID(claims.UserID)
	if err != nil {
		appErr := errors.NotFound("User not found")
		h.writeErrorResponse(w, appErr)
		return
	}

	h.writeSuccessResponse(w, user.ToResponse())
}

// User Management Endpoints (Admin Only)

// CreateUser creates a new user (admin only)
func (h *APIHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username" validate:"required,min=3,max=50"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
		Role     string `json:"role" validate:"required,oneof=admin user"`
	}

	if err := validation.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		appErr := errors.InternalServerError("Failed to process user creation")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Generate new user ID
	userID := uuid.New().String()

	// Create user
	user := &models.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	// Save user
	if err := h.db.CreateUser(user); err != nil {
		if err == models.ErrUserExists {
			appErr := errors.BadRequest("Username already exists")
			h.writeErrorResponse(w, appErr)
			return
		}
		logger.Errorf("Failed to create user: %v", err)
		appErr := errors.InternalServerError("Failed to create user")
		h.writeErrorResponse(w, appErr)
		return
	}

	logger.Infof("User created by admin: %s (Role: %s)", user.Username, user.Role)

	h.writeSuccessResponse(w, user.ToResponse())
}

// ListUsers returns all users (admin only)
func (h *APIHandlers) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.db.ListUsers()

	var response []models.UserResponse
	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	h.writeSuccessResponse(w, response)
}

// UpdateUserByAdmin updates any user (admin only)
func (h *APIHandlers) UpdateUserByAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var req struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty" validate:"omitempty,email"`
		Password string `json:"password,omitempty" validate:"omitempty,min=6"`
		Role     string `json:"role,omitempty" validate:"omitempty,oneof=admin user"`
	}

	if err := validation.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	// Find user
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		appErr := errors.NotFound("User not found")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Prepare updates
	updates := &models.User{}
	if req.Email != "" {
		updates.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := h.authService.HashPassword(req.Password)
		if err != nil {
			logger.Errorf("Failed to hash password: %v", err)
			appErr := errors.InternalServerError("Failed to process password update")
			h.writeErrorResponse(w, appErr)
			return
		}
		updates.Password = hashedPassword
	}
	if req.Role != "" {
		updates.Role = req.Role
	}

	// Update user
	if err := h.db.UpdateUser(user.Username, updates); err != nil {
		logger.Errorf("Failed to update user: %v", err)
		appErr := errors.InternalServerError("Failed to update user")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Get updated user
	updatedUser, _ := h.db.GetUserByID(userID)

	logger.Infof("User updated by admin: %s", updatedUser.Username)

	h.writeSuccessResponse(w, updatedUser.ToResponse())
}

// DeleteUser deletes a user (admin only)
func (h *APIHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Find user
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		appErr := errors.NotFound("User not found")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Delete user
	if err := h.db.DeleteUser(user.Username); err != nil {
		logger.Errorf("Failed to delete user: %v", err)
		appErr := errors.InternalServerError("Failed to delete user")
		h.writeErrorResponse(w, appErr)
		return
	}

	logger.Infof("User deleted by admin: %s", user.Username)

	response := map[string]interface{}{
		"message": "User deleted successfully",
		"user_id": userID,
	}

	h.writeSuccessResponse(w, response)
}

// User Self-Management Endpoints

// UpdateMe allows users to update their own email and password
func (h *APIHandlers) UpdateMe(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())

	var req struct {
		Email    string `json:"email,omitempty" validate:"omitempty,email"`
		Password string `json:"password,omitempty" validate:"omitempty,min=6"`
	}

	if err := validation.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	// Find user
	user, err := h.db.GetUserByID(claims.UserID)
	if err != nil {
		appErr := errors.NotFound("User not found")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Prepare updates
	updates := &models.User{}
	if req.Email != "" {
		updates.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := h.authService.HashPassword(req.Password)
		if err != nil {
			logger.Errorf("Failed to hash password: %v", err)
			appErr := errors.InternalServerError("Failed to process password update")
			h.writeErrorResponse(w, appErr)
			return
		}
		updates.Password = hashedPassword
	}

	// Update user
	if err := h.db.UpdateUser(user.Username, updates); err != nil {
		logger.Errorf("Failed to update user: %v", err)
		appErr := errors.InternalServerError("Failed to update user")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Get updated user
	updatedUser, _ := h.db.GetUserByID(claims.UserID)

	logger.Infof("User updated own profile: %s", updatedUser.Username)

	h.writeSuccessResponse(w, updatedUser.ToResponse())
}

// Resource Management Endpoints

// CreateResource creates a new resource
func (h *APIHandlers) CreateResource(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())

	var req struct {
		Name        string                 `json:"name" validate:"required,min=1,max=100"`
		Description string                 `json:"description,omitempty"`
		Data        map[string]interface{} `json:"data,omitempty"`
	}

	if err := validation.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	// Create resource
	resource := &models.Resource{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Data:        req.Data,
		CreatedBy:   claims.UserID,
	}

	// Save resource
	if err := h.db.CreateResource(resource); err != nil {
		logger.Errorf("Failed to create resource: %v", err)
		appErr := errors.InternalServerError("Failed to create resource")
		h.writeErrorResponse(w, appErr)
		return
	}

	logger.Infof("Resource created: %s by user %s", resource.Name, claims.Username)

	h.writeSuccessResponse(w, resource)
}

// ListResources returns all resources
func (h *APIHandlers) ListResources(w http.ResponseWriter, r *http.Request) {
	resources := h.db.ListResources()
	h.writeSuccessResponse(w, resources)
}

// GetResource returns a specific resource
func (h *APIHandlers) GetResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["id"]

	resource, err := h.db.GetResource(resourceID)
	if err != nil {
		appErr := errors.NotFound("Resource not found")
		h.writeErrorResponse(w, appErr)
		return
	}

	h.writeSuccessResponse(w, resource)
}

// UpdateResource updates a resource (creator or admin only)
func (h *APIHandlers) UpdateResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["id"]
	claims := auth.GetClaimsFromContext(r.Context())

	var req struct {
		Name        string                 `json:"name,omitempty"`
		Description string                 `json:"description,omitempty"`
		Data        map[string]interface{} `json:"data,omitempty"`
	}

	if err := validation.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	// Find resource
	resource, err := h.db.GetResource(resourceID)
	if err != nil {
		appErr := errors.NotFound("Resource not found")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Check permissions (creator or admin)
	if resource.CreatedBy != claims.UserID && !claims.HasRole("admin") {
		appErr := errors.Forbidden("You can only update your own resources")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Prepare updates
	updates := &models.Resource{}
	if req.Name != "" {
		updates.Name = req.Name
	}
	if req.Description != "" {
		updates.Description = req.Description
	}
	if req.Data != nil {
		updates.Data = req.Data
	}

	// Update resource
	if err := h.db.UpdateResource(resourceID, updates); err != nil {
		logger.Errorf("Failed to update resource: %v", err)
		appErr := errors.InternalServerError("Failed to update resource")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Get updated resource
	updatedResource, _ := h.db.GetResource(resourceID)

	logger.Infof("Resource updated: %s by user %s", updatedResource.Name, claims.Username)

	h.writeSuccessResponse(w, updatedResource)
}

// DeleteResource deletes a resource (creator or admin only)
func (h *APIHandlers) DeleteResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["id"]
	claims := auth.GetClaimsFromContext(r.Context())

	// Find resource
	resource, err := h.db.GetResource(resourceID)
	if err != nil {
		appErr := errors.NotFound("Resource not found")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Check permissions (creator or admin)
	if resource.CreatedBy != claims.UserID && !claims.HasRole("admin") {
		appErr := errors.Forbidden("You can only delete your own resources")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Delete resource
	if err := h.db.DeleteResource(resourceID); err != nil {
		logger.Errorf("Failed to delete resource: %v", err)
		appErr := errors.InternalServerError("Failed to delete resource")
		h.writeErrorResponse(w, appErr)
		return
	}

	logger.Infof("Resource deleted: %s by user %s", resource.Name, claims.Username)

	response := map[string]interface{}{
		"message":     "Resource deleted successfully",
		"resource_id": resourceID,
	}

	h.writeSuccessResponse(w, response)
}

// System Endpoints

// GetStatus returns server status
func (h *APIHandlers) GetStatus(w http.ResponseWriter, r *http.Request) {
	stats := h.db.GetStats()

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "2.0.0",
		"database":  stats,
	}

	h.writeSuccessResponse(w, response)
}

// GetHealth returns health check
func (h *APIHandlers) GetHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	h.writeSuccessResponse(w, response)
}

// Helper methods

// writeSuccessResponse writes a successful JSON response with proper headers
func (h *APIHandlers) writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	h.writeSuccessResponseWithStatus(w, http.StatusOK, "", data)
}

// writeSuccessResponseWithStatus writes a successful JSON response with custom status
func (h *APIHandlers) writeSuccessResponseWithStatus(w http.ResponseWriter, statusCode int, resource string, data interface{}) {
	// Set development-friendly headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"http_status_code":    fmt.Sprintf("%d", statusCode),
		"http_status_message": http.StatusText(statusCode),
		"resource":            resource,
		"app":                 "Go REST API Framework",
		"timestamp":           time.Now().Format(time.RFC3339),
		"response":            data,
	}

	json.NewEncoder(w).Encode(response)
}

// writeStandardResponse writes a response in the standard format
func (h *APIHandlers) writeStandardResponse(w http.ResponseWriter, statusCode int, resource string, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"http_status_code":    fmt.Sprintf("%d", statusCode),
		"http_status_message": http.StatusText(statusCode),
		"resource":            resource,
		"app":                 "Go REST API Framework",
		"timestamp":           time.Now().Format(time.RFC3339),
		"response":            data,
	}

	json.NewEncoder(w).Encode(response)
}

// writeStandardError writes an error response in the standard format
func (h *APIHandlers) writeStandardError(w http.ResponseWriter, statusCode int, resource, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"http_status_code":    fmt.Sprintf("%d", statusCode),
		"http_status_message": http.StatusText(statusCode),
		"resource":            resource,
		"app":                 "Go REST API Framework",
		"timestamp":           time.Now().Format(time.RFC3339),
		"response": map[string]interface{}{
			"error": map[string]interface{}{
				"message": message,
			},
		},
	}

	json.NewEncoder(w).Encode(response)
}

// writeErrorResponse writes an error JSON response with proper headers
func (h *APIHandlers) writeErrorResponse(w http.ResponseWriter, appErr *errors.AppError) {
	// Set development-friendly headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	w.WriteHeader(appErr.Status)

	response := map[string]interface{}{
		"success":     false,
		"status_code": appErr.Status,
		"status":      http.StatusText(appErr.Status),
		"error": map[string]interface{}{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// Dynamic Endpoint Management

// AddDynamicEndpoint adds a new dynamic endpoint based on resource data with authentication
func (h *APIHandlers) AddDynamicEndpoint(router *mux.Router, endpoint, method string, response interface{}) {
	// The endpoint path (don't add /v1 prefix since the protected router already has it)
	routePath := endpoint

	// Create the full path for logging and response (with /v1 prefix)
	fullPath := "/v1" + endpoint

	// Create a handler that returns the specified response
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Get authenticated user from context (set by auth middleware)
		claims := auth.GetClaimsFromContext(r.Context())

		// Set development-friendly headers
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("X-API-Framework", "Go-REST-API-v2.0")
		w.Header().Set("X-Dynamic-Endpoint", "true")

		if claims != nil {
			w.Header().Set("X-Authenticated-User", claims.Username)
		}

		w.WriteHeader(http.StatusOK)

		// Return the response with proper format including user info
		apiResponse := map[string]interface{}{
			"success":     true,
			"status_code": http.StatusOK,
			"status":      "OK",
			"response":    response,
			"timestamp":   time.Now().Format(time.RFC3339),
			"endpoint":    fullPath,
			"method":      method,
		}

		// Add user info if authenticated
		if claims != nil {
			apiResponse["user"] = claims.Username
			apiResponse["user_id"] = claims.UserID
		}

		json.NewEncoder(w).Encode(apiResponse)
	}

	// Add the route to the router (use routePath without /v1 prefix)
	router.HandleFunc(routePath, handler).Methods(method)

	logger.Infof("Dynamic endpoint created: %s %s (requires authentication)", method, fullPath)
}

// LoadExistingEndpoints loads all existing resource endpoints on server startup
func (h *APIHandlers) LoadExistingEndpoints(dynamicRouter DynamicRouter) {
	resources := h.db.ListResources()

	for _, resource := range resources {
		if resource.Data != nil {
			// Check if the resource has endpoint data
			if endpoint, ok := resource.Data["endpoint"].(string); ok {
				if method, ok := resource.Data["method"].(string); ok {
					if response, ok := resource.Data["response"]; ok {
						dynamicRouter.AddDynamicEndpoint(endpoint, method, response)
					}
				}
			}
		}
	}

	logger.Infof("Loaded %d existing dynamic endpoints", len(resources))
}

// CreateResourceWithDynamicEndpoint creates a resource and its dynamic endpoint
func (h *APIHandlers) CreateResourceWithDynamicEndpoint(w http.ResponseWriter, r *http.Request, dynamicRouter DynamicRouter) {
	claims := auth.GetClaimsFromContext(r.Context())

	var req struct {
		Name        string                 `json:"name" validate:"required,min=1,max=100"`
		Description string                 `json:"description,omitempty"`
		Data        map[string]interface{} `json:"data,omitempty"`
	}

	if err := validation.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	// Create resource
	resource := &models.Resource{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Data:        req.Data,
		CreatedBy:   claims.UserID,
	}

	// Save resource
	if err := h.db.CreateResource(resource); err != nil {
		logger.Errorf("Failed to create resource: %v", err)
		appErr := errors.InternalServerError("Failed to create resource")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Create dynamic endpoint if resource has endpoint data
	if resource.Data != nil {
		if endpoint, ok := resource.Data["endpoint"].(string); ok {
			if method, ok := resource.Data["method"].(string); ok {
				if response, ok := resource.Data["response"]; ok {
					dynamicRouter.AddDynamicEndpoint(endpoint, method, response)
				}
			}
		}
	}

	logger.Infof("Resource created: %s by user %s", resource.Name, claims.Username)

	h.writeSuccessResponse(w, resource)
}

// UpdateResourceWithDynamicEndpoint updates a resource and its dynamic endpoint
func (h *APIHandlers) UpdateResourceWithDynamicEndpoint(w http.ResponseWriter, r *http.Request, dynamicRouter DynamicRouter) {
	vars := mux.Vars(r)
	resourceID := vars["id"]
	claims := auth.GetClaimsFromContext(r.Context())

	var req struct {
		Name        string                 `json:"name,omitempty"`
		Description string                 `json:"description,omitempty"`
		Data        map[string]interface{} `json:"data,omitempty"`
	}

	if err := validation.ValidateJSON(r, &req); err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	// Find resource
	resource, err := h.db.GetResource(resourceID)
	if err != nil {
		appErr := errors.NotFound("Resource not found")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Check permissions (creator or admin)
	if resource.CreatedBy != claims.UserID && !claims.HasRole("admin") {
		appErr := errors.Forbidden("You can only update your own resources")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Prepare updates
	updates := &models.Resource{}
	if req.Name != "" {
		updates.Name = req.Name
	}
	if req.Description != "" {
		updates.Description = req.Description
	}
	if req.Data != nil {
		updates.Data = req.Data
	}

	// Update resource
	if err := h.db.UpdateResource(resourceID, updates); err != nil {
		logger.Errorf("Failed to update resource: %v", err)
		appErr := errors.InternalServerError("Failed to update resource")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Get updated resource
	updatedResource, _ := h.db.GetResource(resourceID)

	// Update dynamic endpoint if resource has endpoint data
	if updatedResource.Data != nil {
		if endpoint, ok := updatedResource.Data["endpoint"].(string); ok {
			if method, ok := updatedResource.Data["method"].(string); ok {
				if response, ok := updatedResource.Data["response"]; ok {
					// Note: Since Gorilla mux doesn't support removing routes,
					// we just add the new endpoint (it will override the old one)
					dynamicRouter.AddDynamicEndpoint(endpoint, method, response)
				}
			}
		}
	}

	logger.Infof("Resource updated: %s by user %s", updatedResource.Name, claims.Username)

	h.writeSuccessResponse(w, updatedResource)
}

// DeleteResourceWithDynamicEndpoint deletes a resource and its dynamic endpoint
func (h *APIHandlers) DeleteResourceWithDynamicEndpoint(w http.ResponseWriter, r *http.Request, dynamicRouter DynamicRouter) {
	vars := mux.Vars(r)
	resourceID := vars["id"]
	claims := auth.GetClaimsFromContext(r.Context())

	// Find resource
	resource, err := h.db.GetResource(resourceID)
	if err != nil {
		appErr := errors.NotFound("Resource not found")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Check permissions (creator or admin)
	if resource.CreatedBy != claims.UserID && !claims.HasRole("admin") {
		appErr := errors.Forbidden("You can only delete your own resources")
		h.writeErrorResponse(w, appErr)
		return
	}

	// Remove dynamic endpoint if resource has endpoint data
	if resource.Data != nil {
		if endpoint, ok := resource.Data["endpoint"].(string); ok {
			if method, ok := resource.Data["method"].(string); ok {
				dynamicRouter.RemoveDynamicEndpoint(endpoint, method)
			}
		}
	}

	// Delete resource
	if err := h.db.DeleteResource(resourceID); err != nil {
		logger.Errorf("Failed to delete resource: %v", err)
		appErr := errors.InternalServerError("Failed to delete resource")
		h.writeErrorResponse(w, appErr)
		return
	}

	logger.Infof("Resource deleted: %s by user %s", resource.Name, claims.Username)

	response := map[string]interface{}{
		"message":     "Resource deleted successfully",
		"resource_id": resourceID,
	}

	h.writeSuccessResponse(w, response)
}
