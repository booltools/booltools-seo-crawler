package technical

import (
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

var expectedNoindexPatterns = []string{
	"/login", "/signin", "/sign-in", "/signup", "/sign-up",
	"/register", "/auth/", "/forgot-password", "/reset-password",
	"/verify", "/confirm", "/logout", "/sign-out",
	"/admin", "/dashboard", "/account", "/settings",
	"/cart", "/checkout",
}

func isExpectedNoindex(pageURL string) bool {
	lowered := strings.ToLower(pageURL)
	for _, pattern := range expectedNoindexPatterns {
		if strings.Contains(lowered, pattern) {
			return true
		}
	}
	return false
}

type MetaRobotsChecker struct{}

func (c *MetaRobotsChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	metaRobots, _ := page.Document.Find(`meta[name="robots"]`).Attr("content")
	metaRobots = strings.ToLower(strings.TrimSpace(metaRobots))

	xRobotsTag := strings.ToLower(strings.TrimSpace(page.Headers.Get("X-Robots-Tag")))

	noindexRule := valueobject.NewAuditRule("meta_robots_noindex", valueobject.CategoryTechnical, valueobject.SeverityInfo)
	noindexRule.AffectedURL = page.URL

	hasNoindex := strings.Contains(metaRobots, "noindex") || strings.Contains(xRobotsTag, "noindex")
	if hasNoindex {
		if isExpectedNoindex(page.URL) {
			noindexRule.Pass("Page has noindex (expected for auth/utility pages)")
		} else {
			noindexRule.Warn(
				"Page has a noindex directive — it will not appear in search results",
				"This page is intentionally excluded from search engine indexing. If this is unintentional, remove the noindex directive from the meta robots tag or X-Robots-Tag header.",
			)
		}
	} else {
		noindexRule.Pass("Page is indexable (no noindex directive)")
	}
	rules = append(rules, noindexRule)

	conflictRule := valueobject.NewAuditRule("meta_robots_conflict", valueobject.CategoryTechnical, valueobject.SeverityHigh)
	conflictRule.AffectedURL = page.URL
	if metaRobots != "" && xRobotsTag != "" {
		metaHasNoindex := strings.Contains(metaRobots, "noindex")
		headerHasNoindex := strings.Contains(xRobotsTag, "noindex")
		if metaHasNoindex != headerHasNoindex {
			conflictRule.Fail(
				"Conflicting robots directives between meta tag and X-Robots-Tag header",
				"Ensure the meta robots tag and X-Robots-Tag HTTP header have consistent directives.",
			)
		} else {
			conflictRule.Pass("Meta robots and X-Robots-Tag directives are consistent")
		}
	} else {
		conflictRule.Pass("No conflicting robots directives")
	}
	rules = append(rules, conflictRule)

	return rules
}
