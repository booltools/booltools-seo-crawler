package eeat

import (
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type EEATPageChecker struct{}

func (c *EEATPageChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	authorRule := valueobject.NewAuditRule("eeat_author", valueobject.CategoryEEAT, valueobject.SeverityMedium)
	authorRule.AffectedURL = page.URL
	hasAuthor := false

	authorSelectors := []string{
		`[rel="author"]`,
		`[class*="author"]`,
		`[itemprop="author"]`,
		`meta[name="author"]`,
	}
	for _, selector := range authorSelectors {
		if page.Document.Find(selector).Length() > 0 {
			hasAuthor = true
			break
		}
	}

	if !hasAuthor {
		authorRule.Warn(
			"No author attribution found",
			"Add author information to content pages. Include author name, bio, and credentials for better E-E-A-T signals.",
		)
	} else {
		authorRule.Pass("Author attribution is present")
	}
	rules = append(rules, authorRule)

	copyrightRule := valueobject.NewAuditRule("eeat_copyright", valueobject.CategoryEEAT, valueobject.SeverityLow)
	copyrightRule.AffectedURL = page.URL
	bodyHTML := page.HTML
	hasCopyright := strings.Contains(bodyHTML, "©") ||
		strings.Contains(strings.ToLower(bodyHTML), "copyright")
	if !hasCopyright {
		copyrightRule.Warn(
			"No copyright notice found",
			"Add a copyright notice with the current year to establish content ownership.",
		)
	} else {
		copyrightRule.Pass("Copyright notice is present")
	}
	rules = append(rules, copyrightRule)

	return rules
}
