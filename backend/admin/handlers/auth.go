package admin

import (
	"encoding/json"
	"net/http"
	"time"

	adminMW "bestdoctors_service/admin/middleware"
	"bestdoctors_service/internal/session"
)

var adminSessionStore *session.Store

func InitAdminSessionStore(store *session.Store) {
	adminSessionStore = store
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SuperAdminLoginHandler(w http.ResponseWriter, r *http.Request) {
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

	validUsername, validPassword := adminMW.GetSuperAdminCredentials()
	
	if req.Username != validUsername || req.Password != validPassword {
		time.Sleep(200 * time.Millisecond)
				
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid SuperAdmin credentials",
		})
		return
	}

	sessionID, err := adminSessionStore.Create(r.Context(), session.SessionData{
		UserID:   0,      
		Username: req.Username,
		Role:     "superadmin",
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, 
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, 
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Message: "SuperAdmin login successful",
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("admin_session_id")
	if err == nil && cookie.Value != "" {
		adminSessionStore.Delete(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}
