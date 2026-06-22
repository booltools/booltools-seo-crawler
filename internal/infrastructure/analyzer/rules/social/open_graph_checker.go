package social

import (
	"strings"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type OpenGraphChecker struct{}

func (c *OpenGraphChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	ogTags := map[string]struct {
		key      string
		severity valueobject.Severity
		message  string
	}{
		"og:title":       {"og_title", valueobject.SeverityMedium, "Add og:title for better social media sharing previews."},
		"og:description": {"og_description", valueobject.SeverityMedium, "Add og:description for compelling social media previews."},
		"og:image":       {"og_image", valueobject.SeverityMedium, "Add og:image (1200x630px minimum) for visual social media previews."},
		"og:url":         {"og_url", valueobject.SeverityLow, "Add og:url with the canonical URL of the page."},
		"og:type":        {"og_type", valueobject.SeverityLow, "Add og:type (e.g., 'website', 'article') to classify your content."},
		"og:site_name":   {"og_site_name", valueobject.SeverityLow, "Add og:site_name with your brand name."},
		"og:locale":      {"og_locale", valueobject.SeverityLow, "Add og:locale (e.g., 'en_US') to specify the content language."},
	}

	for property, config := range ogTags {
		content, _ := page.Document.Find(`meta[property="` + property + `"]`).Attr("content")
		content = strings.TrimSpace(content)

		rule := valueobject.NewAuditRule(config.key, valueobject.CategorySocial, config.severity)
		rule.AffectedURL = page.URL

		if content == "" {
			rule.Fail(property+" is missing", config.message)
		} else {
			rule.Pass(property + " is present")
		}
		rules = append(rules, rule)
	}

	return rules
}
