package service

import (
	"arjunagl/htmlAnalyzer/domain/port"
	"fmt"
)

type HeadingAnalyzeService struct {
}

func (ta *HeadingAnalyzeService) AnalyzeContent(cr port.ContentReader) (interface{}, error) {
	headings := make(map[string]int)
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		count := cr.Find(tag).Length()
		headings[tag] = count
	}
	return headings, nil
}
