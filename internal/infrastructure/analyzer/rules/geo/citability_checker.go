package geo

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type CitabilityChecker struct{}

var numberPattern = regexp.MustCompile(`\d+[\.,]?\d*\s*(%|percent|million|billion|thousand|users|customers|downloads)`)
var questionPattern = regexp.MustCompile(`(?i)(what|how|why|when|where|which|who|can|does|is|are|should)\s`)

func (c *CitabilityChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	totalPages := 0
	pagesWithStats := 0
	pagesWithFAQ := 0
	pagesWithTables := 0
	pagesWithLists := 0
	pagesWithQuestionHeadings := 0

	for _, page := range result.Pages {
		totalPages++

		if numberPattern.MatchString(page.BodyText) {
			pagesWithStats++
		}

		if page.Document.Find("details, .faq, #faq, [itemtype*='FAQPage']").Length() > 0 {
			pagesWithFAQ++
		}

		if page.Document.Find("table").Length() > 0 {
			pagesWithTables++
		}

		if page.Document.Find("ol").Length() > 0 {
			pagesWithLists++
		}

		page.Document.Find("h2, h3").Each(func(_ int, selection *goquery.Selection) {
			headingText := strings.TrimSpace(selection.Text())
			if questionPattern.MatchString(headingText) || strings.HasSuffix(headingText, "?") {
				pagesWithQuestionHeadings++
			}
		})
	}

	if totalPages == 0 {
		return rules
	}

	statsRule := valueobject.NewAuditRule("geo_citability_statistics", valueobject.CategoryGEO, valueobject.SeverityMedium)
	statsRatio := float64(pagesWithStats) / float64(totalPages) * 100
	if statsRatio < 30 {
		statsRule.Fail(
			fmt.Sprintf("Only %.0f%% of pages contain specific data points or statistics", statsRatio),
			"Add specific numbers, percentages, and data points to your content. AI models prefer citing pages with verifiable, quantifiable information.",
		)
	} else {
		statsRule.Pass(fmt.Sprintf("%.0f%% of pages contain data points and statistics", statsRatio))
	}
	rules = append(rules, statsRule)

	faqRule := valueobject.NewAuditRule("geo_citability_faq", valueobject.CategoryGEO, valueobject.SeverityLow)
	if pagesWithFAQ == 0 {
		faqRule.Warn(
			"No FAQ-style content found on any page",
			"Add FAQ sections with clear question-answer pairs. AI models frequently extract and cite Q&A-formatted content.",
		)
	} else {
		faqRule.Pass(fmt.Sprintf("%d pages have FAQ-style content", pagesWithFAQ))
	}
	rules = append(rules, faqRule)

	tablesRule := valueobject.NewAuditRule("geo_citability_tables", valueobject.CategoryGEO, valueobject.SeverityLow)
	if pagesWithTables == 0 {
		tablesRule.Warn(
			"No comparison tables found",
			"Add comparison tables for \"vs\" content and feature comparisons. Structured data in tables is easily extracted by AI models.",
		)
	} else {
		tablesRule.Pass(fmt.Sprintf("%d pages contain data tables", pagesWithTables))
	}
	rules = append(rules, tablesRule)

	listsRule := valueobject.NewAuditRule("geo_citability_lists", valueobject.CategoryGEO, valueobject.SeverityLow)
	if pagesWithLists == 0 {
		listsRule.Warn(
			"No numbered/ordered lists found",
			"Use numbered lists for step-by-step guides and how-to content. AI models prefer structured, extractable content.",
		)
	} else {
		listsRule.Pass(fmt.Sprintf("%d pages contain ordered lists", pagesWithLists))
	}
	rules = append(rules, listsRule)

	questionRule := valueobject.NewAuditRule("geo_citability_question_headings", valueobject.CategoryGEO, valueobject.SeverityMedium)
	if pagesWithQuestionHeadings == 0 {
		questionRule.Warn(
			"No question-format headings found",
			"Use question keywords in H2/H3 headings (\"What is...\", \"How to...\"). These mirror how users query AI search engines.",
		)
	} else {
		questionRule.Pass(fmt.Sprintf("%d pages use question-format headings", pagesWithQuestionHeadings))
	}
	rules = append(rules, questionRule)

	return rules
}
