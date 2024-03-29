package service

import (
	"arjunagl/htmlAnalyzer/domain/port"
	"strings"
)

type VersionAnalyzeService struct {
}

func (ta *VersionAnalyzeService) AnalyzeContent(cr port.ContentReader) (interface{}, error) {
	sc := strings.TrimSpace(strings.ToLower(cr.Html()))

	switch {
	case strings.Contains(sc, "<!doctype html>"):
		return "HTML5", nil
	case strings.Contains(sc, "-//w3c//dtd html 4.01"):
		return "HTML 4.01", nil
	case strings.Contains(sc, "-//w3c//dtd xhtml 1.0"):
		return "XHTML 1.0", nil
	case strings.Contains(sc, "-//w3c//dtd xhtml 1.1"):
		return "XHTML 1.1", nil
	default:
		return "Unknown", nil
	}
}
