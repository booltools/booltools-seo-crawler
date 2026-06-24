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
	foundAnchors := make(map[string]bool)

	for _, page := range result.Pages {
		loweredURL := strings.ToLower(page.URL)
		for _, pattern := range []string{"/about", "/about-us", "/contact", "/contact-us", "/privacy", "/terms", "/tos"} {
			if strings.Contains(loweredURL, pattern) {
				foundPages[pattern] = true
			}
		}

		for _, link := range page.InternalLinks {
			loweredLink := strings.ToLower(link.URL)
			for _, anchor := range []string{"#about", "#contact"} {
				if strings.Contains(loweredLink, anchor) {
					foundAnchors[anchor] = true
				}
			}
		}

		if page.Document != nil {
			if page.Document.Find("#about, [id='about'], #about-us, [id='about-us']").Length() > 0 {
				foundAnchors["#about"] = true
			}
			if page.Document.Find("#contact, [id='contact'], #contact-us, [id='contact-us']").Length() > 0 {
				foundAnchors["#contact"] = true
			}
		}
	}

	aboutRule := valueobject.NewAuditRule("eeat_about_page", valueobject.CategoryEEAT, valueobject.SeverityLow)
	if foundPages["/about"] || foundPages["/about-us"] {
		aboutRule.Pass("About page exists")
	} else if foundAnchors["#about"] {
		aboutRule.Pass("About section found as anchor on page")
	} else {
		aboutRule.Warn(
			"No About page or section found",
			"Create an About page (/about or /about-us) or add an #about section to your homepage.",
		)
	}
	rules = append(rules, aboutRule)

	contactRule := valueobject.NewAuditRule("eeat_contact_page", valueobject.CategoryEEAT, valueobject.SeverityLow)
	if foundPages["/contact"] || foundPages["/contact-us"] {
		contactRule.Pass("Contact page exists")
	} else if foundAnchors["#contact"] {
		contactRule.Pass("Contact section found as anchor on page")
	} else {
		contactRule.Warn(
			"No Contact page or section found",
			"Create a Contact page (/contact or /contact-us) or add a #contact section to your homepage.",
		)
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
