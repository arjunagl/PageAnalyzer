package controller

import (
	"arjunagl/htmlAnalyzer/domain/model"
	"arjunagl/htmlAnalyzer/domain/service"
	"arjunagl/htmlAnalyzer/infrastructure"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newMockedAnalyzeController() *AnalyzeController {
	app := &infrastructure.Application{
		SchedulerService:  new(service.MockSchedulerService),
		ContentDownloader: new(service.MockContentDownloader),
	}
	return NewAnalyzeController(app)
}

func TestAnalyzeSubmit(t *testing.T) {
	// controller := newMockedAnalyzeController()
	var controller *AnalyzeController

	tests := []struct {
		name         string
		body         string
		expectedCode int
		mockSetup    func()
	}{
		{
			name:         "Valid Request",
			body:         `{"site_url": "http://example.com"}`,
			expectedCode: http.StatusAccepted,
			mockSetup: func() {
				controller = newMockedAnalyzeController()
				mockSchedulerService := controller.app.SchedulerService.(*service.MockSchedulerService)
				mockContentDownloader := controller.app.ContentDownloader.(*service.MockContentDownloader)
				mockContentDownloader.On("IsLinkAccessible", mock.AnythingOfType("string")).Return(true).Once()
				mockSchedulerService.On("Analyze", mock.AnythingOfType("string"), "http://example.com").Return(nil).Once()
			},
		},
		{
			name:         "Invalid Request",
			body:         `{"": ""}`,
			expectedCode: http.StatusBadRequest,
			mockSetup: func() {
				controller = newMockedAnalyzeController()
				mockSchedulerService := controller.app.SchedulerService.(*service.MockSchedulerService)
				mockContentDownloader := controller.app.ContentDownloader.(*service.MockContentDownloader)
				mockContentDownloader.On("IsLinkAccessible", mock.AnythingOfType("string")).Return(true).Once()
				mockSchedulerService.On("Analyze", mock.AnythingOfType("string"), "http://example.com").Return(nil).Once()
			},
		},
		{
			name:         "Unreachable site",
			body:         `{"site_url": "http://example.com"}`,
			expectedCode: http.StatusNotFound,
			mockSetup: func() {
				controller = newMockedAnalyzeController()
				mockSchedulerService := controller.app.SchedulerService.(*service.MockSchedulerService)
				mockContentDownloader := controller.app.ContentDownloader.(*service.MockContentDownloader)
				mockContentDownloader.On("IsLinkAccessible", mock.AnythingOfType("string")).Return(false).Once()
				mockSchedulerService.On("Analyze", mock.AnythingOfType("string"), "http://example.com").Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			reqBody := strings.NewReader(tt.body)
			req, err := http.NewRequest("POST", "/analyze", reqBody)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.Analyze)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedCode == http.StatusAccepted {
				responseBody := rr.Body.String()
				uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}`)
				assert.True(t, uuidRegex.MatchString(responseBody), "Expected response body to be a valid UUID")
				controller.app.SchedulerService.(*service.MockSchedulerService).AssertExpectations(t)
				// controller.app.SchedulerService.AssertExpectations(t)
			}
		})
	}

}

func TestAnalyzeGetResults(t *testing.T) {
	controller := newMockedAnalyzeController()
	mockSchedulerService := controller.app.SchedulerService.(*service.MockSchedulerService)
	// mockContentDownloader := controller.app.ContentDownloader.(*service.MockContentDownloader)

	tests := []struct {
		name           string
		expectedCode   int
		url            string
		expectedResult *model.AnalysisResultWithStatus
		mockSetup      func()
	}{
		{
			name:         "Existing Request",
			url:          "?request_id=1234",
			expectedCode: http.StatusOK,
			expectedResult: &model.AnalysisResultWithStatus{
				Status: model.StatusComplete,
				Result: map[string]*model.AnalyzerResult{"title": {Result: "test1", Err: nil}},
			},
			mockSetup: func() {
				mockSchedulerService.On("GetResult", mock.AnythingOfType("string")).Return(&model.AnalysisResultWithStatus{
					Status: model.StatusComplete,
					Result: map[string]*model.AnalyzerResult{"title": {Result: "test1", Err: nil}},
				}, true).Once()
			},
		},
		{
			name:         "NonExisting Request in the system",
			url:          "?request_id=1234",
			expectedCode: http.StatusBadRequest,
			expectedResult: &model.AnalysisResultWithStatus{
				Status: model.StatusComplete,
				Result: nil,
			},
			mockSetup: func() {
				mockSchedulerService.On("GetResult", mock.AnythingOfType("string")).Return(nil, false).Once()
			},
		},
		{
			name:         "Missing Request id",
			url:          "",
			expectedCode: http.StatusBadRequest,
			expectedResult: &model.AnalysisResultWithStatus{
				Status: model.StatusComplete,
				Result: map[string]*model.AnalyzerResult{"title": {Result: "test1", Err: nil}},
			},
			mockSetup: func() {
				mockSchedulerService.On("GetResult", mock.AnythingOfType("string")).Return(nil, false).Once()
			},
		},
		{
			name:         "Request still being processed",
			url:          "",
			expectedCode: http.StatusBadRequest,
			expectedResult: &model.AnalysisResultWithStatus{
				Status: model.StatusInProgress,
				Result: nil,
			},
			mockSetup: func() {
				mockSchedulerService.On("GetResult", mock.AnythingOfType("string")).Return(nil, true).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, err := http.NewRequest("GET", fmt.Sprintf("/analyze%s", tt.url), nil)
			assert.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.GetResults)

			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedCode, rr.Code)

			if tt.expectedCode == http.StatusOK {
				var actualResult model.AnalysisResultWithStatus
				responseBody := rr.Body.String()
				json.Unmarshal([]byte(responseBody), &actualResult)
				assert.Equal(t, tt.expectedResult, &actualResult)
			}
		})
	}

}
