package analyzer

import (
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/accessibility"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/content"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/duplicate_content"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/eeat"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/geo"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/links"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/mobile"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/on_page"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/performance"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/security"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/social"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/structured_data"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/technical"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/url_structure"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer/rules/internationalization"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type SiteAnalyzer struct {
	registry *Registry
}

func NewSiteAnalyzer() *SiteAnalyzer {
	registry := NewRegistry()

	registry.RegisterPageChecker(&on_page.TitleChecker{})
	registry.RegisterPageChecker(&on_page.MetaDescriptionChecker{})
	registry.RegisterPageChecker(&on_page.HeadingChecker{})
	registry.RegisterPageChecker(&on_page.ImageChecker{})
	registry.RegisterPageChecker(&content.ContentQualityChecker{})
	registry.RegisterPageChecker(&technical.CanonicalChecker{})
	registry.RegisterPageChecker(&technical.MetaRobotsChecker{})
	registry.RegisterPageChecker(&technical.HTTPStatusChecker{})
	registry.RegisterPageChecker(&technical.CrawlDepthChecker{})
	registry.RegisterPageChecker(&links.InternalLinkChecker{})
	registry.RegisterPageChecker(&links.ExternalLinkChecker{})
	registry.RegisterPageChecker(&performance.PageSpeedChecker{})
	registry.RegisterPageChecker(&performance.ResourceChecker{})
	registry.RegisterPageChecker(&structured_data.JsonLdChecker{})
	registry.RegisterPageChecker(&security.HTTPSChecker{})
	registry.RegisterPageChecker(&security.SecurityHeaderChecker{})
	registry.RegisterPageChecker(&accessibility.AccessibilityChecker{})
	registry.RegisterPageChecker(&social.OpenGraphChecker{})
	registry.RegisterPageChecker(&social.TwitterCardChecker{})
	registry.RegisterPageChecker(&mobile.MobileChecker{})
	registry.RegisterPageChecker(&url_structure.URLStructureChecker{})
	registry.RegisterPageChecker(&internationalization.HreflangChecker{})
	registry.RegisterPageChecker(&eeat.EEATPageChecker{})
	registry.RegisterPageChecker(&technical.PaginationChecker{})

	registry.RegisterSiteChecker(&technical.RobotsTxtChecker{})
	registry.RegisterSiteChecker(&technical.SitemapChecker{})
	registry.RegisterSiteChecker(&technical.RedirectChecker{})
	registry.RegisterSiteChecker(&links.BrokenLinkChecker{})
	registry.RegisterSiteChecker(&geo.AICrawlerAccessChecker{})
	registry.RegisterSiteChecker(&geo.LlmsTxtChecker{})
	registry.RegisterSiteChecker(&geo.CitabilityChecker{})
	registry.RegisterSiteChecker(&geo.EntityAuthorityChecker{})
	registry.RegisterSiteChecker(&geo.AIFriendlyChecker{})
	registry.RegisterSiteChecker(&eeat.EEATSiteChecker{})
	registry.RegisterSiteChecker(&duplicate_content.DuplicateContentChecker{})
	registry.RegisterSiteChecker(&performance.BrokenAssetChecker{})

	return &SiteAnalyzer{registry: registry}
}

func (a *SiteAnalyzer) AnalyzePage(page crawler.PageData) []valueobject.AuditRule {
	var allRules []valueobject.AuditRule

	for _, checker := range a.registry.PageCheckers() {
		rules := checker.Check(page)
		allRules = append(allRules, rules...)
	}

	return allRules
}

func (a *SiteAnalyzer) AnalyzeSite(result crawler.CrawlResult) []valueobject.AuditRule {
	var allRules []valueobject.AuditRule

	for _, checker := range a.registry.SiteCheckers() {
		rules := checker.Check(result)
		allRules = append(allRules, rules...)
	}

	return allRules
}
