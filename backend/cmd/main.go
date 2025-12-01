package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	adminHandler "bestdoctors_service/admin/handlers"
	adminMW "bestdoctors_service/admin/middleware"
	"bestdoctors_service/middleware"
	"bestdoctors_service/routes"

	"golang.org/x/time/rate"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\n", "")
	return value
}


func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := getEnv("ALLOWED_ORIGIN")
		if origin == "" {
			origin = "http://localhost" 
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") 
		w.Header().Set("Access-Control-Max-Age", "3600")

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	redisURL := getEnv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable is required")
	}

	if err := routes.InitSessionStore(redisURL); err != nil {
		log.Fatalf("Failed to initialize session store: %v", err)
	}

	loginLimiter := middleware.NewIPRateLimiter(rate.Limit(5.0/60.0), 5)
	apiLimiter := middleware.NewIPRateLimiter(rate.Limit(100.0/60.0), 100)

	authMW := middleware.AuthMiddleware(routes.GetSessionStore())

	mux := http.NewServeMux()

	mux.Handle("/template_bestdoctors/",
		http.StripPrefix("/template_bestdoctors/",
			http.FileServer(http.Dir("./template_bestdoctors"))))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.Handle("/auth/login", middleware.RateLimitMiddleware(loginLimiter)(http.HandlerFunc(routes.LoginHandler)))
	mux.HandleFunc("/auth/logout", routes.LogoutHandler)

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/auth/me", authMW(http.HandlerFunc(routes.MeHandler)))

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/bestdoctors/sessionphone", routes.SessionPhoneHandler)
	protectedMux.HandleFunc("/bestdoctors/sessionphone/active", routes.ToggleAIHandler)
	protectedMux.HandleFunc("/bestdoctors/chathistory", routes.ChatHistoryHandler)
	protectedMux.HandleFunc("/bestdoctors/sessiondelta", routes.SessionDeltaHandler)
	protectedMux.HandleFunc("/bestdoctors/metrics/session", routes.SessionMetricsHandler)
	protectedMux.HandleFunc("/bestdoctors/metrics/abandonment", routes.AbandonmentRateHandler)
	protectedMux.HandleFunc("/bestdoctors/metrics/flowdepth", routes.FlowDepthHandler)
	protectedMux.HandleFunc("/bestdoctors/metrics/reengagement", routes.ReengagementRateHandler)
	protectedMux.HandleFunc("/bestdoctors/sendmessage", routes.SendMessageHandler)
	protectedMux.HandleFunc("/bestdoctors/report", routes.ReportHandler)

	mux.Handle("/bestdoctors/", middleware.RateLimitMiddleware(apiLimiter)(authMW(protectedMux)))

	
	adminHandler.InitAdminSessionStore(routes.GetSessionStore())
	
	adminLimiter := middleware.NewIPRateLimiter(rate.Limit(2.0/60.0), 2) 
	mux.Handle("/admin/auth", middleware.RateLimitMiddleware(adminLimiter)(http.HandlerFunc(adminHandler.SuperAdminLoginHandler)))
	mux.HandleFunc("/admin/logout", adminHandler.LogoutHandler)
	
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/users", adminHandler.UsersHandler)     
	adminMux.HandleFunc("/admin/users/", adminHandler.UserHandler)     
	
	adminAuthMW := adminMW.AdminMiddleware(routes.GetSessionStore())
	mux.Handle("/admin/users", adminAuthMW(adminMux))
	mux.Handle("/admin/users/", adminAuthMW(adminMux))

	port := getEnv("PORT")
	if port == "" {
		port = "9002"
	}

	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(mux)))
}
