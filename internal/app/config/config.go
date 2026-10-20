package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	GinMode        string
	DatabaseURL    string
	RedisURL       string
	AllowedOrigins []string
	TrustedProxies []string
	CookieSecure   bool
	CookieSameSite string
	CookieDomain   string
}

func Load() Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}

	dbURL := os.Getenv("DATABASE_URL")
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	var allowedOrigins []string
	if origins := os.Getenv("CORS_ORIGINS"); origins != "" {
		for _, o := range strings.Split(origins, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	}

	var trustedProxies []string
	if proxies := os.Getenv("TRUSTED_PROXIES"); proxies != "" {
		for _, p := range strings.Split(proxies, ",") {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				trustedProxies = append(trustedProxies, trimmed)
			}
		}
	}

	cookieSecure := os.Getenv("COOKIE_SECURE") == "true" || (os.Getenv("COOKIE_SECURE") == "" && ginMode == "release")

	cookieSameSite := os.Getenv("COOKIE_SAME_SITE")
	if cookieSameSite == "" {
		cookieSameSite = "Lax"
	}

	cookieDomain := os.Getenv("COOKIE_DOMAIN")

	return Config{
		Port:           port,
		GinMode:        ginMode,
		DatabaseURL:    dbURL,
		RedisURL:       redisURL,
		AllowedOrigins: allowedOrigins,
		TrustedProxies: trustedProxies,
		CookieSecure:   cookieSecure,
		CookieSameSite: cookieSameSite,
		CookieDomain:   cookieDomain,
	}
}