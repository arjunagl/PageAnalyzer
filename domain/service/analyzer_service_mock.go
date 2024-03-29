package service

import (
	"arjunagl/htmlAnalyzer/domain/model"

	"github.com/stretchr/testify/mock"
)

type MockAnalyzerService struct {
	mock.Mock
}

func (m *MockAnalyzerService) AnalyzeContent(source string) (map[string]*model.AnalyzerResult, error) {
	args := m.Called(source)
	return args.Get(0).(map[string]*model.AnalyzerResult), args.Error(1)
}
