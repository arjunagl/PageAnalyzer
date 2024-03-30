package adapter

import (
	"arjunagl/htmlAnalyzer/domain/port"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GoqueryAdapter struct {
	sourceContent string
	sourceURL     string
	selection     *goquery.Selection
}

func NewGoqueryAdapter() port.ContentReader {
	return &GoqueryAdapter{}
}

func (g *GoqueryAdapter) LoadContentFromSource(source string) error {
	g.sourceURL = source
	g.sourceContent = source
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(source))
	if err != nil {
		return fmt.Errorf("error converting to document: %w", err)
	}
	g.selection = doc.Selection
	return nil
}

func (g *GoqueryAdapter) Find(selector string) port.ContentReader {
	return &GoqueryAdapter{selection: g.selection.Find(selector)}
}

func (g *GoqueryAdapter) First() port.ContentReader {
	return &GoqueryAdapter{selection: g.selection.First()}
}

func (g *GoqueryAdapter) Each(f func(index int, elem port.ContentReader)) {
	g.selection.Each(func(i int, s *goquery.Selection) {
		f(i, &GoqueryAdapter{selection: s})
	})
}

func (g *GoqueryAdapter) Length() int {
	return g.selection.Length()
}

func (g *GoqueryAdapter) Text() string {
	return g.selection.Text()
}

func (g *GoqueryAdapter) Attr(attr string) (string, bool) {
	return g.selection.Attr(attr)
}

func (g *GoqueryAdapter) Html() string {
	return g.sourceContent
}

func (g *GoqueryAdapter) SourceURL() string {
	return g.sourceURL
}
