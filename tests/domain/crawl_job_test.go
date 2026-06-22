package domain_test

import (
	"testing"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/entity"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
)

func TestNewCrawlJob_ValidDomain(t *testing.T) {
	job, err := entity.NewCrawlJob("https://example.com", 50, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if job.ID == "" {
		t.Error("expected job ID to be set")
	}

	if job.Domain != "https://example.com" {
		t.Errorf("expected domain 'https://example.com', got '%s'", job.Domain)
	}

	if job.Status != entity.CrawlStatusQueued {
		t.Errorf("expected status 'queued', got '%s'", job.Status)
	}

	if job.MaxPages != 50 {
		t.Errorf("expected maxPages 50, got %d", job.MaxPages)
	}
}

func TestNewCrawlJob_EmptyDomain(t *testing.T) {
	_, err := entity.NewCrawlJob("", 50, nil)
	if err == nil {
		t.Error("expected error for empty domain")
	}
}

func TestNewCrawlJob_InvalidMaxPages(t *testing.T) {
	_, err := entity.NewCrawlJob("https://example.com", -1, nil)
	if err == nil {
		t.Error("expected error for negative maxPages")
	}
}

func TestNewCrawlJob_UnlimitedMaxPages(t *testing.T) {
	job, err := entity.NewCrawlJob("https://example.com", 0, nil)
	if err != nil {
		t.Errorf("expected no error for maxPages 0 (unlimited), got %v", err)
	}
	if job.MaxPages != 0 {
		t.Errorf("expected maxPages 0, got %d", job.MaxPages)
	}
}

func TestNewCrawlJob_SelectedRules(t *testing.T) {
	rules := []string{"title_exists", "h1_count"}
	job, err := entity.NewCrawlJob("https://example.com", 50, rules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !job.IsRuleSelected("title_exists") {
		t.Error("expected title_exists to be selected")
	}
	if !job.IsRuleSelected("h1_count") {
		t.Error("expected h1_count to be selected")
	}
	if job.IsRuleSelected("meta_description_exists") {
		t.Error("expected meta_description_exists to NOT be selected")
	}
}

func TestNewCrawlJob_NoFilterSelectsAll(t *testing.T) {
	job, err := entity.NewCrawlJob("https://example.com", 50, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !job.IsRuleSelected("title_exists") {
		t.Error("expected all rules to be selected when no filter is set")
	}
	if !job.IsRuleSelected("anything") {
		t.Error("expected all rules to be selected when no filter is set")
	}
}

func TestCrawlJob_StatusTransitions(t *testing.T) {
	job, _ := entity.NewCrawlJob("https://example.com", 50, nil)

	job.StartCrawling()
	if job.Status != entity.CrawlStatusCrawling {
		t.Errorf("expected 'crawling', got '%s'", job.Status)
	}

	job.StartAnalyzing()
	if job.Status != entity.CrawlStatusAnalyzing {
		t.Errorf("expected 'analyzing', got '%s'", job.Status)
	}

	seoScore := valueobject.NewSeoScore()
	geoScore := valueobject.NewGeoScore()
	job.Complete(seoScore, geoScore, 10)

	if job.Status != entity.CrawlStatusCompleted {
		t.Errorf("expected 'completed', got '%s'", job.Status)
	}

	if job.CompletedAt == nil {
		t.Error("expected completedAt to be set")
	}

	if job.IssuesFound != 10 {
		t.Errorf("expected 10 issues, got %d", job.IssuesFound)
	}
}

func TestCrawlJob_NormalizedDomain(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"example.com", "https://example.com"},
		{"https://example.com", "https://example.com"},
		{"http://example.com", "http://example.com"},
	}

	for _, testCase := range tests {
		job, _ := entity.NewCrawlJob(testCase.input, 10, nil)
		normalized := job.NormalizedDomain()
		if normalized != testCase.expected {
			t.Errorf("NormalizedDomain(%s) = %s, want %s", testCase.input, normalized, testCase.expected)
		}
	}
}
