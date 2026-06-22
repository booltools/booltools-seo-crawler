package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

const (
	ExitCodePass          = 0
	ExitCodeFail          = 1
	ExitCodeConfigError   = 2
)

type ReportOutput struct {
	Mode       string              `json:"mode"`
	TotalPages int                 `json:"totalPages"`
	TotalRules int                 `json:"totalRules"`
	Passed     int                 `json:"passed"`
	Warnings   int                 `json:"warnings"`
	Failures   int                 `json:"failures"`
	Issues     []ReportIssue       `json:"issues"`
	ExitCode   int                 `json:"exitCode"`
	Verdict    string              `json:"verdict"`
}

type ReportIssue struct {
	RuleKey        string `json:"ruleKey"`
	Severity       string `json:"severity"`
	Category       string `json:"category"`
	Message        string `json:"message"`
	Recommendation string `json:"recommendation"`
	AffectedCount  int    `json:"affectedCount"`
	AffectedURLs   []string `json:"affectedUrls,omitempty"`
}

type Reporter struct {
	writer io.Writer
}

func NewReporter(writer io.Writer) *Reporter {
	return &Reporter{writer: writer}
}

func (reporter *Reporter) Report(result *ScanResult, config Config) int {
	rules := FilterRules(result.AllRules, config.Ignore, config.Only)

	grouped := groupRulesByKey(rules)
	exitCode := determineExitCode(grouped, config.FailOn)

	if config.Format == "json" {
		reporter.reportJSON(result, grouped, exitCode, config)
	} else {
		reporter.reportText(result, grouped, exitCode, config)
	}

	if config.Output != "" {
		reporter.writeToFile(result, grouped, exitCode, config)
	}

	return exitCode
}

func (reporter *Reporter) reportText(result *ScanResult, grouped []groupedRule, exitCode int, config Config) {
	writer := reporter.writer

	fmt.Fprintf(writer, "\n%s%sSEO/GEO Check%s — %s mode\n", colorBold, colorCyan, colorReset, result.Mode)
	fmt.Fprintf(writer, "Scanned %d pages", result.TotalPages)
	if config.Mode == "static" {
		fmt.Fprintf(writer, " from %s", config.Dir)
	} else {
		fmt.Fprintf(writer, " at %s", config.URL)
	}
	fmt.Fprintln(writer)
	fmt.Fprintln(writer)

	severityOrder := []string{"critical", "high", "medium", "low", "info"}
	severityColors := map[string]string{
		"critical": colorRed,
		"high":     colorRed,
		"medium":   colorYellow,
		"low":      colorCyan,
		"info":     colorGray,
	}

	totalIssues := 0
	for _, severity := range severityOrder {
		issues := filterBySeverity(grouped, severity)
		if len(issues) == 0 {
			continue
		}

		totalIssues += len(issues)
		color := severityColors[severity]
		fmt.Fprintf(writer, "%s%s%s (%d issue%s)%s\n",
			colorBold, color, strings.ToUpper(severity), len(issues), pluralize(len(issues)), colorReset)

		for _, issue := range issues {
			fmt.Fprintf(writer, "  %s[%s]%s %s\n", colorGray, issue.ruleKey, colorReset, issue.message)
			if issue.affectedCount > 0 {
				fmt.Fprintf(writer, "    %s→ %d page%s affected%s\n", colorGray, issue.affectedCount, pluralize(issue.affectedCount), colorReset)
			}
		}
		fmt.Fprintln(writer)
	}

	passedCount := countByResult(result.AllRules, valueobject.RuleResultPass)
	fmt.Fprintf(writer, "%sPassed: %d%s  |  ", colorGreen, passedCount, colorReset)
	fmt.Fprintf(writer, "%sIssues: %d%s\n", colorYellow, totalIssues, colorReset)

	if exitCode == ExitCodePass {
		fmt.Fprintf(writer, "\n%s%sPASS%s — no issues at or above '%s' severity\n\n", colorBold, colorGreen, colorReset, config.FailOn)
	} else {
		fmt.Fprintf(writer, "\n%s%sFAIL%s — issues found at or above '%s' severity\n\n", colorBold, colorRed, colorReset, config.FailOn)
	}
}

func (reporter *Reporter) reportJSON(result *ScanResult, grouped []groupedRule, exitCode int, config Config) {
	output := buildReportOutput(result, grouped, exitCode)
	encoder := json.NewEncoder(reporter.writer)
	encoder.SetIndent("", "  ")
	encoder.Encode(output)
}

func (reporter *Reporter) writeToFile(result *ScanResult, grouped []groupedRule, exitCode int, config Config) {
	file, err := os.Create(config.Output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not write output file: %v\n", err)
		return
	}
	defer file.Close()

	output := buildReportOutput(result, grouped, exitCode)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(output)
	fmt.Fprintf(os.Stderr, "Report written to %s\n", config.Output)
}

func buildReportOutput(result *ScanResult, grouped []groupedRule, exitCode int) ReportOutput {
	passedCount := countByResult(result.AllRules, valueobject.RuleResultPass)
	warningCount := countByResult(result.AllRules, valueobject.RuleResultWarning)
	failCount := countByResult(result.AllRules, valueobject.RuleResultFail)

	verdict := "PASS"
	if exitCode != ExitCodePass {
		verdict = "FAIL"
	}

	issues := make([]ReportIssue, 0, len(grouped))
	for _, group := range grouped {
		if group.result != valueobject.RuleResultFail && group.result != valueobject.RuleResultWarning {
			continue
		}
		issues = append(issues, ReportIssue{
			RuleKey:        group.ruleKey,
			Severity:       group.severity,
			Category:       group.category,
			Message:        group.message,
			Recommendation: group.recommendation,
			AffectedCount:  group.affectedCount,
			AffectedURLs:   group.affectedURLs,
		})
	}

	return ReportOutput{
		Mode:       result.Mode,
		TotalPages: result.TotalPages,
		TotalRules: len(result.AllRules),
		Passed:     passedCount,
		Warnings:   warningCount,
		Failures:   failCount,
		Issues:     issues,
		ExitCode:   exitCode,
		Verdict:    verdict,
	}
}

type groupedRule struct {
	ruleKey        string
	severity       string
	category       string
	result         valueobject.RuleResult
	message        string
	recommendation string
	affectedCount  int
	affectedURLs   []string
}

func groupRulesByKey(rules []valueobject.AuditRule) []groupedRule {
	ruleMap := make(map[string]*groupedRule)
	var orderedKeys []string

	for _, rule := range rules {
		existing, found := ruleMap[rule.Key]
		if found {
			if rule.AffectedURL != "" && !containsString(existing.affectedURLs, rule.AffectedURL) {
				existing.affectedURLs = append(existing.affectedURLs, rule.AffectedURL)
				existing.affectedCount++
			}
			if severityRank(string(rule.Severity)) > severityRank(existing.severity) {
				existing.severity = string(rule.Severity)
			}
			if rule.Result == valueobject.RuleResultFail {
				existing.result = valueobject.RuleResultFail
			} else if rule.Result == valueobject.RuleResultWarning && existing.result != valueobject.RuleResultFail {
				existing.result = valueobject.RuleResultWarning
			}
		} else {
			entry := &groupedRule{
				ruleKey:        rule.Key,
				severity:       string(rule.Severity),
				category:       string(rule.Category),
				result:         rule.Result,
				message:        rule.Message,
				recommendation: rule.Recommendation,
				affectedCount:  0,
				affectedURLs:   make([]string, 0),
			}
			if rule.AffectedURL != "" {
				entry.affectedURLs = append(entry.affectedURLs, rule.AffectedURL)
				entry.affectedCount = 1
			}
			ruleMap[rule.Key] = entry
			orderedKeys = append(orderedKeys, rule.Key)
		}
	}

	result := make([]groupedRule, 0, len(orderedKeys))
	for _, key := range orderedKeys {
		result = append(result, *ruleMap[key])
	}

	sort.Slice(result, func(i, j int) bool {
		rankI := severityRank(result[i].severity)
		rankJ := severityRank(result[j].severity)
		if rankI != rankJ {
			return rankI > rankJ
		}
		return result[i].affectedCount > result[j].affectedCount
	})

	return result
}

func filterBySeverity(grouped []groupedRule, severity string) []groupedRule {
	var filtered []groupedRule
	for _, group := range grouped {
		if group.severity == severity && (group.result == valueobject.RuleResultFail || group.result == valueobject.RuleResultWarning) {
			filtered = append(filtered, group)
		}
	}
	return filtered
}

func determineExitCode(grouped []groupedRule, failOnSeverity string) int {
	threshold := severityRank(failOnSeverity)
	for _, group := range grouped {
		if group.result != valueobject.RuleResultFail && group.result != valueobject.RuleResultWarning {
			continue
		}
		if severityRank(group.severity) >= threshold {
			return ExitCodeFail
		}
	}
	return ExitCodePass
}

func severityRank(severity string) int {
	ranks := map[string]int{
		"info":     0,
		"low":      1,
		"medium":   2,
		"high":     3,
		"critical": 4,
	}
	if rank, exists := ranks[severity]; exists {
		return rank
	}
	return -1
}

func countByResult(rules []valueobject.AuditRule, result valueobject.RuleResult) int {
	count := 0
	for _, rule := range rules {
		if rule.Result == result {
			count++
		}
	}
	return count
}

func containsString(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
