package crawler

import "strings"

var authPagePatterns = []string{
	"/login", "/signin", "/sign-in", "/signup", "/sign-up",
	"/register", "/auth/", "/forgot-password", "/reset-password",
	"/verify", "/confirm", "/logout", "/sign-out",
}

var nonEditorialPatterns = []string{
	"/login", "/signin", "/sign-in", "/signup", "/sign-up",
	"/register", "/auth/", "/forgot-password", "/reset-password",
	"/verify", "/confirm", "/logout", "/sign-out",
	"/privacy", "/terms", "/tos", "/legal",
	"/cookie-policy", "/cookies", "/disclaimer",
	"/gdpr", "/imprint", "/impressum",
	"/404", "/500", "/error",
}

func IsAuthPage(pageURL string) bool {
	lowered := strings.ToLower(pageURL)
	for _, pattern := range authPagePatterns {
		if strings.Contains(lowered, pattern) {
			return true
		}
	}
	return false
}

func IsNonEditorialPage(pageURL string) bool {
	lowered := strings.ToLower(pageURL)
	for _, pattern := range nonEditorialPatterns {
		if strings.Contains(lowered, pattern) {
			return true
		}
	}
	return false
}
