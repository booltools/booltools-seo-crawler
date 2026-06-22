package handler

import (
	"encoding/json"
	"net/http"

	"github.com/booltools/booltools-seo-crawler/internal/application/dto"
	"github.com/booltools/booltools-seo-crawler/internal/application/usecase"
)

type CrawlHandler struct {
	startCrawl *usecase.StartCrawlUseCase
}

func NewCrawlHandler(startCrawl *usecase.StartCrawlUseCase) *CrawlHandler {
	return &CrawlHandler{startCrawl: startCrawl}
}

func (h *CrawlHandler) StartCrawl(writer http.ResponseWriter, request *http.Request) {
	var crawlRequest dto.CrawlRequest
	if err := json.NewDecoder(request.Body).Decode(&crawlRequest); err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	response, err := h.startCrawl.Execute(request.Context(), crawlRequest)
	if err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(writer, http.StatusAccepted, response)
}

func writeJSON(writer http.ResponseWriter, status int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(data)
}
