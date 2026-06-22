CREATE TABLE IF NOT EXISTS crawl_jobs (
    id TEXT PRIMARY KEY,
    domain TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'queued',
    max_pages INTEGER NOT NULL DEFAULT 50,
    pages_crawled INTEGER NOT NULL DEFAULT 0,
    issues_found INTEGER NOT NULL DEFAULT 0,
    seo_score_overall REAL NOT NULL DEFAULT 0,
    seo_score_data TEXT NOT NULL DEFAULT '{}',
    geo_score_overall REAL NOT NULL DEFAULT 0,
    geo_score_data TEXT NOT NULL DEFAULT '{}',
    error_message TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME
);

CREATE TABLE IF NOT EXISTS page_audits (
    id TEXT PRIMARY KEY,
    crawl_job_id TEXT NOT NULL,
    url TEXT NOT NULL,
    status_code INTEGER NOT NULL DEFAULT 0,
    depth INTEGER NOT NULL DEFAULT 0,
    rules_data TEXT NOT NULL DEFAULT '[]',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (crawl_job_id) REFERENCES crawl_jobs(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_page_audits_crawl_job_id ON page_audits(crawl_job_id);
CREATE INDEX IF NOT EXISTS idx_crawl_jobs_status ON crawl_jobs(status);
CREATE INDEX IF NOT EXISTS idx_crawl_jobs_created_at ON crawl_jobs(created_at);
