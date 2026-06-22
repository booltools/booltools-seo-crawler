package analyzer_test

import (
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer/rules/social"
)

func TestOpenGraphChecker_AllPresent(t *testing.T) {
	page := makePageData(`<html><head>
		<meta property="og:title" content="Test Page">
		<meta property="og:description" content="A test description">
		<meta property="og:image" content="https://example.com/image.jpg">
		<meta property="og:url" content="https://example.com/test">
		<meta property="og:type" content="website">
		<meta property="og:site_name" content="Example">
		<meta property="og:locale" content="en_US">
	</head><body></body></html>`)

	checker := &social.OpenGraphChecker{}
	rules := checker.Check(page)

	for _, rule := range rules {
		if rule.Result != valueobject.RuleResultPass {
			t.Errorf("expected %s to pass, got %s", rule.Key, rule.Result)
		}
	}
}

func TestOpenGraphChecker_AllMissing(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	checker := &social.OpenGraphChecker{}
	rules := checker.Check(page)

	failingKeys := []string{"og_title", "og_description", "og_image"}
	for _, key := range failingKeys {
		rule := findRule(rules, key)
		if rule == nil || rule.Result != valueobject.RuleResultFail {
			t.Errorf("expected %s to fail when OG tags missing", key)
		}
	}
}

func TestTwitterCardChecker_AllPresent(t *testing.T) {
	page := makePageData(`<html><head>
		<meta name="twitter:card" content="summary_large_image">
		<meta name="twitter:title" content="Test">
		<meta name="twitter:description" content="A description">
		<meta name="twitter:image" content="https://example.com/img.jpg">
		<meta name="twitter:site" content="@example">
	</head><body></body></html>`)

	checker := &social.TwitterCardChecker{}
	rules := checker.Check(page)

	for _, rule := range rules {
		if rule.Result != valueobject.RuleResultPass {
			t.Errorf("expected %s to pass, got %s", rule.Key, rule.Result)
		}
	}
}

func TestTwitterCardChecker_AllMissing(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	checker := &social.TwitterCardChecker{}
	rules := checker.Check(page)

	for _, rule := range rules {
		if rule.Result != valueobject.RuleResultWarning {
			t.Errorf("expected %s to warn when missing, got %s", rule.Key, rule.Result)
		}
	}
}
