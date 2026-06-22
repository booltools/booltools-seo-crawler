package internationalization

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type HreflangChecker struct{}

func (c *HreflangChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	hreflangTags := page.Document.Find(`link[rel="alternate"][hreflang]`)
	hreflangCount := hreflangTags.Length()

	if hreflangCount == 0 {
		return rules
	}

	validRule := valueobject.NewAuditRule("hreflang_valid", valueobject.CategoryInternationalization, valueobject.SeverityMedium)
	validRule.AffectedURL = page.URL
	invalidCount := 0
	hasXDefault := false

	hreflangTags.Each(func(_ int, selection *goquery.Selection) {
		hreflang, _ := selection.Attr("hreflang")
		hreflang = strings.TrimSpace(hreflang)

		if hreflang == "x-default" {
			hasXDefault = true
			return
		}

		if !isValidLanguageCode(hreflang) {
			invalidCount++
		}
	})

	if invalidCount > 0 {
		validRule.Fail(
			fmt.Sprintf("%d hreflang tags have invalid language codes", invalidCount),
			"Use valid ISO 639-1 language codes (e.g., 'en', 'pt-BR') in hreflang attributes.",
		)
	} else {
		validRule.Pass("All hreflang language codes are valid")
	}
	rules = append(rules, validRule)

	xDefaultRule := valueobject.NewAuditRule("hreflang_x_default", valueobject.CategoryInternationalization, valueobject.SeverityLow)
	xDefaultRule.AffectedURL = page.URL
	if !hasXDefault {
		xDefaultRule.Warn(
			"Missing x-default hreflang tag",
			"Add an x-default hreflang tag to specify the default or fallback version of the page.",
		)
	} else {
		xDefaultRule.Pass("x-default hreflang tag is present")
	}
	rules = append(rules, xDefaultRule)

	return rules
}

func isValidLanguageCode(code string) bool {
	parts := strings.Split(code, "-")
	if len(parts) < 1 || len(parts) > 3 {
		return false
	}
	lang := parts[0]
	if len(lang) != 2 && len(lang) != 3 {
		return false
	}
	for _, r := range lang {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}
