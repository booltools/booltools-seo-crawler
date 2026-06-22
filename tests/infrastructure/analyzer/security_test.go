package analyzer_test

import (
	"net/http"
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer/rules/security"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

func TestHTTPSChecker_PassHTTPS(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com/page"
	checker := &security.HTTPSChecker{}
	rules := checker.Check(page)

	httpsRule := findRule(rules, "uses_https")
	if httpsRule == nil || httpsRule.Result != valueobject.RuleResultPass {
		t.Error("expected uses_https to pass for HTTPS URL")
	}
}

func TestHTTPSChecker_FailHTTP(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "http://example.com/page"
	checker := &security.HTTPSChecker{}
	rules := checker.Check(page)

	httpsRule := findRule(rules, "uses_https")
	if httpsRule == nil || httpsRule.Result != valueobject.RuleResultFail {
		t.Error("expected uses_https to fail for HTTP URL")
	}
}

func TestHTTPSChecker_MixedContent(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com/page"
	page.Images = []crawler.ImageData{
		{URL: "http://example.com/image.jpg"},
	}

	checker := &security.HTTPSChecker{}
	rules := checker.Check(page)

	mixedRule := findRule(rules, "mixed_content")
	if mixedRule == nil || mixedRule.Result != valueobject.RuleResultFail {
		t.Error("expected mixed_content to fail when HTTP images on HTTPS page")
	}
	if mixedRule.Details == "" {
		t.Error("expected mixed_content to include asset details")
	}
}

func TestHTTPSChecker_NoMixedContent(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com/page"
	page.Images = []crawler.ImageData{
		{URL: "https://example.com/image.jpg"},
	}

	checker := &security.HTTPSChecker{}
	rules := checker.Check(page)

	mixedRule := findRule(rules, "mixed_content")
	if mixedRule == nil || mixedRule.Result != valueobject.RuleResultPass {
		t.Error("expected mixed_content to pass when all resources use HTTPS")
	}
}

func TestHTTPSChecker_HSTSMissing(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com/page"
	checker := &security.HTTPSChecker{}
	rules := checker.Check(page)

	hstsRule := findRule(rules, "hsts_header")
	if hstsRule == nil || hstsRule.Result != valueobject.RuleResultWarning {
		t.Error("expected hsts_header to warn when missing")
	}
}

func TestHTTPSChecker_HSTSPresent(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.URL = "https://example.com/page"
	page.Headers = http.Header{"Strict-Transport-Security": []string{"max-age=31536000"}}
	checker := &security.HTTPSChecker{}
	rules := checker.Check(page)

	hstsRule := findRule(rules, "hsts_header")
	if hstsRule == nil || hstsRule.Result != valueobject.RuleResultPass {
		t.Error("expected hsts_header to pass when present")
	}
}

func TestSecurityHeaderChecker_AllMissing(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	checker := &security.SecurityHeaderChecker{}
	rules := checker.Check(page)

	expectedWarnings := []string{"security_xcto", "security_xfo", "security_csp", "security_referrer", "security_permissions"}
	for _, key := range expectedWarnings {
		rule := findRule(rules, key)
		if rule == nil || rule.Result != valueobject.RuleResultWarning {
			t.Errorf("expected %s to warn when header missing", key)
		}
	}
}

func TestSecurityHeaderChecker_AllPresent(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Headers = http.Header{
		"X-Content-Type-Options":  []string{"nosniff"},
		"X-Frame-Options":        []string{"DENY"},
		"Content-Security-Policy": []string{"default-src 'self'"},
		"Referrer-Policy":        []string{"strict-origin-when-cross-origin"},
		"Permissions-Policy":     []string{"camera=(), microphone=()"},
	}
	checker := &security.SecurityHeaderChecker{}
	rules := checker.Check(page)

	expectedPasses := []string{"security_xcto", "security_xfo", "security_csp", "security_referrer", "security_permissions"}
	for _, key := range expectedPasses {
		rule := findRule(rules, key)
		if rule == nil || rule.Result != valueobject.RuleResultPass {
			t.Errorf("expected %s to pass when header present", key)
		}
	}
}

func TestSecurityHeaderChecker_ServerDisclosure(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Headers = http.Header{"Server": []string{"Apache/2.4.41"}}
	checker := &security.SecurityHeaderChecker{}
	rules := checker.Check(page)

	serverRule := findRule(rules, "security_server_disclosure")
	if serverRule == nil || serverRule.Result != valueobject.RuleResultWarning {
		t.Error("expected server_disclosure to warn when Server header present")
	}
}
