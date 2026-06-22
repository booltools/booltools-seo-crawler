package technical

import (
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type PaginationChecker struct{}

func (c *PaginationChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	hasNext := page.Document.Find(`link[rel="next"]`).Length() > 0
	hasPrev := page.Document.Find(`link[rel="prev"]`).Length() > 0
	hasPagination := page.Document.Find(`nav[aria-label*="pagination"], .pagination, [class*="pager"]`).Length() > 0

	if !hasPagination {
		return rules
	}

	paginationRule := valueobject.NewAuditRule("pagination_rel_tags", valueobject.CategoryTechnical, valueobject.SeverityLow)
	paginationRule.AffectedURL = page.URL

	if !hasNext && !hasPrev {
		paginationRule.Warn(
			"Paginated page is missing rel=\"next\"/rel=\"prev\" tags",
			"Add <link rel=\"next\"> and <link rel=\"prev\"> tags to paginated pages to help search engines understand the series relationship.",
		)
	} else {
		paginationRule.Pass("Pagination uses rel=\"next\"/rel=\"prev\" correctly")
	}
	rules = append(rules, paginationRule)

	return rules
}
