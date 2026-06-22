package sdk_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer"
	"github.com/booltools/booltools-seo-crawler/internal/sdk"
)

func TestStaticScanner_ValidDirectory(t *testing.T) {
	tempDir := t.TempDir()

	validHTML := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Test Page With Good Title Here</title>
	<meta name="description" content="This is a good meta description that is long enough to pass the length check for SEO purposes.">
	<link rel="canonical" href="https://example.com/test">
</head>
<body>
	<nav>Navigation</nav>
	<main>
		<h1>Welcome to the Test Page</h1>
		<p>This is a paragraph with enough content to avoid the thin content warning from the checker.</p>
		<img src="hero.webp" alt="Hero image showing the product" width="800" height="600" loading="lazy">
		<a href="/about">About Us</a>
	</main>
</body>
</html>`

	if err := os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(validHTML), 0644); err != nil {
		t.Fatalf("failed to write test HTML: %v", err)
	}

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if result.TotalPages != 1 {
		t.Errorf("expected 1 page, got %d", result.TotalPages)
	}

	if result.Mode != "static" {
		t.Errorf("expected mode 'static', got '%s'", result.Mode)
	}

	if len(result.AllRules) == 0 {
		t.Error("expected rules to be generated")
	}

	if len(result.Pages) != 1 {
		t.Errorf("expected 1 page result, got %d", len(result.Pages))
	}
}

func TestStaticScanner_InvalidHTML(t *testing.T) {
	tempDir := t.TempDir()

	invalidHTML := `<!DOCTYPE html>
<html>
<head></head>
<body>
	<p>Minimal page with many missing SEO elements.</p>
</body>
</html>`

	if err := os.WriteFile(filepath.Join(tempDir, "bad.html"), []byte(invalidHTML), 0644); err != nil {
		t.Fatalf("failed to write test HTML: %v", err)
	}

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if len(result.Pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(result.Pages))
	}

	pageResult := result.Pages[0]
	if pageResult.Failures == 0 {
		t.Error("expected failures for a page missing many SEO elements")
	}
}

func TestStaticScanner_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	_, err := scanner.Scan(tempDir)
	if err == nil {
		t.Error("expected error for empty directory")
	}
}

func TestStaticScanner_NonexistentDirectory(t *testing.T) {
	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	_, err := scanner.Scan("/nonexistent/path/to/nowhere")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestStaticScanner_NestedHTMLFiles(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "pages")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	html := `<!DOCTYPE html><html lang="en"><head><title>Page</title></head><body><h1>Content</h1></body></html>`
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(html), 0644)
	os.WriteFile(filepath.Join(subDir, "about.html"), []byte(html), 0644)
	os.WriteFile(filepath.Join(subDir, "contact.htm"), []byte(html), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if result.TotalPages != 3 {
		t.Errorf("expected 3 pages (including .htm), got %d", result.TotalPages)
	}
}

func TestStaticScanner_MultiplePages(t *testing.T) {
	tempDir := t.TempDir()

	goodHTML := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Good Page Title for SEO</title>
	<meta name="description" content="A well-crafted meta description for the SEO checker to validate properly.">
</head>
<body><nav>Nav</nav><main><h1>Title</h1><p>Content here.</p></main></body>
</html>`

	badHTML := `<html><body><p>Bad page</p></body></html>`

	os.WriteFile(filepath.Join(tempDir, "good.html"), []byte(goodHTML), 0644)
	os.WriteFile(filepath.Join(tempDir, "bad.html"), []byte(badHTML), 0644)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	scanner := sdk.NewStaticScanner(siteAnalyzer)

	result, err := scanner.Scan(tempDir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if result.TotalPages != 2 {
		t.Errorf("expected 2 pages, got %d", result.TotalPages)
	}

	if len(result.AllRules) == 0 {
		t.Error("expected rules to be populated from both pages")
	}
}
