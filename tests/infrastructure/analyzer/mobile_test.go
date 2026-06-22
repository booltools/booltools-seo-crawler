package analyzer_test

import (
	"testing"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/mobile"
)

func TestMobileChecker_MissingViewport(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	checker := &mobile.MobileChecker{}
	rules := checker.Check(page)

	viewportRule := findRule(rules, "mobile_viewport")
	if viewportRule == nil || viewportRule.Result != valueobject.RuleResultFail {
		t.Error("expected mobile_viewport to fail when missing")
	}
}

func TestMobileChecker_ProperViewport(t *testing.T) {
	page := makePageData(`<html><head><meta name="viewport" content="width=device-width, initial-scale=1"></head><body></body></html>`)
	checker := &mobile.MobileChecker{}
	rules := checker.Check(page)

	viewportRule := findRule(rules, "mobile_viewport")
	if viewportRule == nil || viewportRule.Result != valueobject.RuleResultPass {
		t.Error("expected mobile_viewport to pass when present")
	}

	configRule := findRule(rules, "mobile_viewport_config")
	if configRule == nil || configRule.Result != valueobject.RuleResultPass {
		t.Error("expected mobile_viewport_config to pass with proper config")
	}
}

func TestMobileChecker_PartialViewport(t *testing.T) {
	page := makePageData(`<html><head><meta name="viewport" content="width=500"></head><body></body></html>`)
	checker := &mobile.MobileChecker{}
	rules := checker.Check(page)

	configRule := findRule(rules, "mobile_viewport_config")
	if configRule == nil || configRule.Result != valueobject.RuleResultWarning {
		t.Error("expected mobile_viewport_config to warn with non-optimal config")
	}
}
