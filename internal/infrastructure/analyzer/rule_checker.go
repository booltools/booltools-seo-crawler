package analyzer

import (
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type PageRuleChecker interface {
	Check(page crawler.PageData) []valueobject.AuditRule
}

type SiteRuleChecker interface {
	Check(result crawler.CrawlResult) []valueobject.AuditRule
}
