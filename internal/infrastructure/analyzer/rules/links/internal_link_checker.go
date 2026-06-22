package links

import (
	"fmt"
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type InternalLinkChecker struct{}

func (c *InternalLinkChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	internalLinks := page.InternalLinks

	hasLinksRule := valueobject.NewAuditRule("internal_links_present", valueobject.CategoryLinks, valueobject.SeverityMedium)
	hasLinksRule.AffectedURL = page.URL
	if len(internalLinks) == 0 {
		hasLinksRule.Fail(
			"Page has no internal links",
			"Add internal links to help search engines discover and understand your site structure.",
		)
	} else {
		hasLinksRule.Pass(fmt.Sprintf("Page has %d internal links", len(internalLinks)))
	}
	rules = append(rules, hasLinksRule)

	excessiveRule := valueobject.NewAuditRule("internal_links_count", valueobject.CategoryLinks, valueobject.SeverityLow)
	excessiveRule.AffectedURL = page.URL
	if len(internalLinks) > 100 {
		uniqueTargets := make(map[string]int)
		for _, link := range internalLinks {
			uniqueTargets[link.URL]++
		}
		var linkSummary []string
		linkSummary = append(linkSummary, fmt.Sprintf("Total: %d links to %d unique URLs", len(internalLinks), len(uniqueTargets)))
		for url, count := range uniqueTargets {
			if count > 1 {
				linkSummary = append(linkSummary, fmt.Sprintf("%s (x%d)", url, count))
			}
		}
		excessiveRule.Warn(
			fmt.Sprintf("Page has excessive internal links (%d)", len(internalLinks)),
			"Reduce the number of internal links per page to under 100 to focus link equity on the most important pages.",
		)
		excessiveRule.WithDetails(strings.Join(linkSummary, "\n"))
	} else {
		excessiveRule.Pass(fmt.Sprintf("Internal link count is reasonable (%d)", len(internalLinks)))
	}
	rules = append(rules, excessiveRule)

	var genericAnchorDetails []string
	genericPhrases := []string{"click here", "read more", "learn more", "here", "link", "this"}
	for _, link := range internalLinks {
		anchor := strings.ToLower(strings.TrimSpace(link.AnchorText))
		for _, phrase := range genericPhrases {
			if anchor == phrase {
				genericAnchorDetails = append(genericAnchorDetails, fmt.Sprintf("%s (text: \"%s\")", link.URL, link.AnchorText))
				break
			}
		}
	}

	anchorRule := valueobject.NewAuditRule("internal_links_anchor_text", valueobject.CategoryLinks, valueobject.SeverityLow)
	anchorRule.AffectedURL = page.URL
	if len(genericAnchorDetails) > 0 {
		anchorRule.Warn(
			fmt.Sprintf("%d internal links use generic anchor text", len(genericAnchorDetails)),
			"Replace generic anchor text (\"click here\", \"read more\") with descriptive text that indicates the linked page's content.",
		)
		anchorRule.WithDetails(formatLinkDetailList(genericAnchorDetails))
	} else {
		anchorRule.Pass("Internal links use descriptive anchor text")
	}
	rules = append(rules, anchorRule)

	return rules
}

func formatLinkDetailList(items []string) string {
	return strings.Join(items, "\n")
}
