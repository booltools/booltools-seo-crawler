package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/booltools/booltools-seo-crawler/internal/api"
	"github.com/booltools/booltools-seo-crawler/internal/application/usecase"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/analyzer"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/crawler"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/persistence/sqlite"
	"github.com/booltools/booltools-seo-crawler/internal/infrastructure/queue"
)

func main() {
	database, err := sqlite.NewConnection("data.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := sqlite.RunMigrations(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	crawlJobRepo := sqlite.NewCrawlJobRepository(database)
	pageAuditRepo := sqlite.NewPageAuditRepository(database)

	siteAnalyzer := analyzer.NewSiteAnalyzer()
	siteCrawler := crawler.NewSiteCrawler()

	progressBroker := queue.NewProgressBroker()
	workerPool := queue.NewWorkerPool(3, crawlJobRepo, pageAuditRepo, siteCrawler, siteAnalyzer, progressBroker)

	startCrawl := usecase.NewStartCrawlUseCase(crawlJobRepo, workerPool)
	getReport := usecase.NewGetReportUseCase(crawlJobRepo, pageAuditRepo)
	getProgress := usecase.NewGetProgressUseCase(progressBroker)

	router := api.NewRouter(startCrawl, getReport, getProgress)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	go func() {
		log.Printf("server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	workerPool.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	workerPool.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}
