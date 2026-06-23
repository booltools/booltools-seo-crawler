package analyzer_test

import (
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer"
)

func TestAnalyzePage_NoindexSkipsSeORules(t *testing.T) {
	siteAnalyzer := analyzer.NewSiteAnalyzer()
	siteAnalyzer.ExcludeNoindex = true

	html := `<html><head>
		<meta name="robots" content="noindex">
		<title>Auth Page</title>
	</head><body><h1>Login</h1></body></html>`
	page := makePageData(html)
	page.IsNoindex = true

	rules := siteAnalyzer.AnalyzePage(page)

	for _, rule := range rules {
		if !isAllowedNoindexCategory(rule.Category) {
			t.Errorf("expected only technical/performance/security rules for noindex page, got category %s (rule: %s)", rule.Category, rule.Key)
		}
	}

	if len(rules) == 0 {
		t.Error("expected at least some technical rules for noindex page")
	}
}

func TestAnalyzePage_NoindexKeepsTechnicalRules(t *testing.T) {
	siteAnalyzer := analyzer.NewSiteAnalyzer()
	siteAnalyzer.ExcludeNoindex = true

	html := `<html><head>
		<meta name="robots" content="noindex">
		<title>Auth Page</title>
	</head><body><h1>Login</h1></body></html>`
	page := makePageData(html)
	page.IsNoindex = true

	rules := siteAnalyzer.AnalyzePage(page)

	hasTechnical := false
	for _, rule := range rules {
		if rule.Category == valueobject.CategoryTechnical {
			hasTechnical = true
			break
		}
	}

	if !hasTechnical {
		t.Error("expected technical rules to still run on noindex pages")
	}
}

func TestAnalyzePage_NonNoindexRunsAllRules(t *testing.T) {
	siteAnalyzer := analyzer.NewSiteAnalyzer()
	siteAnalyzer.ExcludeNoindex = true

	html := `<html><head>
		<title>Normal Page With Good Title Here</title>
		<meta name="description" content="A good description for the page that is long enough to pass validation checks.">
	</head><body><h1>Welcome</h1><p>Content here.</p></body></html>`
	page := makePageData(html)

	rules := siteAnalyzer.AnalyzePage(page)

	hasOnPage := false
	for _, rule := range rules {
		if rule.Category == valueobject.CategoryOnPage {
			hasOnPage = true
			break
		}
	}

	if !hasOnPage {
		t.Error("expected on_page rules to run on normal (non-noindex) pages")
	}
}

func TestAnalyzePage_ExcludeNoindexDisabledRunsAllRules(t *testing.T) {
	siteAnalyzer := analyzer.NewSiteAnalyzer()
	siteAnalyzer.ExcludeNoindex = false

	html := `<html><head>
		<meta name="robots" content="noindex">
		<title>Auth Page</title>
	</head><body><h1>Login</h1></body></html>`
	page := makePageData(html)
	page.IsNoindex = true

	rules := siteAnalyzer.AnalyzePage(page)

	hasOnPage := false
	for _, rule := range rules {
		if rule.Category == valueobject.CategoryOnPage {
			hasOnPage = true
			break
		}
	}

	if !hasOnPage {
		t.Error("expected on_page rules to run on noindex pages when ExcludeNoindex is false")
	}
}

func isAllowedNoindexCategory(category valueobject.Category) bool {
	return category == valueobject.CategoryTechnical ||
		category == valueobject.CategoryPerformance ||
		category == valueobject.CategorySecurity
}
