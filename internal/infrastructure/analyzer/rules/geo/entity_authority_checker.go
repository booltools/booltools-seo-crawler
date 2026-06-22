package geo

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type EntityAuthorityChecker struct{}

func (c *EntityAuthorityChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	orgSchemaComplete := false
	hasSocialLinks := false

	for _, page := range result.Pages {
		page.Document.Find(`script[type="application/ld+json"]`).Each(func(_ int, selection *goquery.Selection) {
			content := strings.TrimSpace(selection.Text())
			var parsed map[string]interface{}
			if json.Unmarshal([]byte(content), &parsed) != nil {
				return
			}

			schemaType, _ := parsed["@type"].(string)
			if schemaType == "Organization" {
				requiredFields := []string{"name", "description", "url"}
				allPresent := true
				for _, field := range requiredFields {
					if _, exists := parsed[field]; !exists {
						allPresent = false
						break
					}
				}
				if allPresent {
					orgSchemaComplete = true
				}

				if _, exists := parsed["sameAs"]; exists {
					hasSocialLinks = true
				}
			}
		})
	}

	orgRule := valueobject.NewAuditRule("geo_entity_org_schema", valueobject.CategoryGEO, valueobject.SeverityHigh)
	if !orgSchemaComplete {
		orgRule.Fail(
			"Organization schema is missing or incomplete",
			"Add complete Organization schema to your homepage with name, description, URL, address, phone, email, and social profiles. This is the foundation for AI entity recognition.",
		)
	} else {
		orgRule.Pass("Organization schema is complete")
	}
	rules = append(rules, orgRule)

	socialRule := valueobject.NewAuditRule("geo_entity_social", valueobject.CategoryGEO, valueobject.SeverityMedium)
	if !hasSocialLinks {
		socialRule.Warn(
			"No social profile links found in Organization schema",
			"Add sameAs property to Organization schema with links to your LinkedIn, Twitter, GitHub, and other social profiles.",
		)
	} else {
		socialRule.Pass("Social profile links present in Organization schema")
	}
	rules = append(rules, socialRule)

	return rules
}
