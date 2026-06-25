package geo

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type AIFriendlyChecker struct{}

func (c *AIFriendlyChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	if len(result.Pages) == 0 {
		return rules
	}

	contentPages := 0
	pagesWithDescriptiveHeadings := 0
	pagesWithFreshnessSignals := 0
	totalPages := len(result.Pages)
	pagesWithSemanticHTML := 0

	for _, page := range result.Pages {
		isContent := !crawler.IsNonEditorialPage(page.URL) && !page.IsNoindex

		if isContent {
			contentPages++

			hasDescriptiveH2 := false
			page.Document.Find("h2").Each(func(_ int, selection *goquery.Selection) {
				text := strings.TrimSpace(selection.Text())
				if len(text) > 10 {
					hasDescriptiveH2 = true
				}
			})
			if hasDescriptiveH2 {
				pagesWithDescriptiveHeadings++
			}

			loweredHTML := strings.ToLower(page.HTML)
			hasDateModified := strings.Contains(loweredHTML, "datemodified") || strings.Contains(loweredHTML, "date_modified")
			hasDatePublished := strings.Contains(loweredHTML, "datepublished") || strings.Contains(loweredHTML, "date_published")
			hasLastUpdated := strings.Contains(loweredHTML, "last updated") || strings.Contains(loweredHTML, "updated on") || strings.Contains(loweredHTML, "modified:")
			hasPublishedTime := strings.Contains(loweredHTML, `"article:published_time"`) || strings.Contains(loweredHTML, `"article:modified_time"`)
			if hasDateModified || hasDatePublished || hasLastUpdated || hasPublishedTime {
				pagesWithFreshnessSignals++
			}
		}

		hasArticle := page.Document.Find("article").Length() > 0
		hasSection := page.Document.Find("section").Length() > 0
		hasHeader := page.Document.Find("header").Length() > 0
		hasFooter := page.Document.Find("footer").Length() > 0
		semanticCount := 0
		if hasArticle {
			semanticCount++
		}
		if hasSection {
			semanticCount++
		}
		if hasHeader {
			semanticCount++
		}
		if hasFooter {
			semanticCount++
		}
		if semanticCount >= 2 {
			pagesWithSemanticHTML++
		}
	}

	if contentPages == 0 {
		contentPages = 1
	}

	headingsRule := valueobject.NewAuditRule("geo_ai_descriptive_headings", valueobject.CategoryGEO, valueobject.SeverityMedium)
	headingsRatio := float64(pagesWithDescriptiveHeadings) / float64(contentPages) * 100
	if headingsRatio < 50 {
		headingsRule.Warn(
			fmt.Sprintf("Only %.0f%% of pages have descriptive H2 headings", headingsRatio),
			"Use descriptive H2/H3 headings that mirror user queries. AI models use headings to identify relevant content sections.",
		)
	} else {
		headingsRule.Pass(fmt.Sprintf("%.0f%% of pages have descriptive headings", headingsRatio))
	}
	rules = append(rules, headingsRule)

	freshnessRule := valueobject.NewAuditRule("geo_ai_freshness", valueobject.CategoryGEO, valueobject.SeverityMedium)
	freshnessRatio := float64(pagesWithFreshnessSignals) / float64(contentPages) * 100
	if freshnessRatio < 20 {
		freshnessRule.Warn(
			fmt.Sprintf("Only %.0f%% of pages have content freshness signals", freshnessRatio),
			"Add dateModified schema and visible \"Last Updated\" dates. AI models prefer citing current, recently-updated content.",
		)
	} else {
		freshnessRule.Pass(fmt.Sprintf("%.0f%% of pages have freshness signals", freshnessRatio))
	}
	rules = append(rules, freshnessRule)

	semanticRule := valueobject.NewAuditRule("geo_ai_semantic_html", valueobject.CategoryGEO, valueobject.SeverityMedium)
	semanticRatio := float64(pagesWithSemanticHTML) / float64(totalPages) * 100
	if semanticRatio < 50 {
		semanticRule.Warn(
			fmt.Sprintf("Only %.0f%% of pages use semantic HTML elements", semanticRatio),
			"Use <article>, <section>, <header>, <footer> elements for clean semantic structure. AI crawlers extract content more reliably from semantic HTML.",
		)
	} else {
		semanticRule.Pass(fmt.Sprintf("%.0f%% of pages use semantic HTML", semanticRatio))
	}
	rules = append(rules, semanticRule)

	return rules
}
