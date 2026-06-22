package analyzer_test

import (
	"testing"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/accessibility"
)

func TestAccessibilityChecker_MissingLang(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	checker := &accessibility.AccessibilityChecker{}
	rules := checker.Check(page)

	langRule := findRule(rules, "html_lang")
	if langRule == nil || langRule.Result != valueobject.RuleResultFail {
		t.Error("expected html_lang to fail when missing")
	}
}

func TestAccessibilityChecker_WithLang(t *testing.T) {
	page := makePageData(`<html lang="en"><head></head><body></body></html>`)
	checker := &accessibility.AccessibilityChecker{}
	rules := checker.Check(page)

	langRule := findRule(rules, "html_lang")
	if langRule == nil || langRule.Result != valueobject.RuleResultPass {
		t.Error("expected html_lang to pass when present")
	}
}

func TestAccessibilityChecker_MissingViewport(t *testing.T) {
	page := makePageData(`<html lang="en"><head></head><body></body></html>`)
	checker := &accessibility.AccessibilityChecker{}
	rules := checker.Check(page)

	viewportRule := findRule(rules, "viewport_meta")
	if viewportRule == nil || viewportRule.Result != valueobject.RuleResultFail {
		t.Error("expected viewport_meta to fail when missing")
	}
}

func TestAccessibilityChecker_WithViewport(t *testing.T) {
	page := makePageData(`<html lang="en"><head><meta name="viewport" content="width=device-width, initial-scale=1"></head><body></body></html>`)
	checker := &accessibility.AccessibilityChecker{}
	rules := checker.Check(page)

	viewportRule := findRule(rules, "viewport_meta")
	if viewportRule == nil || viewportRule.Result != valueobject.RuleResultPass {
		t.Error("expected viewport_meta to pass when present")
	}
}

func TestAccessibilityChecker_EmptyLinks(t *testing.T) {
	page := makePageData(`<html lang="en"><head></head><body><a href="/page"></a><a href="/other"></a></body></html>`)
	checker := &accessibility.AccessibilityChecker{}
	rules := checker.Check(page)

	emptyRule := findRule(rules, "empty_links")
	if emptyRule == nil || emptyRule.Result != valueobject.RuleResultWarning {
		t.Error("expected empty_links to warn for links with no text")
	}
	if emptyRule.Details == "" {
		t.Error("expected empty_links to include href details")
	}
}

func TestAccessibilityChecker_LinksWithText(t *testing.T) {
	page := makePageData(`<html lang="en"><head></head><body><a href="/page">Visit Page</a></body></html>`)
	checker := &accessibility.AccessibilityChecker{}
	rules := checker.Check(page)

	emptyRule := findRule(rules, "empty_links")
	if emptyRule == nil || emptyRule.Result != valueobject.RuleResultPass {
		t.Error("expected empty_links to pass when all links have text")
	}
}

func TestAccessibilityChecker_ARIALandmarks(t *testing.T) {
	page := makePageData(`<html lang="en"><head></head><body><nav>Menu</nav><main>Content</main></body></html>`)
	checker := &accessibility.AccessibilityChecker{}
	rules := checker.Check(page)

	ariaRule := findRule(rules, "aria_landmarks")
	if ariaRule == nil || ariaRule.Result != valueobject.RuleResultPass {
		t.Error("expected aria_landmarks to pass when nav and main present")
	}
}

func TestAccessibilityChecker_NoLandmarks(t *testing.T) {
	page := makePageData(`<html lang="en"><head></head><body><div>Content</div></body></html>`)
	checker := &accessibility.AccessibilityChecker{}
	rules := checker.Check(page)

	ariaRule := findRule(rules, "aria_landmarks")
	if ariaRule == nil || ariaRule.Result != valueobject.RuleResultWarning {
		t.Error("expected aria_landmarks to warn when no landmarks")
	}
}

func TestAccessibilityChecker_MissingCharset(t *testing.T) {
	page := makePageData(`<html lang="en"><head></head><body></body></html>`)
	checker := &accessibility.AccessibilityChecker{}
	rules := checker.Check(page)

	charsetRule := findRule(rules, "charset_meta")
	if charsetRule == nil || charsetRule.Result != valueobject.RuleResultWarning {
		t.Error("expected charset_meta to warn when missing")
	}
}

func TestAccessibilityChecker_WithCharset(t *testing.T) {
	page := makePageData(`<html lang="en"><head><meta charset="UTF-8"></head><body></body></html>`)
	checker := &accessibility.AccessibilityChecker{}
	rules := checker.Check(page)

	charsetRule := findRule(rules, "charset_meta")
	if charsetRule == nil || charsetRule.Result != valueobject.RuleResultPass {
		t.Error("expected charset_meta to pass when present")
	}
}
