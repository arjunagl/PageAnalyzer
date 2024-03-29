package service

import (
	"arjunagl/htmlAnalyzer/domain/port"
	"strings"
)

type LinkAnalyzeService struct {
	cd port.ContentDownloader
}

type LinkAnalysis struct {
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks []string
}

func NewLinkAnalyzeService(cd port.ContentDownloader) *LinkAnalyzeService {
	return &LinkAnalyzeService{cd: cd}
}

func (ta *LinkAnalyzeService) AnalyzeContent(cr port.ContentReader) (interface{}, error) {

	analysis := LinkAnalysis{}
	cr.Find("a").Each(func(i int, s port.ContentReader) {
		href, exists := s.Attr("href")
		if !exists || href == "" || strings.HasPrefix(href, "#") {
			return // Ignore anchors or empty hrefs
		}

		if strings.HasPrefix(href, "http") {
			analysis.ExternalLinks++
			// Check accessibility for external links
			if !ta.cd.IsLinkAccessible(href) {
				analysis.InaccessibleLinks = append(analysis.InaccessibleLinks, href)
			}
		} else {
			analysis.InternalLinks++
		}
	})

	return analysis, nil

}
