package accessibility

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type AccessibilityChecker struct{}

func (c *AccessibilityChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	langRule := valueobject.NewAuditRule("html_lang", valueobject.CategoryAccessibility, valueobject.SeverityHigh)
	langRule.AffectedURL = page.URL
	langAttr, langExists := page.Document.Find("html").Attr("lang")
	if !langExists || strings.TrimSpace(langAttr) == "" {
		langRule.Fail(
			"HTML element is missing the lang attribute",
			"Add a lang attribute to the <html> element: <html lang=\"en\">. This helps screen readers and search engines determine the page language.",
		)
	} else {
		langRule.Pass("HTML lang attribute is present: " + langAttr)
	}
	rules = append(rules, langRule)

	viewportRule := valueobject.NewAuditRule("viewport_meta", valueobject.CategoryAccessibility, valueobject.SeverityHigh)
	viewportRule.AffectedURL = page.URL
	viewport, viewportExists := page.Document.Find(`meta[name="viewport"]`).Attr("content")
	if !viewportExists || viewport == "" {
		viewportRule.Fail(
			"Viewport meta tag is missing",
			"Add <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"> for proper mobile rendering.",
		)
	} else {
		viewportRule.Pass("Viewport meta tag is present")
	}
	rules = append(rules, viewportRule)

	charsetRule := valueobject.NewAuditRule("charset_meta", valueobject.CategoryAccessibility, valueobject.SeverityMedium)
	charsetRule.AffectedURL = page.URL
	charset, charsetExists := page.Document.Find(`meta[charset]`).Attr("charset")
	httpEquivCharset, _ := page.Document.Find(`meta[http-equiv="Content-Type"]`).Attr("content")
	if (!charsetExists || charset == "") && httpEquivCharset == "" {
		charsetRule.Warn(
			"Character encoding is not declared",
			"Add <meta charset=\"UTF-8\"> to the <head> to ensure proper text rendering.",
		)
	} else {
		charsetRule.Pass("Character encoding is declared")
	}
	rules = append(rules, charsetRule)

	var emptyLinkDetails []string
	page.Document.Find("a").Each(func(_ int, selection *goquery.Selection) {
		text := strings.TrimSpace(selection.Text())
		ariaLabel := strings.TrimSpace(selection.AttrOr("aria-label", ""))
		title := strings.TrimSpace(selection.AttrOr("title", ""))
		if text == "" && ariaLabel == "" && title == "" && selection.Find("img[alt]").Length() == 0 {
			href := selection.AttrOr("href", "(no href)")
			emptyLinkDetails = append(emptyLinkDetails, href)
		}
	})

	emptyLinksRule := valueobject.NewAuditRule("empty_links", valueobject.CategoryAccessibility, valueobject.SeverityMedium)
	emptyLinksRule.AffectedURL = page.URL
	if len(emptyLinkDetails) > 0 {
		emptyLinksRule.Warn(
			fmt.Sprintf("%d empty links found (no text, aria-label, or title)", len(emptyLinkDetails)),
			"Add descriptive text, aria-label, or title attributes to all links.",
		)
		emptyLinksRule.WithDetails(formatAccessibilityList(emptyLinkDetails))
	} else {
		emptyLinksRule.Pass("All links have accessible text")
	}
	rules = append(rules, emptyLinksRule)

	ariaRule := valueobject.NewAuditRule("aria_landmarks", valueobject.CategoryAccessibility, valueobject.SeverityLow)
	ariaRule.AffectedURL = page.URL
	hasNav := page.Document.Find("nav, [role='navigation']").Length() > 0
	hasMain := page.Document.Find("main, [role='main']").Length() > 0
	if !hasNav && !hasMain {
		ariaRule.Warn(
			"No ARIA landmarks found (<nav>, <main>, or role attributes)",
			"Use semantic HTML elements (<nav>, <main>, <header>, <footer>) or ARIA roles for better accessibility.",
		)
	} else {
		ariaRule.Pass("ARIA landmarks are present")
	}
	rules = append(rules, ariaRule)

	return rules
}

func formatAccessibilityList(items []string) string {
	return strings.Join(items, "\n")
}
