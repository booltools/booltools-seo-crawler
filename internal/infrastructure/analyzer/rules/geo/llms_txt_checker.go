package geo

import (
	"fmt"
	"strings"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type LlmsTxtChecker struct{}

func (c *LlmsTxtChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	existsRule := valueobject.NewAuditRule("geo_llms_txt_exists", valueobject.CategoryGEO, valueobject.SeverityMedium)
	if result.LlmsTxt == "" {
		existsRule.Fail(
			"llms.txt file is missing",
			"Create a /llms.txt file at your domain root to help AI models understand your site. Include an H1 title, blockquote summary, and H2 sections with annotated links.",
		)
		rules = append(rules, existsRule)

		fullRule := valueobject.NewAuditRule("geo_llms_full_txt", valueobject.CategoryGEO, valueobject.SeverityLow)
		fullRule.Warn(
			"llms-full.txt is also missing",
			"Consider creating /llms-full.txt with expanded content for AI models with larger context windows.",
		)
		rules = append(rules, fullRule)

		return rules
	}
	existsRule.Pass("llms.txt file exists")
	rules = append(rules, existsRule)

	content := result.LlmsTxt
	lines := strings.Split(content, "\n")

	h1Rule := valueobject.NewAuditRule("geo_llms_txt_h1", valueobject.CategoryGEO, valueobject.SeverityMedium)
	hasH1 := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") && !strings.HasPrefix(trimmed, "## ") {
			hasH1 = true
			break
		}
	}
	if !hasH1 {
		h1Rule.Fail(
			"llms.txt is missing the required H1 title",
			"Add an H1 (# Title) as the first heading in your llms.txt with your project or site name.",
		)
	} else {
		h1Rule.Pass("llms.txt has the required H1 title")
	}
	rules = append(rules, h1Rule)

	blockquoteRule := valueobject.NewAuditRule("geo_llms_txt_blockquote", valueobject.CategoryGEO, valueobject.SeverityMedium)
	hasBlockquote := false
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "> ") {
			hasBlockquote = true
			break
		}
	}
	if !hasBlockquote {
		blockquoteRule.Fail(
			"llms.txt is missing the blockquote summary",
			"Add a blockquote (> summary) after the H1 title to provide a one-sentence summary of your project.",
		)
	} else {
		blockquoteRule.Pass("llms.txt has a blockquote summary")
	}
	rules = append(rules, blockquoteRule)

	h2Rule := valueobject.NewAuditRule("geo_llms_txt_sections", valueobject.CategoryGEO, valueobject.SeverityLow)
	h2Count := 0
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "## ") {
			h2Count++
		}
	}
	if h2Count == 0 {
		h2Rule.Warn(
			"llms.txt has no H2 sections",
			"Add H2 sections (## Section) with link lists to categorize your site content for AI models.",
		)
	} else {
		h2Rule.Pass(fmt.Sprintf("llms.txt has %d H2 sections", h2Count))
	}
	rules = append(rules, h2Rule)

	linkRule := valueobject.NewAuditRule("geo_llms_txt_links", valueobject.CategoryGEO, valueobject.SeverityLow)
	linkCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [") && strings.Contains(trimmed, "](") {
			linkCount++
		}
	}
	if linkCount == 0 && h2Count > 0 {
		linkRule.Warn(
			"llms.txt sections have no link entries",
			"Add link entries in format: - [Title](URL): Description",
		)
	} else if linkCount > 0 {
		linkRule.Pass(fmt.Sprintf("llms.txt contains %d link entries", linkCount))
	}
	rules = append(rules, linkRule)

	fullRule := valueobject.NewAuditRule("geo_llms_full_txt", valueobject.CategoryGEO, valueobject.SeverityLow)
	if result.LlmsFullTxt == "" {
		fullRule.Warn(
			"llms-full.txt not found",
			"Consider creating /llms-full.txt with expanded content for AI models with larger context windows.",
		)
	} else {
		fullRule.Pass("llms-full.txt exists")
	}
	rules = append(rules, fullRule)

	return rules
}
