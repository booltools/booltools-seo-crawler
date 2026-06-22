package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/MarceloBD/free-seo-crawler/internal/application/usecase"
	"github.com/MarceloBD/free-seo-crawler/internal/domain/entity"
)

type SSEHandler struct {
	getProgress *usecase.GetProgressUseCase
}

func NewSSEHandler(getProgress *usecase.GetProgressUseCase) *SSEHandler {
	return &SSEHandler{getProgress: getProgress}
}

func (h *SSEHandler) StreamProgress(writer http.ResponseWriter, request *http.Request) {
	jobID := chi.URLParam(request, "id")
	if _, parseError := uuid.Parse(jobID); parseError != nil {
		http.Error(writer, "Invalid job ID format", http.StatusBadRequest)
		return
	}

	flusher, ok := writer.(http.Flusher)
	if !ok {
		http.Error(writer, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("X-Accel-Buffering", "no")

	eventChannel := h.getProgress.Subscribe(jobID)
	defer h.getProgress.Unsubscribe(jobID, eventChannel)

	ctx := request.Context()
	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-heartbeat.C:
			fmt.Fprintf(writer, ": heartbeat\n\n")
			flusher.Flush()
		case event, ok := <-eventChannel:
			if !ok {
				return
			}

			data, err := json.Marshal(event)
			if err != nil {
				continue
			}

			fmt.Fprintf(writer, "data: %s\n\n", data)
			flusher.Flush()

			if event.Status == string(entity.CrawlStatusCompleted) ||
				event.Status == string(entity.CrawlStatusFailed) {
				return
			}
		}
	}
}
