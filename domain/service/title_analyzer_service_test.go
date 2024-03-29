package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTitleAnalyzeService_AnalyzeContent tests the AnalyzeContent function of TitleAnalyzeService using a table-driven approach
func TestTitleAnalyzeService_AnalyzeContent(t *testing.T) {
	tests := []struct {
		name        string
		titleText   string // Simulated text of the title tag
		expectTitle string // Expected title text result
	}{
		{
			name:        "empty title",
			titleText:   "",
			expectTitle: "",
		},
		{
			name:        "non-empty title",
			titleText:   "Example Title",
			expectTitle: "Example Title",
		},
		{
			name:        "complex title text",
			titleText:   "Example Title | Site Name",
			expectTitle: "Example Title | Site Name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockContentReader := new(MockContentReader)
			service := TitleAnalyzeService{}

			// Setup mock expectations based on the test case
			firstMock := new(MockContentReader) // Simulate the result of First()
			mockContentReader.On("Find", "title").Return(mockContentReader)
			mockContentReader.On("First").Return(firstMock)
			firstMock.On("Text").Return(tc.titleText)

			// Execute test
			result, err := service.AnalyzeContent(mockContentReader)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectTitle, result)

			// Verify expectations
			mockContentReader.AssertExpectations(t)
			firstMock.AssertExpectations(t)
		})
	}
}
