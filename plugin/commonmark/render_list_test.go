package commonmark_test

import (
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
)

func TestNewCommonmarkPlugin_List(t *testing.T) {
	const nonBreakingSpace = '\u00A0'
	const zeroWidthSpace = '\u200b'

	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "nested code block in list item",
			input:    "<ul><li>list item:<pre>line 1\nline 2</pre></li></ul>",
			expected: "- list item:\n  \n  ```\n  line 1\n  line 2\n  ```",
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			conv := converter.NewConverter(
				converter.WithPlugins(
					base.NewBasePlugin(),
					commonmark.NewCommonmarkPlugin(),
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
