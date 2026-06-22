package links

import (
	"fmt"
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type BrokenLinkChecker struct{}

func (c *BrokenLinkChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	uniqueInternalLinks := make(map[string]string)
	uniqueExternalLinks := make(map[string]string)

	for _, page := range result.Pages {
		for _, link := range page.InternalLinks {
			if _, exists := uniqueInternalLinks[link.URL]; !exists {
				uniqueInternalLinks[link.URL] = page.URL
			}
		}
		for _, link := range page.ExternalLinks {
			if isNonHTTPScheme(link.URL) {
				continue
			}
			if _, exists := uniqueExternalLinks[link.URL]; !exists {
				uniqueExternalLinks[link.URL] = page.URL
			}
		}
	}

	cache := result.URLStatusCache
	brokenInternal, brokenInternalDetails := checkLinksWithCache(cache, uniqueInternalLinks, 50, 3)
	brokenExternal, brokenExternalDetails := checkLinksWithCache(cache, uniqueExternalLinks, 30, 3)

	internalRule := valueobject.NewAuditRule("broken_internal_links", valueobject.CategoryLinks, valueobject.SeverityHigh)
	if brokenInternal > 0 {
		internalRule.Fail(
			fmt.Sprintf("%d broken internal links detected", brokenInternal),
			"Fix or remove broken internal links. They waste crawl budget and create poor user experience.",
		)
		if len(brokenInternalDetails) > 0 {
			internalRule.WithDetails(strings.Join(brokenInternalDetails, "\n"))
		}
	} else {
		internalRule.Pass("No broken internal links detected")
	}
	rules = append(rules, internalRule)

	externalRule := valueobject.NewAuditRule("broken_external_links", valueobject.CategoryLinks, valueobject.SeverityMedium)
	if brokenExternal > 0 {
		externalRule.Warn(
			fmt.Sprintf("%d broken external links detected", brokenExternal),
			"Fix or remove broken external links. They negatively impact user trust and may harm SEO.",
		)
		if len(brokenExternalDetails) > 0 {
			externalRule.WithDetails(strings.Join(brokenExternalDetails, "\n"))
		}
	} else {
		externalRule.Pass("No broken external links detected")
	}
	rules = append(rules, externalRule)

	return rules
}

func checkLinksWithCache(cache *crawler.URLStatusCache, links map[string]string, maxCheck int, concurrency int) (int, []string) {
	statusResults := cache.CheckConcurrent(links, maxCheck, concurrency)

	brokenCount := 0
	var details []string

	for targetURL, statusResult := range statusResults {
		sourceURL := links[targetURL]
		if statusResult.Error != nil {
			brokenCount++
			details = append(details, fmt.Sprintf("%s (from %s, connection error)", targetURL, sourceURL))
		} else if statusResult.StatusCode >= 400 {
			brokenCount++
			details = append(details, fmt.Sprintf("%s (from %s, HTTP %d)", targetURL, sourceURL, statusResult.StatusCode))
		}
	}

	return brokenCount, details
}

func isNonHTTPScheme(targetURL string) bool {
	nonHTTPPrefixes := []string{"mailto:", "tel:", "javascript:", "data:", "ftp:", "file:"}
	for _, prefix := range nonHTTPPrefixes {
		if strings.HasPrefix(targetURL, prefix) {
			return true
		}
	}
	return false
}
