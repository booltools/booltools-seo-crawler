package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
	"github.com/booltools/booltools-seo-crawler/internal/sdk"
)

func main() {
	modeFlag := flag.String("mode", "", "Scan mode: 'static' (local HTML files) or 'full' (live server crawl)")
	dirFlag := flag.String("dir", "", "Directory containing built HTML files (static mode)")
	urlFlag := flag.String("url", "", "Server URL to crawl (full mode)")
	failOnFlag := flag.String("fail-on", "", "Minimum severity to trigger failure: critical, high, medium, low, info (default: high)")
	ignoreFlag := flag.String("ignore", "", "Comma-separated rule keys to skip")
	onlyFlag := flag.String("only", "", "Comma-separated rule keys to run exclusively")
	configFlag := flag.String("config", ".seo-crawler.yml", "Path to config file")
	startCmdFlag := flag.String("start-cmd", "", "Comma-separated commands to start servers before scanning")
	waitForFlag := flag.String("wait-for", "", "Comma-separated URLs to poll until servers are ready")
	waitTimeoutFlag := flag.String("wait-timeout", "", "Max time to wait for each server (e.g. 30s, 1m)")
	formatFlag := flag.String("format", "", "Output format: text or json (default: text)")
	outputFlag := flag.String("output", "", "Write JSON report to file")
	maxPagesFlag := flag.Int("max-pages", 0, "Max pages to crawl in full mode (0 = unlimited)")
	portFlag := flag.Int("port", 0, "Port to set as PORT env var for --start-cmd processes (avoids port conflicts)")
	excludeNoindexFlag := flag.Bool("exclude-noindex", false, "Skip SEO rules on noindex pages (only technical/performance/security rules apply)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: seo-crawler [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Run SEO/GEO checks in CI/CD pipelines.\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  seo-crawler --mode=static --dir=./dist --fail-on=high\n")
		fmt.Fprintf(os.Stderr, "  seo-crawler --mode=full --url=http://localhost:3000\n")
		fmt.Fprintf(os.Stderr, "  seo-crawler --mode=full --start-cmd=\"npm run dev\" --wait-for=http://localhost:3000\n")
		fmt.Fprintf(os.Stderr, "  seo-crawler --mode=full --start-cmd=\"go run ./cmd/server,npm run dev\" --wait-for=\"http://localhost:8080,http://localhost:3000\"\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	config, err := sdk.LoadConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(sdk.ExitCodeConfigError)
	}

	applyFlagOverrides(&config, modeFlag, dirFlag, urlFlag, failOnFlag, ignoreFlag,
		onlyFlag, startCmdFlag, waitForFlag, waitTimeoutFlag, formatFlag, outputFlag, maxPagesFlag, portFlag, excludeNoindexFlag)

	if err := config.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(sdk.ExitCodeConfigError)
	}

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	siteAnalyzer.ExcludeNoindex = config.ExcludeNoindex
	reporter := sdk.NewReporter(os.Stdout)

	var result *sdk.ScanResult

	switch config.Mode {
	case "static":
		scanner := sdk.NewStaticScanner(siteAnalyzer)
		result, err = scanner.Scan(config.Dir)
	case "full":
		siteCrawler := crawler.NewSiteCrawler()
		scanner := sdk.NewLiveScanner(siteAnalyzer, siteCrawler)
		result, err = scanner.Scan(config)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(sdk.ExitCodeConfigError)
	}

	exitCode := reporter.Report(result, config)
	os.Exit(exitCode)
}

func applyFlagOverrides(
	config *sdk.Config,
	modeFlag, dirFlag, urlFlag, failOnFlag, ignoreFlag, onlyFlag,
	startCmdFlag, waitForFlag, waitTimeoutFlag, formatFlag, outputFlag *string,
	maxPagesFlag, portFlag *int,
	excludeNoindexFlag *bool,
) {
	flagWasSet := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		flagWasSet[f.Name] = true
	})

	if flagWasSet["mode"] {
		config.MergeFlag("mode", *modeFlag)
	}
	if flagWasSet["dir"] {
		config.MergeFlag("dir", *dirFlag)
	}
	if flagWasSet["url"] {
		config.MergeFlag("url", *urlFlag)
	}
	if flagWasSet["fail-on"] {
		config.MergeFlag("fail-on", *failOnFlag)
	}
	if flagWasSet["ignore"] {
		config.MergeFlag("ignore", *ignoreFlag)
	}
	if flagWasSet["only"] {
		config.MergeFlag("only", *onlyFlag)
	}
	if flagWasSet["start-cmd"] {
		config.MergeFlag("start-cmd", *startCmdFlag)
	}
	if flagWasSet["wait-for"] {
		config.MergeFlag("wait-for", *waitForFlag)
	}
	if flagWasSet["wait-timeout"] {
		duration, err := time.ParseDuration(*waitTimeoutFlag)
		if err == nil {
			config.WaitTimeout = duration
		}
	}
	if flagWasSet["format"] {
		config.MergeFlag("format", *formatFlag)
	}
	if flagWasSet["output"] {
		config.MergeFlag("output", *outputFlag)
	}
	if flagWasSet["max-pages"] && *maxPagesFlag > 0 {
		config.MaxPages = *maxPagesFlag
	}
	if flagWasSet["port"] && *portFlag > 0 {
		config.Port = *portFlag
	}
	if flagWasSet["exclude-noindex"] {
		config.ExcludeNoindex = *excludeNoindexFlag
	}
}
