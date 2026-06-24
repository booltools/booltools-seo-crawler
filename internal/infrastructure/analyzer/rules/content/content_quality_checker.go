package content

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
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
	minWords := 300
	thinWords := 50
	if isLegalPage(page.URL) {
		minWords = 100
		thinWords = 30
	}
	if wordCount < thinWords {
		wordCountRule.Fail(
			fmt.Sprintf("Page has very thin content (%d words)", wordCount),
			"Add more substantive content. Note: client-rendered SPAs may show low word counts because content renders via JavaScript.",
		)
	} else if wordCount < minWords {
		wordCountRule.Warn(
			fmt.Sprintf("Page content may be too thin (%d words)", wordCount),
			fmt.Sprintf("Consider expanding the content to at least %d words for better search visibility. Note: client-rendered SPAs may show low word counts because content renders via JavaScript.", minWords),
		)
	} else {
		wordCountRule.Pass(fmt.Sprintf("Content length is adequate (%d words)", wordCount))
	}
	rules = append(rules, wordCountRule)

	ratioRule := valueobject.NewAuditRule("content_text_html_ratio", valueobject.CategoryContent, valueobject.SeverityLow)
	ratioRule.AffectedURL = page.URL
	if htmlLength > 0 {
		ratio := float64(textLength) / float64(htmlLength) * 100
		if ratio < 5 && page.IsDevMode {
			ratioRule.Pass(fmt.Sprintf("Low text-to-HTML ratio (%.1f%%) — dev mode detected, HMR/debug scripts inflate HTML size", ratio))
		} else if ratio < 5 {
			ratioRule.Warn(
				fmt.Sprintf("Low text-to-HTML ratio (%.1f%%)", ratio),
				"Increase the amount of visible text content relative to HTML code. Note: client-rendered SPAs (React, Vue, Angular) will show low ratios because content renders via JavaScript.",
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

var legalPagePatterns = []string{
	"/privacy", "/terms", "/tos", "/legal",
	"/cookie-policy", "/cookies", "/disclaimer",
	"/gdpr", "/imprint", "/impressum",
}

func isLegalPage(pageURL string) bool {
	lowered := strings.ToLower(pageURL)
	for _, pattern := range legalPagePatterns {
		if strings.Contains(lowered, pattern) {
			return true
		}
	}
	return false
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
