package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLinkAnalyzeService_AnalyzeContent(t *testing.T) {
	tests := []struct {
		name                  string
		hrefValues            []string
		isLinkAccessibleCalls map[string]bool // Maps href values to whether they're accessible
		expectedAnalysis      LinkAnalysis
	}{
		{
			name: "mixed internal and external links, with inaccessible external link",
			hrefValues: []string{
				"http://example.com", // External, inaccessible
				"/internal-link",     // Internal
			},
			isLinkAccessibleCalls: map[string]bool{
				"http://example.com": false,
			},
			expectedAnalysis: LinkAnalysis{
				InternalLinks:     1,
				ExternalLinks:     1,
				InaccessibleLinks: []string{"http://example.com"},
			},
		},
		{
			name: "mixed internal and external links, with accessible external link",
			hrefValues: []string{
				"http://example.com", // External, inaccessible
				"/internal-link",     // Internal
			},
			isLinkAccessibleCalls: map[string]bool{
				"http://example.com": true,
			},
			expectedAnalysis: LinkAnalysis{
				InternalLinks:     1,
				ExternalLinks:     1,
				InaccessibleLinks: nil,
			},
		},
		{
			name: "all internal links",
			hrefValues: []string{
				"/internal-link1",
				"/internal-link2",
			},
			isLinkAccessibleCalls: map[string]bool{}, // No external links to check
			expectedAnalysis: LinkAnalysis{
				InternalLinks:     2,
				ExternalLinks:     0,
				InaccessibleLinks: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockContentReader := new(MockContentReader)
			mockContentDownloader := new(MockContentDownloader)
			service := NewLinkAnalyzeService(mockContentDownloader)

			mockContentReader.On("Find", "a").Return(mockContentReader).Once()
			mockContentReader.On("Each", mock.Anything).Return(tc.hrefValues).Once()

			for href, accessible := range tc.isLinkAccessibleCalls {
				mockContentDownloader.On("IsLinkAccessible", href).Return(accessible)
			}

			result, err := service.AnalyzeContent(mockContentReader)
			assert.NoError(t, err)
			analysisResult, ok := result.(LinkAnalysis)
			assert.True(t, ok)
			assert.Equal(t, tc.expectedAnalysis, analysisResult)

			mockContentReader.AssertExpectations(t)
			mockContentDownloader.AssertExpectations(t)
		})
	}
}
