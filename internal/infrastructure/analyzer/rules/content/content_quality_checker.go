package content

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type ContentQualityChecker struct{}

func (c *ContentQualityChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	bodyText := page.BodyText
	wordCount := countWords(bodyText)
	htmlLength := len(page.HTML)
	textLength := len(bodyText)

	wordCountRule := valueobject.NewAuditRule("content_word_count", valueobject.CategoryContent, valueobject.SeverityMedium)
	wordCountRule.AffectedURL = page.URL
	if wordCount < 100 {
		wordCountRule.Fail(
			fmt.Sprintf("Page has very thin content (%d words)", wordCount),
			"Add more substantive content. Pages with under 100 words provide little value to users or search engines. Aim for at least 300 words.",
		)
	} else if wordCount < 300 {
		wordCountRule.Warn(
			fmt.Sprintf("Page content may be too thin (%d words)", wordCount),
			"Consider expanding the content to at least 300 words for better search visibility.",
		)
	} else {
		wordCountRule.Pass(fmt.Sprintf("Content length is adequate (%d words)", wordCount))
	}
	rules = append(rules, wordCountRule)

	ratioRule := valueobject.NewAuditRule("content_text_html_ratio", valueobject.CategoryContent, valueobject.SeverityMedium)
	ratioRule.AffectedURL = page.URL
	if htmlLength > 0 {
		ratio := float64(textLength) / float64(htmlLength) * 100
		if ratio < 10 {
			ratioRule.Warn(
				fmt.Sprintf("Low text-to-HTML ratio (%.1f%%)", ratio),
				"Increase the amount of visible text content relative to HTML code. Reduce unnecessary markup, inline styles, and scripts.",
			)
			ratioRule.WithDetails(fmt.Sprintf("Text content: %d bytes, HTML size: %d bytes, Ratio: %.1f%%", textLength, htmlLength, ratio))
		} else {
			ratioRule.Pass(fmt.Sprintf("Text-to-HTML ratio is healthy (%.1f%%)", ratio))
		}
	} else {
		ratioRule.Skip("Could not calculate text-to-HTML ratio")
	}
	rules = append(rules, ratioRule)

	return rules
}

func countWords(text string) int {
	words := 0
	inWord := false
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !inWord {
				words++
				inWord = true
			}
		} else {
			inWord = false
		}
	}
	return words
}

func CountWordsInString(text string) int {
	return len(strings.Fields(text))
}
