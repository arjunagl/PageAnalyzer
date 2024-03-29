package service

import (
	"arjunagl/htmlAnalyzer/domain/port"
)

type TitleAnalyzeService struct {
}

func (ta *TitleAnalyzeService) AnalyzeContent(cr port.ContentReader) (interface{}, error) {
	title := cr.Find("title").First().Text()
	return title, nil
}
