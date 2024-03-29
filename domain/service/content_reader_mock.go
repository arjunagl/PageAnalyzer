package service

import (
	"arjunagl/htmlAnalyzer/domain/port"

	"github.com/stretchr/testify/mock"
)

// MockContentReader mocks the ContentReader interface
type MockContentReader struct {
	mock.Mock
}

func (m *MockContentReader) LoadContentFromSource(source string) error {
	args := m.Called(source)
	return args.Error(0)
}

func (m *MockContentReader) Find(selector string) port.ContentReader {
	args := m.Called(selector)
	return args.Get(0).(port.ContentReader)
}

func (m *MockContentReader) First() port.ContentReader {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(port.ContentReader)
}

func (m *MockContentReader) Each(f func(index int, elem port.ContentReader)) {
	args := m.Called(f)
	hrefValues := args.Get(0).([]string) // Expect a slice of href values to be passed

	for i, href := range hrefValues {
		// Create a new mock for each href and setup Attr
		elemMock := new(MockContentReader)
		elemMock.On("Attr", "href").Return(href, true).Once()

		// Call the callback function with the mock
		f(i, elemMock)
	}
}

func (m *MockContentReader) Length() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockContentReader) Text() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockContentReader) Attr(attr string) (string, bool) {
	args := m.Called(attr)
	return args.String(0), args.Bool(1)
}

func (m *MockContentReader) Html() string {
	args := m.Called()
	return args.String(0)
}
