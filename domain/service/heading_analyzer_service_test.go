package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadingAnalyzeService_AnalyzeContent(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *MockContentReader)
		expectedResult map[string]int
	}{
		{
			name: "all headings present",
			setupMock: func(m *MockContentReader) {
				for i := 1; i <= 6; i++ {
					tag := fmt.Sprintf("h%d", i)
					findResultMock := new(MockContentReader)
					m.On("Find", tag).Return(findResultMock).Once()
					findResultMock.On("Length").Return(i * 2).Once() // Assuming each heading tag appears i*2 times
				}
			},
			expectedResult: map[string]int{
				"h1": 2, "h2": 4, "h3": 6, "h4": 8, "h5": 10, "h6": 12,
			},
		},
		{
			name: "no headings present",
			setupMock: func(m *MockContentReader) {
				for i := 1; i <= 6; i++ {
					tag := fmt.Sprintf("h%d", i)
					findResultMock := new(MockContentReader)
					m.On("Find", tag).Return(findResultMock).Once()
					findResultMock.On("Length").Return(0).Once() // Assuming each heading tag appears i*2 times
				}
			},
			expectedResult: map[string]int{
				"h1": 0, "h2": 0, "h3": 0, "h4": 0, "h5": 0, "h6": 0,
			},
		},
		{
			name: "only h1, and h2 headings present",
			setupMock: func(m *MockContentReader) {
				for i := 1; i <= 6; i++ {
					tag := fmt.Sprintf("h%d", i)
					findResultMock := new(MockContentReader)
					m.On("Find", tag).Return(findResultMock).Once()
					findResultMock.On("Length").Return(func(i int) int {
						if i == 1 || i == 2 {
							return i * 2
						}
						return 0
					}(i)).Once()
				}
			},
			expectedResult: map[string]int{
				"h1": 2, "h2": 4, "h3": 0, "h4": 0, "h5": 0, "h6": 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockContentReader := new(MockContentReader)
			tc.setupMock(mockContentReader)

			service := HeadingAnalyzeService{}
			result, err := service.AnalyzeContent(mockContentReader)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
