package geo

import (
	"strings"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type AICrawlerAccessChecker struct{}

type botConfig struct {
	name           string
	shouldAllow    bool
	recommendation string
}

func (c *AICrawlerAccessChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	robotsTxt := strings.ToLower(result.RobotsTxt)

	searchBots := []botConfig{
		{
			name:           "OAI-SearchBot",
			shouldAllow:    true,
			recommendation: "Allow OAI-SearchBot in robots.txt to be cited in ChatGPT search results.",
		},
		{
			name:           "PerplexityBot",
			shouldAllow:    true,
			recommendation: "Allow PerplexityBot in robots.txt to appear in Perplexity AI search results.",
		},
		{
			name:           "Claude-SearchBot",
			shouldAllow:    true,
			recommendation: "Allow Claude-SearchBot in robots.txt to be referenced in Claude search answers.",
		},
		{
			name:           "Google-Extended",
			shouldAllow:    true,
			recommendation: "Allow Google-Extended in robots.txt to be featured in Gemini and Google AI Overviews.",
		},
		{
			name:           "Applebot",
			shouldAllow:    true,
			recommendation: "Allow Applebot in robots.txt for Apple search and Siri integrations.",
		},
	}

	trainingBots := []botConfig{
		{
			name:           "GPTBot",
			shouldAllow:    false,
			recommendation: "Consider blocking GPTBot (training-only) in robots.txt while allowing OAI-SearchBot (search).",
		},
		{
			name:           "CCBot",
			shouldAllow:    false,
			recommendation: "Consider blocking CCBot (Common Crawl training data) in robots.txt.",
		},
	}

	for _, bot := range searchBots {
		rule := valueobject.NewAuditRule("geo_crawler_"+strings.ToLower(bot.name), valueobject.CategoryGEO, valueobject.SeverityHigh)

		if robotsTxt == "" {
			rule.Warn(
				"robots.txt missing — "+bot.name+" access cannot be verified",
				"Create a robots.txt file and explicitly allow "+bot.name+".",
			)
		} else if isBotBlocked(robotsTxt, bot.name) {
			rule.Fail(
				bot.name+" is blocked in robots.txt",
				bot.recommendation,
			)
		} else {
			rule.Pass(bot.name + " is allowed (not blocked)")
		}
		rules = append(rules, rule)
	}

	for _, bot := range trainingBots {
		rule := valueobject.NewAuditRule("geo_crawler_block_"+strings.ToLower(bot.name), valueobject.CategoryGEO, valueobject.SeverityLow)

		if robotsTxt == "" {
			rule.Warn(
				"robots.txt missing — "+bot.name+" training access cannot be controlled",
				bot.recommendation,
			)
		} else if isBotBlocked(robotsTxt, bot.name) {
			rule.Pass(bot.name + " is blocked (recommended for training-only bots)")
		} else {
			rule.Warn(
				bot.name+" is not blocked — your content may be used for AI model training",
				bot.recommendation,
			)
		}
		rules = append(rules, rule)
	}

	return rules
}

func isBotBlocked(robotsTxt string, botName string) bool {
	lines := strings.Split(robotsTxt, "\n")
	lowerBot := strings.ToLower(botName)
	inBotSection := false

	for _, line := range lines {
		line = strings.TrimSpace(strings.ToLower(line))

		if strings.HasPrefix(line, "user-agent:") {
			agent := strings.TrimSpace(strings.TrimPrefix(line, "user-agent:"))
			inBotSection = agent == lowerBot || agent == "*"
			if agent != lowerBot && agent != "*" {
				inBotSection = false
			}
		}

		if inBotSection && strings.HasPrefix(line, "disallow:") {
			path := strings.TrimSpace(strings.TrimPrefix(line, "disallow:"))
			if path == "/" || path == "/*" {
				return true
			}
		}
	}

	return false
}
