package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginAnalyzeService_AnalyzeContent(t *testing.T) {
	tests := []struct {
		name           string
		passwordInputs int // Number of password input fields to simulate
		expectedResult bool
	}{
		{
			name:           "no password inputs",
			passwordInputs: 0,
			expectedResult: false,
		},
		{
			name:           "one password input",
			passwordInputs: 1,
			expectedResult: true,
		},
		{
			name:           "multiple password inputs",
			passwordInputs: 3,
			expectedResult: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockContentReader := new(MockContentReader)
			service := LoginAnalyzeService{}

			// Setup mock expectations based on the test case
			mockContentReader.On("Find", "input[type='password']").Return(mockContentReader)
			mockContentReader.On("Length").Return(tc.passwordInputs)

			// Execute test
			result, err := service.AnalyzeContent(mockContentReader)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)

			// Verify expectations
			mockContentReader.AssertExpectations(t)
		})
	}
}
