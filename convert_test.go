package htmltomarkdown_test

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"golang.org/x/net/html"
)

func ExampleConvertString() {
	input := `<strong>Bold Text</strong>`

	markdown, err := htmltomarkdown.ConvertString(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markdown)
	// Output: **Bold Text**
}
func ExampleConvertNode() {
	input := `<strong>Bold Text</strong>`

	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		log.Fatal(err)
	}

	markdown, err := htmltomarkdown.ConvertNode(doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(markdown))
	// Output: **Bold Text**
}

func TestConvertString_WindowsCarriageReturn(t *testing.T) {
	testCases := []struct {
		desc string

		input    string
		expected string
	}{
		{
			desc: "just newlines",

			input:    "\r\n\r\n\r\n\r\n",
			expected: "",
		},
		{
			desc: "inside strong",

			input:    "<strong>Bold\r\n\r\n\r\n\r\nText</strong>",
			expected: "**Bold Text**",
		},
		{
			desc: "inside paragraph",

			input:    "<p>Some\r\n\r\n\r\n\r\nText</p>",
			expected: "Some Text",
		},
		{
			desc: "inside list",

			input:    "<ul><li>Some\r\n\r\n\r\n\r\nText</li></ul>",
			expected: "- Some Text",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output, err := htmltomarkdown.ConvertString(tC.input)
			if err != nil {
				log.Fatal(err)
			}
			if output != tC.expected {
				t.Errorf("expected %q but got %q", tC.expected, output)
			}
		})
	}
}

func TestDataRaceDetector(t *testing.T) {
	conv := converter.NewConverter(
		converter.WithPlugins(commonmark.NewCommonmarkPlugin()),
	)

	input := `<i>italic text</i>`

	var wg sync.WaitGroup

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			conv.Register.EscapedChar('~')
			conv.Register.UnEscaper(
				func(chars []byte, index int) int { return -1 },
				converter.PriorityStandard,
			)
			conv.Register.PreRenderer(
				func(ctx converter.Context, doc *html.Node) {},
				converter.PriorityStandard,
			)
			conv.Register.TextTransformer(
				func(ctx converter.Context, content string) string { return content },
				converter.PriorityStandard,
			)
			conv.Register.Renderer(
				func(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
					return converter.RenderTryNext
				},
				converter.PriorityStandard,
			)
			conv.Register.PostRenderer(
				func(ctx converter.Context, content []byte) []byte {
					return content
				},
				converter.PriorityStandard,
			)

			conv.Register.TagStrategy("script", converter.StrategyHTMLBlock)

			output, err := conv.ConvertString(input, converter.WithDomain("example.com"))
			if err != nil {
				t.Error(err)
			}
			_ = output

			output2, err := conv.ConvertString(input)
			if err != nil {
				t.Error(err)
			}
			_ = output2

			wg.Done()
		}()
	}
	wg.Wait()
}
