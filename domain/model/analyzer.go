package model

type AnalyzerResult struct {
	Result interface{} `json:"result"`
	Err    error       `json:"error,omitempty"`
}

type AnalysisStatus string

const (
	StatusInProgress AnalysisStatus = "In Progress"
	StatusComplete   AnalysisStatus = "Complete"
	StatusError      AnalysisStatus = "Error"
)

type AnalysisResultWithStatus struct {
	Result map[string]*AnalyzerResult `json:"result"`
	Status AnalysisStatus             `json:"status"`
}
