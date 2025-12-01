package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"bestdoctors_service/internal/session"
)

type contextKey string

const SessionDataKey contextKey = "session_data"

func AdminMiddleware(sessionStore *session.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("admin_session_id")
			if err != nil {
				http.Error(w, "Unauthorized: Admin access required", http.StatusUnauthorized)
				return
			}

			sessionID := cookie.Value
			if sessionID == "" {
				http.Error(w, "Unauthorized: Admin access required", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			sessionData, err := sessionStore.Get(ctx, sessionID)
			if err != nil {
				http.Error(w, "Unauthorized: Invalid admin session", http.StatusUnauthorized)
				return
			}

			if sessionData.Role != "superadmin" {
				http.Error(w, "Forbidden: SuperAdmin access required", http.StatusForbidden)
				return
			}

			ctx = context.WithValue(ctx, SessionDataKey, sessionData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func IsSuperAdmin(r *http.Request) bool {
	sessionData, ok := r.Context().Value(SessionDataKey).(*session.SessionData)
	if !ok {
		return false
	}
	return sessionData.Role == "superadmin"
}

func RequireSuperAdmin(w http.ResponseWriter, r *http.Request) bool {
	if !IsSuperAdmin(r) {
		http.Error(w, "Forbidden: SuperAdmin access required", http.StatusForbidden)
		return false
	}
	return true
}

func GetSuperAdminCredentials() (username, password string) {
	cleanEnv := func(key string) string {
		v := os.Getenv(key)
		v = strings.TrimSpace(v)
		v = strings.ReplaceAll(v, "\r", "")
		v = strings.ReplaceAll(v, "\n", "")
		return v
	}
	
	username = cleanEnv("SUPERADMIN_USERNAME")
	password = cleanEnv("SUPERADMIN_PASSWORD")
	
	return username, password
}
