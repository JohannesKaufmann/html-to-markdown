package htmltomarkdown_test

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
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

func ExampleWithDomain() {
	input := `<img src="/assets/image.png" />`

	markdown, err := htmltomarkdown.ConvertString(
		input,
		// Provide a different domain for every convert call:
		converter.WithDomain("https://example.com"),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markdown)
	// Output: ![](https://example.com/assets/image.png)
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

func TestConvertString_MultiLineCodeNodes(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "multiple adjacent code tags in pre block",
			input:    "<pre><code>c1</code><code>c2</code></pre>",
			expected: "```\nc1\nc2\n```",
		},
		{
			desc:     "code snippet with data-lang and nested span elements",
			input:    `<section class="code-snippet__fix code-snippet__js"><ul class="code-snippet__line-index code-snippet__js"><li></li><li></li></ul><pre class="code-snippet__js" data-lang="javascript"><code><span leaf="">git config --<span class="code-snippet__variable">global</span>&nbsp;user.<span class="code-snippet__property">name</span>&nbsp;<span class="code-snippet__string">"Your Gitee Name"</span></span></code><code><span leaf="">git config --<span class="code-snippet__variable">global</span>&nbsp;user.<span class="code-snippet__property">email</span>&nbsp;<span class="code-snippet__string">"Your Gitee Email"</span></span></code></pre></section>`,
			expected: "```javascript\ngit config --global\u00a0user.name\u00a0\"Your Gitee Name\"\ngit config --global\u00a0user.email\u00a0\"Your Gitee Email\"\n```",
		},
		{
			desc:     "multiple code tags with empty line in between",
			input:    "<pre><code>line 1</code><code></code><code>line 3</code></pre>",
			expected: "```\nline 1\n\nline 3\n```",
		},
		{
			desc:     "multiple code tags with existing newline separator",
			input:    "<pre><code>line 1</code>\n<code>line 2</code></pre>",
			expected: "```\nline 1\nline 2\n```",
		},
		{
			desc:     "single code tag in pre block remains unchanged",
			input:    "<pre><code>single line</code></pre>",
			expected: "```\nsingle line\n```",
		},
		{
			desc:     "pre block with data-language attribute",
			input:    `<pre data-language="python"><code>def foo():</code><code>    pass</code></pre>`,
			expected: "```python\ndef foo():\n    pass\n```",
		},
		{
			desc:     "div elements inside pre block",
			input:    "<pre><div>line 1</div><div>line 2</div></pre>",
			expected: "```\nline 1\nline 2\n```",
		},
		{
			desc:     "br element inside code block",
			input:    "<pre><code>line 1<br>line 2</code></pre>",
			expected: "```\nline 1\nline 2\n```",
		},
		{
			desc:     "adjacent inline code outside pre block is still merged",
			input:    "<div>before<code>A</code><code>B</code>after</div>",
			expected: "before`AB`after",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output, err := htmltomarkdown.ConvertString(tC.input)
			if err != nil {
				t.Fatal(err)
			}
			if output != tC.expected {
				t.Errorf("expected %q but got %q", tC.expected, output)
			}
		})
	}
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
				t.Fatal(err)
			}
			if output != tC.expected {
				t.Errorf("expected %q but got %q", tC.expected, output)
			}
		})
	}
}

func TestDataRaceDetector(t *testing.T) {
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
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

			conv.Register.TagType("script", converter.TagTypeBlock, converter.PriorityStandard)

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
