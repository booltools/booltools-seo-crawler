package technical

import (
	"fmt"
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type RobotsTxtChecker struct{}

func (c *RobotsTxtChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule
	robotsTxt := result.RobotsTxt

	existsRule := valueobject.NewAuditRule("robots_txt_exists", valueobject.CategoryTechnical, valueobject.SeverityHigh)
	if robotsTxt == "" {
		existsRule.Fail(
			"robots.txt file is missing or inaccessible",
			"Create a robots.txt file at the root of your domain to control how search engines crawl your site.",
		)
		rules = append(rules, existsRule)
		return rules
	}
	existsRule.Pass("robots.txt file exists and is accessible")
	rules = append(rules, existsRule)

	sizeRule := valueobject.NewAuditRule("robots_txt_size", valueobject.CategoryTechnical, valueobject.SeverityLow)
	if len(robotsTxt) > 500*1024 {
		sizeRule.Fail(
			fmt.Sprintf("robots.txt is too large (%d KB)", len(robotsTxt)/1024),
			"Reduce the robots.txt file size to under 500KB. Google may not fully process larger files.",
		)
	} else {
		sizeRule.Pass(fmt.Sprintf("robots.txt size is acceptable (%d bytes)", len(robotsTxt)))
	}
	rules = append(rules, sizeRule)

	sitemapRule := valueobject.NewAuditRule("robots_txt_sitemap", valueobject.CategoryTechnical, valueobject.SeverityMedium)
	if !strings.Contains(strings.ToLower(robotsTxt), "sitemap:") {
		sitemapRule.Fail(
			"robots.txt does not reference a sitemap",
			"Add a Sitemap directive to your robots.txt: Sitemap: https://yourdomain.com/sitemap.xml",
		)
	} else {
		sitemapRule.Pass("robots.txt references a sitemap")
	}
	rules = append(rules, sitemapRule)

	syntaxRule := valueobject.NewAuditRule("robots_txt_syntax", valueobject.CategoryTechnical, valueobject.SeverityMedium)
	syntaxErrors := validateRobotsTxtSyntax(robotsTxt)
	if len(syntaxErrors) > 0 {
		syntaxRule.Fail(
			fmt.Sprintf("robots.txt has %d syntax issues", len(syntaxErrors)),
			"Fix the robots.txt syntax. Valid directives include: User-agent, Disallow, Allow, Sitemap, Crawl-delay.",
		)
		syntaxRule.WithDetails(strings.Join(syntaxErrors, "; "))
	} else {
		syntaxRule.Pass("robots.txt syntax is valid")
	}
	rules = append(rules, syntaxRule)

	return rules
}

func validateRobotsTxtSyntax(content string) []string {
	var errors []string
	validDirectives := map[string]bool{
		"user-agent":  true,
		"disallow":    true,
		"allow":       true,
		"sitemap":     true,
		"crawl-delay": true,
		"host":        true,
	}

	lines := strings.Split(content, "\n")
	for lineNumber, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			errors = append(errors, fmt.Sprintf("line %d: missing colon separator", lineNumber+1))
			continue
		}

		directive := strings.ToLower(strings.TrimSpace(parts[0]))
		if !validDirectives[directive] {
			errors = append(errors, fmt.Sprintf("line %d: unknown directive '%s'", lineNumber+1, parts[0]))
		}
	}

	return errors
}
