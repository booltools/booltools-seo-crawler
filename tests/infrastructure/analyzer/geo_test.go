package analyzer_test

import (
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer/rules/geo"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

func TestAICrawlerAccessChecker_MissingRobotsTxt(t *testing.T) {
	result := crawler.CrawlResult{RobotsTxt: ""}
	checker := &geo.AICrawlerAccessChecker{}
	rules := checker.Check(result)

	if len(rules) == 0 {
		t.Fatal("expected rules")
	}

	for _, rule := range rules {
		if rule.Result != valueobject.RuleResultWarning {
			t.Errorf("expected warning for '%s' when robots.txt missing, got %s", rule.Key, rule.Result)
		}
	}
}

func TestAICrawlerAccessChecker_BlockedSearchBot(t *testing.T) {
	robotsTxt := `User-agent: OAI-SearchBot
Disallow: /

User-agent: PerplexityBot
Disallow: /`

	result := crawler.CrawlResult{RobotsTxt: robotsTxt}
	checker := &geo.AICrawlerAccessChecker{}
	rules := checker.Check(result)

	oaiRule := findRule(rules, "geo_crawler_oai-searchbot")
	if oaiRule == nil || oaiRule.Result != valueobject.RuleResultFail {
		t.Error("expected OAI-SearchBot to be flagged as blocked")
	}
}

func TestAICrawlerAccessChecker_AllowedSearchBot(t *testing.T) {
	robotsTxt := `User-agent: *
Allow: /`

	result := crawler.CrawlResult{RobotsTxt: robotsTxt}
	checker := &geo.AICrawlerAccessChecker{}
	rules := checker.Check(result)

	oaiRule := findRule(rules, "geo_crawler_oai-searchbot")
	if oaiRule == nil || oaiRule.Result != valueobject.RuleResultPass {
		t.Error("expected OAI-SearchBot to pass when not blocked")
	}
}

func TestLlmsTxtChecker_Missing(t *testing.T) {
	result := crawler.CrawlResult{LlmsTxt: ""}
	checker := &geo.LlmsTxtChecker{}
	rules := checker.Check(result)

	existsRule := findRule(rules, "geo_llms_txt_exists")
	if existsRule == nil || existsRule.Result != valueobject.RuleResultFail {
		t.Error("expected geo_llms_txt_exists to fail when missing")
	}
}

func TestLlmsTxtChecker_Valid(t *testing.T) {
	llmsTxt := `# My Site

> A great site for learning stuff

## Docs

- [Getting Started](https://example.com/docs/start): Start here
- [API Reference](https://example.com/docs/api): Full API docs
`
	result := crawler.CrawlResult{LlmsTxt: llmsTxt}
	checker := &geo.LlmsTxtChecker{}
	rules := checker.Check(result)

	existsRule := findRule(rules, "geo_llms_txt_exists")
	if existsRule == nil || existsRule.Result != valueobject.RuleResultPass {
		t.Error("expected exists to pass")
	}

	h1Rule := findRule(rules, "geo_llms_txt_h1")
	if h1Rule == nil || h1Rule.Result != valueobject.RuleResultPass {
		t.Error("expected h1 to pass")
	}

	blockquoteRule := findRule(rules, "geo_llms_txt_blockquote")
	if blockquoteRule == nil || blockquoteRule.Result != valueobject.RuleResultPass {
		t.Error("expected blockquote to pass")
	}

	sectionsRule := findRule(rules, "geo_llms_txt_sections")
	if sectionsRule == nil || sectionsRule.Result != valueobject.RuleResultPass {
		t.Error("expected sections to pass")
	}

	linksRule := findRule(rules, "geo_llms_txt_links")
	if linksRule == nil || linksRule.Result != valueobject.RuleResultPass {
		t.Error("expected links to pass")
	}
}

func TestLlmsTxtChecker_MissingH1(t *testing.T) {
	llmsTxt := `> Summary but no H1

## Section
`
	result := crawler.CrawlResult{LlmsTxt: llmsTxt}
	checker := &geo.LlmsTxtChecker{}
	rules := checker.Check(result)

	h1Rule := findRule(rules, "geo_llms_txt_h1")
	if h1Rule == nil || h1Rule.Result != valueobject.RuleResultFail {
		t.Error("expected h1 to fail when missing")
	}
}
