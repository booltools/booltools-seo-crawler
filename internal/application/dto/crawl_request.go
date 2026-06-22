package dto

import (
	"fmt"
	"net/url"
	"strings"
)

type CrawlRequest struct {
	Domain        string   `json:"domain"`
	MaxPages      int      `json:"maxPages,omitempty"`
	SelectedRules []string `json:"selectedRules,omitempty"`
}

var blockedSchemes = map[string]bool{
	"file":       true,
	"ftp":        true,
	"javascript": true,
	"data":       true,
	"gopher":     true,
}

func (r *CrawlRequest) Validate() error {
	if strings.TrimSpace(r.Domain) == "" {
		return fmt.Errorf("domain is required")
	}

	if len(r.Domain) > 2048 {
		return fmt.Errorf("domain URL is too long")
	}

	domain := r.Domain
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		if isLocalDomain(domain) {
			domain = "http://" + domain
		} else {
			domain = "https://" + domain
		}
	}

	parsedURL, err := url.Parse(domain)
	if err != nil {
		return fmt.Errorf("invalid domain: %w", err)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("domain must have a valid hostname")
	}

	if blockedSchemes[parsedURL.Scheme] {
		return fmt.Errorf("scheme %q is not allowed, use http or https", parsedURL.Scheme)
	}

	hostname := parsedURL.Hostname()
	if isCloudMetadataHost(hostname) {
		return fmt.Errorf("this hostname is not allowed")
	}

	r.Domain = domain

	if r.MaxPages <= 0 {
		r.MaxPages = 0
	}

	if len(r.SelectedRules) > 500 {
		r.SelectedRules = r.SelectedRules[:500]
	}

	return nil
}

func isLocalDomain(domain string) bool {
	host := strings.Split(domain, ":")[0]
	host = strings.ToLower(host)
	return host == "localhost" || host == "127.0.0.1" || host == "0.0.0.0" || host == "::1"
}

func isCloudMetadataHost(hostname string) bool {
	metadataHosts := []string{
		"169.254.169.254",
		"metadata.google.internal",
		"metadata.google.com",
	}
	lower := strings.ToLower(hostname)
	for _, blocked := range metadataHosts {
		if lower == blocked {
			return true
		}
	}
	return false
}
