package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
)

type PageAudit struct {
	ID         string
	CrawlJobID string
	URL        string
	StatusCode int
	Depth      int
	Rules      []valueobject.AuditRule
	CreatedAt  time.Time
}

func NewPageAudit(crawlJobID string, pageURL string, statusCode int, depth int) (*PageAudit, error) {
	audit := &PageAudit{
		ID:         uuid.New().String(),
		CrawlJobID: crawlJobID,
		URL:        pageURL,
		StatusCode: statusCode,
		Depth:      depth,
		Rules:      make([]valueobject.AuditRule, 0),
		CreatedAt:  time.Now().UTC(),
	}

	if err := audit.Validate(); err != nil {
		return nil, err
	}

	return audit, nil
}

func (p *PageAudit) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("page audit ID is required")
	}

	if p.CrawlJobID == "" {
		return fmt.Errorf("crawl job ID is required")
	}

	if p.URL == "" {
		return fmt.Errorf("page URL is required")
	}

	return nil
}

func (p *PageAudit) AddRule(rule valueobject.AuditRule) {
	p.Rules = append(p.Rules, rule)
}

func (p *PageAudit) AddRules(rules []valueobject.AuditRule) {
	p.Rules = append(p.Rules, rules...)
}

func (p *PageAudit) FailingRules() []valueobject.AuditRule {
	var failing []valueobject.AuditRule
	for _, rule := range p.Rules {
		if rule.IsFailing() {
			failing = append(failing, rule)
		}
	}
	return failing
}

func (p *PageAudit) WarningRules() []valueobject.AuditRule {
	var warnings []valueobject.AuditRule
	for _, rule := range p.Rules {
		if rule.IsWarning() {
			warnings = append(warnings, rule)
		}
	}
	return warnings
}

func (p *PageAudit) PassingRules() []valueobject.AuditRule {
	var passing []valueobject.AuditRule
	for _, rule := range p.Rules {
		if rule.IsPassing() {
			passing = append(passing, rule)
		}
	}
	return passing
}

func (p *PageAudit) IssueCount() int {
	count := 0
	for _, rule := range p.Rules {
		if rule.IsFailing() || rule.IsWarning() {
			count++
		}
	}
	return count
}
