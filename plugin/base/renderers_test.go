package base_test

import (
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"golang.org/x/net/html"
)

func TestRenderAsX(t *testing.T) {
	input := `
<h1>heading</h1>
<footer>
	<strong>bold text</strong>
</footer>
	`

	testCases := []struct {
		desc string

		isInline   bool
		renderFunc converter.HandleRenderFunc

		expected string
	}{
		{
			desc: "Default Block",
			renderFunc: func(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
				return converter.RenderTryNext
			},
			expected: "# heading\n\n**bold text**",
		},
		{
			desc:     "Default Inline",
			isInline: true,
			renderFunc: func(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
				return converter.RenderTryNext
			},
			expected: "# heading\n\n**bold text**",
		},

		{
			desc:       "RenderAsHTML Block",
			renderFunc: base.RenderAsHTML,
			expected:   "# heading\n\n<footer><strong>bold text</strong></footer>",
		},
		{
			desc:       "RenderAsHTML Inline",
			isInline:   true,
			renderFunc: base.RenderAsHTML,
			expected:   "# heading\n\n<footer><strong>bold text</strong></footer>",
		},

		{
			desc:       "RenderAsHTMLWrapper Block",
			renderFunc: base.RenderAsHTMLWrapper,
			expected:   "# heading\n\n<footer>\n\n**bold text**\n\n</footer>",
		},
		{
			desc:       "RenderAsHTMLWrapper Inline",
			isInline:   true,
			renderFunc: base.RenderAsHTMLWrapper,
			expected:   "# heading\n\n<footer>\n\n**bold text**\n\n</footer>",
		},

		{
			desc:       "RenderAsPlaintextWrapper Block",
			renderFunc: base.RenderAsPlaintextWrapper,
			expected:   "# heading\n\n**bold text**",
		},
		{
			desc:       "RenderAsPlaintextWrapper Inline",
			isInline:   true,
			renderFunc: base.RenderAsPlaintextWrapper,
			expected:   "# heading\n\n**bold text**",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			conv := converter.NewConverter(
				converter.WithPlugins(
					base.NewBasePlugin(),
					commonmark.NewCommonmarkPlugin(),
				),
			)
			if tC.isInline {
				conv.Register.RendererFor("footer", converter.TagTypeInline, tC.renderFunc, converter.PriorityStandard)
			} else {
				conv.Register.RendererFor("footer", converter.TagTypeBlock, tC.renderFunc, converter.PriorityStandard)
			}

			output, err := conv.ConvertString(input)
			if err != nil {
				t.Fatal(err)
			}
			if output != tC.expected {
				t.Errorf("expected %q but got %q", tC.expected, output)
			}
		})
	}
}
