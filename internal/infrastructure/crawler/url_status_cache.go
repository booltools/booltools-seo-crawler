package crawler

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type URLStatusResult struct {
	StatusCode int
	Error      error
}

type URLStatusCache struct {
	cache  map[string]URLStatusResult
	mutex  sync.RWMutex
	client *http.Client
}

const (
	maxRetries    = 3
	baseBackoffMs = 500
	jitterMs      = 200
	userAgent     = "Mozilla/5.0 (compatible; BoolToolsSEOCrawler/1.0; +https://github.com/booltools/booltools-seo-crawler)"
)

func NewURLStatusCache() *URLStatusCache {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	}

	return &URLStatusCache{
		cache: make(map[string]URLStatusResult),
		client: &http.Client{
			Timeout:   15 * time.Second,
			Transport: transport,
			CheckRedirect: func(request *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

func NewURLStatusCacheNoRedirect() *URLStatusCache {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	}

	return &URLStatusCache{
		cache: make(map[string]URLStatusResult),
		client: &http.Client{
			Timeout:   15 * time.Second,
			Transport: transport,
			CheckRedirect: func(request *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *URLStatusCache) Get(targetURL string) (URLStatusResult, bool) {
	c.mutex.RLock()
	result, exists := c.cache[targetURL]
	c.mutex.RUnlock()
	return result, exists
}

func (c *URLStatusCache) Check(targetURL string) URLStatusResult {
	if cached, exists := c.Get(targetURL); exists {
		return cached
	}

	result := c.fetchWithRetry(targetURL)

	c.mutex.Lock()
	c.cache[targetURL] = result
	c.mutex.Unlock()

	return result
}

func (c *URLStatusCache) fetchWithRetry(targetURL string) URLStatusResult {
	for attempt := range maxRetries {
		request, err := http.NewRequest(http.MethodHead, targetURL, nil)
		if err != nil {
			return URLStatusResult{StatusCode: 0, Error: err}
		}
		setBrowserHeaders(request)

		response, err := c.client.Do(request)
		if err != nil {
			if attempt < maxRetries-1 {
				sleepWithBackoff(attempt)
				continue
			}
			return URLStatusResult{StatusCode: 0, Error: err}
		}
		response.Body.Close()

		if response.StatusCode == 405 {
			return c.fetchGETWithRetry(targetURL, 0)
		}

		if response.StatusCode == 429 || response.StatusCode >= 500 {
			if attempt < maxRetries-1 {
				sleepWithBackoff(attempt)
				continue
			}
		}

		return URLStatusResult{StatusCode: response.StatusCode, Error: nil}
	}

	return URLStatusResult{StatusCode: 0, Error: fmt.Errorf("max retries exceeded")}
}

func (c *URLStatusCache) fetchGETWithRetry(targetURL string, startAttempt int) URLStatusResult {
	for attempt := startAttempt; attempt < maxRetries; attempt++ {
		request, err := http.NewRequest(http.MethodGet, targetURL, nil)
		if err != nil {
			return URLStatusResult{StatusCode: 0, Error: err}
		}
		setBrowserHeaders(request)

		response, err := c.client.Do(request)
		if err != nil {
			if attempt < maxRetries-1 {
				sleepWithBackoff(attempt)
				continue
			}
			return URLStatusResult{StatusCode: 0, Error: err}
		}
		response.Body.Close()

		if response.StatusCode == 429 || response.StatusCode >= 500 {
			if attempt < maxRetries-1 {
				sleepWithBackoff(attempt)
				continue
			}
		}

		return URLStatusResult{StatusCode: response.StatusCode, Error: nil}
	}

	return URLStatusResult{StatusCode: 0, Error: fmt.Errorf("max retries exceeded")}
}

func (c *URLStatusCache) CheckConcurrent(urls map[string]string, maxCheck int, concurrency int) map[string]URLStatusResult {
	results := make(map[string]URLStatusResult)
	var mutex sync.Mutex
	var waitGroup sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)
	checked := 0

	for targetURL := range urls {
		if checked >= maxCheck {
			break
		}
		checked++

		waitGroup.Add(1)
		go func(url string) {
			defer waitGroup.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := c.Check(url)
			mutex.Lock()
			results[url] = result
			mutex.Unlock()
		}(targetURL)
	}

	waitGroup.Wait()
	return results
}

func setBrowserHeaders(request *http.Request) {
	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9")
	request.Header.Set("Connection", "keep-alive")
}

func sleepWithBackoff(attempt int) {
	backoff := time.Duration(baseBackoffMs*(1<<attempt)) * time.Millisecond
	jitter := time.Duration(rand.Intn(jitterMs)) * time.Millisecond
	time.Sleep(backoff + jitter)
}
