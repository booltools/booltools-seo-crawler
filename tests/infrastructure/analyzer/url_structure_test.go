package analyzer_test

import (
	"strings"
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer/rules/url_structure"
)

func TestURLStructureChecker_CleanURL(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com/about-us"
	checker := &url_structure.URLStructureChecker{}
	rules := checker.Check(page)

	lowercaseRule := findRule(rules, "url_lowercase")
	if lowercaseRule == nil || lowercaseRule.Result != valueobject.RuleResultPass {
		t.Error("expected url_lowercase to pass for lowercase URL")
	}

	hyphensRule := findRule(rules, "url_hyphens")
	if hyphensRule == nil || hyphensRule.Result != valueobject.RuleResultPass {
		t.Error("expected url_hyphens to pass for URL with hyphens")
	}

	doubleSlashRule := findRule(rules, "url_double_slash")
	if doubleSlashRule == nil || doubleSlashRule.Result != valueobject.RuleResultPass {
		t.Error("expected url_double_slash to pass for clean URL")
	}
}

func TestURLStructureChecker_UppercaseURL(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com/About-Us"
	checker := &url_structure.URLStructureChecker{}
	rules := checker.Check(page)

	lowercaseRule := findRule(rules, "url_lowercase")
	if lowercaseRule == nil || lowercaseRule.Result != valueobject.RuleResultWarning {
		t.Error("expected url_lowercase to warn for uppercase URL")
	}
}

func TestURLStructureChecker_UnderscoresInURL(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com/about_us"
	checker := &url_structure.URLStructureChecker{}
	rules := checker.Check(page)

	hyphensRule := findRule(rules, "url_hyphens")
	if hyphensRule == nil || hyphensRule.Result != valueobject.RuleResultWarning {
		t.Error("expected url_hyphens to warn for underscores")
	}
}

func TestURLStructureChecker_LongURL(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com/" + strings.Repeat("very-long-path-segment/", 10)
	checker := &url_structure.URLStructureChecker{}
	rules := checker.Check(page)

	lengthRule := findRule(rules, "url_length")
	if lengthRule == nil || lengthRule.Result != valueobject.RuleResultWarning {
		t.Error("expected url_length to warn for long URL")
	}
}

func TestURLStructureChecker_DoubleSlash(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com//about"
	checker := &url_structure.URLStructureChecker{}
	rules := checker.Check(page)

	doubleSlashRule := findRule(rules, "url_double_slash")
	if doubleSlashRule == nil || doubleSlashRule.Result != valueobject.RuleResultFail {
		t.Error("expected url_double_slash to fail for double slashes")
	}
}

func TestURLStructureChecker_ManyParameters(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com/search?q=test&page=1&sort=asc&filter=new"
	checker := &url_structure.URLStructureChecker{}
	rules := checker.Check(page)

	paramsRule := findRule(rules, "url_parameters")
	if paramsRule == nil || paramsRule.Result != valueobject.RuleResultWarning {
		t.Error("expected url_parameters to warn for many query params")
	}
}
