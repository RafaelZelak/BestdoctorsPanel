package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	adminMW "bestdoctors_service/admin/middleware"
	"bestdoctors_service/admin/validators"
	"bestdoctors_service/internal/db"
	"bestdoctors_service/models"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if !adminMW.RequireSuperAdmin(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		listUsers(w, r)
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if !adminMW.RequireSuperAdmin(w, r) {
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/admin/users/")
	parts := strings.Split(path, "/")
	
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	userIDStr := parts[0]
	
	if len(parts) == 2 && parts[1] == "password" {
		resetPassword(w, r, userIDStr)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getUser(w, r, userID)
	case http.MethodPut:
		updateUser(w, r, userID)
	case http.MethodDelete:
		deleteUser(w, r, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /admin/users - List all users
func listUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"users":   users,
	})
}

// POST /admin/users - Create user
func createUser(w http.ResponseWriter, r *http.Request) {
	var req validators.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var existingUser models.User
	if err := db.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Username already exists",
		})
		return
	}

	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Email already exists",
		})
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Role:     "user",
		IsActive: req.IsActive,
	}

	if err := user.SetPassword(req.Password); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User created successfully",
		"user":    user,
	})
}

// GET /admin/users/:id - Get user details
func getUser(w http.ResponseWriter, r *http.Request, userID int) {
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

// PUT /admin/users/:id - Update user
func updateUser(w http.ResponseWriter, r *http.Request, userID int) {
	var req validators.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User updated successfully",
		"user":    user,
	})
}

// DELETE /admin/users/:id - Delete user (soft delete via is_active)
func deleteUser(w http.ResponseWriter, r *http.Request, userID int) {
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.IsActive = false
	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}

// PATCH /admin/users/:id/password - Reset password
func resetPassword(w http.ResponseWriter, r *http.Request, userIDStr string) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req validators.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := user.SetPassword(req.NewPassword); err != nil {
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Password reset successfully",
	})
}
