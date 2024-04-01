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
	analyzers []*ItemAnalyzeService
}

func NewAnalyzerService(ias []*ItemAnalyzeService, cr port.ContentReader, cd port.ContentDownloader) *AnalyzerService {
	return &AnalyzerService{analyzers: ias, cr: cr, cd: cd}
}

func (c *AnalyzerService) AnalyzeContent(source string) (map[string]*model.AnalyzerResult, error) {
	// download content from source
	content, err := c.cd.DownloadContent(source)
	if err != nil {
		return nil, fmt.Errorf("failed to download content from %s : %w", source, err)
	}

	// load the content into the reader
	var rdr port.ContentReader
	rdr, err = c.cr.LoadContentFromSource(source, content)
	if err != nil {
		return nil, fmt.Errorf("failed to read content from %s : %w", source, err)
	}

	// Start analysis of all registered analyzers on the loaded content
	results := make(map[*ItemAnalyzeService]*model.AnalyzerResult)
	for _, as := range c.analyzers {
		r, err := as.ContentAnalyzer.AnalyzeContent(rdr)
		results[as] = &model.AnalyzerResult{
			Result: r,
			Err:    err,
		}
	}
	return c.getAnalysisResults(results), nil
}

func (c *AnalyzerService) getAnalysisResults(results map[*ItemAnalyzeService]*model.AnalyzerResult) map[string]*model.AnalyzerResult {
	aResults := make(map[string]*model.AnalyzerResult, len(results))
	for service, result := range results {
		aResults[service.Title] = result
	}
	return aResults
}
