package table

import (
	"bytes"
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
)

func TestGoldenFiles(t *testing.T) {
	goldenFileConvert := func(htmlInput []byte) ([]byte, error) {
		conv := converter.NewConverter(
			converter.WithPlugins(
				base.NewBasePlugin(),
				commonmark.NewCommonmarkPlugin(),
				NewTablePlugin(),
			),
		)

		return conv.ConvertReader(bytes.NewReader(htmlInput))
	}

	// TODO: Options (e.g. PromoteFirstRowToHeader)
	// TODO: Option: what to do with <br />

	tester.GoldenFiles(t, goldenFileConvert, goldenFileConvert)
}
