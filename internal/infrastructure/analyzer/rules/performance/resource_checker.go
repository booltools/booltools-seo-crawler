package performance

import (
	"fmt"
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type ResourceChecker struct{}

func (c *ResourceChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	jsCountRule := valueobject.NewAuditRule("js_file_count", valueobject.CategoryPerformance, valueobject.SeverityMedium)
	jsCountRule.AffectedURL = page.URL
	jsCount := len(page.Scripts)
	if jsCount > 10 {
		jsCountRule.Warn(
			fmt.Sprintf("Page loads %d JavaScript files", jsCount),
			"Reduce the number of JavaScript files by bundling them. Fewer HTTP requests improve load time.",
		)
		jsCountRule.WithDetails(formatResourceURLs(scriptURLs(page.Scripts)))
	} else {
		jsCountRule.Pass(fmt.Sprintf("JavaScript file count is reasonable (%d)", jsCount))
	}
	rules = append(rules, jsCountRule)

	cssCountRule := valueobject.NewAuditRule("css_file_count", valueobject.CategoryPerformance, valueobject.SeverityMedium)
	cssCountRule.AffectedURL = page.URL
	cssCount := len(page.Stylesheets)
	if cssCount > 5 {
		cssCountRule.Warn(
			fmt.Sprintf("Page loads %d CSS files", cssCount),
			"Reduce the number of CSS files by combining them. Fewer HTTP requests improve load time.",
		)
		cssCountRule.WithDetails(formatResourceURLs(stylesheetURLs(page.Stylesheets)))
	} else {
		cssCountRule.Pass(fmt.Sprintf("CSS file count is reasonable (%d)", cssCount))
	}
	rules = append(rules, cssCountRule)

	var renderBlockingScripts []string
	for _, script := range page.Scripts {
		if script.Location == "head" && !script.IsAsync && !script.IsDefer {
			renderBlockingScripts = append(renderBlockingScripts, script.URL)
		}
	}

	renderBlockingRule := valueobject.NewAuditRule("render_blocking", valueobject.CategoryPerformance, valueobject.SeverityHigh)
	renderBlockingRule.AffectedURL = page.URL
	if len(renderBlockingScripts) > 0 {
		renderBlockingRule.Fail(
			fmt.Sprintf("%d render-blocking scripts found in <head>", len(renderBlockingScripts)),
			"Add 'async' or 'defer' attributes to scripts in the <head>, or move them to the end of <body>.",
		)
		renderBlockingRule.WithDetails(formatResourceURLs(renderBlockingScripts))
	} else {
		renderBlockingRule.Pass("No render-blocking scripts in <head>")
	}
	rules = append(rules, renderBlockingRule)

	totalRequests := jsCount + cssCount + len(page.Images)
	requestRule := valueobject.NewAuditRule("total_requests", valueobject.CategoryPerformance, valueobject.SeverityMedium)
	requestRule.AffectedURL = page.URL
	if totalRequests > 50 {
		requestRule.Warn(
			fmt.Sprintf("High number of HTTP requests (%d)", totalRequests),
			"Reduce the total number of HTTP requests by combining files, using sprites, and lazy-loading resources.",
		)
		var requestDetails []string
		requestDetails = append(requestDetails, fmt.Sprintf("Total: %d (JS: %d, CSS: %d, Images: %d)", totalRequests, jsCount, cssCount, len(page.Images)))
		jsURLs := scriptURLs(page.Scripts)
		if len(jsURLs) > 0 {
			requestDetails = append(requestDetails, fmt.Sprintf("\nJavaScript files (%d):", len(jsURLs)))
			requestDetails = append(requestDetails, jsURLs...)
		}
		cssURLs := stylesheetURLs(page.Stylesheets)
		if len(cssURLs) > 0 {
			requestDetails = append(requestDetails, fmt.Sprintf("\nCSS files (%d):", len(cssURLs)))
			requestDetails = append(requestDetails, cssURLs...)
		}
		if len(page.Images) > 0 {
			requestDetails = append(requestDetails, fmt.Sprintf("\nImages (%d):", len(page.Images)))
			for _, image := range page.Images {
				if image.URL != "" {
					requestDetails = append(requestDetails, image.URL)
				}
			}
		}
		requestRule.WithDetails(strings.Join(requestDetails, "\n"))
	} else {
		requestRule.Pass(fmt.Sprintf("Total HTTP requests are reasonable (%d)", totalRequests))
	}
	rules = append(rules, requestRule)

	cacheRule := valueobject.NewAuditRule("cache_headers", valueobject.CategoryPerformance, valueobject.SeverityLow)
	cacheRule.AffectedURL = page.URL
	cacheControl := page.Headers.Get("Cache-Control")
	if cacheControl == "" {
		cacheRule.Warn(
			"No Cache-Control header found",
			"Add Cache-Control headers to enable browser caching and reduce repeated downloads.",
		)
	} else {
		cacheRule.Pass("Cache-Control header is present")
	}
	rules = append(rules, cacheRule)

	return rules
}

func scriptURLs(scripts []crawler.ResourceData) []string {
	urls := make([]string, 0, len(scripts))
	for _, script := range scripts {
		if script.URL != "" {
			urls = append(urls, script.URL)
		}
	}
	return urls
}

func stylesheetURLs(stylesheets []crawler.ResourceData) []string {
	urls := make([]string, 0, len(stylesheets))
	for _, stylesheet := range stylesheets {
		if stylesheet.URL != "" {
			urls = append(urls, stylesheet.URL)
		}
	}
	return urls
}

func formatResourceURLs(urls []string) string {
	return strings.Join(urls, "\n")
}
