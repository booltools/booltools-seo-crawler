package dto

import "time"

type CrawlResponse struct {
	ID           string     `json:"id"`
	Domain       string     `json:"domain"`
	Status       string     `json:"status"`
	MaxPages     int        `json:"maxPages"`
	PagesCrawled int        `json:"pagesCrawled"`
	IssuesFound  int        `json:"issuesFound"`
	ErrorMessage string     `json:"errorMessage,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
}
