package controller

import (
	"arjunagl/htmlAnalyzer/infrastructure"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type AnalyzeRequest struct {
	SiteURL string `json:"site_url"`
}

type AnalyzeController struct {
	app *infrastructure.Application
}

func (h *AnalyzeController) Analyze(w http.ResponseWriter, r *http.Request) {
	var analyzeReq AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&analyzeReq); err != nil || strings.TrimSpace(analyzeReq.SiteURL) == "" {
		// Handle error if the request body cannot be decoded
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// do a quick check to see if the site can be reached
	if isAcceesible := h.app.ContentDownloader.IsLinkAccessible(analyzeReq.SiteURL); !isAcceesible {
		http.Error(w, "404 Site is inaccessible", http.StatusNotFound)
		return
	}

	id := uuid.New().String()
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, id)
	h.app.SchedulerService.Analyze(id, analyzeReq.SiteURL)
}

func (h *AnalyzeController) GetResults(w http.ResponseWriter, r *http.Request) {
	requestID := strings.TrimSpace(r.URL.Query().Get("request_id"))
	if requestID == "" {
		// Handle the case where "request_id" is not provided or is empty
		http.Error(w, "request_id is required", http.StatusBadRequest)
		return
	}

	result, exists := h.app.SchedulerService.GetResult(requestID)
	if !exists {
		http.Error(w, "invalid request_id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewAnalyzeController(app *infrastructure.Application) *AnalyzeController {
	return &AnalyzeController{
		app: app,
	}
}
