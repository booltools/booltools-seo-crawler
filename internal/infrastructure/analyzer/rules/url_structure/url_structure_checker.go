package url_structure

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type URLStructureChecker struct{}

func (c *URLStructureChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	parsedURL, err := url.Parse(page.URL)
	if err != nil {
		return rules
	}

	lowercaseRule := valueobject.NewAuditRule("url_lowercase", valueobject.CategoryURLStructure, valueobject.SeverityLow)
	lowercaseRule.AffectedURL = page.URL
	if parsedURL.Path != strings.ToLower(parsedURL.Path) {
		lowercaseRule.Warn(
			"URL contains uppercase characters",
			"Use lowercase URLs to avoid duplicate content issues. Set up redirects from uppercase to lowercase versions.",
		)
	} else {
		lowercaseRule.Pass("URL uses lowercase characters")
	}
	rules = append(rules, lowercaseRule)

	hyphensRule := valueobject.NewAuditRule("url_hyphens", valueobject.CategoryURLStructure, valueobject.SeverityLow)
	hyphensRule.AffectedURL = page.URL
	if strings.Contains(parsedURL.Path, "_") {
		hyphensRule.Warn(
			"URL uses underscores instead of hyphens",
			"Use hyphens (-) instead of underscores (_) in URLs. Google treats hyphens as word separators.",
		)
	} else {
		hyphensRule.Pass("URL uses hyphens for word separation")
	}
	rules = append(rules, hyphensRule)

	lengthRule := valueobject.NewAuditRule("url_length", valueobject.CategoryURLStructure, valueobject.SeverityLow)
	lengthRule.AffectedURL = page.URL
	urlLength := len(page.URL)
	if urlLength > 100 {
		lengthRule.Warn(
			fmt.Sprintf("URL is long (%d characters)", urlLength),
			"Keep URLs under 100 characters for readability and shareability.",
		)
	} else {
		lengthRule.Pass(fmt.Sprintf("URL length is good (%d characters)", urlLength))
	}
	rules = append(rules, lengthRule)

	specialCharsRule := valueobject.NewAuditRule("url_special_chars", valueobject.CategoryURLStructure, valueobject.SeverityLow)
	specialCharsRule.AffectedURL = page.URL
	hasSpecialChars := false
	for _, r := range parsedURL.Path {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '/' && r != '-' && r != '.' {
			hasSpecialChars = true
			break
		}
	}
	if hasSpecialChars {
		specialCharsRule.Warn(
			"URL contains special characters or encoded spaces",
			"Use only alphanumeric characters, hyphens, and forward slashes in URLs.",
		)
	} else {
		specialCharsRule.Pass("URL has clean character usage")
	}
	rules = append(rules, specialCharsRule)

	doubleSlashRule := valueobject.NewAuditRule("url_double_slash", valueobject.CategoryURLStructure, valueobject.SeverityMedium)
	doubleSlashRule.AffectedURL = page.URL
	if strings.Contains(parsedURL.Path, "//") {
		doubleSlashRule.Fail(
			"URL contains double slashes in path",
			"Remove double slashes from the URL path. These can cause duplicate content issues.",
		)
	} else {
		doubleSlashRule.Pass("URL path has no double slashes")
	}
	rules = append(rules, doubleSlashRule)

	paramsRule := valueobject.NewAuditRule("url_parameters", valueobject.CategoryURLStructure, valueobject.SeverityLow)
	paramsRule.AffectedURL = page.URL
	queryParams := parsedURL.Query()
	if len(queryParams) > 3 {
		paramsRule.Warn(
			fmt.Sprintf("URL has %d query parameters", len(queryParams)),
			"Minimize URL parameters. Consider using clean URL paths instead of query strings.",
		)
	} else {
		paramsRule.Pass("URL parameters are minimal")
	}
	rules = append(rules, paramsRule)

	return rules
}
