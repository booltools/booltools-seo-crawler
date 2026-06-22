package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.Mutex
	limit    int
	window   time.Duration
}

func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	limiter := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    requestsPerMinute,
		window:   time.Minute,
	}

	go limiter.cleanup()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			clientIP := extractClientIP(request)

			if !limiter.allow(clientIP) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusTooManyRequests)
				writer.Write([]byte(`{"error":"rate limit exceeded, try again later"}`))
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func (rl *rateLimiter) allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	timestamps := rl.requests[clientIP]
	validTimestamps := make([]time.Time, 0)
	for _, timestamp := range timestamps {
		if timestamp.After(windowStart) {
			validTimestamps = append(validTimestamps, timestamp)
		}
	}

	if len(validTimestamps) >= rl.limit {
		rl.requests[clientIP] = validTimestamps
		return false
	}

	rl.requests[clientIP] = append(validTimestamps, now)
	return true
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		cutoff := time.Now().Add(-rl.window)
		for ip, timestamps := range rl.requests {
			valid := make([]time.Time, 0)
			for _, timestamp := range timestamps {
				if timestamp.After(cutoff) {
					valid = append(valid, timestamp)
				}
			}
			if len(valid) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = valid
			}
		}
		rl.mutex.Unlock()
	}
}

func extractClientIP(request *http.Request) string {
	if forwarded := request.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.SplitN(forwarded, ",", 2)
		return strings.TrimSpace(parts[0])
	}

	if realIP := request.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return request.RemoteAddr
	}
	return host
}
