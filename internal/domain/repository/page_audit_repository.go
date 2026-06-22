package repository

import (
	"context"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/entity"
)

type PageAuditRepository interface {
	Create(ctx context.Context, audit *entity.PageAudit) error
	GetByCrawlJobID(ctx context.Context, crawlJobID string) ([]*entity.PageAudit, error)
	GetByID(ctx context.Context, id string) (*entity.PageAudit, error)
}
