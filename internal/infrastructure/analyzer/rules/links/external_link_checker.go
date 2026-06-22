package links

import (
	"fmt"
	"strings"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type ExternalLinkChecker struct{}

func (c *ExternalLinkChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	externalLinks := page.ExternalLinks
	if len(externalLinks) == 0 {
		return rules
	}

	var missingRelLinks []string
	for _, link := range externalLinks {
		rel := strings.ToLower(link.Rel)
		if !strings.Contains(rel, "noopener") || !strings.Contains(rel, "noreferrer") {
			if link.Target == "_blank" {
				missingRelLinks = append(missingRelLinks, link.URL)
			}
		}
	}

	securityRule := valueobject.NewAuditRule("external_links_rel", valueobject.CategoryLinks, valueobject.SeverityMedium)
	securityRule.AffectedURL = page.URL
	if len(missingRelLinks) > 0 {
		securityRule.Warn(
			fmt.Sprintf("%d external links with target=\"_blank\" are missing rel=\"noopener noreferrer\"", len(missingRelLinks)),
			"Add rel=\"noopener noreferrer\" to all external links that open in a new tab for security.",
		)
		securityRule.WithDetails(formatLinkDetailList(missingRelLinks))
	} else {
		securityRule.Pass("External links have proper rel attributes")
	}
	rules = append(rules, securityRule)

	return rules
}
