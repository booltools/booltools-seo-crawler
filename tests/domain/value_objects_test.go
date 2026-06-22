package domain_test

import (
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
)

func TestSeverity_IsValid(t *testing.T) {
	validSeverities := []valueobject.Severity{
		valueobject.SeverityCritical,
		valueobject.SeverityHigh,
		valueobject.SeverityMedium,
		valueobject.SeverityLow,
		valueobject.SeverityInfo,
	}

	for _, severity := range validSeverities {
		if !severity.IsValid() {
			t.Errorf("expected '%s' to be valid", severity)
		}
	}

	invalid := valueobject.Severity("unknown")
	if invalid.IsValid() {
		t.Error("expected 'unknown' to be invalid")
	}
}

func TestCategory_IsValid(t *testing.T) {
	for _, category := range valueobject.AllCategories() {
		if !category.IsValid() {
			t.Errorf("expected '%s' to be valid", category)
		}
	}
}

func TestCategory_Label(t *testing.T) {
	label := valueobject.CategoryTechnical.Label()
	if label != "Technical SEO" {
		t.Errorf("expected 'Technical SEO', got '%s'", label)
	}
}

func TestAuditRule_PassFailWarn(t *testing.T) {
	rule := valueobject.NewAuditRule("test_rule", valueobject.CategoryTechnical, valueobject.SeverityHigh)

	if rule.Result != valueobject.RuleResultSkipped {
		t.Errorf("new rule should be skipped, got '%s'", rule.Result)
	}

	rule.Pass("all good")
	if !rule.IsPassing() {
		t.Error("expected rule to be passing")
	}

	rule.Fail("broken", "fix it")
	if !rule.IsFailing() {
		t.Error("expected rule to be failing")
	}
	if rule.Recommendation != "fix it" {
		t.Errorf("expected recommendation 'fix it', got '%s'", rule.Recommendation)
	}

	rule.Warn("maybe broken", "check it")
	if !rule.IsWarning() {
		t.Error("expected rule to be warning")
	}
}

func TestSeoScore_Calculate(t *testing.T) {
	rules := []valueobject.AuditRule{
		{Key: "title_exists", Category: valueobject.CategoryOnPage, Severity: valueobject.SeverityCritical, Result: valueobject.RuleResultPass},
		{Key: "title_length", Category: valueobject.CategoryOnPage, Severity: valueobject.SeverityMedium, Result: valueobject.RuleResultFail},
		{Key: "robots_txt", Category: valueobject.CategoryTechnical, Severity: valueobject.SeverityHigh, Result: valueobject.RuleResultPass},
		{Key: "geo_test", Category: valueobject.CategoryGEO, Severity: valueobject.SeverityHigh, Result: valueobject.RuleResultFail},
	}

	score := valueobject.NewSeoScore()
	score.Calculate(rules)

	if score.TotalRules != 3 {
		t.Errorf("expected 3 total rules (GEO excluded), got %d", score.TotalRules)
	}

	if score.PassedRules != 2 {
		t.Errorf("expected 2 passed, got %d", score.PassedRules)
	}

	if score.FailedRules != 1 {
		t.Errorf("expected 1 failed, got %d", score.FailedRules)
	}

	if score.Overall <= 0 {
		t.Error("expected positive overall score")
	}
}

func TestGeoScore_Calculate(t *testing.T) {
	rules := []valueobject.AuditRule{
		{Key: "geo_crawler_oai", Category: valueobject.CategoryGEO, Severity: valueobject.SeverityHigh, Result: valueobject.RuleResultPass},
		{Key: "geo_crawler_perplexity", Category: valueobject.CategoryGEO, Severity: valueobject.SeverityHigh, Result: valueobject.RuleResultFail},
		{Key: "geo_llms_txt_exists", Category: valueobject.CategoryGEO, Severity: valueobject.SeverityMedium, Result: valueobject.RuleResultPass},
		{Key: "title_exists", Category: valueobject.CategoryOnPage, Severity: valueobject.SeverityHigh, Result: valueobject.RuleResultPass},
	}

	score := valueobject.NewGeoScore()
	score.Calculate(rules)

	if score.TotalRules != 3 {
		t.Errorf("expected 3 GEO rules, got %d", score.TotalRules)
	}

	if score.PassedRules != 2 {
		t.Errorf("expected 2 passed, got %d", score.PassedRules)
	}

	if score.Overall <= 0 {
		t.Error("expected positive GEO score")
	}
}

func TestSeoScore_Grade(t *testing.T) {
	score := valueobject.NewSeoScore()

	score.Overall = 95
	if score.Grade() != "A" {
		t.Errorf("expected A, got %s", score.Grade())
	}

	score.Overall = 85
	if score.Grade() != "B" {
		t.Errorf("expected B, got %s", score.Grade())
	}

	score.Overall = 55
	if score.Grade() != "F" {
		t.Errorf("expected F, got %s", score.Grade())
	}
}
