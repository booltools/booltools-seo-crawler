package analyzer_test

import (
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer/rules/structured_data"
)

func TestJsonLdChecker_NoStructuredData(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	checker := &structured_data.JsonLdChecker{}
	rules := checker.Check(page)

	existsRule := findRule(rules, "jsonld_exists")
	if existsRule == nil || existsRule.Result != valueobject.RuleResultFail {
		t.Error("expected jsonld_exists to fail when no JSON-LD found")
	}
}

func TestJsonLdChecker_ValidJSONLD(t *testing.T) {
	page := makePageData(`<html><head>
		<script type="application/ld+json">{"@type":"Organization","name":"Example","url":"https://example.com"}</script>
	</head><body></body></html>`)
	checker := &structured_data.JsonLdChecker{}
	rules := checker.Check(page)

	existsRule := findRule(rules, "jsonld_exists")
	if existsRule == nil || existsRule.Result != valueobject.RuleResultPass {
		t.Error("expected jsonld_exists to pass")
	}

	validRule := findRule(rules, "jsonld_valid")
	if validRule == nil || validRule.Result != valueobject.RuleResultPass {
		t.Error("expected jsonld_valid to pass for valid JSON")
	}
}

func TestJsonLdChecker_InvalidJSON(t *testing.T) {
	page := makePageData(`<html><head>
		<script type="application/ld+json">{invalid json here}</script>
	</head><body></body></html>`)
	checker := &structured_data.JsonLdChecker{}
	rules := checker.Check(page)

	validRule := findRule(rules, "jsonld_valid")
	if validRule == nil || validRule.Result != valueobject.RuleResultFail {
		t.Error("expected jsonld_valid to fail for invalid JSON")
	}
}

func TestJsonLdChecker_WithBreadcrumb(t *testing.T) {
	page := makePageData(`<html><head>
		<script type="application/ld+json">{"@type":"BreadcrumbList","itemListElement":[]}</script>
	</head><body></body></html>`)
	checker := &structured_data.JsonLdChecker{}
	rules := checker.Check(page)

	breadcrumbRule := findRule(rules, "jsonld_breadcrumb")
	if breadcrumbRule == nil || breadcrumbRule.Result != valueobject.RuleResultPass {
		t.Error("expected jsonld_breadcrumb to pass when BreadcrumbList present")
	}
}

func TestJsonLdChecker_NoBreadcrumb(t *testing.T) {
	page := makePageData(`<html><head>
		<script type="application/ld+json">{"@type":"Organization","name":"Example"}</script>
	</head><body></body></html>`)
	checker := &structured_data.JsonLdChecker{}
	rules := checker.Check(page)

	breadcrumbRule := findRule(rules, "jsonld_breadcrumb")
	if breadcrumbRule == nil || breadcrumbRule.Result != valueobject.RuleResultWarning {
		t.Error("expected jsonld_breadcrumb to warn when missing")
	}
}
