package service

import (
	"arjunagl/htmlAnalyzer/domain/port"
)

type LoginAnalyzeService struct {
}

func (ta *LoginAnalyzeService) AnalyzeContent(cr port.ContentReader) (interface{}, error) {
	contains := cr.Find("input[type='password']").Length() > 0
	return contains, nil
}
