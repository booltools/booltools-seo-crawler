package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/entity"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
)

type CrawlJobRepository struct {
	database *sql.DB
}

func NewCrawlJobRepository(database *sql.DB) *CrawlJobRepository {
	return &CrawlJobRepository{database: database}
}

func (r *CrawlJobRepository) Create(ctx context.Context, job *entity.CrawlJob) error {
	seoData, err := json.Marshal(job.SeoScore)
	if err != nil {
		return fmt.Errorf("failed to marshal seo score: %w", err)
	}

	geoData, err := json.Marshal(job.GeoScore)
	if err != nil {
		return fmt.Errorf("failed to marshal geo score: %w", err)
	}

	query := `INSERT INTO crawl_jobs (id, domain, status, max_pages, pages_crawled, issues_found, seo_score_overall, seo_score_data, geo_score_overall, geo_score_data, error_message, created_at, updated_at, completed_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = r.database.ExecContext(ctx, query,
		job.ID, job.Domain, job.Status, job.MaxPages, job.PagesCrawled,
		job.IssuesFound, job.SeoScore.Overall, string(seoData),
		job.GeoScore.Overall, string(geoData), job.ErrorMessage,
		job.CreatedAt, job.UpdatedAt, job.CompletedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create crawl job: %w", err)
	}

	return nil
}

func (r *CrawlJobRepository) GetByID(ctx context.Context, id string) (*entity.CrawlJob, error) {
	query := `SELECT id, domain, status, max_pages, pages_crawled, issues_found, seo_score_overall, seo_score_data, geo_score_overall, geo_score_data, error_message, created_at, updated_at, completed_at FROM crawl_jobs WHERE id = ?`

	var job entity.CrawlJob
	var seoData, geoData string
	var completedAt sql.NullTime

	err := r.database.QueryRowContext(ctx, query, id).Scan(
		&job.ID, &job.Domain, &job.Status, &job.MaxPages, &job.PagesCrawled,
		&job.IssuesFound, &job.SeoScore.Overall, &seoData,
		&job.GeoScore.Overall, &geoData, &job.ErrorMessage,
		&job.CreatedAt, &job.UpdatedAt, &completedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get crawl job: %w", err)
	}

	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}

	if err := json.Unmarshal([]byte(seoData), &job.SeoScore); err != nil {
		return nil, fmt.Errorf("failed to unmarshal seo score: %w", err)
	}

	if err := json.Unmarshal([]byte(geoData), &job.GeoScore); err != nil {
		return nil, fmt.Errorf("failed to unmarshal geo score: %w", err)
	}

	return &job, nil
}

func (r *CrawlJobRepository) Update(ctx context.Context, job *entity.CrawlJob) error {
	seoData, err := json.Marshal(job.SeoScore)
	if err != nil {
		return fmt.Errorf("failed to marshal seo score: %w", err)
	}

	geoData, err := json.Marshal(job.GeoScore)
	if err != nil {
		return fmt.Errorf("failed to marshal geo score: %w", err)
	}

	query := `UPDATE crawl_jobs SET domain = ?, status = ?, max_pages = ?, pages_crawled = ?, issues_found = ?, seo_score_overall = ?, seo_score_data = ?, geo_score_overall = ?, geo_score_data = ?, error_message = ?, updated_at = ?, completed_at = ? WHERE id = ?`

	_, err = r.database.ExecContext(ctx, query,
		job.Domain, job.Status, job.MaxPages, job.PagesCrawled,
		job.IssuesFound, job.SeoScore.Overall, string(seoData),
		job.GeoScore.Overall, string(geoData), job.ErrorMessage,
		job.UpdatedAt, job.CompletedAt, job.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update crawl job: %w", err)
	}

	return nil
}

func (r *CrawlJobRepository) List(ctx context.Context, limit int, offset int) ([]*entity.CrawlJob, error) {
	query := `SELECT id, domain, status, max_pages, pages_crawled, issues_found, seo_score_overall, seo_score_data, geo_score_overall, geo_score_data, error_message, created_at, updated_at, completed_at FROM crawl_jobs ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.database.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list crawl jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*entity.CrawlJob
	for rows.Next() {
		var job entity.CrawlJob
		var seoData, geoData string
		var completedAt sql.NullTime

		if err := rows.Scan(
			&job.ID, &job.Domain, &job.Status, &job.MaxPages, &job.PagesCrawled,
			&job.IssuesFound, &job.SeoScore.Overall, &seoData,
			&job.GeoScore.Overall, &geoData, &job.ErrorMessage,
			&job.CreatedAt, &job.UpdatedAt, &completedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan crawl job: %w", err)
		}

		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}

		var seoScore valueobject.SeoScore
		if err := json.Unmarshal([]byte(seoData), &seoScore); err == nil {
			job.SeoScore = seoScore
		}

		var geoScore valueobject.GeoScore
		if err := json.Unmarshal([]byte(geoData), &geoScore); err == nil {
			job.GeoScore = geoScore
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}
