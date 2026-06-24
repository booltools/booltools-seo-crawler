package sdk_test

import (
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/sdk"
)

func makeRules() []valueobject.AuditRule {
	return []valueobject.AuditRule{
		{Key: "title_exists", Category: valueobject.CategoryOnPage, Severity: valueobject.SeverityHigh, Result: valueobject.RuleResultFail},
		{Key: "h1_count", Category: valueobject.CategoryOnPage, Severity: valueobject.SeverityHigh, Result: valueobject.RuleResultPass},
		{Key: "og_locale", Category: valueobject.CategorySocial, Severity: valueobject.SeverityLow, Result: valueobject.RuleResultFail},
		{Key: "twitter_site", Category: valueobject.CategorySocial, Severity: valueobject.SeverityLow, Result: valueobject.RuleResultWarning},
		{Key: "mobile_viewport", Category: valueobject.CategoryMobile, Severity: valueobject.SeverityCritical, Result: valueobject.RuleResultPass},
	}
}

func TestFilterRules_NoFilters(t *testing.T) {
	rules := makeRules()
	filtered := sdk.FilterRules(rules, nil, nil)

	if len(filtered) != 5 {
		t.Errorf("expected all 5 rules, got %d", len(filtered))
	}
}

func TestFilterRules_Ignore(t *testing.T) {
	rules := makeRules()
	filtered := sdk.FilterRules(rules, []string{"og_locale", "twitter_site"}, nil)

	if len(filtered) != 3 {
		t.Errorf("expected 3 rules after ignoring 2, got %d", len(filtered))
	}
	for _, rule := range filtered {
		if rule.Key == "og_locale" || rule.Key == "twitter_site" {
			t.Errorf("rule %s should have been filtered out", rule.Key)
		}
	}
}

func TestFilterRules_Only(t *testing.T) {
	rules := makeRules()
	filtered := sdk.FilterRules(rules, nil, []string{"title_exists", "h1_count"})

	if len(filtered) != 2 {
		t.Errorf("expected 2 rules with only filter, got %d", len(filtered))
	}
}

func TestFilterRules_IgnoreAndOnly(t *testing.T) {
	rules := makeRules()
	filtered := sdk.FilterRules(rules, []string{"h1_count"}, []string{"title_exists", "h1_count"})

	if len(filtered) != 1 {
		t.Errorf("expected 1 rule (title_exists), got %d", len(filtered))
	}
	if filtered[0].Key != "title_exists" {
		t.Errorf("expected title_exists, got %s", filtered[0].Key)
	}
}

func TestShouldAnalyzeURL_NoFilters(t *testing.T) {
	if !sdk.ShouldAnalyzeURL("http://localhost:3000/page", nil, nil) {
		t.Error("expected all URLs to pass with no filters")
	}
}

func TestShouldAnalyzeURL_ExcludeExactPath(t *testing.T) {
	excludeURLs := []string{"/admin"}
	if sdk.ShouldAnalyzeURL("http://localhost:3000/admin", excludeURLs, nil) {
		t.Error("expected /admin to be excluded")
	}
	if !sdk.ShouldAnalyzeURL("http://localhost:3000/blog", excludeURLs, nil) {
		t.Error("expected /blog to pass")
	}
}

func TestShouldAnalyzeURL_ExcludeWildcard(t *testing.T) {
	excludeURLs := []string{"/admin/*"}
	if sdk.ShouldAnalyzeURL("http://localhost:3000/admin/settings", excludeURLs, nil) {
		t.Error("expected /admin/settings to be excluded")
	}
	if sdk.ShouldAnalyzeURL("http://localhost:3000/admin/users/123", excludeURLs, nil) {
		t.Error("expected /admin/users/123 to be excluded")
	}
	if !sdk.ShouldAnalyzeURL("http://localhost:3000/blog/post", excludeURLs, nil) {
		t.Error("expected /blog/post to pass")
	}
}

func TestShouldAnalyzeURL_OnlyExactPath(t *testing.T) {
	onlyURLs := []string{"/marketplace"}
	if !sdk.ShouldAnalyzeURL("http://localhost:3000/marketplace", nil, onlyURLs) {
		t.Error("expected /marketplace to be included")
	}
	if sdk.ShouldAnalyzeURL("http://localhost:3000/blog", nil, onlyURLs) {
		t.Error("expected /blog to be excluded")
	}
}

func TestShouldAnalyzeURL_OnlyWildcard(t *testing.T) {
	onlyURLs := []string{"/marketplace/*"}
	if !sdk.ShouldAnalyzeURL("http://localhost:3000/marketplace/shoes", nil, onlyURLs) {
		t.Error("expected /marketplace/shoes to be included")
	}
	if !sdk.ShouldAnalyzeURL("http://localhost:3000/marketplace/shoes/nike", nil, onlyURLs) {
		t.Error("expected /marketplace/shoes/nike to be included")
	}
	if sdk.ShouldAnalyzeURL("http://localhost:3000/blog/post", nil, onlyURLs) {
		t.Error("expected /blog/post to be excluded")
	}
}

func TestShouldAnalyzeURL_MultiplePatterns(t *testing.T) {
	onlyURLs := []string{"/marketplace/*", "/blog/*", "/"}
	if !sdk.ShouldAnalyzeURL("http://example.com/", nil, onlyURLs) {
		t.Error("expected / to be included")
	}
	if !sdk.ShouldAnalyzeURL("http://example.com/marketplace/item", nil, onlyURLs) {
		t.Error("expected /marketplace/item to be included")
	}
	if !sdk.ShouldAnalyzeURL("http://example.com/blog/post-1", nil, onlyURLs) {
		t.Error("expected /blog/post-1 to be included")
	}
	if sdk.ShouldAnalyzeURL("http://example.com/admin/panel", nil, onlyURLs) {
		t.Error("expected /admin/panel to be excluded")
	}
}

func TestShouldAnalyzeURL_TrailingSlashNormalization(t *testing.T) {
	onlyURLs := []string{"/privacy"}
	if !sdk.ShouldAnalyzeURL("http://localhost:3000/privacy/", nil, onlyURLs) {
		t.Error("expected /privacy/ to match /privacy pattern")
	}
}

func TestShouldAnalyzeURL_CaseInsensitive(t *testing.T) {
	onlyURLs := []string{"/Blog/*"}
	if !sdk.ShouldAnalyzeURL("http://localhost:3000/blog/post", nil, onlyURLs) {
		t.Error("expected case-insensitive matching")
	}
}
