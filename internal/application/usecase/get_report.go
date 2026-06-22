package usecase

import (
	"context"
	"fmt"

	"github.com/MarceloBD/free-seo-crawler/internal/application/dto"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/repository"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
)

type GetReportUseCase struct {
	crawlJobRepo  repository.CrawlJobRepository
	pageAuditRepo repository.PageAuditRepository
}

func NewGetReportUseCase(
	crawlJobRepo repository.CrawlJobRepository,
	pageAuditRepo repository.PageAuditRepository,
) *GetReportUseCase {
	return &GetReportUseCase{
		crawlJobRepo:  crawlJobRepo,
		pageAuditRepo: pageAuditRepo,
	}
}

func (uc *GetReportUseCase) Execute(ctx context.Context, jobID string) (*dto.AuditReport, error) {
	job, err := uc.crawlJobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get crawl job: %w", err)
	}
	if job == nil {
		return nil, fmt.Errorf("crawl job not found")
	}

	pageAudits, err := uc.pageAuditRepo.GetByCrawlJobID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get page audits: %w", err)
	}

	report := &dto.AuditReport{
		Job: dto.CrawlResponse{
			ID:           job.ID,
			Domain:       job.Domain,
			Status:       string(job.Status),
			MaxPages:     job.MaxPages,
			PagesCrawled: job.PagesCrawled,
			IssuesFound:  job.IssuesFound,
			ErrorMessage: job.ErrorMessage,
			CreatedAt:    job.CreatedAt,
			CompletedAt:  job.CompletedAt,
		},
		SeoScore: dto.ScoreDetail{
			Overall:        job.SeoScore.Overall,
			Grade:          job.SeoScore.Grade(),
			CategoryScores: convertCategoryScores(job.SeoScore.CategoryScores),
			TotalRules:     job.SeoScore.TotalRules,
			PassedRules:    job.SeoScore.PassedRules,
			FailedRules:    job.SeoScore.FailedRules,
			WarningRules:   job.SeoScore.WarningRules,
			CriticalIssues: job.SeoScore.CriticalIssues,
		},
		GeoScore: dto.GeoScoreDetail{
			Overall:       job.GeoScore.Overall,
			Grade:         job.GeoScore.Grade(),
			CrawlerAccess: job.GeoScore.CrawlerAccess,
			LlmsTxt:       job.GeoScore.LlmsTxt,
			Citability:    job.GeoScore.Citability,
			EntityAuth:    job.GeoScore.EntityAuth,
			AIFriendly:    job.GeoScore.AIFriendly,
			TotalRules:    job.GeoScore.TotalRules,
			PassedRules:   job.GeoScore.PassedRules,
			FailedRules:   job.GeoScore.FailedRules,
		},
		IssuesBySeverity: make(map[string][]dto.IssueItem),
		LinksToChange:    make([]dto.LinkChange, 0),
	}

	linkChangesMap := make(map[string]*dto.LinkChange)

	for _, audit := range pageAudits {
		pageReport := dto.PageReport{
			URL:        audit.URL,
			StatusCode: audit.StatusCode,
			Depth:      audit.Depth,
			Issues:     make([]dto.IssueItem, 0),
		}

		for _, rule := range audit.Rules {
			item := dto.IssueItem{
				RuleKey:        rule.Key,
				Category:       string(rule.Category),
				CategoryLabel:  rule.Category.Label(),
				Severity:       string(rule.Severity),
				Result:         string(rule.Result),
				Message:        rule.Message,
				Recommendation: rule.Recommendation,
				AffectedURL:    rule.AffectedURL,
				Details:        rule.Details,
			}

			switch rule.Result {
			case valueobject.RuleResultFail:
				pageReport.Failures++
				pageReport.Issues = append(pageReport.Issues, item)
				report.IssuesBySeverity[string(rule.Severity)] = append(
					report.IssuesBySeverity[string(rule.Severity)], item,
				)
			case valueobject.RuleResultWarning:
				pageReport.Warnings++
				pageReport.Issues = append(pageReport.Issues, item)
				report.IssuesBySeverity[string(rule.Severity)] = append(
					report.IssuesBySeverity[string(rule.Severity)], item,
				)
			case valueobject.RuleResultPass:
				pageReport.Passes++
			}

			if (rule.IsFailing() || rule.IsWarning()) && rule.AffectedURL != "" {
				affectedURL := rule.AffectedURL
				linkChange, exists := linkChangesMap[affectedURL]
				if !exists {
					linkChange = &dto.LinkChange{
						URL:     affectedURL,
						Changes: make([]dto.ChangeItem, 0),
					}
					linkChangesMap[affectedURL] = linkChange
				}
				linkChange.Changes = append(linkChange.Changes, dto.ChangeItem{
					Category:       rule.Category.Label(),
					Severity:       string(rule.Severity),
					Issue:          rule.Message,
					Recommendation: rule.Recommendation,
				})
			}
		}

		report.Pages = append(report.Pages, pageReport)
	}

	for _, linkChange := range linkChangesMap {
		report.LinksToChange = append(report.LinksToChange, *linkChange)
	}

	report.Summary = buildSummary(report)

	return report, nil
}

func convertCategoryScores(scores map[valueobject.Category]float64) map[string]float64 {
	result := make(map[string]float64)
	for category, score := range scores {
		result[string(category)] = score
	}
	return result
}

func buildSummary(report *dto.AuditReport) dto.ReportSummary {
	summary := dto.ReportSummary{
		TotalPages: len(report.Pages),
	}

	for severity, issues := range report.IssuesBySeverity {
		count := len(issues)
		summary.TotalIssues += count
		switch severity {
		case "critical":
			summary.CriticalIssues = count
		case "high":
			summary.HighIssues = count
		case "medium":
			summary.MediumIssues = count
		case "low":
			summary.LowIssues = count
		case "info":
			summary.InfoIssues = count
		}
	}

	for _, page := range report.Pages {
		summary.PassedChecks += page.Passes
	}

	return summary
}
