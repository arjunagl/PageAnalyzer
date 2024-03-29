package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

			firstMock := new(MockContentReader)
			mockContentReader.On("Find", "title").Return(mockContentReader)
			mockContentReader.On("First").Return(firstMock)
			firstMock.On("Text").Return(tc.titleText)

			result, err := service.AnalyzeContent(mockContentReader)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectTitle, result)

			mockContentReader.AssertExpectations(t)
			firstMock.AssertExpectations(t)
		})
	}
}
