package eeat

import (
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type EEATSiteChecker struct{}

func (c *EEATSiteChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	foundPages := make(map[string]bool)
	for _, page := range result.Pages {
		loweredURL := strings.ToLower(page.URL)
		for _, pattern := range []string{"/about", "/about-us", "/contact", "/contact-us", "/privacy", "/terms", "/tos"} {
			if strings.Contains(loweredURL, pattern) {
				foundPages[pattern] = true
			}
		}
	}

	aboutRule := valueobject.NewAuditRule("eeat_about_page", valueobject.CategoryEEAT, valueobject.SeverityMedium)
	if !foundPages["/about"] && !foundPages["/about-us"] {
		aboutRule.Fail(
			"No About page found",
			"Create an About page (/about or /about-us) that describes your organization, mission, and team. This is critical for E-E-A-T.",
		)
	} else {
		aboutRule.Pass("About page exists")
	}
	rules = append(rules, aboutRule)

	contactRule := valueobject.NewAuditRule("eeat_contact_page", valueobject.CategoryEEAT, valueobject.SeverityMedium)
	if !foundPages["/contact"] && !foundPages["/contact-us"] {
		contactRule.Fail(
			"No Contact page found",
			"Create a Contact page (/contact or /contact-us) with contact information. This builds trust with both users and search engines.",
		)
	} else {
		contactRule.Pass("Contact page exists")
	}
	rules = append(rules, contactRule)

	privacyRule := valueobject.NewAuditRule("eeat_privacy_policy", valueobject.CategoryEEAT, valueobject.SeverityMedium)
	if !foundPages["/privacy"] {
		privacyRule.Warn(
			"No Privacy Policy page found",
			"Create a Privacy Policy page (/privacy or /privacy-policy) as required by most privacy regulations.",
		)
	} else {
		privacyRule.Pass("Privacy Policy page exists")
	}
	rules = append(rules, privacyRule)

	termsRule := valueobject.NewAuditRule("eeat_terms", valueobject.CategoryEEAT, valueobject.SeverityLow)
	if !foundPages["/terms"] && !foundPages["/tos"] {
		termsRule.Warn(
			"No Terms of Service page found",
			"Create a Terms of Service page to establish legal trust and transparency.",
		)
	} else {
		termsRule.Pass("Terms of Service page exists")
	}
	rules = append(rules, termsRule)

	return rules
}
