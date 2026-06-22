package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/booltools/booltools-seo-crawler/internal/api/handler"
	"github.com/booltools/booltools-seo-crawler/internal/api/middleware"
	"github.com/booltools/booltools-seo-crawler/internal/application/usecase"
)

func NewRouter(
	startCrawl *usecase.StartCrawlUseCase,
	getReport *usecase.GetReportUseCase,
	getProgress *usecase.GetProgressUseCase,
) http.Handler {
	router := chi.NewRouter()

	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RealIP)
	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.SecurityHeadersMiddleware)
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.BodyLimitMiddleware)
	router.Use(middleware.CompressMiddleware)
	router.Use(middleware.TimeoutMiddleware(60 * time.Second))

	crawlHandler := handler.NewCrawlHandler(startCrawl)
	reportHandler := handler.NewReportHandler(getReport)
	sseHandler := handler.NewSSEHandler(getProgress)
	exportHandler := handler.NewExportHandler(getReport)
	rulesHandler := handler.NewRulesHandler()

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(middleware.RateLimitMiddleware(30))

		apiRouter.Get("/rules", rulesHandler.ListRules)
		apiRouter.Post("/crawl", crawlHandler.StartCrawl)
		apiRouter.Get("/report/{id}", reportHandler.GetReport)
		apiRouter.Get("/report/{id}/export/csv", exportHandler.ExportCSV)
		apiRouter.Get("/report/{id}/export/md", exportHandler.ExportMarkdown)
		apiRouter.Get("/crawl/{id}/progress", sseHandler.StreamProgress)

		apiRouter.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			writer.Write([]byte(`{"status":"ok"}`))
		})
	})

	return router
}
