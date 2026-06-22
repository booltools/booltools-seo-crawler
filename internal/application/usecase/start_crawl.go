package usecase

import (
	"context"
	"fmt"

	"github.com/MarceloBD/free-seo-crawler/internal/application/dto"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/entity"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/repository"
	"github.com/MarceloBD/free-seo-crawler/internal/infrastructure/queue"
)

type StartCrawlUseCase struct {
	crawlJobRepo repository.CrawlJobRepository
	workerPool   *queue.WorkerPool
}

func NewStartCrawlUseCase(
	crawlJobRepo repository.CrawlJobRepository,
	workerPool *queue.WorkerPool,
) *StartCrawlUseCase {
	return &StartCrawlUseCase{
		crawlJobRepo: crawlJobRepo,
		workerPool:   workerPool,
	}
}

func (uc *StartCrawlUseCase) Execute(ctx context.Context, request dto.CrawlRequest) (*dto.CrawlResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	job, err := entity.NewCrawlJob(request.Domain, request.MaxPages, request.SelectedRules)
	if err != nil {
		return nil, fmt.Errorf("failed to create crawl job: %w", err)
	}

	if err := uc.crawlJobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to save crawl job: %w", err)
	}

	uc.workerPool.Enqueue(job)

	return &dto.CrawlResponse{
		ID:           job.ID,
		Domain:       job.Domain,
		Status:       string(job.Status),
		MaxPages:     job.MaxPages,
		PagesCrawled: job.PagesCrawled,
		IssuesFound:  job.IssuesFound,
		CreatedAt:    job.CreatedAt,
	}, nil
}
