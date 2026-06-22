package technical

import (
	"fmt"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type HTTPStatusChecker struct{}

func (c *HTTPStatusChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	statusRule := valueobject.NewAuditRule("http_status_ok", valueobject.CategoryTechnical, valueobject.SeverityCritical)
	statusRule.AffectedURL = page.URL

	switch {
	case page.StatusCode == 200:
		statusRule.Pass("Page returns 200 OK")
	case page.StatusCode >= 400 && page.StatusCode < 500:
		statusRule.Fail(
			fmt.Sprintf("Page returns %d client error", page.StatusCode),
			"Fix or remove this page. If the content has moved, set up a 301 redirect to the new URL.",
		)
	case page.StatusCode >= 500:
		statusRule.Fail(
			fmt.Sprintf("Page returns %d server error", page.StatusCode),
			"Investigate the server error. Check server logs and fix the underlying issue.",
		)
	case page.StatusCode >= 300 && page.StatusCode < 400:
		statusRule.Warn(
			fmt.Sprintf("Page returns %d redirect", page.StatusCode),
			"Update internal links to point directly to the final destination URL.",
		)
	default:
		statusRule.Warn(
			fmt.Sprintf("Unexpected status code: %d", page.StatusCode),
			"Investigate this unusual HTTP status code.",
		)
	}
	rules = append(rules, statusRule)

	return rules
}
