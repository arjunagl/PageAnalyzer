package service

import (
	"arjunagl/htmlAnalyzer/domain/model"
	"arjunagl/htmlAnalyzer/domain/port"
	"fmt"
)

type Analyzer interface {
	AnalyzeContent(source string) (map[string]*model.AnalyzerResult, error)
}

type AnalyzerService struct {
	cr        port.ContentReader
	cd        port.ContentDownloader
	analyzers map[*ItemAnalyzeService]*model.AnalyzerResult
}

func NewAnalyzerService(ias []*ItemAnalyzeService, cr port.ContentReader, cd port.ContentDownloader) *AnalyzerService {
	as := make(map[*ItemAnalyzeService]*model.AnalyzerResult)
	for _, a := range ias {
		as[a] = nil
	}
	return &AnalyzerService{analyzers: as, cr: cr, cd: cd}
}

func (c *AnalyzerService) AnalyzeContent(source string) (map[string]*model.AnalyzerResult, error) {
	// download content from source
	content, err := c.cd.DownloadContent(source)
	if err != nil {
		return nil, fmt.Errorf("failed to download content from %s : %w", source, err)
	}

	// load the content into the reader
	if err := c.cr.LoadContentFromSource(content); err != nil {
		return nil, fmt.Errorf("failed to read content from %s : %w", source, err)
	}

	// Start analysis of all registered analyzers on the loaded content
	for a := range c.analyzers {
		r, err := a.ContentAnalyzer.AnalyzeContent(c.cr)
		c.analyzers[a] = &model.AnalyzerResult{
			Result: r,
			Err:    err,
		}
	}
	return c.getAnalysisResults(), nil
}

func (c *AnalyzerService) getAnalysisResults() map[string]*model.AnalyzerResult {
	results := make(map[string]*model.AnalyzerResult, len(c.analyzers))
	for service, result := range c.analyzers {
		results[service.Title] = result
	}
	return results
}
