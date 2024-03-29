package service

import (
	"arjunagl/htmlAnalyzer/domain/model"
	"arjunagl/htmlAnalyzer/domain/port"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ExpectedResults struct {
	result map[string]*model.AnalyzerResult
	err    error
}

func TestAnalyzeContent(t *testing.T) {
	tests := []struct {
		name                 string
		prepareReader        func() port.ContentReader
		prepareDownloader    func() port.ContentDownloader
		prepareItemAnalyzers func() []*ItemAnalyzeService
		expectedResults      ExpectedResults
	}{
		{
			name: "successful analysis",
			prepareReader: func() port.ContentReader {
				cr := &MockContentReader{}
				cr.On("LoadContentFromSource", mock.AnythingOfType("string")).Return(nil)
				return cr
			},
			prepareDownloader: func() port.ContentDownloader {
				dl := &MockContentDownloader{} // success scenario
				dl.On("DownloadContent", mock.AnythingOfType("string")).Return("", nil).Once()
				return dl
			},
			prepareItemAnalyzers: func() []*ItemAnalyzeService {
				m1 := &MockItemAnalyzerService{}
				m1.On("AnalyzeContent", mock.Anything).Return("Hello World", nil).Once() // title
				m1.On("AnalyzeContent", mock.Anything).Return("HTML5", nil).Once()       // version
				m1.On("AnalyzeContent", mock.Anything).Return(false, nil).Once()         // login
				m1.On("AnalyzeContent", mock.Anything).Return(map[string]interface{}{
					"InternalLinks":     6,
					"ExternalLinks":     13,
					"InaccessibleLinks": nil,
				}, nil).Once() // link
				m1.On("AnalyzeContent", mock.Anything).Return(map[string]interface{}{
					"h1": 1,
					"h2": 2,
					"h3": 3,
					"h4": 4,
					"h5": 5,
					"h6": 6,
				}, nil).Once() // heading

				return []*ItemAnalyzeService{
					{Title: "Title", ContentAnalyzer: m1},
					{Title: "Version", ContentAnalyzer: m1},
					{Title: "Login", ContentAnalyzer: m1},
					{Title: "Link", ContentAnalyzer: m1},
					{Title: "Heading", ContentAnalyzer: m1},
				}
			},
			expectedResults: ExpectedResults{
				result: map[string]*model.AnalyzerResult{
					"Title":   {Result: "Hello World", Err: nil},
					"Version": {Result: "HTML5", Err: nil},
					"Login":   {Result: false, Err: nil},
					"Link": {Result: map[string]interface{}{
						"InternalLinks":     6,
						"ExternalLinks":     13,
						"InaccessibleLinks": nil,
					}, Err: nil},
					"Heading": {Result: map[string]interface{}{
						"h1": 1,
						"h2": 2,
						"h3": 3,
						"h4": 4,
						"h5": 5,
						"h6": 6,
					}, Err: nil},
				},
				err: nil,
			},
		},
		{
			name: "failure to download content",
			prepareReader: func() port.ContentReader {
				cr := &MockContentReader{}
				cr.On("LoadContentFromSource", mock.AnythingOfType("string")).Return(nil)
				return cr
			},
			prepareDownloader: func() port.ContentDownloader {
				dl := &MockContentDownloader{} // success scenario
				dl.On("DownloadContent", mock.AnythingOfType("string")).Return("", fmt.Errorf("error downloading content")).Once()
				return dl
			},
			prepareItemAnalyzers: func() []*ItemAnalyzeService {
				m1 := &MockItemAnalyzerService{}
				m1.On("AnalyzeContent", mock.Anything).Return("Hello World", nil).Once() // title

				return []*ItemAnalyzeService{
					{Title: "Title", ContentAnalyzer: m1},
				}
			},
			expectedResults: ExpectedResults{
				result: nil,
				err:    fmt.Errorf("failed to download content from %s : %w", "source", fmt.Errorf("error downloading content")),
			},
		},
		{
			name: "should run successfully without any analyzers",
			prepareReader: func() port.ContentReader {
				cr := &MockContentReader{}
				cr.On("LoadContentFromSource", mock.AnythingOfType("string")).Return(nil)
				return cr
			},
			prepareDownloader: func() port.ContentDownloader {
				dl := &MockContentDownloader{} // success scenario
				dl.On("DownloadContent", mock.AnythingOfType("string")).Return("html content", nil).Once()
				return dl
			},
			prepareItemAnalyzers: func() []*ItemAnalyzeService {
				return []*ItemAnalyzeService{}
			},
			expectedResults: ExpectedResults{
				result: map[string]*model.AnalyzerResult{},
				err:    nil,
			},
		},
		{
			name: "should capture errors correctly to download content",
			prepareReader: func() port.ContentReader {
				cr := &MockContentReader{}
				cr.On("LoadContentFromSource", mock.AnythingOfType("string")).Return(nil)
				return cr
			},
			prepareDownloader: func() port.ContentDownloader {
				dl := &MockContentDownloader{} // success scenario
				dl.On("DownloadContent", mock.AnythingOfType("string")).Return("content", nil).Once()
				return dl
			},
			prepareItemAnalyzers: func() []*ItemAnalyzeService {
				m1 := &MockItemAnalyzerService{}
				m1.On("AnalyzeContent", mock.Anything).Return("Hello World", nil).Once()                         // title
				m1.On("AnalyzeContent", mock.Anything).Return(nil, fmt.Errorf("error analyzing content")).Once() // error

				return []*ItemAnalyzeService{
					{Title: "Title", ContentAnalyzer: m1},
					{Title: "ErrorAnalyzer", ContentAnalyzer: m1},
				}
			},
			expectedResults: ExpectedResults{
				result: map[string]*model.AnalyzerResult{
					"Title":         {Result: "Hello World", Err: nil},
					"ErrorAnalyzer": {Result: nil, Err: fmt.Errorf("error analyzing content")},
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := tc.prepareReader()
			downloader := tc.prepareDownloader()
			service := NewAnalyzerService(tc.prepareItemAnalyzers(), reader, downloader)

			r, err := service.AnalyzeContent("source")

			assert.Equal(t, tc.expectedResults.result, r)
			assert.Equal(t, tc.expectedResults.err, err)
		})
	}
}
