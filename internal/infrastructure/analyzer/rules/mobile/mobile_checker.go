package mobile

import (
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type MobileChecker struct{}

func (c *MobileChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	viewportRule := valueobject.NewAuditRule("mobile_viewport", valueobject.CategoryMobile, valueobject.SeverityCritical)
	viewportRule.AffectedURL = page.URL
	viewport, viewportExists := page.Document.Find(`meta[name="viewport"]`).Attr("content")
	if !viewportExists || viewport == "" {
		viewportRule.Fail(
			"Viewport meta tag is missing",
			"Add <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"> for responsive design.",
		)
	} else {
		viewportRule.Pass("Viewport meta tag is present")
	}
	rules = append(rules, viewportRule)

	if viewport != "" {
		viewportConfigRule := valueobject.NewAuditRule("mobile_viewport_config", valueobject.CategoryMobile, valueobject.SeverityMedium)
		viewportConfigRule.AffectedURL = page.URL
		hasDeviceWidth := strings.Contains(viewport, "width=device-width")
		hasInitialScale := strings.Contains(viewport, "initial-scale=1")

		if !hasDeviceWidth || !hasInitialScale {
			viewportConfigRule.Warn(
				"Viewport is not optimally configured",
				"Set viewport to: width=device-width, initial-scale=1",
			)
		} else {
			viewportConfigRule.Pass("Viewport is properly configured for mobile")
		}
		rules = append(rules, viewportConfigRule)
	}

	return rules
}
