package repository

import (
	"context"

	"github.com/booltools/booltools-seo-crawler/internal/domain/entity"
)

type CrawlJobRepository interface {
	Create(ctx context.Context, job *entity.CrawlJob) error
	GetByID(ctx context.Context, id string) (*entity.CrawlJob, error)
	Update(ctx context.Context, job *entity.CrawlJob) error
	List(ctx context.Context, limit int, offset int) ([]*entity.CrawlJob, error)
}
