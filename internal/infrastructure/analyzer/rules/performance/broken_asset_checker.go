package performance

import (
	"fmt"
	"strings"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type BrokenAssetChecker struct{}

func (c *BrokenAssetChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	uniqueScripts := make(map[string]string)
	uniqueStylesheets := make(map[string]string)
	uniqueImages := make(map[string]string)

	for _, page := range result.Pages {
		for _, script := range page.Scripts {
			if script.URL != "" {
				if _, exists := uniqueScripts[script.URL]; !exists {
					uniqueScripts[script.URL] = page.URL
				}
			}
		}
		for _, stylesheet := range page.Stylesheets {
			if stylesheet.URL != "" {
				if _, exists := uniqueStylesheets[stylesheet.URL]; !exists {
					uniqueStylesheets[stylesheet.URL] = page.URL
				}
			}
		}
		for _, image := range page.Images {
			if image.URL != "" && !strings.HasPrefix(image.URL, "data:") {
				if _, exists := uniqueImages[image.URL]; !exists {
					uniqueImages[image.URL] = page.URL
				}
			}
		}
	}

	cache := result.URLStatusCache

	scriptResults := cache.CheckConcurrent(uniqueScripts, 30, 8)
	cssResults := cache.CheckConcurrent(uniqueStylesheets, 20, 8)
	imageResults := cache.CheckConcurrent(uniqueImages, 40, 8)

	brokenScripts, brokenScriptDetails := countBrokenAssets(scriptResults, uniqueScripts)
	brokenCSS, brokenCSSDetails := countBrokenAssets(cssResults, uniqueStylesheets)
	brokenImages, brokenImageDetails := countBrokenAssets(imageResults, uniqueImages)

	scriptRule := valueobject.NewAuditRule("broken_scripts", valueobject.CategoryPerformance, valueobject.SeverityHigh)
	if brokenScripts > 0 {
		scriptRule.Fail(
			fmt.Sprintf("%d JavaScript files return errors (checked %d)", brokenScripts, len(uniqueScripts)),
			"Fix or remove broken script references. Broken JS files cause functionality failures and degrade user experience.",
		)
		scriptRule.WithDetails(strings.Join(brokenScriptDetails, "\n"))
	} else if len(uniqueScripts) > 0 {
		scriptRule.Pass(fmt.Sprintf("All %d script files are reachable", len(uniqueScripts)))
	} else {
		scriptRule.Pass("No external scripts to check")
	}
	rules = append(rules, scriptRule)

	cssRule := valueobject.NewAuditRule("broken_stylesheets", valueobject.CategoryPerformance, valueobject.SeverityHigh)
	if brokenCSS > 0 {
		cssRule.Fail(
			fmt.Sprintf("%d CSS files return errors (checked %d)", brokenCSS, len(uniqueStylesheets)),
			"Fix or remove broken stylesheet references. Missing CSS causes unstyled content and poor user experience.",
		)
		cssRule.WithDetails(strings.Join(brokenCSSDetails, "\n"))
	} else if len(uniqueStylesheets) > 0 {
		cssRule.Pass(fmt.Sprintf("All %d stylesheets are reachable", len(uniqueStylesheets)))
	} else {
		cssRule.Pass("No external stylesheets to check")
	}
	rules = append(rules, cssRule)

	imageRule := valueobject.NewAuditRule("broken_images", valueobject.CategoryPerformance, valueobject.SeverityMedium)
	if brokenImages > 0 {
		imageRule.Warn(
			fmt.Sprintf("%d images return errors (checked %d)", brokenImages, len(uniqueImages)),
			"Fix or remove broken image references. Missing images hurt user experience and waste bandwidth.",
		)
		imageRule.WithDetails(strings.Join(brokenImageDetails, "\n"))
	} else if len(uniqueImages) > 0 {
		imageRule.Pass(fmt.Sprintf("All %d checked images are reachable", min(len(uniqueImages), 40)))
	} else {
		imageRule.Pass("No images to check")
	}
	rules = append(rules, imageRule)

	return rules
}

func countBrokenAssets(results map[string]crawler.URLStatusResult, sourceMap map[string]string) (int, []string) {
	brokenCount := 0
	var details []string

	for targetURL, statusResult := range results {
		sourceURL := sourceMap[targetURL]
		if statusResult.Error != nil {
			brokenCount++
			details = append(details, fmt.Sprintf("%s (from %s, connection error)", targetURL, sourceURL))
		} else if statusResult.StatusCode >= 400 {
			brokenCount++
			details = append(details, fmt.Sprintf("%s (from %s, HTTP %d)", targetURL, sourceURL, statusResult.StatusCode))
		}
	}

	return brokenCount, details
}
