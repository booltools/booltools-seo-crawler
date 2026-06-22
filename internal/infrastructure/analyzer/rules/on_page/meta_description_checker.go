package on_page

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type MetaDescriptionChecker struct{}

func (c *MetaDescriptionChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	descriptions := page.Document.Find(`meta[name="description"]`)
	descCount := descriptions.Length()
	descContent, _ := descriptions.First().Attr("content")
	descContent = strings.TrimSpace(descContent)

	existsRule := valueobject.NewAuditRule("meta_description_exists", valueobject.CategoryOnPage, valueobject.SeverityHigh)
	existsRule.AffectedURL = page.URL
	if descCount == 0 || descContent == "" {
		existsRule.Fail(
			"Page is missing a meta description",
			"Add a meta description between 120-160 characters that summarizes the page content and includes target keywords.",
		)
	} else {
		existsRule.Pass("Meta description is present")
	}
	rules = append(rules, existsRule)

	if descContent == "" {
		return rules
	}

	lengthRule := valueobject.NewAuditRule("meta_description_length", valueobject.CategoryOnPage, valueobject.SeverityMedium)
	lengthRule.AffectedURL = page.URL
	descLen := len(descContent)
	if descLen < 120 {
		lengthRule.Warn(
			fmt.Sprintf("Meta description is short (%d characters)", descLen),
			"Expand your meta description to 120-160 characters to maximize SERP real estate.",
		)
		lengthRule.WithDetails(fmt.Sprintf("Current value (%d chars): %s", descLen, descContent))
	} else if descLen > 160 {
		lengthRule.Warn(
			fmt.Sprintf("Meta description may be truncated (%d characters)", descLen),
			"Shorten your meta description to under 160 characters to prevent truncation.",
		)
		lengthRule.WithDetails(fmt.Sprintf("Current value (%d chars): %s", descLen, descContent))
	} else {
		lengthRule.Pass(fmt.Sprintf("Meta description length is optimal (%d characters)", descLen))
	}
	rules = append(rules, lengthRule)

	multipleRule := valueobject.NewAuditRule("meta_description_multiple", valueobject.CategoryOnPage, valueobject.SeverityMedium)
	multipleRule.AffectedURL = page.URL
	if descCount > 1 {
		var descTexts []string
		descriptions.Each(func(index int, selection *goquery.Selection) {
			content, _ := selection.Attr("content")
			content = strings.TrimSpace(content)
			if content == "" {
				content = "(empty)"
			}
			descTexts = append(descTexts, fmt.Sprintf("#%d: %s", index+1, content))
		})
		multipleRule.Fail(
			fmt.Sprintf("Page has %d meta description tags", descCount),
			"Remove duplicate meta description tags. Each page should have exactly one.",
		)
		multipleRule.WithDetails(strings.Join(descTexts, "\n"))
	} else {
		multipleRule.Pass("Page has exactly one meta description")
	}
	rules = append(rules, multipleRule)

	return rules
}
