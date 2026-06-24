package sdk

import (
	"net/url"
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
)

func FilterRules(rules []valueobject.AuditRule, ignore []string, only []string) []valueobject.AuditRule {
	if len(ignore) == 0 && len(only) == 0 {
		return rules
	}

	ignoreSet := toSet(ignore)
	onlySet := toSet(only)

	filtered := make([]valueobject.AuditRule, 0, len(rules))
	for _, rule := range rules {
		if len(onlySet) > 0 {
			if _, included := onlySet[rule.Key]; !included {
				continue
			}
		}

		if _, excluded := ignoreSet[rule.Key]; excluded {
			continue
		}

		filtered = append(filtered, rule)
	}

	return filtered
}

func ShouldAnalyzeURL(pageURL string, excludeURLs []string, onlyURLs []string) bool {
	pagePath := extractPath(pageURL)

	if len(onlyURLs) > 0 {
		for _, pattern := range onlyURLs {
			if matchURLPattern(pagePath, pattern) {
				return true
			}
		}
		return false
	}

	for _, pattern := range excludeURLs {
		if matchURLPattern(pagePath, pattern) {
			return false
		}
	}

	return true
}

func extractPath(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	path := parsed.Path
	if path == "" {
		return "/"
	}
	return path
}

func matchURLPattern(pagePath string, pattern string) bool {
	pagePath = strings.ToLower(strings.TrimSuffix(pagePath, "/"))
	pattern = strings.ToLower(strings.TrimSuffix(pattern, "/"))

	if pagePath == "" {
		pagePath = "/"
	}
	if pattern == "" {
		pattern = "/"
	}

	if !strings.Contains(pattern, "*") {
		return pagePath == pattern
	}

	return globMatch(pagePath, pattern)
}

func globMatch(text string, pattern string) bool {
	textIndex := 0
	patternIndex := 0
	lastStarText := -1
	lastStarPattern := -1

	for textIndex < len(text) {
		if patternIndex < len(pattern) && pattern[patternIndex] == '*' {
			lastStarPattern = patternIndex
			lastStarText = textIndex
			patternIndex++
			continue
		}

		if patternIndex < len(pattern) && text[textIndex] == pattern[patternIndex] {
			textIndex++
			patternIndex++
			continue
		}

		if lastStarPattern >= 0 {
			patternIndex = lastStarPattern + 1
			lastStarText++
			textIndex = lastStarText
			continue
		}

		return false
	}

	for patternIndex < len(pattern) && pattern[patternIndex] == '*' {
		patternIndex++
	}

	return patternIndex == len(pattern)
}

func toSet(items []string) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}
