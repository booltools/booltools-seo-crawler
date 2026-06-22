package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if strings.Contains(request.URL.Path, "/progress") {
				next.ServeHTTP(writer, request)
				return
			}

			ctx, cancel := context.WithTimeout(request.Context(), timeout)
			defer cancel()

			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
