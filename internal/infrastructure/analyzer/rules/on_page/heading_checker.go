package on_page

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type HeadingChecker struct{}

func (c *HeadingChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	h1Elements := page.Document.Find("h1")
	h1Count := h1Elements.Length()
	h1Text := strings.TrimSpace(h1Elements.First().Text())

	h1CountRule := valueobject.NewAuditRule("h1_count", valueobject.CategoryOnPage, valueobject.SeverityHigh)
	h1CountRule.AffectedURL = page.URL
	if h1Count == 0 {
		h1CountRule.Fail("Page is missing an H1 heading", "Add exactly one H1 heading that describes the main topic of the page and includes your primary keyword.")
	} else if h1Count > 1 {
		var h1Texts []string
		h1Elements.Each(func(index int, selection *goquery.Selection) {
			text := strings.TrimSpace(selection.Text())
			if text == "" {
				text = "(empty)"
			}
			h1Texts = append(h1Texts, fmt.Sprintf("H1 #%d: %s", index+1, text))
		})
		h1CountRule.Warn(
			fmt.Sprintf("Page has %d H1 headings", h1Count),
			"Use exactly one H1 heading per page. Convert extra H1 tags to H2 or lower.",
		)
		h1CountRule.WithDetails(strings.Join(h1Texts, "\n"))
	} else {
		h1CountRule.Pass("Page has exactly one H1 heading")
	}
	rules = append(rules, h1CountRule)

	h1EmptyRule := valueobject.NewAuditRule("h1_not_empty", valueobject.CategoryOnPage, valueobject.SeverityHigh)
	h1EmptyRule.AffectedURL = page.URL
	if h1Count > 0 && h1Text == "" {
		h1EmptyRule.Fail("H1 heading is empty", "Add descriptive text to your H1 heading that includes the primary keyword.")
	} else if h1Count > 0 {
		h1EmptyRule.Pass("H1 heading has content")
	} else {
		h1EmptyRule.Skip("No H1 found")
	}
	rules = append(rules, h1EmptyRule)

	hierarchyRule := valueobject.NewAuditRule("heading_hierarchy", valueobject.CategoryOnPage, valueobject.SeverityMedium)
	hierarchyRule.AffectedURL = page.URL
	headingSequence := getHeadingSequence(page.Document)
	hierarchyValid := checkHeadingHierarchy(page.Document)
	if !hierarchyValid {
		hierarchyRule.Warn(
			"Heading hierarchy is not sequential",
			"Ensure headings follow a logical order (H1 > H2 > H3). Do not skip heading levels (e.g., H1 directly to H3).",
		)
		hierarchyRule.WithDetails(strings.Join(headingSequence, "\n"))
	} else {
		hierarchyRule.Pass("Heading hierarchy is properly sequential")
	}
	rules = append(rules, hierarchyRule)

	emptyHeadingsRule := valueobject.NewAuditRule("heading_not_empty", valueobject.CategoryOnPage, valueobject.SeverityLow)
	emptyHeadingsRule.AffectedURL = page.URL
	emptyHeadingCount := 0
	var emptyHeadingTags []string
	page.Document.Find("h1, h2, h3, h4, h5, h6").Each(func(_ int, selection *goquery.Selection) {
		if strings.TrimSpace(selection.Text()) == "" {
			emptyHeadingCount++
			tagName := goquery.NodeName(selection)
			emptyHeadingTags = append(emptyHeadingTags, fmt.Sprintf("<%s> (empty)", tagName))
		}
	})
	if emptyHeadingCount > 0 {
		emptyHeadingsRule.Warn(
			fmt.Sprintf("Found %d empty heading tags", emptyHeadingCount),
			"Remove or fill in empty heading tags. Empty headings provide no value to users or search engines.",
		)
		emptyHeadingsRule.WithDetails(strings.Join(emptyHeadingTags, "\n"))
	} else {
		emptyHeadingsRule.Pass("All heading tags have content")
	}
	rules = append(rules, emptyHeadingsRule)

	return rules
}

func getHeadingSequence(document *goquery.Document) []string {
	var sequence []string
	document.Find("h1, h2, h3, h4, h5, h6").Each(func(_ int, selection *goquery.Selection) {
		tagName := strings.ToUpper(goquery.NodeName(selection))
		text := strings.TrimSpace(selection.Text())
		if len(text) > 80 {
			text = text[:80] + "..."
		}
		if text == "" {
			text = "(empty)"
		}
		sequence = append(sequence, fmt.Sprintf("%s: %s", tagName, text))
	})
	return sequence
}

func checkHeadingHierarchy(document *goquery.Document) bool {
	previousLevel := 0
	valid := true

	document.Find("h1, h2, h3, h4, h5, h6").Each(func(_ int, selection *goquery.Selection) {
		if !valid {
			return
		}

		tagName := goquery.NodeName(selection)
		level := int(tagName[1] - '0')

		if previousLevel > 0 && level > previousLevel+1 {
			valid = false
		}
		previousLevel = level
	})

	return valid
}
