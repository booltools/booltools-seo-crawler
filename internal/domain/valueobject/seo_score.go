package valueobject

import "fmt"

type SeoScore struct {
	Overall        float64
	CategoryScores map[Category]float64
	TotalRules     int
	PassedRules    int
	FailedRules    int
	WarningRules   int
	CriticalIssues int
}

func NewSeoScore() SeoScore {
	return SeoScore{
		CategoryScores: make(map[Category]float64),
	}
}

func (s *SeoScore) Calculate(rules []AuditRule) {
	type ruleScore struct {
		maxPoints float64
		earned    float64
	}

	categoryRules := make(map[Category][]ruleScore)

	for _, rule := range rules {
		if rule.Category == CategoryGEO {
			continue
		}
		if rule.Result == RuleResultSkipped {
			continue
		}

		s.TotalRules++
		severityWeight := float64(rule.Severity.Weight())
		if severityWeight == 0 {
			severityWeight = 1
		}

		var earned float64
		switch rule.Result {
		case RuleResultPass:
			s.PassedRules++
			earned = severityWeight
		case RuleResultFail:
			s.FailedRules++
			earned = 0
			if rule.Severity == SeverityCritical {
				s.CriticalIssues++
			}
		case RuleResultWarning:
			s.WarningRules++
			earned = severityWeight * 0.5
		}

		categoryRules[rule.Category] = append(categoryRules[rule.Category], ruleScore{
			maxPoints: severityWeight,
			earned:    earned,
		})
	}

	weightedSum := 0.0
	totalWeight := 0.0

	for category, scores := range categoryRules {
		var maxPoints, earnedPoints float64
		for _, score := range scores {
			maxPoints += score.maxPoints
			earnedPoints += score.earned
		}

		var categoryScore float64
		if maxPoints > 0 {
			categoryScore = (earnedPoints / maxPoints) * 100
		}
		s.CategoryScores[category] = categoryScore

		weight := CategoryWeights[category]
		weightedSum += categoryScore * weight
		totalWeight += weight
	}

	if totalWeight > 0 {
		s.Overall = weightedSum / totalWeight
	}

	if s.CriticalIssues > 0 {
		cap := 60.0 - float64(s.CriticalIssues-1)*5.0
		if cap < 20 {
			cap = 20
		}
		if s.Overall > cap {
			s.Overall = cap
		}
	}
}

func (s SeoScore) Validate() error {
	if s.Overall < 0 || s.Overall > 100 {
		return fmt.Errorf("SEO score must be between 0 and 100, got %.2f", s.Overall)
	}
	return nil
}

func (s SeoScore) Grade() string {
	switch {
	case s.Overall >= 90:
		return "A"
	case s.Overall >= 80:
		return "B"
	case s.Overall >= 70:
		return "C"
	case s.Overall >= 60:
		return "D"
	default:
		return "F"
	}
}
