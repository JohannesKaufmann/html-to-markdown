package htmltomarkdown

import (
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"golang.org/x/net/html"
)

// ConvertString converts a html-string to a markdown-string.
//
// Under the hood `html.Parse()` is used to parse the HTML.
func ConvertString(htmlInput string) (string, error) {
	conv := converter.NewConverter(
		converter.WithPlugins(commonmark.NewCommonmarkPlugin()),
	)

	return conv.ConvertString(htmlInput)
}

// ConvertNode converts a `*html.Node` to a markdown byte slice.
//
// If you have already parsed an HTML page using the `html.Parse()` function
// from the "golang.org/x/net/html" package then you can pass this node
// directly to the converter.
func ConvertNode(doc *html.Node) ([]byte, error) {
	conv := converter.NewConverter(
		converter.WithPlugins(commonmark.NewCommonmarkPlugin()),
	)

	return conv.ConvertNode(doc)
}
