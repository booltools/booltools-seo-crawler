package analyzer_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/on_page"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

func makePageData(html string) crawler.PageData {
	document, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	return crawler.PageData{
		URL:        "https://example.com/test",
		StatusCode: 200,
		Headers:    http.Header{},
		Document:   document,
		HTML:       html,
		BodyText:   document.Find("body").Text(),
	}
}

func TestTitleChecker_MissingTitle(t *testing.T) {
	page := makePageData(`<html><head></head><body>Hello</body></html>`)
	checker := &on_page.TitleChecker{}
	rules := checker.Check(page)

	if len(rules) == 0 {
		t.Fatal("expected at least one rule")
	}

	if rules[0].Result != valueobject.RuleResultFail {
		t.Errorf("expected fail for missing title, got %s", rules[0].Result)
	}
}

func TestTitleChecker_ValidTitle(t *testing.T) {
	page := makePageData(`<html><head><title>A Good Page Title For SEO Testing</title></head><body>Hello</body></html>`)
	checker := &on_page.TitleChecker{}
	rules := checker.Check(page)

	titleExists := findRule(rules, "title_exists")
	if titleExists == nil || titleExists.Result != valueobject.RuleResultPass {
		t.Error("expected title_exists to pass")
	}

	titleLength := findRule(rules, "title_length")
	if titleLength == nil || titleLength.Result != valueobject.RuleResultPass {
		t.Errorf("expected title_length to pass, got %v", titleLength)
	}
}

func TestTitleChecker_TooLong(t *testing.T) {
	longTitle := strings.Repeat("A", 70)
	page := makePageData(`<html><head><title>` + longTitle + `</title></head><body>Hello</body></html>`)
	checker := &on_page.TitleChecker{}
	rules := checker.Check(page)

	titleLength := findRule(rules, "title_length")
	if titleLength == nil || titleLength.Result != valueobject.RuleResultWarning {
		t.Error("expected title_length to warn for long title")
	}
}

func TestMetaDescriptionChecker_Missing(t *testing.T) {
	page := makePageData(`<html><head><title>Test</title></head><body>Hello</body></html>`)
	checker := &on_page.MetaDescriptionChecker{}
	rules := checker.Check(page)

	descExists := findRule(rules, "meta_description_exists")
	if descExists == nil || descExists.Result != valueobject.RuleResultFail {
		t.Error("expected meta_description_exists to fail")
	}
}

func TestMetaDescriptionChecker_Valid(t *testing.T) {
	desc := strings.Repeat("A good description. ", 8)
	page := makePageData(`<html><head><meta name="description" content="` + desc + `"></head><body>Hello</body></html>`)
	checker := &on_page.MetaDescriptionChecker{}
	rules := checker.Check(page)

	descExists := findRule(rules, "meta_description_exists")
	if descExists == nil || descExists.Result != valueobject.RuleResultPass {
		t.Error("expected meta_description_exists to pass")
	}
}

func TestHeadingChecker_NoH1(t *testing.T) {
	page := makePageData(`<html><head><title>Test</title></head><body><h2>Sub</h2></body></html>`)
	checker := &on_page.HeadingChecker{}
	rules := checker.Check(page)

	h1Count := findRule(rules, "h1_count")
	if h1Count == nil || h1Count.Result != valueobject.RuleResultFail {
		t.Error("expected h1_count to fail when no H1")
	}
}

func TestHeadingChecker_MultipleH1(t *testing.T) {
	page := makePageData(`<html><head></head><body><h1>First</h1><h1>Second</h1></body></html>`)
	checker := &on_page.HeadingChecker{}
	rules := checker.Check(page)

	h1Count := findRule(rules, "h1_count")
	if h1Count == nil || h1Count.Result != valueobject.RuleResultWarning {
		t.Error("expected h1_count to warn for multiple H1s")
	}
}

func TestHeadingChecker_CorrectHierarchy(t *testing.T) {
	page := makePageData(`<html><head></head><body><h1>Title</h1><h2>Sub</h2><h3>Sub Sub</h3></body></html>`)
	checker := &on_page.HeadingChecker{}
	rules := checker.Check(page)

	hierarchy := findRule(rules, "heading_hierarchy")
	if hierarchy == nil || hierarchy.Result != valueobject.RuleResultPass {
		t.Error("expected heading_hierarchy to pass")
	}
}

func TestHeadingChecker_SkippedHierarchy(t *testing.T) {
	page := makePageData(`<html><head></head><body><h1>Title</h1><h3>Skipped H2</h3></body></html>`)
	checker := &on_page.HeadingChecker{}
	rules := checker.Check(page)

	hierarchy := findRule(rules, "heading_hierarchy")
	if hierarchy == nil || hierarchy.Result != valueobject.RuleResultWarning {
		t.Error("expected heading_hierarchy to warn for skipped levels")
	}
}

func TestImageChecker_MissingAlt(t *testing.T) {
	page := makePageData(`<html><head></head><body><img src="test.jpg"></body></html>`)
	page.Images = []crawler.ImageData{
		{URL: "https://example.com/test.jpg", Alt: ""},
	}

	checker := &on_page.ImageChecker{}
	rules := checker.Check(page)

	altRule := findRule(rules, "images_alt_text")
	if altRule == nil || altRule.Result != valueobject.RuleResultFail {
		t.Error("expected images_alt_text to fail for missing alt")
	}
	if altRule.Details == "" {
		t.Error("expected images_alt_text to include image URL in details")
	}
}

func TestImageChecker_AllAltPresent(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Images = []crawler.ImageData{
		{URL: "https://example.com/a.jpg", Alt: "A photo of nature"},
		{URL: "https://example.com/b.jpg", Alt: "Profile picture"},
	}
	checker := &on_page.ImageChecker{}
	rules := checker.Check(page)

	altRule := findRule(rules, "images_alt_text")
	if altRule == nil || altRule.Result != valueobject.RuleResultPass {
		t.Error("expected images_alt_text to pass when all have alt")
	}
}

func TestImageChecker_LegacyFormat(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Images = []crawler.ImageData{
		{URL: "https://example.com/photo.jpg", Alt: "Photo"},
	}
	checker := &on_page.ImageChecker{}
	rules := checker.Check(page)

	formatRule := findRule(rules, "images_modern_format")
	if formatRule == nil || formatRule.Result != valueobject.RuleResultWarning {
		t.Error("expected images_modern_format to warn for .jpg format")
	}
	if formatRule.Details == "" {
		t.Error("expected images_modern_format to include image URL in details")
	}
}

func TestImageChecker_ModernFormat(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Images = []crawler.ImageData{
		{URL: "https://example.com/photo.webp", Alt: "Photo"},
	}
	checker := &on_page.ImageChecker{}
	rules := checker.Check(page)

	formatRule := findRule(rules, "images_modern_format")
	if formatRule == nil || formatRule.Result != valueobject.RuleResultPass {
		t.Error("expected images_modern_format to pass for .webp format")
	}
}

func TestImageChecker_MissingDimensions(t *testing.T) {
	page := makePageData(`<html><head></head><body></body></html>`)
	page.Images = []crawler.ImageData{
		{URL: "https://example.com/img.jpg", Alt: "Test", Width: "", Height: ""},
	}
	checker := &on_page.ImageChecker{}
	rules := checker.Check(page)

	dimRule := findRule(rules, "images_dimensions")
	if dimRule == nil || dimRule.Result != valueobject.RuleResultWarning {
		t.Error("expected images_dimensions to warn when missing width/height")
	}
}

func findRule(rules []valueobject.AuditRule, key string) *valueobject.AuditRule {
	for _, rule := range rules {
		if rule.Key == key {
			return &rule
		}
	}
	return nil
}
