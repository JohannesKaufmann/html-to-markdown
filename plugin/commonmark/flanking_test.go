package commonmark_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/yuin/goldmark"
	"golang.org/x/net/html"
)

// renderedText parses the markdown and returns the text that a reader
// would end up seeing.
func renderedText(t *testing.T, markdown string) string {
	t.Helper()

	var buf bytes.Buffer
	if err := goldmark.New().Convert([]byte(markdown), &buf); err != nil {
		t.Fatal(err)
	}
	doc, err := html.Parse(&buf)
	if err != nil {
		t.Fatal(err)
	}

	var text strings.Builder
	for _, node := range dom.AllNodes(doc) {
		if node.Type == html.TextNode {
			text.WriteString(node.Data)
		}
	}
	return strings.Join(strings.Fields(text.String()), " ")
}

var flankingContents = []string{
	"x", "xy", ".x", "x.", ".x.", "!x", "x!", "?x", "x?", "-x", "x-",
	"(x)", "[x]", `"x"`, "'x'", "...x", "x...", ",x", "x,", ":x", "x:",
	"/x", "x/", "&x", "x&", "+x", "x+", "=x", "x=", "~x", "x~", "|x", "x|",
	".", "!", "-", "...", "€x", "x€", "日x", " x", "x ", " x ",
}

var flankingContexts = []struct{ before, after string }{
	{"", ""},
	{"before", "after"},
	{"before", ""},
	{"", "after"},
	{"before ", " after"},
	{"before ", "after"},
	{"before", " after"},
	{"before-", "-after"},
	{"before.", ".after"},
	{"1", "2"},
}

// A delimiter that does not flank its content is not emphasis - it is a
// literal asterisk that the reader gets to see. Whatever we emit, the text
// itself has to survive unchanged.
func TestNewCommonmarkPlugin_DelimitersNeverBecomeLiteralText(t *testing.T) {
	for _, tag := range []string{"em", "i", "strong", "b"} {
		for _, content := range flankingContents {
			for _, ctx := range flankingContexts {
				input := fmt.Sprintf("<p>%s<%s>%s</%s>%s</p>", ctx.before, tag, content, tag, ctx.after)

				t.Run(input, func(t *testing.T) {
					conv := converter.NewConverter(
						converter.WithPlugins(
							base.NewBasePlugin(),
							commonmark.NewCommonmarkPlugin(),
						),
					)

					out, err := conv.ConvertString(input)
					if err != nil {
						t.Fatal(err)
					}

					expected := strings.Join(strings.Fields(ctx.before+content+ctx.after), " ")
					if got := renderedText(t, out); got != expected {
						t.Errorf("markdown %q renders as %q but the text should stay %q", out, got, expected)
					}
				})
			}
		}
	}
}

func TestNewCommonmarkPlugin_PunctuationAtEdges(t *testing.T) {
	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "punctuation at the start",
			input:    `<p>before<em>.middle</em>after</p>`,
			expected: `before.*middle*after`,
		},
		{
			desc:     "punctuation at the end",
			input:    `<p>before<em>middle.</em>after</p>`,
			expected: `before*middle*.after`,
		},
		{
			desc:     "punctuation at both edges",
			input:    `<p>before<strong>.middle.</strong>after</p>`,
			expected: `before.**middle**.after`,
		},
		{
			desc:     "moving one character is enough",
			input:    `<p>before<em>?!!middle?!!</em>after</p>`,
			expected: `before?*!!middle?!*!after`,
		},
		{
			desc:     "a space already does the job",
			input:    `<p>before <em>.middle.</em> after</p>`,
			expected: `before *.middle.* after`,
		},
		{
			desc:     "surrounding punctuation already does the job",
			input:    `<p>before-<em>.middle.</em>-after</p>`,
			expected: `before-*.middle.*-after`,
		},
		{
			desc:     "the start of the line counts as whitespace",
			input:    `<p><em>.middle.</em></p>`,
			expected: `*.middle.*`,
		},
		{
			desc:     "a currency sign counts as punctuation",
			input:    `<p>before<em>€x</em>after</p>`,
			expected: `before€*x*after`,
		},
		{
			desc:     "a cjk ideograph does not",
			input:    `<p>before<em>日x</em>after</p>`,
			expected: `before*日x*after`,
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
