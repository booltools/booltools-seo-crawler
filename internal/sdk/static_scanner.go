package sdk

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/analyzer"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/crawler"
)

type StaticScanner struct {
	analyzer *analyzer.SiteAnalyzer
}

func NewStaticScanner(siteAnalyzer *analyzer.SiteAnalyzer) *StaticScanner {
	return &StaticScanner{analyzer: siteAnalyzer}
}

func (scanner *StaticScanner) Scan(directory string) (*ScanResult, error) {
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

	totalFiles := len(htmlFiles)
	for index, htmlFile := range htmlFiles {
		relativePath, _ := filepath.Rel(absDirectory, htmlFile)
		fmt.Fprintf(os.Stderr, "\r  [%d/%d] %s", index+1, totalFiles, filepath.ToSlash(relativePath))

		pageData, parseError := parseHTMLFile(htmlFile, absDirectory)
		if parseError != nil {
			fmt.Fprintf(os.Stderr, " (skipped: %v)\n", parseError)
			continue
		}

		rules := scanner.analyzer.AnalyzePage(pageData)
		rules = filterNetworkDependentRules(rules)
		pageResult := buildPageScanResult(pageData.URL, rules)
		result.Pages = append(result.Pages, pageResult)
		result.AllRules = append(result.AllRules, rules...)
	}
	fmt.Fprintln(os.Stderr)

	return result, nil
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
	pageURL := "file://" + filepath.ToSlash(relativePath)

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
	}
	crawler.ExtractResources(document, &pageData)

	return pageData, nil
}

var networkDependentRules = map[string]struct{}{
	"http_status_ok":              {},
	"uses_https":                  {},
	"mixed_content":               {},
	"hsts_header":                 {},
	"security_xcto":               {},
	"security_xfo":                {},
	"security_csp":                {},
	"security_referrer":           {},
	"security_permissions":        {},
	"security_server_disclosure":  {},
	"cache_headers":               {},
	"ttfb":                        {},
	"compression":                 {},
	"page_size":                   {},
	"crawl_depth":                 {},
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
