package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/cors"
)

func CORSMiddleware() func(http.Handler) http.Handler {
	allowedOrigins := []string{"http://localhost:*", "https://localhost:*"}

	if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
		allowedOrigins = strings.Split(envOrigins, ",")
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
