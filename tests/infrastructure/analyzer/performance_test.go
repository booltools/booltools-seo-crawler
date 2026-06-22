package analyzer_test

import (
	"net/http"
	"testing"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/performance"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
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
	scripts := make([]crawler.ResourceData, 15)
	for i := range scripts {
		scripts[i] = crawler.ResourceData{URL: "https://example.com/script.js"}
	}
	page.Scripts = scripts
	checker := &performance.ResourceChecker{}
	rules := checker.Check(page)

	jsRule := findRule(rules, "js_file_count")
	if jsRule == nil || jsRule.Result != valueobject.RuleResultWarning {
		t.Error("expected js_file_count to warn for 15 scripts")
	}
	if jsRule.Details == "" {
		t.Error("expected js_file_count to include script URL details")
	}
}

func TestResourceChecker_RenderBlocking(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Scripts = []crawler.ResourceData{
		{URL: "https://example.com/blocking.js", Location: "head", IsAsync: false, IsDefer: false},
	}
	checker := &performance.ResourceChecker{}
	rules := checker.Check(page)

	blockingRule := findRule(rules, "render_blocking")
	if blockingRule == nil || blockingRule.Result != valueobject.RuleResultFail {
		t.Error("expected render_blocking to fail for sync script in head")
	}
	if blockingRule.Details == "" {
		t.Error("expected render_blocking to include script URL details")
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
