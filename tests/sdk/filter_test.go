package sdk_test

import (
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/sdk"
)

func makeRules() []valueobject.AuditRule {
	return []valueobject.AuditRule{
		{Key: "title_exists", Category: valueobject.CategoryOnPage, Severity: valueobject.SeverityHigh, Result: valueobject.RuleResultFail},
		{Key: "h1_count", Category: valueobject.CategoryOnPage, Severity: valueobject.SeverityHigh, Result: valueobject.RuleResultPass},
		{Key: "og_locale", Category: valueobject.CategorySocial, Severity: valueobject.SeverityLow, Result: valueobject.RuleResultFail},
		{Key: "twitter_site", Category: valueobject.CategorySocial, Severity: valueobject.SeverityLow, Result: valueobject.RuleResultWarning},
		{Key: "mobile_viewport", Category: valueobject.CategoryMobile, Severity: valueobject.SeverityCritical, Result: valueobject.RuleResultPass},
	}
}

func TestFilterRules_NoFilters(t *testing.T) {
	rules := makeRules()
	filtered := sdk.FilterRules(rules, nil, nil)

	if len(filtered) != 5 {
		t.Errorf("expected all 5 rules, got %d", len(filtered))
	}
}

func TestFilterRules_Ignore(t *testing.T) {
	rules := makeRules()
	filtered := sdk.FilterRules(rules, []string{"og_locale", "twitter_site"}, nil)

	if len(filtered) != 3 {
		t.Errorf("expected 3 rules after ignoring 2, got %d", len(filtered))
	}
	for _, rule := range filtered {
		if rule.Key == "og_locale" || rule.Key == "twitter_site" {
			t.Errorf("rule %s should have been filtered out", rule.Key)
		}
	}
}

func TestFilterRules_Only(t *testing.T) {
	rules := makeRules()
	filtered := sdk.FilterRules(rules, nil, []string{"title_exists", "h1_count"})

	if len(filtered) != 2 {
		t.Errorf("expected 2 rules with only filter, got %d", len(filtered))
	}
}

func TestFilterRules_IgnoreAndOnly(t *testing.T) {
	rules := makeRules()
	filtered := sdk.FilterRules(rules, []string{"h1_count"}, []string{"title_exists", "h1_count"})

	if len(filtered) != 1 {
		t.Errorf("expected 1 rule (title_exists), got %d", len(filtered))
	}
	if filtered[0].Key != "title_exists" {
		t.Errorf("expected title_exists, got %s", filtered[0].Key)
	}
}
