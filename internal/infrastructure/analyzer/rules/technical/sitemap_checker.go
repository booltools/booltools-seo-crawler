package technical

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type SitemapChecker struct{}

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc        string `xml:"loc"`
	Lastmod    string `xml:"lastmod"`
	Changefreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

type sitemapIndex struct {
	XMLName  xml.Name          `xml:"sitemapindex"`
	Sitemaps []sitemapLocation `xml:"sitemap"`
}

type sitemapLocation struct {
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod"`
}

func (c *SitemapChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	existsRule := valueobject.NewAuditRule("sitemap_exists", valueobject.CategoryTechnical, valueobject.SeverityHigh)
	if result.SitemapXML == "" {
		existsRule.Fail(
			"XML sitemap is missing or inaccessible",
			"Create a sitemap.xml file and submit it to Google Search Console and Bing Webmaster Tools.",
		)
		rules = append(rules, existsRule)
		return rules
	}
	existsRule.Pass("XML sitemap exists and is accessible")
	rules = append(rules, existsRule)

	allURLs := c.parseSitemapURLs(result.SitemapXML, &rules)

	c.checkSitemapSize(allURLs, result.SitemapXML, &rules)
	c.checkLastmodDates(allURLs, &rules)
	c.checkCrawledPageCoverage(allURLs, result, &rules)
	c.checkSitemapURLsReachable(allURLs, result, &rules)
	c.checkSitemapURLsNotBlocked(allURLs, result, &rules)
	c.checkImageSitemap(result.SitemapXML, &rules)
	c.checkVideoSitemap(result.SitemapXML, &rules)

	return rules
}

func (c *SitemapChecker) parseSitemapURLs(sitemapContent string, rules *[]valueobject.AuditRule) []sitemapURL {
	validXMLRule := valueobject.NewAuditRule("sitemap_valid_xml", valueobject.CategoryTechnical, valueobject.SeverityHigh)

	var urlSet sitemapURLSet
	if err := xml.Unmarshal([]byte(sitemapContent), &urlSet); err == nil {
		validXMLRule.Pass("Sitemap XML is valid (urlset format)")
		*rules = append(*rules, validXMLRule)
		return urlSet.URLs
	}

	var sitemapIdx sitemapIndex
	if err := xml.Unmarshal([]byte(sitemapContent), &sitemapIdx); err == nil && len(sitemapIdx.Sitemaps) > 0 {
		validXMLRule.Pass(fmt.Sprintf("Sitemap index found with %d child sitemaps", len(sitemapIdx.Sitemaps)))
		*rules = append(*rules, validXMLRule)

		indexRule := valueobject.NewAuditRule("sitemap_index", valueobject.CategoryTechnical, valueobject.SeverityInfo)
		indexRule.Pass(fmt.Sprintf("Site uses a sitemap index with %d sitemaps", len(sitemapIdx.Sitemaps)))
		*rules = append(*rules, indexRule)

		var allURLs []sitemapURL
		for _, loc := range sitemapIdx.Sitemaps {
			allURLs = append(allURLs, sitemapURL{Loc: loc.Loc, Lastmod: loc.Lastmod})
		}
		return allURLs
	}

	validXMLRule.Fail(
		"Sitemap is not valid XML (neither urlset nor sitemapindex)",
		"Fix the XML syntax in your sitemap. Ensure proper encoding and well-formed tags.",
	)
	*rules = append(*rules, validXMLRule)
	return nil
}

func (c *SitemapChecker) checkSitemapSize(urls []sitemapURL, rawContent string, rules *[]valueobject.AuditRule) {
	sizeRule := valueobject.NewAuditRule("sitemap_size", valueobject.CategoryTechnical, valueobject.SeverityMedium)
	urlCount := len(urls)
	fileSizeKB := len(rawContent) / 1024

	if urlCount > 50000 {
		sizeRule.Fail(
			fmt.Sprintf("Sitemap exceeds 50,000 URL limit (%d URLs)", urlCount),
			"Split your sitemap into multiple files using a sitemap index. Each file should have at most 50,000 URLs.",
		)
	} else if fileSizeKB > 50*1024 {
		sizeRule.Fail(
			fmt.Sprintf("Sitemap exceeds 50MB size limit (%d KB)", fileSizeKB),
			"Reduce the sitemap file size by splitting into multiple sitemaps using a sitemap index.",
		)
	} else {
		sizeRule.Pass(fmt.Sprintf("Sitemap contains %d URLs (%d KB)", urlCount, fileSizeKB))
	}
	*rules = append(*rules, sizeRule)
}

func (c *SitemapChecker) checkLastmodDates(urls []sitemapURL, rules *[]valueobject.AuditRule) {
	if len(urls) == 0 {
		return
	}

	freshnessRule := valueobject.NewAuditRule("sitemap_freshness", valueobject.CategoryTechnical, valueobject.SeverityLow)
	hasLastmod := false
	for _, entry := range urls {
		if entry.Lastmod != "" {
			hasLastmod = true
			break
		}
	}

	if !hasLastmod {
		freshnessRule.Warn(
			"Sitemap URLs do not include lastmod dates",
			"Add <lastmod> dates to sitemap entries to help search engines prioritize recently updated content.",
		)
	} else {
		freshnessRule.Pass("Sitemap includes lastmod dates")
	}
	*rules = append(*rules, freshnessRule)
}

func (c *SitemapChecker) checkCrawledPageCoverage(sitemapURLs []sitemapURL, result crawler.CrawlResult, rules *[]valueobject.AuditRule) {
	sitemapNormalized := make(map[string]bool)
	for _, entry := range sitemapURLs {
		normalized := strings.TrimSuffix(strings.ToLower(entry.Loc), "/")
		sitemapNormalized[normalized] = true
	}

	missingFromSitemap := 0
	var missingURLs []string
	for _, page := range result.Pages {
		normalized := strings.TrimSuffix(strings.ToLower(page.URL), "/")
		if !sitemapNormalized[normalized] {
			missingFromSitemap++
			missingURLs = append(missingURLs, page.URL)
		}
	}

	coverageRule := valueobject.NewAuditRule("sitemap_coverage", valueobject.CategoryTechnical, valueobject.SeverityMedium)
	if missingFromSitemap > 0 && len(result.Pages) > 0 {
		coverageRule.Warn(
			fmt.Sprintf("%d of %d crawled pages are missing from the sitemap", missingFromSitemap, len(result.Pages)),
			"Add all important indexable pages to your sitemap to ensure search engines discover them.",
		)
		coverageRule.WithDetails(strings.Join(missingURLs, "\n"))
	} else {
		coverageRule.Pass("All crawled pages are present in the sitemap")
	}
	*rules = append(*rules, coverageRule)

	orphanRule := valueobject.NewAuditRule("sitemap_orphan_urls", valueobject.CategoryTechnical, valueobject.SeverityMedium)
	crawledNormalized := make(map[string]bool)
	for _, page := range result.Pages {
		normalized := strings.TrimSuffix(strings.ToLower(page.URL), "/")
		crawledNormalized[normalized] = true
	}

	orphanCount := 0
	var orphanURLs []string
	for _, entry := range sitemapURLs {
		normalized := strings.TrimSuffix(strings.ToLower(entry.Loc), "/")
		if !crawledNormalized[normalized] {
			orphanCount++
			orphanURLs = append(orphanURLs, entry.Loc)
		}
	}

	if orphanCount > 0 {
		orphanRule.Warn(
			fmt.Sprintf("%d sitemap URLs were not discovered during crawling (possible orphan pages)", orphanCount),
			"Orphan pages in the sitemap are not linked from your site. Either add internal links to them or remove them from the sitemap.",
		)
		orphanRule.WithDetails(strings.Join(orphanURLs, "\n"))
	} else {
		orphanRule.Pass("All sitemap URLs are internally linked")
	}
	*rules = append(*rules, orphanRule)
}

func (c *SitemapChecker) checkSitemapURLsReachable(urls []sitemapURL, result crawler.CrawlResult, rules *[]valueobject.AuditRule) {
	if len(urls) == 0 {
		return
	}

	maxToCheck := min(len(urls), 30)
	noRedirectCache := crawler.NewURLStatusCacheNoRedirect()

	urlMap := make(map[string]string, maxToCheck)
	for i, entry := range urls {
		if i >= maxToCheck {
			break
		}
		urlMap[entry.Loc] = entry.Loc
	}

	statusResults := noRedirectCache.CheckConcurrent(urlMap, maxToCheck, 5)

	brokenCount := 0
	redirectCount := 0
	var brokenExamples []string
	var redirectExamples []string

	for urlToCheck, statusResult := range statusResults {
		if statusResult.Error != nil {
			brokenCount++
			brokenExamples = append(brokenExamples, urlToCheck+" (connection error)")
		} else if statusResult.StatusCode >= 400 {
			brokenCount++
			brokenExamples = append(brokenExamples, fmt.Sprintf("%s (HTTP %d)", urlToCheck, statusResult.StatusCode))
		} else if statusResult.StatusCode >= 300 {
			redirectCount++
			redirectExamples = append(redirectExamples, fmt.Sprintf("%s -> %d", urlToCheck, statusResult.StatusCode))
		}
	}

	brokenRule := valueobject.NewAuditRule("sitemap_broken_urls", valueobject.CategoryTechnical, valueobject.SeverityHigh)
	if brokenCount > 0 {
		brokenRule.Fail(
			fmt.Sprintf("%d sitemap URLs return errors (checked %d of %d)", brokenCount, maxToCheck, len(urls)),
			"Remove broken URLs from the sitemap or fix them. Sitemaps should only contain URLs returning 200 OK.",
		)
		if len(brokenExamples) > 0 {
			brokenRule.WithDetails(strings.Join(brokenExamples, "; "))
		}
	} else {
		brokenRule.Pass(fmt.Sprintf("All %d checked sitemap URLs are reachable", maxToCheck))
	}
	*rules = append(*rules, brokenRule)

	redirectRule := valueobject.NewAuditRule("sitemap_redirect_urls", valueobject.CategoryTechnical, valueobject.SeverityMedium)
	if redirectCount > 0 {
		redirectRule.Warn(
			fmt.Sprintf("%d sitemap URLs redirect to another location", redirectCount),
			"Update the sitemap to use the final destination URLs instead of redirecting URLs.",
		)
		if len(redirectExamples) > 0 {
			redirectRule.WithDetails(strings.Join(redirectExamples, "; "))
		}
	} else {
		redirectRule.Pass("No sitemap URLs redirect")
	}
	*rules = append(*rules, redirectRule)
}

func (c *SitemapChecker) checkSitemapURLsNotBlocked(urls []sitemapURL, result crawler.CrawlResult, rules *[]valueobject.AuditRule) {
	if result.RobotsTxt == "" || len(urls) == 0 {
		return
	}

	blockedRule := valueobject.NewAuditRule("sitemap_robots_conflict", valueobject.CategoryTechnical, valueobject.SeverityHigh)

	disallowedPaths := parseDisallowedPaths(result.RobotsTxt)
	if len(disallowedPaths) == 0 {
		blockedRule.Pass("No sitemap/robots.txt conflicts")
		*rules = append(*rules, blockedRule)
		return
	}

	blockedCount := 0
	var blockedExamples []string

	for _, entry := range urls {
		for _, disallowed := range disallowedPaths {
			if strings.Contains(strings.ToLower(entry.Loc), strings.ToLower(disallowed)) {
				blockedCount++
				if len(blockedExamples) < 5 {
					blockedExamples = append(blockedExamples, fmt.Sprintf("%s (blocked by Disallow: %s)", entry.Loc, disallowed))
				}
				break
			}
		}
	}

	if blockedCount > 0 {
		blockedRule.Fail(
			fmt.Sprintf("%d sitemap URLs are blocked by robots.txt", blockedCount),
			"Remove blocked URLs from the sitemap, or update robots.txt to allow them. URLs in the sitemap should always be crawlable.",
		)
		if len(blockedExamples) > 0 {
			blockedRule.WithDetails(strings.Join(blockedExamples, "; "))
		}
	} else {
		blockedRule.Pass("No sitemap URLs are blocked by robots.txt")
	}
	*rules = append(*rules, blockedRule)
}

func (c *SitemapChecker) checkImageSitemap(content string, rules *[]valueobject.AuditRule) {
	imageRule := valueobject.NewAuditRule("sitemap_image", valueobject.CategoryTechnical, valueobject.SeverityLow)
	if strings.Contains(content, "image:image") || strings.Contains(content, "image:loc") {
		imageRule.Pass("Sitemap includes image entries for better image indexing")
	} else {
		imageRule.Warn(
			"Sitemap does not include image entries",
			"Add <image:image> tags to your sitemap for better image discovery and indexing in Google Images.",
		)
	}
	*rules = append(*rules, imageRule)
}

func (c *SitemapChecker) checkVideoSitemap(content string, rules *[]valueobject.AuditRule) {
	videoRule := valueobject.NewAuditRule("sitemap_video", valueobject.CategoryTechnical, valueobject.SeverityInfo)
	if strings.Contains(content, "video:video") || strings.Contains(content, "video:content_loc") {
		videoRule.Pass("Sitemap includes video entries for better video indexing")
	} else {
		videoRule.Warn(
			"Sitemap does not include video entries (if applicable)",
			"If your site has videos, add <video:video> tags to your sitemap for better video discovery in search results.",
		)
	}
	*rules = append(*rules, videoRule)
}

func parseDisallowedPaths(robotsTxt string) []string {
	var paths []string
	lines := strings.Split(robotsTxt, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "disallow:") {
			path := strings.TrimSpace(strings.TrimPrefix(line, strings.SplitN(line, ":", 2)[0]+":"))
			if path != "" && path != "/" {
				paths = append(paths, path)
			}
		}
	}

	return paths
}
