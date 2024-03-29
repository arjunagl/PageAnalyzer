package service

import (
	"arjunagl/htmlAnalyzer/domain/model"
	"sync"
)

type Scheduler interface {
	Analyze(id string, source string)
	GetResult(id string) (*model.AnalysisResultWithStatus, bool)
}

type SchedulerService struct {
	analyzerService Analyzer
	results         map[string]*model.AnalysisResultWithStatus
	mu              sync.RWMutex // Protects results
}

func NewSchedulerService(analyzerService Analyzer) *SchedulerService {
	return &SchedulerService{
		analyzerService: analyzerService,
		results:         make(map[string]*model.AnalysisResultWithStatus),
	}
}

func (s *SchedulerService) Analyze(id string, source string) {
	s.mu.Lock()
	s.results[id] = &model.AnalysisResultWithStatus{
		Result: nil,
		Status: model.StatusInProgress,
	}
	s.mu.Unlock()
	go s.runAnalysis(id, source)
}

func (s *SchedulerService) runAnalysis(id string, source string) {
	r, err := s.analyzerService.AnalyzeContent(source)
	s.mu.Lock()
	if result, exists := s.results[id]; exists {
		result.Result = r
		if err != nil {
			result.Status = model.StatusError
		} else {
			result.Status = model.StatusComplete
		}
	}
	s.mu.Unlock()
}

func (s *SchedulerService) GetResult(id string) (*model.AnalysisResultWithStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, exists := s.results[id]
	return result, exists
}
