package service

import (
	"arjunagl/htmlAnalyzer/domain/model"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type expectedReturn struct {
	Result map[string]*model.AnalyzerResult
	Err    error
}

func TestSchedulerService(t *testing.T) {
	mockAnalyzerService := new(MockAnalyzerService)
	scheduler := NewSchedulerService(mockAnalyzerService)

	testCases := []struct {
		name           string
		id             string
		source         string
		mockReturn     expectedReturn
		expectedStatus model.AnalysisStatus
	}{
		{
			name:   "Valid Analysis",
			id:     "test1",
			source: "https://example.com",
			mockReturn: expectedReturn{
				Result: map[string]*model.AnalyzerResult{"Title": {Result: "Hello World", Err: nil}},
				Err:    nil,
			},
			expectedStatus: model.StatusComplete,
		},
		{
			name:   "Error Analysis",
			id:     "test1",
			source: "https://example.com",
			mockReturn: expectedReturn{
				Result: map[string]*model.AnalyzerResult{"Title": {Result: "Hello World", Err: nil}},
				Err:    fmt.Errorf("error analyzing"),
			},
			expectedStatus: model.StatusError,
		},
		{
			name:   "Not Existant Analysis",
			id:     "test1",
			source: "https://example.com",
			mockReturn: expectedReturn{
				Result: map[string]*model.AnalyzerResult{"Title": {Result: "Hello World", Err: nil}},
				Err:    nil,
			},
			expectedStatus: model.StatusComplete,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Setup mock expectation
			mockAnalyzerService.On("AnalyzeContent", tc.source).Return(tc.mockReturn.Result, tc.mockReturn.Err).Once()

			// Test the results
			_, exists := scheduler.GetResult(tc.id)
			assert.Equal(t, false, exists)

		})
	}
}
