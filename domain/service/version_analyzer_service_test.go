package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestVersionAnalyzeService_AnalyzeContent corrected for nil interface issue
func TestVersionAnalyzeService_AnalyzeContent(t *testing.T) {
	tests := []struct {
		name           string
		doctypePresent bool // Simulates whether the <!doctype html> is present
		sourceContent  string
		expectedResult string
	}{
		{
			name:           "HTML5",
			sourceContent:  "<!doctype html><html></html>",
			expectedResult: "HTML5",
		},
		{
			name:           "doctype not present",
			sourceContent:  "<html></html>",
			expectedResult: "Unknown",
		},
		{
			name:           "HTML1",
			sourceContent:  "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\">",
			expectedResult: "HTML 4.01",
		},
		{
			name:           "XHTML1",
			sourceContent:  `<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">`,
			expectedResult: "XHTML 1.0",
		},
		{
			name:           "XHTML1.1",
			sourceContent:  `<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">`,
			expectedResult: "XHTML 1.1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockContentReader := new(MockContentReader)
			mockContentReader.On("Html").Return(tc.sourceContent).Once()
			service := VersionAnalyzeService{}

			// Execute test
			result, err := service.AnalyzeContent(mockContentReader)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)

			// Verify expectations
			mockContentReader.AssertExpectations(t)
		})
	}
}
