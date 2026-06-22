package middleware

import "net/http"

const maxBodySize = 1 << 20 // 1 MB

func BodyLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Body != nil {
			request.Body = http.MaxBytesReader(writer, request.Body, maxBodySize)
		}
		next.ServeHTTP(writer, request)
	})
}
