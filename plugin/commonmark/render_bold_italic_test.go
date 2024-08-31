package commonmark_test

import (
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
)

func TestNewCommonmarkPlugin_Italic(t *testing.T) {
	const nonBreakingSpace = '\u00A0'
	const zeroWidthSpace = '\u200b'

	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "simple",
			input:    `<p><em>Text</em></p>`,
			expected: `*Text*`,
		},
		{
			desc:     "normal text surrounded by italic",
			input:    `<em>Italic</em>Normal<em>Italic</em>`,
			expected: `*Italic*Normal*Italic*`,
		},
		{
			desc:     "italic text surrounded by normal",
			input:    `Normal<em>Italic</em>Normal`,
			expected: `Normal*Italic*Normal`,
		},
		{
			desc:     "with spaces inside",
			input:    `<p><em>  Text  </em></p>`,
			expected: `*Text*`,
		},
		{
			desc:     "with delimiter inside",
			input:    `<p><em>*A*B*</em></p>`,
			expected: `*\*A\*B\**`,
		},
		{
			desc:     "adjacent",
			input:    `<em>A</em><em>B</em> <em>C</em>`,
			expected: `*AB* *C*`,
		},
		{
			desc:     "adjacent and lots of spaces",
			input:    `<em>  A  </em><em>  B  </em>  <em>  C  </em>`,
			expected: `*A B* *C*`,
		},
		{
			desc:     "nested",
			input:    `<em>A <em>B</em> C</em>`,
			expected: `*A B C*`,
		},
		{
			desc:     "nested and lots of spaces",
			input:    `<em>  A  <em>  B  </em>  C  </em>`,
			expected: `*A B C*`,
		},
		{
			desc:     "mixed nested 1",
			input:    `<em>A <strong>B</strong> C</em>`,
			expected: `*A **B** C*`,
		},
		{
			desc:     "mixed nested 2",
			input:    `<strong>A <em>B</em> C</strong>`,
			expected: `**A *B* C**`,
		},
		{
			desc:     "mixed different italic",
			input:    `<i>A<em>B</em>C</i>`,
			expected: `*ABC*`,
		},

		{
			desc: "next to each other in other containers",
			input: `<div>
	<em>A</em>
	<article><em>B</em></article>
	<em>C</em>
</div>`,
			expected: "*A*\n\n*B*\n\n*C*",
		},

		// - - - - //
		{
			desc:     "empty italic #1",
			input:    `before<i></i>after`,
			expected: `beforeafter`,
		},
		{
			desc:     "empty italic #2",
			input:    `before<i> </i>after`,
			expected: `before after`,
		},
		{
			desc:     "empty italic #3",
			input:    `before <i> </i> after`,
			expected: `before after`,
		},
		{
			desc:     "italic with non-breaking-space",
			input:    `before<i>` + string(nonBreakingSpace) + `</i>after`,
			expected: `before` + string(nonBreakingSpace) + `after`,
		},
		{
			desc:     "italic with zero-width-space",
			input:    `before<i>` + string(zeroWidthSpace) + `</i>after`,
			expected: `before*` + string(zeroWidthSpace) + `*after`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			conv := converter.NewConverter(
				converter.WithPlugins(commonmark.NewCommonmarkPlugin()),
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
