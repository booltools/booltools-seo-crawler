package sdk_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer"
	"github.com/booltools/booltools-seo-crawler/internal/sdk"
)

func TestStaticScanner_ValidDirectory(t *testing.T) {
	tempDir := t.TempDir()

	validHTML := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Test Page With Good Title Here</title>
	<meta name="description" content="This is a good meta description that is long enough to pass the length check for SEO purposes.">
	<link rel="canonical" href="https://example.com/test">
</head>
<body>
	<nav>Navigation</nav>
	<main>
		<h1>Welcome to the Test Page</h1>
		<p>This is a paragraph with enough content to avoid the thin content warning from the checker.</p>
		<img src="hero.webp" alt="Hero image showing the product" width="800" height="600" loading="lazy">
		<a href="/about">About Us</a>
	</main>
</body>
</html>`

	if err := os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(validHTML), 0644); err != nil {
		t.Fatalf("failed to write test HTML: %v", err)
	}

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if result.TotalPages != 1 {
		t.Errorf("expected 1 page, got %d", result.TotalPages)
	}

	if result.Mode != "static" {
		t.Errorf("expected mode 'static', got '%s'", result.Mode)
	}

	if len(result.AllRules) == 0 {
		t.Error("expected rules to be generated")
	}

	if len(result.Pages) != 1 {
		t.Errorf("expected 1 page result, got %d", len(result.Pages))
	}
}

func TestStaticScanner_InvalidHTML(t *testing.T) {
	tempDir := t.TempDir()

	invalidHTML := `<!DOCTYPE html>
<html>
<head></head>
<body>
	<p>Minimal page with many missing SEO elements.</p>
</body>
</html>`

	if err := os.WriteFile(filepath.Join(tempDir, "bad.html"), []byte(invalidHTML), 0644); err != nil {
		t.Fatalf("failed to write test HTML: %v", err)
	}

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if len(result.Pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(result.Pages))
	}

	pageResult := result.Pages[0]
	if pageResult.Failures == 0 {
		t.Error("expected failures for a page missing many SEO elements")
	}
}

func TestStaticScanner_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	_, err := scanner.Scan(tempDir, nil, nil)
	if err == nil {
		t.Error("expected error for empty directory")
	}
}

func TestStaticScanner_NonexistentDirectory(t *testing.T) {
	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	_, err := scanner.Scan("/nonexistent/path/to/nowhere", nil, nil)
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestStaticScanner_NestedHTMLFiles(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "pages")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Content</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)
	os.WriteFile(filepath.Join(subDir, "about.html"), []byte(html), 0644)
	os.WriteFile(filepath.Join(subDir, "contact.htm"), []byte(html), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if result.TotalPages != 3 {
		t.Errorf("expected 3 pages (including .htm), got %d", result.TotalPages)
	}
}

func TestStaticScanner_MultiplePages(t *testing.T) {
	tempDir := t.TempDir()

	goodHTML := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Good Page Title for SEO</title>
	<meta name="description" content="A well-crafted meta description for the SEO checker to validate properly.">
</head>
<body><nav>Nav</nav><main><h1>Title</h1><p>Content here.</p></main></body>
</html>`

	badHTML := `<html><body><p>Bad page</p></body></html>`

	os.WriteFile(filepath.Join(tempDir, "good.html"), []byte(goodHTML), 0644)
	os.WriteFile(filepath.Join(tempDir, "bad.html"), []byte(badHTML), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if result.TotalPages != 2 {
		t.Errorf("expected 2 pages, got %d", result.TotalPages)
	}

	if len(result.AllRules) == 0 {
		t.Error("expected rules to be populated from both pages")
	}
}

func TestStaticScanner_DetectsRobotsTxt(t *testing.T) {
	tempDir := t.TempDir()

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Hello</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)

	robotsTxt := `User-agent: *
Allow: /

Sitemap: https://example.com/sitemap.xml`
	os.WriteFile(filepath.Join(tempDir, "robots.txt"), []byte(robotsTxt), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	foundRobotsTxtExists := false
	robotsTxtPassed := false
	for _, rule := range result.AllRules {
		if rule.Key == "robots_txt_exists" {
			foundRobotsTxtExists = true
			if rule.Result == "pass" {
				robotsTxtPassed = true
			}
		}
	}

	if !foundRobotsTxtExists {
		t.Error("expected robots_txt_exists rule to be evaluated in static mode")
	}
	if !robotsTxtPassed {
		t.Error("expected robots_txt_exists to pass when robots.txt file is present")
	}
}

func TestStaticScanner_DetectsMissingRobotsTxt(t *testing.T) {
	tempDir := t.TempDir()

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Hello</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	foundRobotsTxtExists := false
	robotsTxtFailed := false
	for _, rule := range result.AllRules {
		if rule.Key == "robots_txt_exists" {
			foundRobotsTxtExists = true
			if rule.Result == "fail" {
				robotsTxtFailed = true
			}
		}
	}

	if !foundRobotsTxtExists {
		t.Error("expected robots_txt_exists rule to be evaluated in static mode")
	}
	if !robotsTxtFailed {
		t.Error("expected robots_txt_exists to fail when robots.txt is missing")
	}
}

func TestStaticScanner_DetectsSitemapXml(t *testing.T) {
	tempDir := t.TempDir()

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Hello</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)

	sitemapXml := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
    <lastmod>2026-06-22</lastmod>
  </url>
</urlset>`
	os.WriteFile(filepath.Join(tempDir, "sitemap.xml"), []byte(sitemapXml), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	foundSitemapExists := false
	sitemapPassed := false
	for _, rule := range result.AllRules {
		if rule.Key == "sitemap_exists" {
			foundSitemapExists = true
			if rule.Result == "pass" {
				sitemapPassed = true
			}
		}
	}

	if !foundSitemapExists {
		t.Error("expected sitemap_exists rule to be evaluated in static mode")
	}
	if !sitemapPassed {
		t.Error("expected sitemap_exists to pass when sitemap.xml file is present")
	}
}

func TestStaticScanner_DetectsLlmsTxt(t *testing.T) {
	tempDir := t.TempDir()

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Hello</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)

	llmsTxt := `# My Site

> A brief description of the site.

## Docs

- [Getting Started](https://example.com/docs/getting-started): Setup guide
`
	os.WriteFile(filepath.Join(tempDir, "llms.txt"), []byte(llmsTxt), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	foundLlmsTxtExists := false
	llmsTxtPassed := false
	for _, rule := range result.AllRules {
		if rule.Key == "geo_llms_txt_exists" {
			foundLlmsTxtExists = true
			if rule.Result == "pass" {
				llmsTxtPassed = true
			}
		}
	}

	if !foundLlmsTxtExists {
		t.Error("expected geo_llms_txt_exists rule to be evaluated in static mode")
	}
	if !llmsTxtPassed {
		t.Error("expected geo_llms_txt_exists to pass when llms.txt file is present")
	}
}

func TestStaticScanner_DetectsMissingLlmsTxt(t *testing.T) {
	tempDir := t.TempDir()

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Hello</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	foundLlmsTxtExists := false
	llmsTxtFailed := false
	for _, rule := range result.AllRules {
		if rule.Key == "geo_llms_txt_exists" {
			foundLlmsTxtExists = true
			if rule.Result == "fail" {
				llmsTxtFailed = true
			}
		}
	}

	if !foundLlmsTxtExists {
		t.Error("expected geo_llms_txt_exists rule to be evaluated in static mode")
	}
	if !llmsTxtFailed {
		t.Error("expected geo_llms_txt_exists to fail when llms.txt is missing")
	}
}

func TestStaticScanner_AICrawlerAccessChecked(t *testing.T) {
	tempDir := t.TempDir()

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Hello</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)

	robotsTxt := `User-agent: *
Allow: /

User-agent: OAI-SearchBot
Allow: /

User-agent: GPTBot
Disallow: /

Sitemap: https://example.com/sitemap.xml`
	os.WriteFile(filepath.Join(tempDir, "robots.txt"), []byte(robotsTxt), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	foundAICrawlerRule := false
	for _, rule := range result.AllRules {
		if strings.HasPrefix(rule.Key, "geo_crawler_") {
			foundAICrawlerRule = true
			break
		}
	}

	if !foundAICrawlerRule {
		t.Error("expected AI crawler access rules (geo_crawler_*) to be evaluated in static mode when robots.txt is present")
	}
}

func TestStaticScanner_NetworkRulesExcluded(t *testing.T) {
	tempDir := t.TempDir()

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Hello</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	networkRules := map[string]bool{
		"broken_external_links": true,
		"broken_internal_links": true,
		"broken_scripts":        true,
		"broken_stylesheets":    true,
		"broken_images":         true,
		"http_status_ok":        true,
		"uses_https":            true,
		"ttfb":                  true,
		"compression":           true,
		"sitemap_broken_urls":   true,
	}

	for _, rule := range result.AllRules {
		if networkRules[rule.Key] {
			t.Errorf("network-dependent rule %q should be excluded in static mode", rule.Key)
		}
	}
}

func TestStaticScanner_RobotsTxtSyntaxChecked(t *testing.T) {
	tempDir := t.TempDir()

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Hello</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)

	robotsTxt := `User-agent: *
Allow: /
InvalidDirective: something
Sitemap: https://example.com/sitemap.xml`
	os.WriteFile(filepath.Join(tempDir, "robots.txt"), []byte(robotsTxt), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	foundSyntaxRule := false
	syntaxFailed := false
	for _, rule := range result.AllRules {
		if rule.Key == "robots_txt_syntax" {
			foundSyntaxRule = true
			if rule.Result == "fail" {
				syntaxFailed = true
			}
		}
	}

	if !foundSyntaxRule {
		t.Error("expected robots_txt_syntax rule to be evaluated")
	}
	if !syntaxFailed {
		t.Error("expected robots_txt_syntax to fail when robots.txt has invalid directives")
	}
}

func TestStaticScanner_SitemapValidXml(t *testing.T) {
	tempDir := t.TempDir()

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Hello</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)

	invalidSitemap := `not valid xml at all`
	os.WriteFile(filepath.Join(tempDir, "sitemap.xml"), []byte(invalidSitemap), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir, nil, nil)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	foundValidXml := false
	validXmlFailed := false
	for _, rule := range result.AllRules {
		if rule.Key == "sitemap_valid_xml" {
			foundValidXml = true
			if rule.Result == "fail" {
				validXmlFailed = true
			}
		}
	}

	if !foundValidXml {
		t.Error("expected sitemap_valid_xml rule to be evaluated")
	}
	if !validXmlFailed {
		t.Error("expected sitemap_valid_xml to fail for invalid XML content")
	}
}
