package db

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	SupabaseDB *gorm.DB
	PostgresDB *gorm.DB
	DB         *gorm.DB
)

func init() {
	cfg := &gorm.Config{
		PrepareStmt: false,
		Logger:      logger.Default.LogMode(logger.Silent),
	}

	cleanEnv := func(key string) string {
		v := os.Getenv(key)
		v = strings.TrimSpace(v)
		v = strings.ReplaceAll(v, "\r", "")
		v = strings.ReplaceAll(v, "\n", "")
		return v
	}
	
	pgHost := cleanEnv("PG_HOST")
	pgPort := cleanEnv("PG_PORT")
	pgUser := cleanEnv("PG_USER")
	pgPassword := cleanEnv("PG_PASSWORD")
	pgDatabase := cleanEnv("PG_DATABASE")
	pgSSLMode := cleanEnv("PG_SSLMODE")
	
	if pgSSLMode == "" {
		pgSSLMode = "disable"
	}
	
	supaDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgHost, pgPort, pgUser, pgPassword, pgDatabase, pgSSLMode)

	supaConn, err := gorm.Open(postgres.Open(supaDSN), cfg)
	if err != nil {
		log.Fatalf("failed to connect to Supabase: %v", err)
	}
	sqlSupa, err := supaConn.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlSupa.SetMaxIdleConns(10)
	sqlSupa.SetMaxOpenConns(100)
	sqlSupa.SetConnMaxLifetime(5 * time.Minute)
	SupabaseDB = supaConn
	DB = supaConn

}

func toURLWithParams(raw string, kv map[string]string) string {
	if strings.HasPrefix(raw, "postgres://") || strings.HasPrefix(raw, "postgresql://") {
		u, err := url.Parse(raw)
		if err != nil {
			log.Fatalf("invalid SUPRABASE_PGSQL URL: %v", err)
		}
		q := u.Query()
		for k, v := range kv {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		return u.String()
	}

	parts := fieldsToMap(raw)
	host := parts["host"]
	user := parts["user"]
	pass := parts["password"]
	db := parts["dbname"]
	port := parts["port"]
	if host == "" || user == "" || db == "" {
		log.Fatal("SUPRABASE_PGSQL must be a postgres URL or include host= user= dbname= [password=] [port=]")
	}

	if port == "" {
		port = "5432"
	}
	u := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   "/" + db,
	}
	q := u.Query()
	for k, v := range kv {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func fieldsToMap(s string) map[string]string {
	m := map[string]string{}
	for _, f := range strings.Fields(s) {
		if i := strings.IndexByte(f, '='); i > 0 && i < len(f)-1 {
			k := strings.TrimSpace(f[:i])
			v := strings.TrimSpace(f[i+1:])
			m[k] = v
		}
	}
	return m
}

func RetryForever(initialDelay time.Duration, fn func() error) {
	delay := initialDelay
	for {
		if err := fn(); err != nil {
			SupabaseDB.Exec("DEALLOCATE ALL")
			time.Sleep(delay)
			if delay < 5*time.Second {
				delay *= 2
			}
			continue
		}
		return
	}
}
