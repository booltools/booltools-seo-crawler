package sdk_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/booltools/booltools-seo-crawler/internal/sdk"
)

func TestDefaultConfig(t *testing.T) {
	config := sdk.DefaultConfig()

	if config.FailOn != "high" {
		t.Errorf("expected default FailOn to be 'high', got '%s'", config.FailOn)
	}
	if config.Format != "text" {
		t.Errorf("expected default Format to be 'text', got '%s'", config.Format)
	}
	if config.MaxPages != 1000 {
		t.Errorf("expected default MaxPages to be 1000, got %d", config.MaxPages)
	}
	if config.WaitTimeout != 30*time.Second {
		t.Errorf("expected default WaitTimeout to be 30s, got %s", config.WaitTimeout)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	config, err := sdk.LoadConfig("nonexistent.yml")
	if err != nil {
		t.Fatalf("expected no error for missing config, got %v", err)
	}
	if config.FailOn != "high" {
		t.Error("expected defaults to be applied when file not found")
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".seo-crawler.yml")
	configContent := `mode: static
dir: ./build
fail_on: medium
ignore:
  - og_locale
  - twitter_site
format: json
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	config, err := sdk.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.Mode != "static" {
		t.Errorf("expected mode 'static', got '%s'", config.Mode)
	}
	if config.Dir != "./build" {
		t.Errorf("expected dir './build', got '%s'", config.Dir)
	}
	if config.FailOn != "medium" {
		t.Errorf("expected fail_on 'medium', got '%s'", config.FailOn)
	}
	if len(config.Ignore) != 2 {
		t.Errorf("expected 2 ignore rules, got %d", len(config.Ignore))
	}
	if config.Format != "json" {
		t.Errorf("expected format 'json', got '%s'", config.Format)
	}
}

func TestConfig_Validate_StaticMode(t *testing.T) {
	config := sdk.DefaultConfig()
	config.Mode = "static"
	config.Dir = "./dist"

	if err := config.Validate(); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}
}

func TestConfig_Validate_StaticMissingDir(t *testing.T) {
	config := sdk.DefaultConfig()
	config.Mode = "static"

	if err := config.Validate(); err == nil {
		t.Error("expected error for static mode without dir")
	}
}

func TestConfig_Validate_FullMode(t *testing.T) {
	config := sdk.DefaultConfig()
	config.Mode = "full"
	config.URL = "http://localhost:3000"

	if err := config.Validate(); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}
}

func TestConfig_Validate_FullMissingURL(t *testing.T) {
	config := sdk.DefaultConfig()
	config.Mode = "full"

	if err := config.Validate(); err == nil {
		t.Error("expected error for full mode without url")
	}
}

func TestConfig_Validate_InvalidMode(t *testing.T) {
	config := sdk.DefaultConfig()
	config.Mode = "invalid"

	if err := config.Validate(); err == nil {
		t.Error("expected error for invalid mode")
	}
}

func TestConfig_Validate_MissingMode(t *testing.T) {
	config := sdk.DefaultConfig()

	if err := config.Validate(); err == nil {
		t.Error("expected error for missing mode")
	}
}

func TestConfig_Validate_InvalidFailOn(t *testing.T) {
	config := sdk.DefaultConfig()
	config.Mode = "static"
	config.Dir = "./dist"
	config.FailOn = "super"

	if err := config.Validate(); err == nil {
		t.Error("expected error for invalid fail-on severity")
	}
}

func TestConfig_MergeFlag(t *testing.T) {
	config := sdk.DefaultConfig()

	config.MergeFlag("mode", "full")
	if config.Mode != "full" {
		t.Errorf("expected mode 'full' after merge, got '%s'", config.Mode)
	}

	config.MergeFlag("url", "http://localhost:4000")
	if config.URL != "http://localhost:4000" {
		t.Errorf("expected url to be set, got '%s'", config.URL)
	}

	config.MergeFlag("ignore", "rule_a, rule_b, rule_c")
	if len(config.Ignore) != 3 {
		t.Errorf("expected 3 ignore rules, got %d", len(config.Ignore))
	}
}
