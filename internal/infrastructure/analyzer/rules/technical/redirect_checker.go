package technical

import (
	"fmt"
	"net/http"
	"time"

	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type RedirectChecker struct{}

func (c *RedirectChecker) Check(result crawler.CrawlResult) []valueobject.AuditRule {
	var rules []valueobject.AuditRule

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	redirectChains := 0
	temporaryRedirects := 0
	maxChecked := min(len(result.Pages), 50)

	for _, page := range result.Pages[:maxChecked] {
		response, err := client.Head(page.URL)
		if err != nil {
			continue
		}
		response.Body.Close()

		if response.Request.URL.String() != page.URL {
			chainLength := 0
			tempRedirect := false

			checkClient := &http.Client{
				Timeout: 10 * time.Second,
				CheckRedirect: func(request *http.Request, via []*http.Request) error {
					chainLength = len(via)
					return nil
				},
			}

			headResp, headErr := checkClient.Head(page.URL)
			if headErr == nil {
				headResp.Body.Close()
			}

			if chainLength > 1 {
				redirectChains++
			}

			initialResp, initialErr := http.DefaultTransport.RoundTrip(&http.Request{
				Method: "HEAD",
				URL:    response.Request.URL,
			})
			if initialErr == nil {
				if initialResp.StatusCode == 302 || initialResp.StatusCode == 307 {
					tempRedirect = true
					temporaryRedirects++
				}
				initialResp.Body.Close()
			}
			_ = tempRedirect
		}
	}

	chainRule := valueobject.NewAuditRule("redirect_chains", valueobject.CategoryTechnical, valueobject.SeverityMedium)
	if redirectChains > 0 {
		chainRule.Warn(
			fmt.Sprintf("%d redirect chains detected (more than 1 hop)", redirectChains),
			"Eliminate redirect chains by updating links to point directly to the final destination URL.",
		)
	} else {
		chainRule.Pass("No redirect chains detected")
	}
	rules = append(rules, chainRule)

	tempRedirectRule := valueobject.NewAuditRule("temporary_redirects", valueobject.CategoryTechnical, valueobject.SeverityLow)
	if temporaryRedirects > 0 {
		tempRedirectRule.Warn(
			fmt.Sprintf("%d temporary (302/307) redirects found that may need to be permanent (301)", temporaryRedirects),
			"Convert temporary redirects to permanent 301 redirects if the move is permanent. Search engines pass more link equity through 301 redirects.",
		)
	} else {
		tempRedirectRule.Pass("No unnecessary temporary redirects")
	}
	rules = append(rules, tempRedirectRule)

	return rules
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
