package technical

import (
	"net/url"
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type CanonicalChecker struct{}

func (c *CanonicalChecker) Check(page crawler.PageData) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	canonicals := page.Document.Find(`link[rel="canonical"]`)
	canonicalCount := canonicals.Length()
	canonicalHref, _ := canonicals.First().Attr("href")
	canonicalHref = strings.TrimSpace(canonicalHref)

	headerCanonical := ""
	linkHeader := page.Headers.Get("Link")
	if strings.Contains(linkHeader, `rel="canonical"`) {
		parts := strings.Split(linkHeader, ";")
		if len(parts) > 0 {
			headerCanonical = strings.Trim(strings.TrimSpace(parts[0]), "<>")
		}
	}

	existsRule := valueobject.NewAuditRule("canonical_exists", valueobject.CategoryTechnical, valueobject.SeverityHigh)
	existsRule.AffectedURL = page.URL
	if canonicalCount == 0 && headerCanonical == "" {
		existsRule.Fail(
			"Page is missing a canonical tag",
			"Add a <link rel=\"canonical\" href=\"...\"> tag to specify the preferred URL for this page.",
		)
	} else {
		existsRule.Pass("Canonical tag is present")
	}
	rules = append(rules, existsRule)

	if canonicalHref == "" && headerCanonical == "" {
		return rules
	}

	effectiveCanonical := canonicalHref
	if effectiveCanonical == "" {
		effectiveCanonical = headerCanonical
	}

	absoluteRule := valueobject.NewAuditRule("canonical_absolute", valueobject.CategoryTechnical, valueobject.SeverityMedium)
	absoluteRule.AffectedURL = page.URL
	parsedCanonical, err := url.Parse(effectiveCanonical)
	if err != nil || !parsedCanonical.IsAbs() {
		absoluteRule.Fail(
			"Canonical URL is not absolute",
			"Use an absolute URL (including https://) in your canonical tag.",
		)
	} else {
		absoluteRule.Pass("Canonical URL is absolute")
	}
	rules = append(rules, absoluteRule)

	selfRefRule := valueobject.NewAuditRule("canonical_self_ref", valueobject.CategoryTechnical, valueobject.SeverityLow)
	selfRefRule.AffectedURL = page.URL
	normalizedPageURL := normalizeURL(page.URL)
	normalizedCanonical := normalizeURL(effectiveCanonical)
	if normalizedCanonical != normalizedPageURL {
		selfRefRule.Warn(
			"Canonical URL points to a different page",
			"Verify that the canonical URL is intentional. If this is the preferred version of the page, use a self-referencing canonical.",
		)
		selfRefRule.WithDetails("Canonical: " + effectiveCanonical)
	} else {
		selfRefRule.Pass("Canonical is self-referencing")
	}
	rules = append(rules, selfRefRule)

	if headerCanonical != "" && canonicalHref != "" {
		conflictRule := valueobject.NewAuditRule("canonical_conflict", valueobject.CategoryTechnical, valueobject.SeverityHigh)
		conflictRule.AffectedURL = page.URL
		if normalizeURL(headerCanonical) != normalizeURL(canonicalHref) {
			conflictRule.Fail(
				"Conflicting canonical signals between HTTP header and HTML tag",
				"Ensure the canonical URL in the HTTP Link header matches the HTML <link rel=\"canonical\"> tag.",
			)
		} else {
			conflictRule.Pass("HTTP header and HTML canonical tags are consistent")
		}
		rules = append(rules, conflictRule)
	}

	return rules
}

func normalizeURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	parsed.Fragment = ""
	result := parsed.String()
	result = strings.TrimSuffix(result, "/")
	return strings.ToLower(result)
}
