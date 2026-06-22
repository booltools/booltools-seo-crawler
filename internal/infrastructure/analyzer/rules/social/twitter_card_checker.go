package social

import (
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type TwitterCardChecker struct{}

func (c *TwitterCardChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	twitterTags := map[string]struct {
		key      string
		severity valueobject.Severity
		message  string
	}{
		"twitter:card":        {"twitter_card", valueobject.SeverityMedium, "Add twitter:card (use 'summary_large_image' for best visibility)."},
		"twitter:title":       {"twitter_title", valueobject.SeverityLow, "Add twitter:title for Twitter sharing previews."},
		"twitter:description": {"twitter_description", valueobject.SeverityLow, "Add twitter:description for Twitter sharing previews."},
		"twitter:image":       {"twitter_image", valueobject.SeverityLow, "Add twitter:image for visual Twitter previews."},
		"twitter:site":        {"twitter_site", valueobject.SeverityLow, "Add twitter:site with your @username for attribution."},
	}

	for name, config := range twitterTags {
		content, _ := page.Document.Find(`meta[name="` + name + `"]`).Attr("content")
		content = strings.TrimSpace(content)

		rule := valueobject.NewAuditRule(config.key, valueobject.CategorySocial, config.severity)
		rule.AffectedURL = page.URL

		if content == "" {
			rule.Warn(name+" is missing", config.message)
		} else {
			rule.Pass(name + " is present")
		}
		rules = append(rules, rule)
	}

	return rules
}
