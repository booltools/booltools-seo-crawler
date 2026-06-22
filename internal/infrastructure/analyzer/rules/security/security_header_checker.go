package security

import (
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type SecurityHeaderChecker struct{}

func (c *SecurityHeaderChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	headers := map[string]struct {
		key            string
		recommendation string
	}{
		"X-Content-Type-Options": {
			key:            "security_xcto",
			recommendation: "Add the header: X-Content-Type-Options: nosniff",
		},
		"X-Frame-Options": {
			key:            "security_xfo",
			recommendation: "Add the header: X-Frame-Options: DENY (or SAMEORIGIN)",
		},
		"Content-Security-Policy": {
			key:            "security_csp",
			recommendation: "Add a Content-Security-Policy header to prevent XSS and injection attacks.",
		},
		"Referrer-Policy": {
			key:            "security_referrer",
			recommendation: "Add the header: Referrer-Policy: strict-origin-when-cross-origin",
		},
		"Permissions-Policy": {
			key:            "security_permissions",
			recommendation: "Add a Permissions-Policy header to control browser feature access.",
		},
	}

	for headerName, config := range headers {
		rule := valueobject.NewAuditRule(config.key, valueobject.CategorySecurity, valueobject.SeverityLow)
		rule.AffectedURL = page.URL
		if page.Headers.Get(headerName) == "" {
			rule.Warn(
				headerName+" header is missing",
				config.recommendation,
			)
		} else {
			rule.Pass(headerName + " header is present")
		}
		rules = append(rules, rule)
	}

	serverRule := valueobject.NewAuditRule("security_server_disclosure", valueobject.CategorySecurity, valueobject.SeverityLow)
	serverRule.AffectedURL = page.URL
	serverHeader := page.Headers.Get("Server")
	if serverHeader != "" {
		serverRule.Warn(
			"Server version is disclosed: "+serverHeader,
			"Remove or obfuscate the Server header to avoid exposing server software information to attackers.",
		)
	} else {
		serverRule.Pass("Server version is not disclosed")
	}
	rules = append(rules, serverRule)

	return rules
}
