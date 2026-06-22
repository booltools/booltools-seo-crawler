package valueobject

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

func (s Severity) IsValid() bool {
	switch s {
	case SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo:
		return true
	}
	return false
}

func (s Severity) Weight() int {
	switch s {
	case SeverityCritical:
		return 10
	case SeverityHigh:
		return 7
	case SeverityMedium:
		return 4
	case SeverityLow:
		return 2
	case SeverityInfo:
		return 0
	}
	return 0
}
