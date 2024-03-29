package service

import (
	"arjunagl/htmlAnalyzer/domain/model"

	"github.com/stretchr/testify/mock"
)

type MockSchedulerService struct {
	mock.Mock
}

func (m *MockSchedulerService) Analyze(id, siteURL string) {
	m.Called(id, siteURL)
}

func (m *MockSchedulerService) GetResult(requestID string) (*model.AnalysisResultWithStatus, bool) {
	args := m.Called(requestID)
	var result *model.AnalysisResultWithStatus
	if args.Get(0) != nil {
		result = args.Get(0).(*model.AnalysisResultWithStatus)
	}
	return result, args.Bool(1)
}
