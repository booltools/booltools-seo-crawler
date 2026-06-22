package sdk_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/MarceloBD/free-seo-crawler/internal/domain/valueobject"
	"github.com/MarceloBD/free-seo-crawler/internal/sdk"
)

func TestReporter_ExitCodePass(t *testing.T) {
	result := &sdk.ScanResult{
		Mode:       "static",
		TotalPages: 1,
		AllRules: []valueobject.AuditRule{
			{Key: "title_exists", Result: valueobject.RuleResultPass, Severity: valueobject.SeverityHigh},
			{Key: "h1_count", Result: valueobject.RuleResultPass, Severity: valueobject.SeverityHigh},
		},
	}

	config := sdk.DefaultConfig()
	config.Mode = "static"
	config.Dir = "./dist"

	var buffer bytes.Buffer
	reporter := sdk.NewReporter(&buffer)
	exitCode := reporter.Report(result, config)

	if exitCode != sdk.ExitCodePass {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	output := buffer.String()
	if !strings.Contains(output, "PASS") {
		t.Error("expected PASS verdict in output")
	}
}

func TestReporter_ExitCodeFail(t *testing.T) {
	result := &sdk.ScanResult{
		Mode:       "static",
		TotalPages: 1,
		AllRules: []valueobject.AuditRule{
			{Key: "title_exists", Result: valueobject.RuleResultFail, Severity: valueobject.SeverityHigh, Message: "Title is missing"},
		},
	}

	config := sdk.DefaultConfig()
	config.Mode = "static"
	config.Dir = "./dist"

	var buffer bytes.Buffer
	reporter := sdk.NewReporter(&buffer)
	exitCode := reporter.Report(result, config)

	if exitCode != sdk.ExitCodeFail {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	output := buffer.String()
	if !strings.Contains(output, "FAIL") {
		t.Error("expected FAIL verdict in output")
	}
}

func TestReporter_FailOnThreshold(t *testing.T) {
	result := &sdk.ScanResult{
		Mode:       "static",
		TotalPages: 1,
		AllRules: []valueobject.AuditRule{
			{Key: "og_locale", Result: valueobject.RuleResultFail, Severity: valueobject.SeverityLow, Message: "Missing og:locale"},
		},
	}

	config := sdk.DefaultConfig()
	config.Mode = "static"
	config.Dir = "./dist"
	config.FailOn = "high"

	var buffer bytes.Buffer
	reporter := sdk.NewReporter(&buffer)
	exitCode := reporter.Report(result, config)

	if exitCode != sdk.ExitCodePass {
		t.Errorf("expected pass (low severity issue with high threshold), got %d", exitCode)
	}
}

func TestReporter_JSONFormat(t *testing.T) {
	result := &sdk.ScanResult{
		Mode:       "static",
		TotalPages: 2,
		AllRules: []valueobject.AuditRule{
			{Key: "title_exists", Result: valueobject.RuleResultFail, Severity: valueobject.SeverityHigh, Message: "Missing title"},
			{Key: "h1_count", Result: valueobject.RuleResultPass, Severity: valueobject.SeverityHigh},
		},
	}

	config := sdk.DefaultConfig()
	config.Mode = "static"
	config.Dir = "./dist"
	config.Format = "json"

	var buffer bytes.Buffer
	reporter := sdk.NewReporter(&buffer)
	reporter.Report(result, config)

	output := buffer.String()
	if !strings.Contains(output, `"verdict"`) {
		t.Error("expected JSON output with verdict field")
	}
	if !strings.Contains(output, `"totalPages"`) {
		t.Error("expected JSON output with totalPages field")
	}
}

func TestReporter_IgnoreFilter(t *testing.T) {
	result := &sdk.ScanResult{
		Mode:       "static",
		TotalPages: 1,
		AllRules: []valueobject.AuditRule{
			{Key: "title_exists", Result: valueobject.RuleResultFail, Severity: valueobject.SeverityHigh, Message: "Missing title"},
		},
	}

	config := sdk.DefaultConfig()
	config.Mode = "static"
	config.Dir = "./dist"
	config.Ignore = []string{"title_exists"}

	var buffer bytes.Buffer
	reporter := sdk.NewReporter(&buffer)
	exitCode := reporter.Report(result, config)

	if exitCode != sdk.ExitCodePass {
		t.Errorf("expected pass after ignoring the only failing rule, got %d", exitCode)
	}
}
