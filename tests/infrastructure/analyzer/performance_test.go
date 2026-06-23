package analyzer_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer/rules/performance"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

func TestResourceChecker_FewScripts(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Scripts = make([]crawler.ResourceData, 3)
	checker := &performance.ResourceChecker{}
	rules := checker.Check(page)

	jsRule := findRule(rules, "js_file_count")
	if jsRule == nil || jsRule.Result != valueobject.RuleResultPass {
		t.Error("expected js_file_count to pass for 3 scripts")
	}
}

func TestResourceChecker_TooManyScripts(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	scripts := make([]crawler.ResourceData, 35)
	for i := range scripts {
		scripts[i] = crawler.ResourceData{URL: "https://example.com/script.js"}
	}
	page.Scripts = scripts
	checker := &performance.ResourceChecker{}
	rules := checker.Check(page)

	jsRule := findRule(rules, "js_file_count")
	if jsRule == nil || jsRule.Result != valueobject.RuleResultWarning {
		t.Error("expected js_file_count to warn for 35 scripts")
	}
	if jsRule.Details == "" {
		t.Error("expected js_file_count to include script URL details")
	}
}

func TestResourceChecker_RenderBlocking(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Scripts = []crawler.ResourceData{
		{URL: "https://example.com/blocking1.js", Location: "head", IsAsync: false, IsDefer: false},
		{URL: "https://example.com/blocking2.js", Location: "head", IsAsync: false, IsDefer: false},
		{URL: "https://example.com/blocking3.js", Location: "head", IsAsync: false, IsDefer: false},
		{URL: "https://example.com/blocking4.js", Location: "head", IsAsync: false, IsDefer: false},
	}
	checker := &performance.ResourceChecker{}
	rules := checker.Check(page)

	blockingRule := findRule(rules, "render_blocking")
	if blockingRule == nil || blockingRule.Result != valueobject.RuleResultWarning {
		t.Error("expected render_blocking to warn for 4+ sync scripts in head")
	}
	if blockingRule.Details == "" {
		t.Error("expected render_blocking to include script URL details")
	}
}

func TestResourceChecker_RenderBlockingSingleFrameworkScript(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Scripts = []crawler.ResourceData{
		{URL: "https://example.com/_next/static/chunks/bootstrap.js", Location: "head", IsAsync: false, IsDefer: false},
	}
	checker := &performance.ResourceChecker{}
	rules := checker.Check(page)

	blockingRule := findRule(rules, "render_blocking")
	if blockingRule == nil || blockingRule.Result != valueobject.RuleResultPass {
		t.Error("expected render_blocking to pass for 1-3 framework bootstrap scripts")
	}
}

func TestResourceChecker_NoRenderBlocking(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Scripts = []crawler.ResourceData{
		{URL: "https://example.com/async.js", Location: "head", IsAsync: true, IsDefer: false},
		{URL: "https://example.com/defer.js", Location: "head", IsAsync: false, IsDefer: true},
	}
	checker := &performance.ResourceChecker{}
	rules := checker.Check(page)

	blockingRule := findRule(rules, "render_blocking")
	if blockingRule == nil || blockingRule.Result != valueobject.RuleResultPass {
		t.Error("expected render_blocking to pass when all scripts are async/defer")
	}
}

func TestResourceChecker_CacheHeaderMissing(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	checker := &performance.ResourceChecker{}
	rules := checker.Check(page)

	cacheRule := findRule(rules, "cache_headers")
	if cacheRule == nil || cacheRule.Result != valueobject.RuleResultWarning {
		t.Error("expected cache_headers to warn when missing")
	}
}

func TestResourceChecker_CacheHeaderPresent(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Headers = http.Header{"Cache-Control": []string{"max-age=3600"}}
	checker := &performance.ResourceChecker{}
	rules := checker.Check(page)

	cacheRule := findRule(rules, "cache_headers")
	if cacheRule == nil || cacheRule.Result != valueobject.RuleResultPass {
		t.Error("expected cache_headers to pass when present")
	}
}

func TestPageSpeedChecker_TTFBFast(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.ResponseTime = 300 * time.Millisecond
	checker := &performance.PageSpeedChecker{}
	rules := checker.Check(page)

	ttfbRule := findRule(rules, "ttfb")
	if ttfbRule == nil || ttfbRule.Result != valueobject.RuleResultPass {
		t.Error("expected ttfb to pass for 300ms response time")
	}
}

func TestPageSpeedChecker_TTFBWarn(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.ResponseTime = 900 * time.Millisecond
	checker := &performance.PageSpeedChecker{}
	rules := checker.Check(page)

	ttfbRule := findRule(rules, "ttfb")
	if ttfbRule == nil || ttfbRule.Result != valueobject.RuleResultWarning {
		t.Errorf("expected ttfb to warn for 900ms response time, got %v", ttfbRule)
	}
}

func TestPageSpeedChecker_TTFBFail(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.ResponseTime = 1500 * time.Millisecond
	checker := &performance.PageSpeedChecker{}
	rules := checker.Check(page)

	ttfbRule := findRule(rules, "ttfb")
	if ttfbRule == nil || ttfbRule.Result != valueobject.RuleResultFail {
		t.Errorf("expected ttfb to fail for 1500ms response time, got %v", ttfbRule)
	}
}

func TestPageSpeedChecker_TTFBOldThresholdNowPasses(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.ResponseTime = 500 * time.Millisecond
	checker := &performance.PageSpeedChecker{}
	rules := checker.Check(page)

	ttfbRule := findRule(rules, "ttfb")
	if ttfbRule == nil || ttfbRule.Result != valueobject.RuleResultPass {
		t.Error("expected ttfb to pass for 500ms (previously would have warned at 400ms)")
	}
}
