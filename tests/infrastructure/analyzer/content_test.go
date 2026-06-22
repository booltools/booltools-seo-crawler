package analyzer_test

import (
	"strings"
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer/rules/content"
)

func TestContentQualityChecker_ThinContent(t *testing.T) {
	page := makePageData(`<html><head></head><body><p>Very short text.</p></body></html>`)
	checker := &content.ContentQualityChecker{}
	rules := checker.Check(page)

	wordRule := findRule(rules, "content_word_count")
	if wordRule == nil || wordRule.Result != valueobject.RuleResultFail {
		t.Error("expected content_word_count to fail for very thin content")
	}
}

func TestContentQualityChecker_AdequateContent(t *testing.T) {
	longText := strings.Repeat("This is a meaningful sentence with real content. ", 40)
	page := makePageData(`<html><head></head><body><p>` + longText + `</p></body></html>`)
	checker := &content.ContentQualityChecker{}
	rules := checker.Check(page)

	wordRule := findRule(rules, "content_word_count")
	if wordRule == nil || wordRule.Result != valueobject.RuleResultPass {
		t.Error("expected content_word_count to pass for adequate content")
	}
}

func TestContentQualityChecker_LowTextHTMLRatio(t *testing.T) {
	hugeMarkup := strings.Repeat(`<div class="wrapper"><span class="style"></span></div>`, 100)
	page := makePageData(`<html><head></head><body>` + hugeMarkup + `<p>tiny</p></body></html>`)
	checker := &content.ContentQualityChecker{}
	rules := checker.Check(page)

	ratioRule := findRule(rules, "content_text_html_ratio")
	if ratioRule == nil || ratioRule.Result != valueobject.RuleResultWarning {
		t.Error("expected content_text_html_ratio to warn for low ratio")
	}
}

func TestContentQualityChecker_HealthyRatio(t *testing.T) {
	longText := strings.Repeat("Good content text that provides value. ", 50)
	page := makePageData(`<html><head></head><body><p>` + longText + `</p></body></html>`)
	checker := &content.ContentQualityChecker{}
	rules := checker.Check(page)

	ratioRule := findRule(rules, "content_text_html_ratio")
	if ratioRule == nil || ratioRule.Result != valueobject.RuleResultPass {
		t.Error("expected content_text_html_ratio to pass for healthy ratio")
	}
}
