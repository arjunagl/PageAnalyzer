package port

type ContentDownloader interface {
	DownloadContent(url string) (string, error)
	IsLinkAccessible(url string) bool
}
