package middleware

import (
	"net/http"

	"bestdoctors_service/internal/session"
)

// SystemMiddleware validates if the logged user has access to a specific system
func SystemMiddleware(requiredSystem string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionData, ok := r.Context().Value(SessionDataKey).(*session.SessionData)
			if !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// If the required system is not empty, check if it's in the user's systems list
			if requiredSystem != "" {
				hasAccess := false

				for _, sys := range sessionData.System {
					if sys == requiredSystem {
						hasAccess = true
						break
					}
				}

				if !hasAccess {
					http.Error(w, "Forbidden: Missing required system access", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
