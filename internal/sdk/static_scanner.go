package sdk

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type StaticScanner struct {
	analyzer *analyzer.SiteAnalyzer
}

func NewStaticScanner(siteAnalyzer *analyzer.SiteAnalyzer) *StaticScanner {
	return &StaticScanner{analyzer: siteAnalyzer}
}

func (scanner *StaticScanner) Scan(directory string, excludeURLs []string, onlyURLs []string) (*ScanResult, error) {
	absDirectory, err := filepath.Abs(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve directory: %w", err)
	}

	info, err := os.Stat(absDirectory)
	if err != nil {
		return nil, fmt.Errorf("directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", absDirectory)
	}

	htmlFiles, err := findHTMLFiles(absDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to find HTML files: %w", err)
	}

	if len(htmlFiles) == 0 {
		return nil, fmt.Errorf("no .html files found in '%s'", absDirectory)
	}

	result := &ScanResult{
		Mode:       "static",
		TotalPages: len(htmlFiles),
		Pages:      make([]PageScanResult, 0, len(htmlFiles)),
	}

	var allPages []crawler.PageData

	totalFiles := len(htmlFiles)
	for index, htmlFile := range htmlFiles {
		relativePath, _ := filepath.Rel(absDirectory, htmlFile)
		fmt.Fprintf(os.Stderr, "\r  [%d/%d] %s", index+1, totalFiles, filepath.ToSlash(relativePath))

		pageData, parseError := parseHTMLFile(htmlFile, absDirectory)
		if parseError != nil {
			fmt.Fprintf(os.Stderr, " (skipped: %v)\n", parseError)
			continue
		}

		allPages = append(allPages, pageData)

		if !ShouldAnalyzeURL(pageData.URL, excludeURLs, onlyURLs) {
			continue
		}

		rules := scanner.analyzer.AnalyzePage(pageData)
		rules = filterNetworkDependentRules(rules)
		pageResult := buildPageScanResult(pageData.URL, rules)
		result.Pages = append(result.Pages, pageResult)
		result.AllRules = append(result.AllRules, rules...)
	}
	fmt.Fprintln(os.Stderr)

	crawlResult := buildStaticCrawlResult(absDirectory, allPages)
	siteRules := scanner.analyzer.AnalyzeSite(crawlResult)
	siteRules = filterNetworkDependentSiteRules(siteRules)
	result.AllRules = append(result.AllRules, siteRules...)

	return result, nil
}

func buildStaticCrawlResult(directory string, pages []crawler.PageData) crawler.CrawlResult {
	crawlResult := crawler.CrawlResult{
		Pages:          pages,
		Domain:         "file:///" + filepath.ToSlash(directory),
		URLStatusCache: crawler.NewURLStatusCache(),
	}

	crawlResult.RobotsTxt = readFileIfExists(filepath.Join(directory, "robots.txt"))
	crawlResult.SitemapXML = readFileIfExists(filepath.Join(directory, "sitemap.xml"))
	crawlResult.LlmsTxt = readFileIfExists(filepath.Join(directory, "llms.txt"))
	crawlResult.LlmsFullTxt = readFileIfExists(filepath.Join(directory, "llms-full.txt"))

	return crawlResult
}

func readFileIfExists(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

func findHTMLFiles(directory string) ([]string, error) {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, walkError error) error {
		if walkError != nil {
			return walkError
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".html") || strings.HasSuffix(strings.ToLower(info.Name()), ".htm") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func parseHTMLFile(filePath string, baseDirectory string) (crawler.PageData, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return crawler.PageData{}, fmt.Errorf("failed to read file: %w", err)
	}

	html := string(content)
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return crawler.PageData{}, fmt.Errorf("failed to parse HTML: %w", err)
	}

	relativePath, _ := filepath.Rel(baseDirectory, filePath)
	pageURL := "file:///" + filepath.ToSlash(relativePath)

	parsedURL, _ := url.Parse(pageURL)

	bodyClone := document.Clone()
	bodyClone.Find("script, style, noscript").Remove()
	bodyText := strings.TrimSpace(bodyClone.Find("body").Text())

	pageData := crawler.PageData{
		URL:      pageURL,
		FinalURL: pageURL,
		Document: document,
		HTML:     html,
		BodyText: bodyText,
	}

	if parsedURL != nil {
		crawler.ExtractLinks(document, parsedURL, "", &pageData)
		crawler.ExtractImages(document, parsedURL, &pageData)
		crawler.ExtractResources(document, parsedURL, &pageData)
	}

	return pageData, nil
}

var networkDependentRules = map[string]struct{}{
	"http_status_ok":             {},
	"uses_https":                 {},
	"mixed_content":              {},
	"hsts_header":                {},
	"security_xcto":              {},
	"security_xfo":               {},
	"security_csp":               {},
	"security_referrer":          {},
	"security_permissions":       {},
	"security_server_disclosure": {},
	"cache_headers":              {},
	"ttfb":                       {},
	"compression":                {},
	"page_size":                  {},
	"crawl_depth":                {},
	"canonical_self_ref":         {},
}

var networkDependentSiteRules = map[string]struct{}{
	"broken_external_links": {},
	"broken_internal_links": {},
	"broken_scripts":        {},
	"broken_stylesheets":    {},
	"broken_images":         {},
	"redirect_chains":       {},
	"temporary_redirects":   {},
	"sitemap_broken_urls":   {},
	"sitemap_redirect_urls": {},
	"sitemap_coverage":      {},
	"sitemap_orphan_urls":   {},
}

func filterNetworkDependentRules(rules []valueobject.AuditRule) []valueobject.AuditRule {
	filtered := make([]valueobject.AuditRule, 0, len(rules))
	for _, rule := range rules {
		if _, isNetworkRule := networkDependentRules[rule.Key]; isNetworkRule {
			continue
		}
		filtered = append(filtered, rule)
	}
	return filtered
}

func filterNetworkDependentSiteRules(rules []valueobject.AuditRule) []valueobject.AuditRule {
	filtered := make([]valueobject.AuditRule, 0, len(rules))
	for _, rule := range rules {
		if _, isNetworkRule := networkDependentSiteRules[rule.Key]; isNetworkRule {
			continue
		}
		filtered = append(filtered, rule)
	}
	return filtered
}
