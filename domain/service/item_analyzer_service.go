package service

import "arjunagl/htmlAnalyzer/domain/port"

type ContentAnalyzer interface {
	AnalyzeContent(content port.ContentReader) (interface{}, error)
}

type ItemAnalyzeService struct {
	Title           string
	ContentAnalyzer ContentAnalyzer
}
