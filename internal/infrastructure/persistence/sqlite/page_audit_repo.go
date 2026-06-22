package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/booltools/booltools-seo-crawler/internal/domain/entity"
	"github.com/booltools/booltools-seo-crawler/internal/domain/valueobject"
)

type PageAuditRepository struct {
	database *sql.DB
}

func NewPageAuditRepository(database *sql.DB) *PageAuditRepository {
	return &PageAuditRepository{database: database}
}

func (r *PageAuditRepository) Create(ctx context.Context, audit *entity.PageAudit) error {
	rulesData, err := json.Marshal(audit.Rules)
	if err != nil {
		return fmt.Errorf("failed to marshal rules: %w", err)
	}

	query := `INSERT INTO page_audits (id, crawl_job_id, url, status_code, depth, rules_data, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = r.database.ExecContext(ctx, query,
		audit.ID, audit.CrawlJobID, audit.URL, audit.StatusCode,
		audit.Depth, string(rulesData), audit.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create page audit: %w", err)
	}

	return nil
}

func (r *PageAuditRepository) GetByCrawlJobID(ctx context.Context, crawlJobID string) ([]*entity.PageAudit, error) {
	query := `SELECT id, crawl_job_id, url, status_code, depth, rules_data, created_at FROM page_audits WHERE crawl_job_id = ? ORDER BY depth ASC, url ASC`

	rows, err := r.database.QueryContext(ctx, query, crawlJobID)
	if err != nil {
		return nil, fmt.Errorf("failed to query page audits: %w", err)
	}
	defer rows.Close()

	var audits []*entity.PageAudit
	for rows.Next() {
		var audit entity.PageAudit
		var rulesData string

		if err := rows.Scan(
			&audit.ID, &audit.CrawlJobID, &audit.URL, &audit.StatusCode,
			&audit.Depth, &rulesData, &audit.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan page audit: %w", err)
		}

		var rules []valueobject.AuditRule
		if err := json.Unmarshal([]byte(rulesData), &rules); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rules: %w", err)
		}
		audit.Rules = rules

		audits = append(audits, &audit)
	}

	return audits, nil
}

func (r *PageAuditRepository) GetByID(ctx context.Context, id string) (*entity.PageAudit, error) {
	query := `SELECT id, crawl_job_id, url, status_code, depth, rules_data, created_at FROM page_audits WHERE id = ?`

	var audit entity.PageAudit
	var rulesData string

	err := r.database.QueryRowContext(ctx, query, id).Scan(
		&audit.ID, &audit.CrawlJobID, &audit.URL, &audit.StatusCode,
		&audit.Depth, &rulesData, &audit.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get page audit: %w", err)
	}

	var rules []valueobject.AuditRule
	if err := json.Unmarshal([]byte(rulesData), &rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rules: %w", err)
	}
	audit.Rules = rules

	return &audit, nil
}
