package queue

import (
	"context"
	"log"
	"sync"

	"github.com/booltools/booltools-seo-crawler/internal/domain/entity"
	"github.com/booltools/booltools-seo-crawler/internal/domain/repository"
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type WorkerPool struct {
	workerCount    int
	jobChannel     chan *entity.CrawlJob
	stopChannel    chan struct{}
	waitGroup      sync.WaitGroup
	crawlJobRepo   repository.CrawlJobRepository
	pageAuditRepo  repository.PageAuditRepository
	siteCrawler    *crawler.SiteCrawler
	siteAnalyzer   *analyzer.SiteAnalyzer
	progressBroker *ProgressBroker
}

func NewWorkerPool(
	workerCount int,
	crawlJobRepo repository.CrawlJobRepository,
	pageAuditRepo repository.PageAuditRepository,
	siteCrawler *crawler.SiteCrawler,
	siteAnalyzer *analyzer.SiteAnalyzer,
	progressBroker *ProgressBroker,
) *WorkerPool {
	return &WorkerPool{
		workerCount:    workerCount,
		jobChannel:     make(chan *entity.CrawlJob, 100),
		stopChannel:    make(chan struct{}),
		crawlJobRepo:   crawlJobRepo,
		pageAuditRepo:  pageAuditRepo,
		siteCrawler:    siteCrawler,
		siteAnalyzer:   siteAnalyzer,
		progressBroker: progressBroker,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.waitGroup.Add(1)
		go wp.worker(i)
	}
	log.Printf("worker pool started with %d workers", wp.workerCount)
}

func (wp *WorkerPool) Stop() {
	close(wp.stopChannel)
	wp.waitGroup.Wait()
	log.Println("worker pool stopped")
}

func (wp *WorkerPool) Enqueue(job *entity.CrawlJob) {
	wp.jobChannel <- job
}

func (wp *WorkerPool) worker(workerID int) {
	defer wp.waitGroup.Done()

	for {
		select {
		case <-wp.stopChannel:
			return
		case job := <-wp.jobChannel:
			wp.processJob(workerID, job)
		}
	}
}

func (wp *WorkerPool) processJob(workerID int, job *entity.CrawlJob) {
	ctx := context.Background()
	log.Printf("[worker %d] processing job %s for domain %s", workerID, job.ID, job.Domain)

	job.StartCrawling()
	if err := wp.crawlJobRepo.Update(ctx, job); err != nil {
		log.Printf("[worker %d] failed to update job status: %v", workerID, err)
		return
	}

	wp.progressBroker.Publish(ProgressEvent{
		JobID:   job.ID,
		Status:  string(entity.CrawlStatusCrawling),
		Message: "Starting crawl...",
	})

	targetURL := job.NormalizedDomain()
	var allRules []valueobject.AuditRule
	totalIssues := 0

	crawlResult, err := wp.siteCrawler.Crawl(targetURL, job.MaxPages, func(page crawler.PageData, pagesCompleted int, totalDiscovered int) {
		pageRules := filterRulesBySelection(wp.siteAnalyzer.AnalyzePage(page), job)

		pageAudit, auditErr := entity.NewPageAudit(job.ID, page.URL, page.StatusCode, page.Depth)
		if auditErr != nil {
			log.Printf("[worker %d] failed to create page audit: %v", workerID, auditErr)
			return
		}
		pageAudit.AddRules(pageRules)
		allRules = append(allRules, pageRules...)

		if auditErr := wp.pageAuditRepo.Create(ctx, pageAudit); auditErr != nil {
			log.Printf("[worker %d] failed to save page audit: %v", workerID, auditErr)
		}

		job.PagesCrawled = pagesCompleted
		totalIssues += pageAudit.IssueCount()

		var pageIssues []PageIssue
		for _, rule := range pageRules {
			if rule.Result == valueobject.RuleResultFail || rule.Result == valueobject.RuleResultWarning {
				truncatedDetails := rule.Details
				if len(truncatedDetails) > 200 {
					truncatedDetails = truncatedDetails[:200] + "..."
				}
				pageIssues = append(pageIssues, PageIssue{
					RuleKey:  rule.Key,
					Severity: string(rule.Severity),
					Result:   string(rule.Result),
					Message:  rule.Message,
					Details:  truncatedDetails,
				})
			}
		}

		wp.progressBroker.Publish(ProgressEvent{
			JobID:           job.ID,
			Status:          string(entity.CrawlStatusCrawling),
			PagesCrawled:    pagesCompleted,
			TotalDiscovered: totalDiscovered,
			IssuesFound:     totalIssues,
			CurrentURL:      page.URL,
			PageIssues:      pageIssues,
		})
	})

	if err != nil {
		job.Fail(err.Error())
		wp.crawlJobRepo.Update(ctx, job)
		wp.progressBroker.Publish(ProgressEvent{
			JobID:   job.ID,
			Status:  string(entity.CrawlStatusFailed),
			Message: err.Error(),
		})
		return
	}

	job.StartAnalyzing()
	wp.crawlJobRepo.Update(ctx, job)

	wp.progressBroker.Publish(ProgressEvent{
		JobID:   job.ID,
		Status:  string(entity.CrawlStatusAnalyzing),
		Message: "Running site-level analysis...",
	})

	siteRules := filterRulesBySelection(wp.siteAnalyzer.AnalyzeSite(*crawlResult), job)
	allRules = append(allRules, siteRules...)

	siteAudit, _ := entity.NewPageAudit(job.ID, targetURL, 200, 0)
	if siteAudit != nil {
		siteAudit.AddRules(siteRules)
		wp.pageAuditRepo.Create(ctx, siteAudit)
		totalIssues += siteAudit.IssueCount()
	}

	seoScore := valueobject.NewSeoScore()
	seoScore.Calculate(allRules)

	geoScore := valueobject.NewGeoScore()
	geoScore.Calculate(allRules)

	job.Complete(seoScore, geoScore, totalIssues)
	wp.crawlJobRepo.Update(ctx, job)

	wp.progressBroker.Publish(ProgressEvent{
		JobID:        job.ID,
		Status:       string(entity.CrawlStatusCompleted),
		PagesCrawled: job.PagesCrawled,
		IssuesFound:  totalIssues,
		Message:      "Crawl complete!",
	})

	log.Printf("[worker %d] completed job %s: %d pages, %d issues, SEO: %.1f, GEO: %.1f",
		workerID, job.ID, job.PagesCrawled, totalIssues, seoScore.Overall, geoScore.Overall)
}

func filterRulesBySelection(rules []valueobject.AuditRule, job *entity.CrawlJob) []valueobject.AuditRule {
	if !job.HasRuleFilter() {
		return rules
	}
	filtered := make([]valueobject.AuditRule, 0, len(rules))
	for _, rule := range rules {
		if job.IsRuleSelected(rule.Key) {
			filtered = append(filtered, rule)
		}
	}
	return filtered
}
