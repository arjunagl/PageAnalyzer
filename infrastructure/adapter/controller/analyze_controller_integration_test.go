package controller

import (
	"arjunagl/htmlAnalyzer/domain/model"
	"arjunagl/htmlAnalyzer/domain/service"
	"arjunagl/htmlAnalyzer/infrastructure"
	"arjunagl/htmlAnalyzer/infrastructure/adapter"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newIntegratedMockedAnalyzeController() *AnalyzeController {
	wcd := new(service.MockContentDownloader)
	tas := &service.TitleAnalyzeService{}
	vas := &service.VersionAnalyzeService{}
	has := &service.HeadingAnalyzeService{}
	las := service.NewLinkAnalyzeService(wcd)
	lgas := &service.LoginAnalyzeService{}
	cr := adapter.NewGoqueryAdapter()
	as := service.NewAnalyzerService([]*service.ItemAnalyzeService{
		{Title: "Title", ContentAnalyzer: tas},
		{Title: "Version", ContentAnalyzer: vas},
		{Title: "Link", ContentAnalyzer: las},
		{Title: "Login", ContentAnalyzer: lgas},
		{Title: "Heading", ContentAnalyzer: has}}, cr, wcd)
	ss := service.NewSchedulerService(as)

	app := infrastructure.NewApplication(wcd, ss, cr)
	return NewAnalyzeController(app)
}

func ReadTestContent(file string) string {
	// Read the entire file
	fileContent, err := os.ReadFile(fmt.Sprintf("./test/%s", file))
	if err != nil {
		fmt.Println(err)
	}
	return string(fileContent)
}

func TestAnalyze(t *testing.T) {
	var controller *AnalyzeController
	tests := []struct {
		name           string
		expectedCode   int
		expectedResult string
		mockSetup      func()
	}{
		{
			name:           "Plain html",
			expectedCode:   http.StatusOK,
			expectedResult: ReadTestContent("1_result.json"),
			mockSetup: func() {
				controller = newIntegratedMockedAnalyzeController()
				mockContentDownloader := controller.app.ContentDownloader.(*service.MockContentDownloader)
				mockContentDownloader.On("IsLinkAccessible", mock.AnythingOfType("string")).Return(true)
				mockContentDownloader.On("DownloadContent", mock.AnythingOfType("string")).Return(ReadTestContent("1_html.html"), nil).Once()
			},
		},
		{
			name:           "All elements inside",
			expectedCode:   http.StatusOK,
			expectedResult: ReadTestContent("2_result.json"),
			mockSetup: func() {
				controller = newIntegratedMockedAnalyzeController()
				mockContentDownloader := controller.app.ContentDownloader.(*service.MockContentDownloader)
				mockContentDownloader.On("IsLinkAccessible", mock.AnythingOfType("string")).Return(true)
				mockContentDownloader.On("DownloadContent", mock.AnythingOfType("string")).Return(ReadTestContent("2_html.html"), nil).Once()
			},
		},
		{
			name:           "Actual go site",
			expectedCode:   http.StatusOK,
			expectedResult: ReadTestContent("go_site_result.json"),
			mockSetup: func() {
				controller = newIntegratedMockedAnalyzeController()
				mockContentDownloader := controller.app.ContentDownloader.(*service.MockContentDownloader)
				mockContentDownloader.On("IsLinkAccessible", mock.AnythingOfType("string")).Return(true)
				mockContentDownloader.On("DownloadContent", mock.AnythingOfType("string")).Return(ReadTestContent("go_site_html.html"), nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			//POST
			reqBody := strings.NewReader(`{"site_url": "https://example.com"}`)
			req, err := http.NewRequest("POST", "/analyze", reqBody)
			assert.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.Analyze)
			handler.ServeHTTP(rr, req)
			result := rr.Body.String()

			// sleep for 3 seconds
			time.Sleep(3 * time.Second)

			//GET
			req, err = http.NewRequest("GET", fmt.Sprintf("/analyze?request_id=%s", result), reqBody)
			assert.NoError(t, err)
			rr = httptest.NewRecorder()
			handler = http.HandlerFunc(controller.GetResults)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedCode, rr.Code)
			result = rr.Body.String()

			var analysisResult model.AnalysisResultWithStatus
			err = json.Unmarshal([]byte(result), &analysisResult)
			if err != nil {
				log.Fatalf("Error unmarshalling JSON: %v", err)
			}

			var expectedResult model.AnalysisResultWithStatus
			err = json.Unmarshal([]byte(tt.expectedResult), &expectedResult)
			if err != nil {
				log.Fatalf("Error unmarshalling JSON: %v", err)
			}
			assert.Equal(t, expectedResult, analysisResult)
		})
	}
}
