package port

type ContentReader interface {
	LoadContentFromSource(string) error

	Find(selector string) ContentReader

	First() ContentReader

	Each(func(index int, elem ContentReader))

	Length() int

	Text() string

	Attr(attr string) (string, bool)

	Html() string
}
