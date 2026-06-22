package crawler

import (
	"fmt"
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

func NewURLStatusCache() *URLStatusCache {
	return &URLStatusCache{
		cache: make(map[string]URLStatusResult),
		client: &http.Client{
			Timeout: 10 * time.Second,
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
	return &URLStatusCache{
		cache: make(map[string]URLStatusResult),
		client: &http.Client{
			Timeout: 10 * time.Second,
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

	response, err := c.client.Head(targetURL)
	var result URLStatusResult
	if err != nil {
		result = URLStatusResult{StatusCode: 0, Error: err}
	} else {
		response.Body.Close()
		result = URLStatusResult{StatusCode: response.StatusCode, Error: nil}
	}

	c.mutex.Lock()
	c.cache[targetURL] = result
	c.mutex.Unlock()

	return result
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
