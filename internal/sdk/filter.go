package sdk

import "github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"

func FilterRules(rules []valueobject.AuditRule, ignore []string, only []string) []valueobject.AuditRule {
	if len(ignore) == 0 && len(only) == 0 {
		return rules
	}

	ignoreSet := toSet(ignore)
	onlySet := toSet(only)

	filtered := make([]valueobject.AuditRule, 0, len(rules))
	for _, rule := range rules {
		if len(onlySet) > 0 {
			if _, included := onlySet[rule.Key]; !included {
				continue
			}
		}

		if _, excluded := ignoreSet[rule.Key]; excluded {
			continue
		}

		filtered = append(filtered, rule)
	}

	return filtered
}

func toSet(items []string) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}
