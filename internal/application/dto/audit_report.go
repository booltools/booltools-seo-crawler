package dto

type AuditReport struct {
	Job              CrawlResponse          `json:"job"`
	SeoScore         ScoreDetail            `json:"seoScore"`
	GeoScore         GeoScoreDetail         `json:"geoScore"`
	Pages            []PageReport           `json:"pages"`
	IssuesBySeverity map[string][]IssueItem `json:"issuesBySeverity"`
	LinksToChange    []LinkChange           `json:"linksToChange"`
	Summary          ReportSummary          `json:"summary"`
}

type ScoreDetail struct {
	Overall        float64            `json:"overall"`
	Grade          string             `json:"grade"`
	CategoryScores map[string]float64 `json:"categoryScores"`
	TotalRules     int                `json:"totalRules"`
	PassedRules    int                `json:"passedRules"`
	FailedRules    int                `json:"failedRules"`
	WarningRules   int                `json:"warningRules"`
	CriticalIssues int                `json:"criticalIssues"`
}

type GeoScoreDetail struct {
	Overall       float64 `json:"overall"`
	Grade         string  `json:"grade"`
	CrawlerAccess float64 `json:"crawlerAccess"`
	LlmsTxt       float64 `json:"llmsTxt"`
	Citability    float64 `json:"citability"`
	EntityAuth    float64 `json:"entityAuth"`
	AIFriendly    float64 `json:"aiFriendly"`
	TotalRules    int     `json:"totalRules"`
	PassedRules   int     `json:"passedRules"`
	FailedRules   int     `json:"failedRules"`
}

type PageReport struct {
	URL        string      `json:"url"`
	StatusCode int         `json:"statusCode"`
	Depth      int         `json:"depth"`
	Issues     []IssueItem `json:"issues"`
	Passes     int         `json:"passes"`
	Warnings   int         `json:"warnings"`
	Failures   int         `json:"failures"`
}

type IssueItem struct {
	RuleKey        string `json:"ruleKey"`
	Category       string `json:"category"`
	CategoryLabel  string `json:"categoryLabel"`
	Severity       string `json:"severity"`
	Result         string `json:"result"`
	Message        string `json:"message"`
	Recommendation string `json:"recommendation"`
	AffectedURL    string `json:"affectedUrl,omitempty"`
	Details        string `json:"details,omitempty"`
}

type LinkChange struct {
	URL     string       `json:"url"`
	Changes []ChangeItem `json:"changes"`
}

type ChangeItem struct {
	Category       string `json:"category"`
	Severity       string `json:"severity"`
	Issue          string `json:"issue"`
	Recommendation string `json:"recommendation"`
}

type ReportSummary struct {
	TotalPages       int `json:"totalPages"`
	TotalIssues      int `json:"totalIssues"`
	CriticalIssues   int `json:"criticalIssues"`
	HighIssues       int `json:"highIssues"`
	MediumIssues     int `json:"mediumIssues"`
	LowIssues        int `json:"lowIssues"`
	InfoIssues       int `json:"infoIssues"`
	PassedChecks     int `json:"passedChecks"`
}
