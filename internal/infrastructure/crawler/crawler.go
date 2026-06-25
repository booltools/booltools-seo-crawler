package crawler

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type sitemapURLSet struct {
	XMLName xml.Name          `xml:"urlset"`
	URLs    []sitemapURLEntry `xml:"url"`
}

type sitemapURLEntry struct {
	Loc string `xml:"loc"`
}

type sitemapIndexSet struct {
	XMLName  xml.Name              `xml:"sitemapindex"`
	Sitemaps []sitemapIndexEntry   `xml:"sitemap"`
}

type sitemapIndexEntry struct {
	Loc string `xml:"loc"`
}

type CrawlResult struct {
	Pages          []PageData
	RobotsTxt      string
	SitemapXML     string
	LlmsTxt        string
	LlmsFullTxt    string
	Domain         string
	BaseURL        string
	URLStatusCache *URLStatusCache
}

type SiteCrawler struct{}

func NewSiteCrawler() *SiteCrawler {
	return &SiteCrawler{}
}

type OnPageCallback func(page PageData, pagesCompleted int, totalDiscovered int)

func (sc *SiteCrawler) Crawl(targetDomain string, maxPages int, onPage OnPageCallback) (*CrawlResult, error) {
	targetDomain = normalizeURL(targetDomain)

	parsedURL, err := url.Parse(targetDomain)
	if err != nil {
		return nil, fmt.Errorf("invalid domain URL: %w", err)
	}

	hostname := parsedURL.Hostname()

	result := &CrawlResult{
		Pages:          make([]PageData, 0),
		Domain:         hostname,
		BaseURL:        parsedURL.Scheme + "://" + parsedURL.Host,
		URLStatusCache: NewURLStatusCache(),
	}

	var mutex sync.Mutex
	pagesCompleted := 0
	totalDiscovered := 0

	collector := colly.NewCollector(
		colly.AllowedDomains(hostname, "www."+hostname),
		colly.MaxDepth(10),
		colly.Async(true),
		colly.UserAgent(userAgent),
	)


	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
		Delay:       200 * time.Millisecond,
	})

	collector.SetRequestTimeout(30 * time.Second)

	collector.OnRequest(func(request *colly.Request) {
		mutex.Lock()
		totalDiscovered++
		completed := pagesCompleted
		mutex.Unlock()

		if maxPages > 0 && completed >= maxPages {
			request.Abort()
			return
		}

		request.Ctx.Put("start_time", time.Now().Format(time.RFC3339Nano))
	})

	collector.OnResponse(func(response *colly.Response) {
		mutex.Lock()
		if maxPages > 0 && pagesCompleted >= maxPages {
			mutex.Unlock()
			return
		}
		mutex.Unlock()

		contentType := response.Headers.Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			return
		}

		startTimeStr := response.Ctx.Get("start_time")
		startTime, _ := time.Parse(time.RFC3339Nano, startTimeStr)
		responseTime := time.Since(startTime)

		document, err := goquery.NewDocumentFromReader(strings.NewReader(string(response.Body)))
		if err != nil {
			return
		}

		pageData := PageData{
			URL:           response.Request.URL.String(),
			FinalURL:      response.Request.URL.String(),
			StatusCode:    response.StatusCode,
			Headers:       cloneHeaders(response.Headers),
			Document:      document,
			HTML:          string(response.Body),
			BodyText:      ExtractBodyText(document),
			Depth:         response.Request.Depth,
			ResponseTime:  responseTime,
			ContentLength: int64(len(response.Body)),
			ContentType:   contentType,
		}

		pageData.IsNoindex = detectNoindex(document, response.Headers)
		pageData.IsDevMode = detectDevMode(string(response.Body))

		ExtractLinks(document, response.Request.URL, hostname, &pageData)
		ExtractImages(document, response.Request.URL, &pageData)
		ExtractResources(document, response.Request.URL, &pageData)

		mutex.Lock()
		pagesCompleted++
		currentCompleted := pagesCompleted
		currentTotal := totalDiscovered
		result.Pages = append(result.Pages, pageData)
		mutex.Unlock()

		if onPage != nil {
			onPage(pageData, currentCompleted, currentTotal)
		}
	})

	collector.OnError(func(response *colly.Response, crawlError error) {
		log.Printf("[crawler] error on %s (status %d): %v", response.Request.URL, response.StatusCode, crawlError)
	})

	collector.OnHTML("a[href]", func(element *colly.HTMLElement) {
		link := element.Attr("href")
		element.Request.Visit(link)
	})

	fetchSiteFiles(targetDomain, result)

	visitError := collector.Visit(targetDomain)
	if visitError != nil {
		log.Printf("[crawler] initial visit to %s failed: %v", targetDomain, visitError)
		if parsedURL.Scheme == "https" {
			httpFallback := "http://" + parsedURL.Host + parsedURL.Path
			log.Printf("[crawler] falling back to HTTP: %s", httpFallback)
			visitError = collector.Visit(httpFallback)
			if visitError != nil {
				log.Printf("[crawler] HTTP fallback also failed: %v", visitError)
			}
		}
	}

	if visitError != nil {
		return nil, fmt.Errorf("failed to start crawl for %s: %w", hostname, visitError)
	}

	sitemapURLs := parseSitemapLocations(result.SitemapXML, targetDomain)
	if isLocalhostHost(hostname) {
		sitemapURLs = rewriteURLsToBase(sitemapURLs, result.BaseURL)
	}
	for _, sitemapURL := range sitemapURLs {
		collector.Visit(sitemapURL)
	}
	if len(sitemapURLs) > 0 {
		log.Printf("[crawler] seeded %d URLs from sitemap", len(sitemapURLs))
	}

	collector.Wait()

	if len(result.Pages) == 0 {
		return result, fmt.Errorf("no pages could be crawled from %s — the site may be unreachable, blocking bots, or returning non-HTML content", hostname)
	}

	return result, nil
}

func parseSitemapLocations(content string, baseURL string) []string {
	if content == "" {
		return nil
	}

	var locations []string

	var urlSet sitemapURLSet
	if err := xml.Unmarshal([]byte(content), &urlSet); err == nil && len(urlSet.URLs) > 0 {
		for _, entry := range urlSet.URLs {
			if entry.Loc != "" {
				locations = append(locations, entry.Loc)
			}
		}
		return locations
	}

	var indexSet sitemapIndexSet
	if err := xml.Unmarshal([]byte(content), &indexSet); err == nil && len(indexSet.Sitemaps) > 0 {
		client := &http.Client{Timeout: 10 * time.Second}
		for _, entry := range indexSet.Sitemaps {
			if entry.Loc == "" {
				continue
			}
			response, err := client.Get(entry.Loc)
			if err != nil {
				continue
			}
			body, err := io.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				continue
			}
			childURLs := parseSitemapLocations(string(body), baseURL)
			locations = append(locations, childURLs...)
		}
	}

	return locations
}

func ExtractLinks(document *goquery.Document, baseURL *url.URL, hostname string, pageData *PageData) {
	document.Find("a[href]").Each(func(_ int, selection *goquery.Selection) {
		href, exists := selection.Attr("href")
		if !exists || href == "" {
			return
		}

		resolvedURL, err := baseURL.Parse(href)
		if err != nil {
			return
		}

		linkData := LinkData{
			URL:        resolvedURL.String(),
			AnchorText: strings.TrimSpace(selection.Text()),
			Rel:        selection.AttrOr("rel", ""),
			Target:     selection.AttrOr("target", ""),
			IsInternal: resolvedURL.Hostname() == hostname || resolvedURL.Hostname() == "www."+hostname,
		}

		if linkData.IsInternal {
			pageData.InternalLinks = append(pageData.InternalLinks, linkData)
		} else {
			pageData.ExternalLinks = append(pageData.ExternalLinks, linkData)
		}
	})
}

func ExtractImages(document *goquery.Document, baseURL *url.URL, pageData *PageData) {
	document.Find("img").Each(func(_ int, selection *goquery.Selection) {
		src := selection.AttrOr("src", "")
		if src == "" {
			src = selection.AttrOr("data-src", "")
		}

		if src != "" {
			if resolvedURL, err := baseURL.Parse(src); err == nil {
				src = resolvedURL.String()
			}
		}

		hasPictureSource := false
		pictureParent := selection.ParentsFiltered("picture")
		if pictureParent.Length() > 0 {
			pictureParent.Find("source[srcset]").Each(func(_ int, source *goquery.Selection) {
				srcset := strings.ToLower(source.AttrOr("srcset", ""))
				if strings.Contains(srcset, ".webp") || strings.Contains(srcset, ".avif") {
					hasPictureSource = true
				}
			})
		}

		imageData := ImageData{
			URL:              src,
			Alt:              selection.AttrOr("alt", ""),
			Width:            selection.AttrOr("width", ""),
			Height:           selection.AttrOr("height", ""),
			Loading:          selection.AttrOr("loading", ""),
			FetchPriority:    selection.AttrOr("fetchpriority", ""),
			HasPictureSource: hasPictureSource,
		}

		pageData.Images = append(pageData.Images, imageData)
	})
}

func ExtractResources(document *goquery.Document, baseURL *url.URL, pageData *PageData) {
	document.Find("script[src]").Each(func(_ int, selection *goquery.Selection) {
		src := selection.AttrOr("src", "")
		if src != "" {
			if resolvedURL, err := baseURL.Parse(src); err == nil {
				src = resolvedURL.String()
			}
		}

		_, isAsync := selection.Attr("async")
		_, isDefer := selection.Attr("defer")

		location := "body"
		if selection.ParentsFiltered("head").Length() > 0 {
			location = "head"
		}

		pageData.Scripts = append(pageData.Scripts, ResourceData{
			URL:      src,
			Type:     "javascript",
			IsAsync:  isAsync,
			IsDefer:  isDefer,
			Location: location,
		})
	})

	document.Find("link[rel='stylesheet']").Each(func(_ int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")
		if href != "" {
			if resolvedURL, err := baseURL.Parse(href); err == nil {
				href = resolvedURL.String()
			}
		}
		pageData.Stylesheets = append(pageData.Stylesheets, ResourceData{
			URL:  href,
			Type: "css",
		})
	})
}

func ExtractBodyText(document *goquery.Document) string {
	cloned := document.Clone()
	cloned.Find("script, style, noscript").Remove()
	return strings.TrimSpace(cloned.Find("body").Text())
}

var devModeIndicators = []string{
	"__turbopack",
	"__webpack_hmr",
	"hot-update.js",
	"webpack-dev-server",
	"/@vite/client",
	"/@react-refresh",
	"_next/static/development",
	"__next_error__",
	"__NEXT_DATA__.*\"dev\":true",
	"turbopack-chunk",
}

func detectDevMode(html string) bool {
	lowered := strings.ToLower(html)
	for _, indicator := range devModeIndicators {
		if strings.Contains(lowered, strings.ToLower(indicator)) {
			return true
		}
	}
	return false
}

func detectNoindex(document *goquery.Document, headers *http.Header) bool {
	robotsMeta := document.Find(`meta[name="robots"]`)
	if robotsMeta.Length() > 0 {
		content, exists := robotsMeta.First().Attr("content")
		if exists && strings.Contains(strings.ToLower(content), "noindex") {
			return true
		}
	}

	if headers != nil {
		xRobotsTag := headers.Get("X-Robots-Tag")
		if strings.Contains(strings.ToLower(xRobotsTag), "noindex") {
			return true
		}
	}

	return false
}

func cloneHeaders(headers *http.Header) http.Header {
	if headers == nil {
		return http.Header{}
	}
	clone := make(http.Header)
	for key, values := range *headers {
		clone[key] = append([]string{}, values...)
	}
	return clone
}

func fetchSiteFiles(baseURL string, result *CrawlResult) {
	client := &http.Client{Timeout: 10 * time.Second}

	filesToFetch := map[string]*string{
		"/robots.txt":    &result.RobotsTxt,
		"/sitemap.xml":   &result.SitemapXML,
		"/llms.txt":      &result.LlmsTxt,
		"/llms-full.txt": &result.LlmsFullTxt,
	}

	var waitGroup sync.WaitGroup
	var mutex sync.Mutex

	for path, target := range filesToFetch {
		waitGroup.Add(1)
		go func(filePath string, targetPtr *string) {
			defer waitGroup.Done()

			response, err := client.Get(baseURL + filePath)
			if err != nil {
				return
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				body, err := io.ReadAll(response.Body)
				if err != nil {
					return
				}
				mutex.Lock()
				*targetPtr = string(body)
				mutex.Unlock()
			}
		}(path, target)
	}

	waitGroup.Wait()
}

func normalizeURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return rawURL
	}

	if strings.Contains(rawURL, "://") {
		return rawURL
	}

	host := rawURL
	if slashIndex := strings.Index(rawURL, "/"); slashIndex > 0 {
		host = rawURL[:slashIndex]
	}

	if isLocalhostHost(strings.Split(host, ":")[0]) {
		return "http://" + rawURL
	}

	return "https://" + rawURL
}

func isLocalhostHost(hostname string) bool {
	return hostname == "localhost" || hostname == "127.0.0.1" || hostname == "0.0.0.0" || hostname == "::1"
}

func rewriteURLsToBase(urls []string, baseURL string) []string {
	rewritten := make([]string, 0, len(urls))
	for _, rawURL := range urls {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			rewritten = append(rewritten, rawURL)
			continue
		}
		newURL := baseURL + parsed.Path
		if parsed.RawQuery != "" {
			newURL += "?" + parsed.RawQuery
		}
		rewritten = append(rewritten, newURL)
	}
	return rewritten
}
