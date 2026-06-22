package entity

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
)

type CrawlStatus string

const (
	CrawlStatusQueued    CrawlStatus = "queued"
	CrawlStatusCrawling  CrawlStatus = "crawling"
	CrawlStatusAnalyzing CrawlStatus = "analyzing"
	CrawlStatusCompleted CrawlStatus = "completed"
	CrawlStatusFailed    CrawlStatus = "failed"
)

type CrawlJob struct {
	ID             string
	Domain         string
	Status         CrawlStatus
	MaxPages       int
	PagesCrawled   int
	IssuesFound    int
	SeoScore       valueobject.SeoScore
	GeoScore       valueobject.GeoScore
	ErrorMessage   string
	SelectedRules  []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CompletedAt    *time.Time
}

func NewCrawlJob(domain string, maxPages int, selectedRules []string) (*CrawlJob, error) {
	job := &CrawlJob{
		ID:            uuid.New().String(),
		Domain:        domain,
		Status:        CrawlStatusQueued,
		MaxPages:      maxPages,
		SelectedRules: selectedRules,
		SeoScore:      valueobject.NewSeoScore(),
		GeoScore:      valueobject.NewGeoScore(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := job.Validate(); err != nil {
		return nil, err
	}

	return job, nil
}

func (j *CrawlJob) Validate() error {
	if j.ID == "" {
		return fmt.Errorf("crawl job ID is required")
	}

	if j.Domain == "" {
		return fmt.Errorf("domain is required")
	}

	parsedURL, err := url.Parse(j.Domain)
	if err != nil {
		return fmt.Errorf("invalid domain: %w", err)
	}

	if parsedURL.Host == "" && parsedURL.Path == "" {
		return fmt.Errorf("domain must have a valid host")
	}

	if j.MaxPages < 0 {
		return fmt.Errorf("max pages must be zero (unlimited) or positive")
	}

	return nil
}

func (j *CrawlJob) StartCrawling() {
	j.Status = CrawlStatusCrawling
	j.UpdatedAt = time.Now().UTC()
}

func (j *CrawlJob) StartAnalyzing() {
	j.Status = CrawlStatusAnalyzing
	j.UpdatedAt = time.Now().UTC()
}

func (j *CrawlJob) Complete(seoScore valueobject.SeoScore, geoScore valueobject.GeoScore, issuesFound int) {
	now := time.Now().UTC()
	j.Status = CrawlStatusCompleted
	j.SeoScore = seoScore
	j.GeoScore = geoScore
	j.IssuesFound = issuesFound
	j.CompletedAt = &now
	j.UpdatedAt = now
}

func (j *CrawlJob) Fail(errorMessage string) {
	now := time.Now().UTC()
	j.Status = CrawlStatusFailed
	j.ErrorMessage = errorMessage
	j.CompletedAt = &now
	j.UpdatedAt = now
}

func (j *CrawlJob) IncrementPagesCrawled() {
	j.PagesCrawled++
	j.UpdatedAt = time.Now().UTC()
}

func (j *CrawlJob) HasRuleFilter() bool {
	return len(j.SelectedRules) > 0
}

func (j *CrawlJob) IsRuleSelected(ruleKey string) bool {
	if !j.HasRuleFilter() {
		return true
	}
	for _, selected := range j.SelectedRules {
		if selected == ruleKey {
			return true
		}
	}
	return false
}

func (j *CrawlJob) NormalizedDomain() string {
	parsed, err := url.Parse(j.Domain)
	if err != nil {
		return j.Domain
	}
	if parsed.Scheme == "" {
		parsed.Scheme = "https"
	}
	if parsed.Host == "" {
		parsed.Host = parsed.Path
		parsed.Path = ""
	}
	return parsed.String()
}
