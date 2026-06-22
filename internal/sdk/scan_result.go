package sdk

import "github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"

type ScanResult struct {
	Mode       string
	TotalPages int
	Pages      []PageScanResult
	AllRules   []valueobject.AuditRule
}

type PageScanResult struct {
	URL      string
	Issues   []valueobject.AuditRule
	Passes   int
	Warnings int
	Failures int
}

func buildPageScanResult(pageURL string, rules []valueobject.AuditRule) PageScanResult {
	result := PageScanResult{URL: pageURL}
	for _, rule := range rules {
		switch rule.Result {
		case valueobject.RuleResultPass:
			result.Passes++
		case valueobject.RuleResultWarning:
			result.Warnings++
			result.Issues = append(result.Issues, rule)
		case valueobject.RuleResultFail:
			result.Failures++
			result.Issues = append(result.Issues, rule)
		}
	}
	return result
}
