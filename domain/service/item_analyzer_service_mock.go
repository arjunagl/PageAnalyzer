package service

import (
	"arjunagl/htmlAnalyzer/domain/port"

	"github.com/stretchr/testify/mock"
)

type MockItemAnalyzerService struct {
	mock.Mock
}

func (m *MockItemAnalyzerService) AnalyzeContent(content port.ContentReader) (interface{}, error) {
	args := m.Called(content)

	result := args.Get(0)
	var err error
	if args.Get(1) != nil {
		err, _ = args.Get(1).(error)
	}

	return result, err
}
