package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"bestdoctors_service/internal/db"
	"bestdoctors_service/internal/session"
	"bestdoctors_service/middleware"
	"bestdoctors_service/models"
)

var sessionStore *session.Store

func InitSessionStore(redisURL string) error {
	var err error
	sessionStore, err = session.NewStore(redisURL)
	if err != nil {
		return err
	}
	return nil
}

func GetSessionStore() *session.Store {
	return sessionStore
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	User    *UserInfo `json:"user,omitempty"`
}

type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Username and password are required",
		})
		return
	}

	var user models.User
	if err := db.DB.Where("username = ? AND is_active = true", req.Username).First(&user).Error; err != nil {
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	now := time.Now()
	user.LastLogin = &now
	db.DB.Save(&user)

	sessionID, err := sessionStore.Create(r.Context(), session.SessionData{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,                      
		Secure:   true,                      
		SameSite: http.SameSiteStrictMode,   
		MaxAge:   86400,                    
	})


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Message: "Login successful",
		User: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Role:     user.Role,
		},
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		sessionStore.Delete(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, 
	})


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionData, ok := r.Context().Value(middleware.SessionDataKey).(*session.SessionData)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := db.DB.First(&user, sessionData.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserInfo{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Role:     user.Role,
	})
}
