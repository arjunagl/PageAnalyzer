package adapter

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type WebContentDownloader struct {
}

func (wd *WebContentDownloader) DownloadContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error downloading content from %s %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading html content %w", err)
	}

	return string(body), nil
}

func (wd *WebContentDownloader) IsLinkAccessible(url string) bool {
	client := http.Client{
		Timeout: 5 * time.Second, // Set a timeout to prevent hanging on a request
	}
	resp, err := client.Head(url)
	if err != nil {
		return false // Assume inaccessible if any error occurs
	}
	defer resp.Body.Close()

	// Consider HTTP status codes outside of 200-399 range as inaccessible
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}
