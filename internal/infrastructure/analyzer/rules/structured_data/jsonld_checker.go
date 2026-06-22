package structured_data

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type JsonLdChecker struct{}

func (c *JsonLdChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	jsonLdScripts := page.Document.Find(`script[type="application/ld+json"]`)
	jsonLdCount := jsonLdScripts.Length()

	existsRule := valueobject.NewAuditRule("jsonld_exists", valueobject.CategoryStructuredData, valueobject.SeverityMedium)
	existsRule.AffectedURL = page.URL
	if jsonLdCount == 0 {
		existsRule.Fail(
			"No structured data (JSON-LD) found",
			"Add JSON-LD structured data to help search engines understand your content. Use appropriate Schema.org types (Article, Organization, Product, etc.).",
		)
		rules = append(rules, existsRule)
		return rules
	}
	existsRule.Pass(fmt.Sprintf("Found %d JSON-LD blocks", jsonLdCount))
	rules = append(rules, existsRule)

	validJSONCount := 0
	schemaTypes := make([]string, 0)

	jsonLdScripts.Each(func(_ int, selection *goquery.Selection) {
		content := strings.TrimSpace(selection.Text())
		if content == "" {
			return
		}

		var raw interface{}
		if err := json.Unmarshal([]byte(content), &raw); err != nil {
			return
		}
		validJSONCount++

		extractTypes(raw, &schemaTypes)
	})

	validRule := valueobject.NewAuditRule("jsonld_valid", valueobject.CategoryStructuredData, valueobject.SeverityHigh)
	validRule.AffectedURL = page.URL
	if validJSONCount < jsonLdCount {
		validRule.Fail(
			fmt.Sprintf("%d of %d JSON-LD blocks contain invalid JSON", jsonLdCount-validJSONCount, jsonLdCount),
			"Fix the invalid JSON in your structured data blocks. Validate at https://validator.schema.org/",
		)
	} else {
		validRule.Pass("All JSON-LD blocks contain valid JSON")
	}
	rules = append(rules, validRule)

	breadcrumbRule := valueobject.NewAuditRule("jsonld_breadcrumb", valueobject.CategoryStructuredData, valueobject.SeverityLow)
	breadcrumbRule.AffectedURL = page.URL
	hasBreadcrumb := containsType(schemaTypes, "BreadcrumbList")
	if !hasBreadcrumb {
		breadcrumbRule.Warn(
			"BreadcrumbList schema is missing",
			"Add BreadcrumbList structured data to help search engines understand your site hierarchy and display breadcrumbs in results.",
		)
	} else {
		breadcrumbRule.Pass("BreadcrumbList schema is present")
	}
	rules = append(rules, breadcrumbRule)

	return rules
}

func extractTypes(raw interface{}, schemaTypes *[]string) {
	switch value := raw.(type) {
	case map[string]interface{}:
		if schemaType, exists := value["@type"]; exists {
			if typeStr, ok := schemaType.(string); ok {
				*schemaTypes = append(*schemaTypes, typeStr)
			}
		}
	case []interface{}:
		for _, item := range value {
			extractTypes(item, schemaTypes)
		}
	}
}

func containsType(types []string, target string) bool {
	for _, schemaType := range types {
		if strings.EqualFold(schemaType, target) {
			return true
		}
	}
	return false
}
