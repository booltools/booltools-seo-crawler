package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/MarceloBD/free-seo-crawler/internal/application/dto"
	"github.com/MarceloBD/free-seo-crawler/internal/application/usecase"
)

type ExportHandler struct {
	getReport *usecase.GetReportUseCase
}

func NewExportHandler(getReport *usecase.GetReportUseCase) *ExportHandler {
	return &ExportHandler{getReport: getReport}
}

type affectedEntry struct {
	URL     string
	Details string
}

type groupedIssue struct {
	RuleKey        string
	Category       string
	Severity       string
	Result         string
	Message        string
	Recommendation string
	Affected       []affectedEntry
}

func groupIssuesByRule(report *dto.AuditReport) []groupedIssue {
	ruleMap := make(map[string]*groupedIssue)
	var orderedKeys []string

	for _, page := range report.Pages {
		for _, issue := range page.Issues {
			key := issue.RuleKey
			existing, found := ruleMap[key]
			if found {
				if issue.AffectedURL != "" && !containsAffectedURL(existing.Affected, issue.AffectedURL) {
					existing.Affected = append(existing.Affected, affectedEntry{
						URL:     issue.AffectedURL,
						Details: issue.Details,
					})
				}
			} else {
				affected := make([]affectedEntry, 0)
				if issue.AffectedURL != "" {
					affected = append(affected, affectedEntry{
						URL:     issue.AffectedURL,
						Details: issue.Details,
					})
				}
				ruleMap[key] = &groupedIssue{
					RuleKey:        issue.RuleKey,
					Category:       issue.CategoryLabel,
					Severity:       issue.Severity,
					Result:         issue.Result,
					Message:        issue.Message,
					Recommendation: issue.Recommendation,
					Affected:       affected,
				}
				orderedKeys = append(orderedKeys, key)
			}
		}
	}

	severityRank := map[string]int{
		"critical": 0,
		"high":     1,
		"medium":   2,
		"low":      3,
		"info":     4,
	}

	sort.SliceStable(orderedKeys, func(i, j int) bool {
		a := ruleMap[orderedKeys[i]]
		b := ruleMap[orderedKeys[j]]
		rankA := severityRank[a.Severity]
		rankB := severityRank[b.Severity]
		if rankA != rankB {
			return rankA < rankB
		}
		return len(a.Affected) > len(b.Affected)
	})

	result := make([]groupedIssue, 0, len(orderedKeys))
	for _, key := range orderedKeys {
		result = append(result, *ruleMap[key])
	}
	return result
}

func containsAffectedURL(entries []affectedEntry, targetURL string) bool {
	for _, entry := range entries {
		if entry.URL == targetURL {
			return true
		}
	}
	return false
}

func parseExcludedRules(request *http.Request) map[string]bool {
	excludeParam := request.URL.Query().Get("exclude")
	if excludeParam == "" {
		return nil
	}
	excluded := make(map[string]bool)
	for _, key := range strings.Split(excludeParam, ",") {
		trimmed := strings.TrimSpace(key)
		if trimmed != "" {
			excluded[trimmed] = true
		}
	}
	return excluded
}

func filterGroupedIssues(groups []groupedIssue, excluded map[string]bool) []groupedIssue {
	if len(excluded) == 0 {
		return groups
	}
	filtered := make([]groupedIssue, 0, len(groups))
	for _, group := range groups {
		if !excluded[group.RuleKey] {
			filtered = append(filtered, group)
		}
	}
	return filtered
}

func (h *ExportHandler) ExportCSV(writer http.ResponseWriter, request *http.Request) {
	report, err := h.resolveReport(writer, request)
	if err != nil {
		return
	}

	filename := sanitizeFilename(report.Job.Domain) + "-seo-audit.csv"
	writer.Header().Set("Content-Type", "text/csv; charset=utf-8")
	writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	csvWriter.Write([]string{
		"Severity", "Category", "Issue", "Affected URLs", "URL Count", "Recommendation", "Asset Details",
	})

	excluded := parseExcludedRules(request)
	grouped := filterGroupedIssues(groupIssuesByRule(report), excluded)
	for _, group := range grouped {
		var urls []string
		var detailParts []string
		for _, entry := range group.Affected {
			urls = append(urls, entry.URL)
			if entry.Details != "" {
				detailParts = append(detailParts, fmt.Sprintf("[%s] %s", entry.URL, strings.ReplaceAll(entry.Details, "\n", " ; ")))
			}
		}
		csvWriter.Write([]string{
			group.Severity,
			group.Category,
			group.Message,
			strings.Join(urls, " | "),
			fmt.Sprintf("%d", len(urls)),
			group.Recommendation,
			strings.Join(detailParts, " || "),
		})
	}
}

func (h *ExportHandler) ExportMarkdown(writer http.ResponseWriter, request *http.Request) {
	report, err := h.resolveReport(writer, request)
	if err != nil {
		return
	}

	filename := sanitizeFilename(report.Job.Domain) + "-ai-fix-instructions.md"
	writer.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("# SEO/GEO Audit Fix Instructions for %s\n\n", report.Job.Domain))
	builder.WriteString("Use this file to instruct an AI agent to fix all remaining SEO and GEO issues.\n\n")

	builder.WriteString("## Summary\n\n")
	builder.WriteString(fmt.Sprintf("- **Domain**: %s\n", report.Job.Domain))
	builder.WriteString(fmt.Sprintf("- **Pages Crawled**: %d\n", report.Summary.TotalPages))
	builder.WriteString(fmt.Sprintf("- **SEO Score**: %.0f/100 (Grade %s)\n", report.SeoScore.Overall, report.SeoScore.Grade))
	builder.WriteString(fmt.Sprintf("- **GEO Score**: %.0f/100 (Grade %s)\n", report.GeoScore.Overall, report.GeoScore.Grade))
	builder.WriteString(fmt.Sprintf("- **Total Issues**: %d\n", report.Summary.TotalIssues))
	builder.WriteString(fmt.Sprintf("- **Critical**: %d | **High**: %d | **Medium**: %d | **Low**: %d\n\n",
		report.Summary.CriticalIssues, report.Summary.HighIssues,
		report.Summary.MediumIssues, report.Summary.LowIssues))

	builder.WriteString("---\n\n")
	builder.WriteString("## Instructions for AI Agent\n\n")
	builder.WriteString("Fix each issue below in priority order (critical first, then high, medium, low).\n")
	builder.WriteString("Each issue is grouped with all affected URLs and their specific assets listed.\n\n")

	excludedMd := parseExcludedRules(request)
	grouped := filterGroupedIssues(groupIssuesByRule(report), excludedMd)

	currentSeverity := ""
	for _, group := range grouped {
		if group.Severity != currentSeverity {
			currentSeverity = group.Severity
			builder.WriteString(fmt.Sprintf("## %s Severity\n\n", capitalizeFirst(currentSeverity)))
		}

		builder.WriteString(fmt.Sprintf("### [%s] %s\n\n", group.Category, group.Message))
		builder.WriteString(fmt.Sprintf("- **Fix**: %s\n", group.Recommendation))
		builder.WriteString(fmt.Sprintf("- **Affected URLs** (%d):\n", len(group.Affected)))
		for _, entry := range group.Affected {
			builder.WriteString(fmt.Sprintf("  - %s\n", entry.URL))
			if entry.Details != "" {
				for _, line := range strings.Split(entry.Details, "\n") {
					trimmed := strings.TrimSpace(line)
					if trimmed != "" {
						builder.WriteString(fmt.Sprintf("    - `%s`\n", trimmed))
					}
				}
			}
		}
		builder.WriteString("\n")
	}

	writer.Write([]byte(builder.String()))
}

func (h *ExportHandler) resolveReport(writer http.ResponseWriter, request *http.Request) (*dto.AuditReport, error) {
	jobID := chi.URLParam(request, "id")
	if _, parseError := uuid.Parse(jobID); parseError != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "Invalid job ID format"})
		return nil, parseError
	}

	report, err := h.getReport.Execute(request.Context(), jobID)
	if err != nil {
		writeJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return nil, err
	}

	return report, nil
}

func capitalizeFirst(value string) string {
	if len(value) == 0 {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func sanitizeFilename(domain string) string {
	replacer := strings.NewReplacer(
		"https://", "",
		"http://", "",
		"/", "-",
		":", "",
		"?", "",
		"&", "",
		"\\", "",
		"..", "",
		"\n", "",
		"\r", "",
		"\"", "",
	)
	result := replacer.Replace(domain)
	if len(result) > 100 {
		result = result[:100]
	}
	return result
}
