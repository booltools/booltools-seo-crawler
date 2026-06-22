package technical

import (
	"fmt"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type CrawlDepthChecker struct{}

func (c *CrawlDepthChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	depthRule := valueobject.NewAuditRule("crawl_depth", valueobject.CategoryTechnical, valueobject.SeverityMedium)
	depthRule.AffectedURL = page.URL
	if page.Depth > 3 {
		depthRule.Warn(
			fmt.Sprintf("Page is %d clicks deep from the homepage", page.Depth),
			"Move important pages closer to the homepage (within 3 clicks). Deep pages are harder for search engines to discover and index.",
		)
	} else {
		depthRule.Pass(fmt.Sprintf("Page is reachable in %d clicks from the homepage", page.Depth))
	}
	rules = append(rules, depthRule)

	return rules
}
