package analyzer_test

import (
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer/rules/links"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

func TestInternalLinkChecker_NoLinks(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.InternalLinks = []crawler.LinkData{}
	checker := &links.InternalLinkChecker{}
	rules := checker.Check(page)

	presentRule := findRule(rules, "internal_links_present")
	if presentRule == nil || presentRule.Result != valueobject.RuleResultFail {
		t.Error("expected internal_links_present to fail when no links")
	}
}

func TestInternalLinkChecker_WithLinks(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.InternalLinks = []crawler.LinkData{
		{URL: "https://example.com/about", AnchorText: "About Us"},
		{URL: "https://example.com/contact", AnchorText: "Contact"},
	}
	checker := &links.InternalLinkChecker{}
	rules := checker.Check(page)

	presentRule := findRule(rules, "internal_links_present")
	if presentRule == nil || presentRule.Result != valueobject.RuleResultPass {
		t.Error("expected internal_links_present to pass with links")
	}
}

func TestInternalLinkChecker_GenericAnchorText(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.InternalLinks = []crawler.LinkData{
		{URL: "https://example.com/page1", AnchorText: "click here"},
		{URL: "https://example.com/page2", AnchorText: "read more"},
	}
	checker := &links.InternalLinkChecker{}
	rules := checker.Check(page)

	anchorRule := findRule(rules, "internal_links_anchor_text")
	if anchorRule == nil || anchorRule.Result != valueobject.RuleResultWarning {
		t.Error("expected anchor_text to warn for generic text")
	}
	if anchorRule.Details == "" {
		t.Error("expected anchor_text to include details about which links")
	}
}

func TestInternalLinkChecker_DescriptiveAnchorText(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.InternalLinks = []crawler.LinkData{
		{URL: "https://example.com/pricing", AnchorText: "View our pricing plans"},
	}
	checker := &links.InternalLinkChecker{}
	rules := checker.Check(page)

	anchorRule := findRule(rules, "internal_links_anchor_text")
	if anchorRule == nil || anchorRule.Result != valueobject.RuleResultPass {
		t.Error("expected anchor_text to pass for descriptive text")
	}
}

func TestExternalLinkChecker_MissingRel(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.ExternalLinks = []crawler.LinkData{
		{URL: "https://other.com", Target: "_blank", Rel: ""},
	}
	checker := &links.ExternalLinkChecker{}
	rules := checker.Check(page)

	relRule := findRule(rules, "external_links_rel")
	if relRule == nil || relRule.Result != valueobject.RuleResultWarning {
		t.Error("expected external_links_rel to warn when rel missing")
	}
	if relRule.Details == "" {
		t.Error("expected external_links_rel to include URL details")
	}
}

func TestExternalLinkChecker_ProperRel(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.ExternalLinks = []crawler.LinkData{
		{URL: "https://other.com", Target: "_blank", Rel: "noopener noreferrer"},
	}
	checker := &links.ExternalLinkChecker{}
	rules := checker.Check(page)

	relRule := findRule(rules, "external_links_rel")
	if relRule == nil || relRule.Result != valueobject.RuleResultPass {
		t.Error("expected external_links_rel to pass with proper rel")
	}
}
