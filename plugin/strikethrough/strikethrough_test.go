package strikethrough_test

import (
	"bytes"
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/strikethrough"
)

func TestNewStrikethroughPlugin(t *testing.T) {
	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "simple",
			input:    `<p><s>Text</s></p>`,
			expected: `~~Text~~`,
		},
		{
			desc:     "with spaces inside",
			input:    `<p><s>  Text  </s></p>`,
			expected: `~~Text~~`,
		},
		{
			desc:     "with spaces inside",
			input:    `<p><s>~~A~~B~~</s></p>`,
			expected: `~~\~\~A\~\~B\~\~~~`,
		},
		{
			desc:     "nested",
			input:    `<p><s>A <s>B</s> C</s></p>`,
			expected: `~~A B C~~`,
		},
		{
			desc:     "adjacent",
			input:    `<p><s>A</s><s>B</s> <s>C</s></p>`,
			expected: `~~AB~~ ~~C~~`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			conv := converter.NewConverter(
				converter.WithPlugins(
					base.NewBasePlugin(),
					strikethrough.NewStrikethroughPlugin(),
				),
			)

			out, err := conv.ConvertString(run.input)
			if err != nil {
				t.Error(err)
			}
			if out != run.expected {
				t.Errorf("expected %q but got %q", run.expected, out)
			}
		})
	}
}
func TestWithDelimiter(t *testing.T) {
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			strikethrough.NewStrikethroughPlugin(
				strikethrough.WithDelimiter("=="),
			),
		),
	)

	input := `<p><s>Text</s></p>`
	expected := `==Text==`

	out, err := conv.ConvertString(input)
	if err != nil {
		t.Error(err)
	}
	if out != expected {
		t.Errorf("expected %q but got %q", expected, out)
	}
}

func TestGoldenFiles(t *testing.T) {
	goldenFileConvert := func(htmlInput []byte) ([]byte, error) {
		conv := converter.NewConverter(
			converter.WithPlugins(
				base.NewBasePlugin(),
				commonmark.NewCommonmarkPlugin(),
				strikethrough.NewStrikethroughPlugin(),
			),
		)

		return conv.ConvertReader(bytes.NewReader(htmlInput))
	}

	tester.GoldenFiles(t, goldenFileConvert, goldenFileConvert)
}
