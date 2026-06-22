package analyzer_test

import (
	"testing"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/technical"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

func TestRobotsTxtChecker_Missing(t *testing.T) {
	result := crawler.CrawlResult{RobotsTxt: ""}
	checker := &technical.RobotsTxtChecker{}
	rules := checker.Check(result)

	existsRule := findRule(rules, "robots_txt_exists")
	if existsRule == nil || existsRule.Result != valueobject.RuleResultFail {
		t.Error("expected robots_txt_exists to fail when missing")
	}
}

func TestRobotsTxtChecker_Valid(t *testing.T) {
	robotsTxt := `User-agent: *
Disallow: /admin/
Allow: /
Sitemap: https://example.com/sitemap.xml`

	result := crawler.CrawlResult{RobotsTxt: robotsTxt}
	checker := &technical.RobotsTxtChecker{}
	rules := checker.Check(result)

	existsRule := findRule(rules, "robots_txt_exists")
	if existsRule == nil || existsRule.Result != valueobject.RuleResultPass {
		t.Error("expected robots_txt_exists to pass")
	}

	sitemapRule := findRule(rules, "robots_txt_sitemap")
	if sitemapRule == nil || sitemapRule.Result != valueobject.RuleResultPass {
		t.Error("expected robots_txt_sitemap to pass")
	}

	syntaxRule := findRule(rules, "robots_txt_syntax")
	if syntaxRule == nil || syntaxRule.Result != valueobject.RuleResultPass {
		t.Error("expected robots_txt_syntax to pass")
	}
}

func TestRobotsTxtChecker_NoSitemap(t *testing.T) {
	robotsTxt := `User-agent: *
Disallow: /admin/`

	result := crawler.CrawlResult{RobotsTxt: robotsTxt}
	checker := &technical.RobotsTxtChecker{}
	rules := checker.Check(result)

	sitemapRule := findRule(rules, "robots_txt_sitemap")
	if sitemapRule == nil || sitemapRule.Result != valueobject.RuleResultFail {
		t.Error("expected robots_txt_sitemap to fail when no sitemap reference")
	}
}

func TestSitemapChecker_Missing(t *testing.T) {
	result := crawler.CrawlResult{SitemapXML: ""}
	checker := &technical.SitemapChecker{}
	rules := checker.Check(result)

	existsRule := findRule(rules, "sitemap_exists")
	if existsRule == nil || existsRule.Result != valueobject.RuleResultFail {
		t.Error("expected sitemap_exists to fail when missing")
	}
}

func TestSitemapChecker_ValidXML(t *testing.T) {
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
    <lastmod>2024-01-01</lastmod>
  </url>
</urlset>`

	result := crawler.CrawlResult{SitemapXML: sitemap}
	checker := &technical.SitemapChecker{}
	rules := checker.Check(result)

	existsRule := findRule(rules, "sitemap_exists")
	if existsRule == nil || existsRule.Result != valueobject.RuleResultPass {
		t.Error("expected sitemap_exists to pass")
	}

	validRule := findRule(rules, "sitemap_valid_xml")
	if validRule == nil || validRule.Result != valueobject.RuleResultPass {
		t.Error("expected sitemap_valid_xml to pass")
	}
}

func TestCanonicalChecker_Missing(t *testing.T) {
	page := makePageData(`<html><head><title>Test</title></head><body></body></html>`)
	checker := &technical.CanonicalChecker{}
	rules := checker.Check(page)

	existsRule := findRule(rules, "canonical_exists")
	if existsRule == nil || existsRule.Result != valueobject.RuleResultFail {
		t.Error("expected canonical_exists to fail when missing")
	}
}

func TestCanonicalChecker_Present(t *testing.T) {
	page := makePageData(`<html><head><link rel="canonical" href="https://example.com/test"></head><body></body></html>`)
	checker := &technical.CanonicalChecker{}
	rules := checker.Check(page)

	existsRule := findRule(rules, "canonical_exists")
	if existsRule == nil || existsRule.Result != valueobject.RuleResultPass {
		t.Error("expected canonical_exists to pass")
	}
}

func TestHTTPStatusChecker_200(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.StatusCode = 200
	checker := &technical.HTTPStatusChecker{}
	rules := checker.Check(page)

	statusRule := findRule(rules, "http_status_ok")
	if statusRule == nil || statusRule.Result != valueobject.RuleResultPass {
		t.Error("expected http_status_ok to pass for 200")
	}
}

func TestHTTPStatusChecker_404(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.StatusCode = 404
	checker := &technical.HTTPStatusChecker{}
	rules := checker.Check(page)

	statusRule := findRule(rules, "http_status_ok")
	if statusRule == nil || statusRule.Result != valueobject.RuleResultFail {
		t.Error("expected http_status_ok to fail for 404")
	}
}
