package valueobject

import "fmt"

type RuleResult string

const (
	RuleResultPass    RuleResult = "pass"
	RuleResultFail    RuleResult = "fail"
	RuleResultWarning RuleResult = "warning"
	RuleResultSkipped RuleResult = "skipped"
)

type AuditRule struct {
	Key            string
	Category       Category
	Severity       Severity
	Result         RuleResult
	Message        string
	Recommendation string
	AffectedURL    string
	Details        string
}

func NewAuditRule(key string, category Category, severity Severity) AuditRule {
	return AuditRule{
		Key:      key,
		Category: category,
		Severity: severity,
		Result:   RuleResultSkipped,
	}
}

func (r *AuditRule) Pass(message string) {
	r.Result = RuleResultPass
	r.Message = message
}

func (r *AuditRule) Fail(message string, recommendation string) {
	r.Result = RuleResultFail
	r.Message = message
	r.Recommendation = recommendation
}

func (r *AuditRule) Warn(message string, recommendation string) {
	r.Result = RuleResultWarning
	r.Message = message
	r.Recommendation = recommendation
}

func (r *AuditRule) Skip(reason string) {
	r.Result = RuleResultSkipped
	r.Message = reason
}

func (r *AuditRule) WithURL(url string) *AuditRule {
	r.AffectedURL = url
	return r
}

func (r *AuditRule) WithDetails(details string) *AuditRule {
	r.Details = details
	return r
}

func (r AuditRule) Validate() error {
	if r.Key == "" {
		return fmt.Errorf("audit rule key is required")
	}
	if !r.Category.IsValid() {
		return fmt.Errorf("invalid category: %s", r.Category)
	}
	if !r.Severity.IsValid() {
		return fmt.Errorf("invalid severity: %s", r.Severity)
	}
	return nil
}

func (r AuditRule) IsFailing() bool {
	return r.Result == RuleResultFail
}

func (r AuditRule) IsWarning() bool {
	return r.Result == RuleResultWarning
}

func (r AuditRule) IsPassing() bool {
	return r.Result == RuleResultPass
}
