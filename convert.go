package htmltomarkdown

import (
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"golang.org/x/net/html"
)

func ConvertString(htmlInput string) (string, error) {
	conv := converter.NewConverter(
		converter.WithPlugins(commonmark.NewCommonmarkPlugin()),
	)

	return conv.ConvertString(htmlInput)
}

func ConvertNode(doc *html.Node) ([]byte, error) {
	conv := converter.NewConverter(
		converter.WithPlugins(commonmark.NewCommonmarkPlugin()),
	)

	return conv.ConvertNode(doc)
}
