package duplicate_content

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type DuplicateContentChecker struct{}

func (c *DuplicateContentChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	contentHashes := make(map[string][]string)
	titleMap := make(map[string][]string)
	descriptionMap := make(map[string][]string)

	for _, page := range result.Pages {
		bodyText := strings.TrimSpace(page.BodyText)
		if bodyText != "" {
			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(bodyText)))
			contentHashes[hash] = append(contentHashes[hash], page.URL)
		}

		title := strings.TrimSpace(page.Document.Find("title").Text())
		if title != "" {
			titleMap[title] = append(titleMap[title], page.URL)
		}

		desc, _ := page.Document.Find(`meta[name="description"]`).Attr("content")
		desc = strings.TrimSpace(desc)
		if desc != "" {
			descriptionMap[desc] = append(descriptionMap[desc], page.URL)
		}
	}

	duplicateContentCount := 0
	for _, urls := range contentHashes {
		if len(urls) > 1 {
			duplicateContentCount++
		}
	}

	contentRule := valueobject.NewAuditRule("duplicate_content_body", valueobject.CategoryDuplicateContent, valueobject.SeverityHigh)
	if duplicateContentCount > 0 {
		var duplicateDetails []string
		groupIndex := 1
		for _, urls := range contentHashes {
			if len(urls) > 1 {
				duplicateDetails = append(duplicateDetails, fmt.Sprintf("Group %d (%d pages):", groupIndex, len(urls)))
				for _, url := range urls {
					duplicateDetails = append(duplicateDetails, "  "+url)
				}
				groupIndex++
			}
		}
		contentRule.Fail(
			fmt.Sprintf("%d groups of pages have identical content", duplicateContentCount),
			"Consolidate duplicate pages using canonical tags, 301 redirects, or by differentiating the content.",
		)
		contentRule.WithDetails(strings.Join(duplicateDetails, "\n"))
	} else {
		contentRule.Pass("No exact duplicate content detected")
	}
	rules = append(rules, contentRule)

	duplicateTitleCount := 0
	for _, urls := range titleMap {
		if len(urls) > 1 {
			duplicateTitleCount++
		}
	}

	titleRule := valueobject.NewAuditRule("duplicate_titles", valueobject.CategoryDuplicateContent, valueobject.SeverityMedium)
	if duplicateTitleCount > 0 {
		var titleDetails []string
		for title, urls := range titleMap {
			if len(urls) > 1 {
				displayTitle := title
				if len(displayTitle) > 100 {
					displayTitle = displayTitle[:100] + "..."
				}
				titleDetails = append(titleDetails, fmt.Sprintf("\"%s\" (%d pages):", displayTitle, len(urls)))
				for _, url := range urls {
					titleDetails = append(titleDetails, "  "+url)
				}
			}
		}
		titleRule.Fail(
			fmt.Sprintf("%d groups of pages share identical title tags", duplicateTitleCount),
			"Write unique, descriptive title tags for each page targeting different keywords.",
		)
		titleRule.WithDetails(strings.Join(titleDetails, "\n"))
	} else {
		titleRule.Pass("All pages have unique title tags")
	}
	rules = append(rules, titleRule)

	duplicateDescCount := 0
	for _, urls := range descriptionMap {
		if len(urls) > 1 {
			duplicateDescCount++
		}
	}

	descRule := valueobject.NewAuditRule("duplicate_descriptions", valueobject.CategoryDuplicateContent, valueobject.SeverityMedium)
	if duplicateDescCount > 0 {
		var descDetails []string
		for desc, urls := range descriptionMap {
			if len(urls) > 1 {
				displayDesc := desc
				if len(displayDesc) > 100 {
					displayDesc = displayDesc[:100] + "..."
				}
				descDetails = append(descDetails, fmt.Sprintf("\"%s\" (%d pages):", displayDesc, len(urls)))
				for _, url := range urls {
					descDetails = append(descDetails, "  "+url)
				}
			}
		}
		descRule.Fail(
			fmt.Sprintf("%d groups of pages share identical meta descriptions", duplicateDescCount),
			"Write unique meta descriptions for each page to improve click-through rates.",
		)
		descRule.WithDetails(strings.Join(descDetails, "\n"))
	} else {
		descRule.Pass("All pages have unique meta descriptions")
	}
	rules = append(rules, descRule)

	return rules
}
