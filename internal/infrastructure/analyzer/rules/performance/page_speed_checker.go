package performance

import (
	"fmt"
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type PageSpeedChecker struct{}

func (c *PageSpeedChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	pageSizeRule := valueobject.NewAuditRule("page_size", valueobject.CategoryPerformance, valueobject.SeverityMedium)
	pageSizeRule.AffectedURL = page.URL
	pageSizeKB := page.ContentLength / 1024
	if pageSizeKB > 3000 {
		pageSizeRule.Fail(
			fmt.Sprintf("Page is too large (%d KB)", pageSizeKB),
			"Reduce the total page size to under 3MB by optimizing images, minifying code, and removing unused resources.",
		)
	} else if pageSizeKB > 1500 {
		pageSizeRule.Warn(
			fmt.Sprintf("Page size is large (%d KB)", pageSizeKB),
			"Consider reducing page size for faster loading on slower connections.",
		)
	} else {
		pageSizeRule.Pass(fmt.Sprintf("Page size is acceptable (%d KB)", pageSizeKB))
	}
	rules = append(rules, pageSizeRule)

	htmlSizeRule := valueobject.NewAuditRule("html_size", valueobject.CategoryPerformance, valueobject.SeverityLow)
	htmlSizeRule.AffectedURL = page.URL
	htmlSizeKB := len(page.HTML) / 1024
	if htmlSizeKB > 150 {
		htmlSizeRule.Warn(
			fmt.Sprintf("HTML document is large (%d KB)", htmlSizeKB),
			"Reduce HTML document size to under 150KB. Remove inline scripts/styles and use external files. Note: i18n pages or apps with inline data may exceed this threshold acceptably.",
		)
	} else {
		htmlSizeRule.Pass(fmt.Sprintf("HTML document size is acceptable (%d KB)", htmlSizeKB))
	}
	rules = append(rules, htmlSizeRule)

	ttfbRule := valueobject.NewAuditRule("ttfb", valueobject.CategoryPerformance, valueobject.SeverityHigh)
	ttfbRule.AffectedURL = page.URL
	ttfbMs := page.ResponseTime.Milliseconds()
	if ttfbMs > 600 {
		ttfbRule.Fail(
			fmt.Sprintf("Slow server response time (TTFB: %dms)", ttfbMs),
			"Improve server response time to under 600ms. Consider caching, CDN, and server optimization.",
		)
	} else if ttfbMs > 400 {
		ttfbRule.Warn(
			fmt.Sprintf("Server response time could be improved (TTFB: %dms)", ttfbMs),
			"Consider optimizing server response time for better user experience.",
		)
	} else {
		ttfbRule.Pass(fmt.Sprintf("Server response time is fast (TTFB: %dms)", ttfbMs))
	}
	rules = append(rules, ttfbRule)

	compressionRule := valueobject.NewAuditRule("compression", valueobject.CategoryPerformance, valueobject.SeverityLow)
	compressionRule.AffectedURL = page.URL
	contentEncoding := page.Headers.Get("Content-Encoding")
	vary := page.Headers.Get("Vary")
	transferEncoding := page.Headers.Get("Transfer-Encoding")
	hasCompressionIndicator := contentEncoding != "" ||
		strings.Contains(strings.ToLower(vary), "accept-encoding") ||
		strings.Contains(strings.ToLower(transferEncoding), "gzip") ||
		strings.Contains(strings.ToLower(transferEncoding), "br")
	if hasCompressionIndicator {
		encoding := contentEncoding
		if encoding == "" {
			encoding = "detected via Vary/Transfer-Encoding headers"
		}
		compressionRule.Pass(fmt.Sprintf("Response uses compression (%s)", encoding))
	} else if page.ContentLength > 0 && page.ContentLength < 1024 {
		compressionRule.Pass("Response is small enough that compression is unnecessary")
	} else {
		compressionRule.Warn(
			"Response may not be compressed",
			"Enable GZIP or Brotli compression on your server. Note: CDN edge layers (Vercel, Cloudflare, Netlify) often apply compression transparently — this check may not detect it.",
		)
		compressionRule.WithDetails(fmt.Sprintf("Content-Encoding: %q, Vary: %q, Content-Length: %d bytes", contentEncoding, vary, page.ContentLength))
	}
	rules = append(rules, compressionRule)

	return rules
}
