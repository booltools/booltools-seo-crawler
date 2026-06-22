package on_page

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type TitleChecker struct{}

func (c *TitleChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	titles := page.Document.Find("title")
	titleCount := titles.Length()
	titleText := strings.TrimSpace(titles.First().Text())

	existsRule := valueobject.NewAuditRule("title_exists", valueobject.CategoryOnPage, valueobject.SeverityCritical)
	existsRule.AffectedURL = page.URL
	if titleCount == 0 || titleText == "" {
		existsRule.Fail("Page is missing a title tag", "Add a descriptive <title> tag between 30-60 characters that includes your primary keyword.")
	} else {
		existsRule.Pass("Title tag is present")
	}
	rules = append(rules, existsRule)

	if titleText == "" {
		return rules
	}

	lengthRule := valueobject.NewAuditRule("title_length", valueobject.CategoryOnPage, valueobject.SeverityMedium)
	lengthRule.AffectedURL = page.URL
	titleLen := len(titleText)
	if titleLen < 30 {
		lengthRule.Fail(
			fmt.Sprintf("Title is too short (%d characters)", titleLen),
			"Expand your title to 30-60 characters for optimal search display. Include your primary keyword and a compelling description.",
		)
		lengthRule.WithDetails(fmt.Sprintf("Current title (%d chars): %s", titleLen, titleText))
	} else if titleLen > 60 {
		lengthRule.Warn(
			fmt.Sprintf("Title may be truncated in search results (%d characters)", titleLen),
			"Shorten your title to under 60 characters to prevent truncation in search engine results pages.",
		)
		lengthRule.WithDetails(fmt.Sprintf("Current title (%d chars): %s", titleLen, titleText))
	} else {
		lengthRule.Pass(fmt.Sprintf("Title length is optimal (%d characters)", titleLen))
	}
	rules = append(rules, lengthRule)

	multipleRule := valueobject.NewAuditRule("title_multiple", valueobject.CategoryOnPage, valueobject.SeverityHigh)
	multipleRule.AffectedURL = page.URL
	if titleCount > 1 {
		var titleTexts []string
		titles.Each(func(_ int, selection *goquery.Selection) {
			text := strings.TrimSpace(selection.Text())
			if text == "" {
				text = "(empty)"
			}
			titleTexts = append(titleTexts, fmt.Sprintf("<title>%s</title>", text))
		})
		multipleRule.Fail(
			fmt.Sprintf("Page has %d title tags", titleCount),
			"Remove duplicate <title> tags. Each page should have exactly one title tag.",
		)
		multipleRule.WithDetails(strings.Join(titleTexts, "\n"))
	} else {
		multipleRule.Pass("Page has exactly one title tag")
	}
	rules = append(rules, multipleRule)

	return rules
}
