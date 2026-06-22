# Contributing to Booltools Seo Crawler

Thank you for your interest in contributing! This guide will help you get started.

## Getting Started

### Prerequisites

- **Go** 1.21+
- **Node.js** 20+
- **Make** (optional, for convenience commands)

### Setup

```bash
# Clone the repository
git clone https://github.com/booltools/booltools-seo-crawler.git
cd booltools-seo-crawler

# Install Go dependencies
go mod download

# Install frontend dependencies
cd web && npm install && cd ..

# Run tests to verify setup
make test
```

### Running Locally

```bash
# Terminal 1 — Backend
make dev

# Terminal 2 — Frontend
make web-dev
```

The frontend runs at `http://localhost:4321` and proxies API calls to the backend at `:8080`.

## Project Structure

```
├── cmd/                    # Entry points
│   ├── server/             # Web server
│   └── seo-crawler/        # CLI/SDK binary
├── internal/               # Application code
│   ├── api/                # HTTP handlers, router, middleware
│   ├── application/        # Use cases and DTOs
│   ├── domain/             # Entities and value objects
│   ├── infrastructure/     # Crawler, analyzer, persistence
│   └── sdk/                # CI/CD SDK logic
├── tests/                  # Test files
├── web/                    # Astro frontend
│   ├── src/pages/          # Pages (index, report, docs)
│   └── src/components/     # Astro components
└── Makefile
```

## Adding a New Rule

1. **Create the checker** in `internal/infrastructure/analyzer/rules/<category>/`
2. **Register it** in the appropriate analyzer (page or site level)
3. **Add it to the manifest** in `internal/infrastructure/analyzer/rule_manifest.go`
4. **Write tests** in `tests/infrastructure/analyzer/`
5. **Add documentation** in `web/src/data/rule-descriptions.ts`

### Rule Checker Pattern

```go
type MyChecker struct{}

func (c *MyChecker) Check(page crawler.PageData) []valueobject.AuditRule {
    var rules []valueobject.AuditRule

    rule := valueobject.NewAuditRule("my_rule_key", valueobject.CategoryOnPage, valueobject.SeverityMedium)
    rule.AffectedURL = page.URL

    if /* condition fails */ {
        rule.Fail("Description of the issue", "How to fix it")
        rule.WithDetails("Additional context")
    } else {
        rule.Pass("Everything looks good")
    }
    rules = append(rules, rule)

    return rules
}
```

## Running Tests

```bash
# All tests
make test

# Specific package
go test ./tests/infrastructure/analyzer/... -v

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Code Style

- Use descriptive variable names (no abbreviations)
- Prefer early returns over nested if/else
- Add `WithDetails()` to failing rules for better context
- Keep functions focused — split long files semantically
- Use dependency injection
- Follow Clean Architecture with DDD principles

## Pull Request Process

1. Fork the repository and create your branch from `main`
2. Make your changes following the code style above
3. Add or update tests as needed
4. Ensure all tests pass: `make test`
5. Ensure Go vet passes: `make lint`
6. Submit a pull request using the PR template

## Reporting Issues

Use the issue templates provided:

- **Bug Report** — for bugs and unexpected behavior
- **Feature Request** — for new features or improvements
- **New Rule** — for proposing new SEO/GEO audit rules
