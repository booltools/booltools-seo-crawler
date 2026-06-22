package analyzer_test

import (
	"testing"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/internationalization"
)

func TestHreflangChecker_NoHreflang(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	checker := &internationalization.HreflangChecker{}
	rules := checker.Check(page)

	if len(rules) != 0 {
		t.Error("expected no rules when no hreflang tags present")
	}
}

func TestHreflangChecker_ValidHreflang(t *testing.T) {
	page := makePageData(`<html><head>
		<link rel="alternate" hreflang="en" href="https://example.com/en/">
		<link rel="alternate" hreflang="pt-BR" href="https://example.com/pt-br/">
		<link rel="alternate" hreflang="x-default" href="https://example.com/">
	</head><body></body></html>`)
	checker := &internationalization.HreflangChecker{}
	rules := checker.Check(page)

	validRule := findRule(rules, "hreflang_valid")
	if validRule == nil || validRule.Result != valueobject.RuleResultPass {
		t.Error("expected hreflang_valid to pass for valid codes")
	}

	xDefaultRule := findRule(rules, "hreflang_x_default")
	if xDefaultRule == nil || xDefaultRule.Result != valueobject.RuleResultPass {
		t.Error("expected hreflang_x_default to pass when present")
	}
}

func TestHreflangChecker_InvalidCode(t *testing.T) {
	page := makePageData(`<html><head>
		<link rel="alternate" hreflang="INVALID" href="https://example.com/invalid/">
	</head><body></body></html>`)
	checker := &internationalization.HreflangChecker{}
	rules := checker.Check(page)

	validRule := findRule(rules, "hreflang_valid")
	if validRule == nil || validRule.Result != valueobject.RuleResultFail {
		t.Error("expected hreflang_valid to fail for invalid language code")
	}
}

func TestHreflangChecker_MissingXDefault(t *testing.T) {
	page := makePageData(`<html><head>
		<link rel="alternate" hreflang="en" href="https://example.com/en/">
	</head><body></body></html>`)
	checker := &internationalization.HreflangChecker{}
	rules := checker.Check(page)

	xDefaultRule := findRule(rules, "hreflang_x_default")
	if xDefaultRule == nil || xDefaultRule.Result != valueobject.RuleResultWarning {
		t.Error("expected hreflang_x_default to warn when missing")
	}
}
