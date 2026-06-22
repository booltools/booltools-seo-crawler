package valueobject

import "fmt"

type GeoScore struct {
	Overall        float64
	CrawlerAccess  float64
	LlmsTxt        float64
	Citability     float64
	EntityAuth     float64
	AIFriendly     float64
	TotalRules     int
	PassedRules    int
	FailedRules    int
	WarningRules   int
}

func NewGeoScore() GeoScore {
	return GeoScore{}
}

var geoSubcategoryWeights = map[string]float64{
	"crawler_access":    0.25,
	"llms_txt":          0.15,
	"citability":        0.25,
	"entity_authority":  0.15,
	"ai_friendly":       0.20,
}

func (g *GeoScore) Calculate(rules []AuditRule) {
	type ruleScore struct {
		maxPoints float64
		earned    float64
	}

	subcategoryRules := map[string][]ruleScore{
		"crawler_access":   {},
		"llms_txt":         {},
		"citability":       {},
		"entity_authority": {},
		"ai_friendly":      {},
	}

	for _, rule := range rules {
		if rule.Category != CategoryGEO {
			continue
		}
		if rule.Result == RuleResultSkipped {
			continue
		}

		g.TotalRules++
		subcat := classifyGeoRule(rule.Key)
		severityWeight := float64(rule.Severity.Weight())
		if severityWeight == 0 {
			severityWeight = 1
		}

		var earned float64
		switch rule.Result {
		case RuleResultPass:
			g.PassedRules++
			earned = severityWeight
		case RuleResultFail:
			g.FailedRules++
			earned = 0
		case RuleResultWarning:
			g.WarningRules++
			earned = severityWeight * 0.5
		}

		subcategoryRules[subcat] = append(subcategoryRules[subcat], ruleScore{
			maxPoints: severityWeight,
			earned:    earned,
		})
	}

	weightedSum := 0.0
	totalWeight := 0.0

	for subcat, scores := range subcategoryRules {
		if len(scores) == 0 {
			continue
		}

		var maxPoints, earnedPoints float64
		for _, score := range scores {
			maxPoints += score.maxPoints
			earnedPoints += score.earned
		}

		var subcatScore float64
		if maxPoints > 0 {
			subcatScore = (earnedPoints / maxPoints) * 100
		}

		weight := geoSubcategoryWeights[subcat]
		weightedSum += subcatScore * weight
		totalWeight += weight

		switch subcat {
		case "crawler_access":
			g.CrawlerAccess = subcatScore
		case "llms_txt":
			g.LlmsTxt = subcatScore
		case "citability":
			g.Citability = subcatScore
		case "entity_authority":
			g.EntityAuth = subcatScore
		case "ai_friendly":
			g.AIFriendly = subcatScore
		}
	}

	if totalWeight > 0 {
		g.Overall = weightedSum / totalWeight
	}
}

func (g GeoScore) Validate() error {
	if g.Overall < 0 || g.Overall > 100 {
		return fmt.Errorf("GEO score must be between 0 and 100, got %.2f", g.Overall)
	}
	return nil
}

func (g GeoScore) Grade() string {
	switch {
	case g.Overall >= 90:
		return "A"
	case g.Overall >= 80:
		return "B"
	case g.Overall >= 70:
		return "C"
	case g.Overall >= 60:
		return "D"
	default:
		return "F"
	}
}

func classifyGeoRule(key string) string {
	prefixes := map[string]string{
		"geo_crawler_":   "crawler_access",
		"geo_llms_":      "llms_txt",
		"geo_citability_": "citability",
		"geo_entity_":    "entity_authority",
		"geo_ai_":        "ai_friendly",
	}
	for prefix, subcat := range prefixes {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			return subcat
		}
	}
	return "ai_friendly"
}
