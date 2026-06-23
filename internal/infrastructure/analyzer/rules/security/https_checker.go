package security

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type HTTPSChecker struct{}

func (c *HTTPSChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	parsedURL, _ := url.Parse(page.URL)

	httpsRule := valueobject.NewAuditRule("uses_https", valueobject.CategorySecurity, valueobject.SeverityCritical)
	httpsRule.AffectedURL = page.URL
	if parsedURL != nil && isLocalhostHost(parsedURL.Hostname()) {
		httpsRule.Pass("HTTPS check skipped for localhost")
	} else if parsedURL != nil && parsedURL.Scheme != "https" {
		httpsRule.Fail(
			"Page is not served over HTTPS",
			"Migrate your site to HTTPS. HTTPS is a confirmed ranking factor and required for user trust.",
		)
	} else {
		httpsRule.Pass("Page is served over HTTPS")
	}
	rules = append(rules, httpsRule)

	mixedContentRule := valueobject.NewAuditRule("mixed_content", valueobject.CategorySecurity, valueobject.SeverityHigh)
	mixedContentRule.AffectedURL = page.URL
	var mixedContentAssets []string

	if parsedURL != nil && parsedURL.Scheme == "https" {
		for _, image := range page.Images {
			if strings.HasPrefix(image.URL, "http://") {
				mixedContentAssets = append(mixedContentAssets, fmt.Sprintf("[img] %s", image.URL))
			}
		}
		for _, script := range page.Scripts {
			if strings.HasPrefix(script.URL, "http://") {
				mixedContentAssets = append(mixedContentAssets, fmt.Sprintf("[js] %s", script.URL))
			}
		}
		for _, stylesheet := range page.Stylesheets {
			if strings.HasPrefix(stylesheet.URL, "http://") {
				mixedContentAssets = append(mixedContentAssets, fmt.Sprintf("[css] %s", stylesheet.URL))
			}
		}
	}

	if len(mixedContentAssets) > 0 {
		mixedContentRule.Fail(
			fmt.Sprintf("Page has %d mixed content resources (HTTP on HTTPS page)", len(mixedContentAssets)),
			"Update all resource URLs to use HTTPS. Mixed content blocks can prevent resources from loading and trigger browser warnings.",
		)
		mixedContentRule.WithDetails(formatMixedContentList(mixedContentAssets))
	} else {
		mixedContentRule.Pass("No mixed content detected")
	}
	rules = append(rules, mixedContentRule)

	hstsRule := valueobject.NewAuditRule("hsts_header", valueobject.CategorySecurity, valueobject.SeverityMedium)
	hstsRule.AffectedURL = page.URL
	if parsedURL != nil && isLocalhostHost(parsedURL.Hostname()) {
		hstsRule.Pass("HSTS check skipped for localhost")
	} else {
		hsts := page.Headers.Get("Strict-Transport-Security")
		if hsts == "" {
			hstsRule.Warn(
				"HSTS header is missing",
				"Add a Strict-Transport-Security header to enforce HTTPS connections: Strict-Transport-Security: max-age=31536000; includeSubDomains",
			)
		} else {
			hstsRule.Pass("HSTS header is present")
		}
	}
	rules = append(rules, hstsRule)

	return rules
}

func isLocalhostHost(hostname string) bool {
	return hostname == "localhost" || hostname == "127.0.0.1" || hostname == "0.0.0.0" || hostname == "::1"
}

func formatMixedContentList(assets []string) string {
	return strings.Join(assets, "\n")
}
