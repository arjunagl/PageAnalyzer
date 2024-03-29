package service

import (
	"github.com/stretchr/testify/mock"
)

type MockContentDownloader struct {
	mock.Mock
}

func (m *MockContentDownloader) IsLinkAccessible(url string) bool {
	args := m.Called(url)
	return args.Bool(0)
}

func (m *MockContentDownloader) DownloadContent(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}
