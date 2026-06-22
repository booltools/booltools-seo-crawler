package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/booltools/booltools-seo-crawler/internal/application/usecase"
)

type ReportHandler struct {
	getReport *usecase.GetReportUseCase
}

func NewReportHandler(getReport *usecase.GetReportUseCase) *ReportHandler {
	return &ReportHandler{getReport: getReport}
}

func (h *ReportHandler) GetReport(writer http.ResponseWriter, request *http.Request) {
	jobID := chi.URLParam(request, "id")
	if _, parseError := uuid.Parse(jobID); parseError != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{
			"error": "Invalid job ID format",
		})
		return
	}

	report, err := h.getReport.Execute(request.Context(), jobID)
	if err != nil {
		writeJSON(writer, http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(writer, http.StatusOK, report)
}
