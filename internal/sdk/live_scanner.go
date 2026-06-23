package sdk

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
)

type LiveScanner struct {
	analyzer    *analyzer.SiteAnalyzer
	siteCrawler *crawler.SiteCrawler
}

func NewLiveScanner(siteAnalyzer *analyzer.SiteAnalyzer, siteCrawler *crawler.SiteCrawler) *LiveScanner {
	return &LiveScanner{
		analyzer:    siteAnalyzer,
		siteCrawler: siteCrawler,
	}
}

func (scanner *LiveScanner) Scan(config Config) (*ScanResult, error) {
	var serverProcesses []*exec.Cmd

	var extraEnv []string
	if config.Port != 0 {
		extraEnv = append(extraEnv, fmt.Sprintf("PORT=%d", config.Port))
	}

	for _, command := range config.StartCmd {
		process, err := startServerProcess(command, extraEnv)
		if err != nil {
			stopAllServerProcesses(serverProcesses)
			return nil, fmt.Errorf("failed to start server with '%s': %w", command, err)
		}
		serverProcesses = append(serverProcesses, process)
	}
	defer stopAllServerProcesses(serverProcesses)

	for _, waitURL := range config.WaitFor {
		fmt.Fprintf(os.Stderr, "Waiting for server at %s...\n", waitURL)
		if err := waitForServer(waitURL, config.WaitTimeout); err != nil {
			return nil, fmt.Errorf("server not ready at %s: %w", waitURL, err)
		}
		fmt.Fprintf(os.Stderr, "Server at %s is ready.\n", waitURL)
	}

	fmt.Fprintf(os.Stderr, "Crawling %s (max %d pages)...\n", config.URL, config.MaxPages)

	crawlStartTime := time.Now()
	crawlResult, err := scanner.siteCrawler.Crawl(config.URL, config.MaxPages, func(page crawler.PageData, pagesCompleted int, totalDiscovered int) {
		elapsed := time.Since(crawlStartTime).Seconds()
		eta := ""
		if pagesCompleted > 2 && elapsed > 0 {
			rate := float64(pagesCompleted) / elapsed
			remaining := float64(totalDiscovered-pagesCompleted) / rate
			if remaining > 60 {
				eta = fmt.Sprintf(" (~%.0fm remaining)", remaining/60)
			} else if remaining > 0 {
				eta = fmt.Sprintf(" (~%.0fs remaining)", remaining)
			}
		}
		fmt.Fprintf(os.Stderr, "\r  [%d/%d] %s%s", pagesCompleted, totalDiscovered, page.URL, eta)
	})
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("crawl failed: %w", err)
	}

	result := &ScanResult{
		Mode:       "full",
		TotalPages: len(crawlResult.Pages),
		Pages:      make([]PageScanResult, 0, len(crawlResult.Pages)),
	}

	for _, page := range crawlResult.Pages {
		rules := scanner.analyzer.AnalyzePage(page)
		pageResult := buildPageScanResult(page.URL, rules)
		result.Pages = append(result.Pages, pageResult)
		result.AllRules = append(result.AllRules, rules...)
	}

	siteRules := scanner.analyzer.AnalyzeSite(*crawlResult)
	result.AllRules = append(result.AllRules, siteRules...)

	return result, nil
}

func startServerProcess(command string, extraEnv []string) (*exec.Cmd, error) {
	fmt.Fprintf(os.Stderr, "Starting server: %s\n", command)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if len(extraEnv) > 0 {
		cmd.Env = append(os.Environ(), extraEnv...)
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}

func stopServerProcess(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Stopping server (PID %d)...\n", cmd.Process.Pid)
	cmd.Process.Kill()
	cmd.Wait()
}

func stopAllServerProcesses(processes []*exec.Cmd) {
	for _, process := range processes {
		stopServerProcess(process)
	}
}

func waitForServer(targetURL string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if !strings.HasPrefix(targetURL, "http") {
		targetURL = "http://" + targetURL
	}

	client := &http.Client{Timeout: 2 * time.Second}
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout after %s waiting for %s", timeout, targetURL)
		case <-ticker.C:
			response, err := client.Get(targetURL)
			if err != nil {
				continue
			}
			response.Body.Close()
			if response.StatusCode < 500 {
				return nil
			}
		}
	}
}
