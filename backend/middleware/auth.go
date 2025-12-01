package middleware

import (
	"context"
	"net/http"

	"bestdoctors_service/internal/session"
)

type contextKey string

const SessionDataKey contextKey = "session_data"

func AuthMiddleware(sessionStore *session.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sessionID := cookie.Value
			if sessionID == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			sessionData, err := sessionStore.Get(ctx, sessionID)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, SessionDataKey, sessionData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionData, ok := r.Context().Value(SessionDataKey).(*session.SessionData)
			if !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			allowed := false
			for _, role := range allowedRoles {
				if sessionData.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
