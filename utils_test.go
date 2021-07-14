package md

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func getNodeFromString(t *testing.T, rawHTML string) *html.Node {
	docNode, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		t.Error(err)
		return nil
	}

	//  -> #document -> body -> actuall content
	return docNode.FirstChild.LastChild.FirstChild
}

func TestAddSpaceIfNessesary(t *testing.T) {
	var tests = []struct {
		Name string

		Prev     string
		Next     string
		Markdown string

		Expect string
	}{
		{
			Name: "dont count comment",
			Prev: `<html><head></head>
				<!--some comment-->
			<body></body></html>`,
			Next: `<html><head></head>
			<!--another comment-->
		<body></body></html>`,
			Markdown: `_Comment Content_`,
			Expect:   `_Comment Content_`,
		},
		{

			Name:     "bold with break",
			Prev:     `<br />`,
			Next:     `<br />`,
			Markdown: `**Bold**`,
			Expect:   `**Bold**`,
		},
		{
			Name:     "italic with no space",
			Prev:     ``,
			Next:     `and no space afterward.`, // #text
			Markdown: `_Content_`,
			Expect:   `_Content_ `,
		},
		{
			Name:     "bold with no space",
			Prev:     `Some`,
			Next:     `Text`,
			Markdown: `**Bold**`,
			Expect:   ` **Bold** `,
		},
		{
			Name:     "bold with no space in span",
			Prev:     `<span>Some</span>`,
			Next:     `<span>Text</span>`,
			Markdown: `**Bold**`,
			Expect:   ` **Bold** `,
		},
		{
			Name:     "italic with no space",
			Prev:     ``,
			Next:     `and no space afterward.`,
			Markdown: `_Content_`,
			Expect:   `_Content_ `,
		},
		{
			Name:     "github example without new lines",
			Prev:     `<a>go</a>`,
			Next:     `<a>html</a>`,
			Markdown: `[golang](http://example.com/topics/golang "Topic: golang")`,
			Expect:   ` [golang](http://example.com/topics/golang "Topic: golang")`,
		},
		{
			Name: "github example",
			Prev: `<a class="topic-tag topic-tag-link " data-ga-click="Topic, repository page" data-octo-click="topic_click" data-octo-dimensions="topic:go" href="/topics/go" title="Topic: go">
			go
			</a>`,
			Next: `<a class="topic-tag topic-tag-link " data-ga-click="Topic, repository page" data-octo-click="topic_click" data-octo-dimensions="topic:html" href="/topics/html" title="Topic: html">
			html
			</a>`,
			Markdown: `[golang](http://example.com/topics/golang "Topic: golang")`,
			Expect:   ` [golang](http://example.com/topics/golang "Topic: golang")`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			// build a selection for goquery with siblings
			selec := &goquery.Selection{
				Nodes: []*html.Node{
					{
						Data:        "a",
						PrevSibling: getNodeFromString(t, test.Prev),
						NextSibling: getNodeFromString(t, test.Next),
					},
				},
			}
			output := AddSpaceIfNessesary(selec, test.Markdown)

			if output != test.Expect {
				t.Errorf("expected '%s' but got '%s'", test.Expect, output)
			}
		})
	}
}

func TestTrimpLeadingSpaces(t *testing.T) {
	var tests = []struct {
		Name   string
		Text   string
		Expect string
	}{
		{
			Name: "trim normal text",
			Text: `
This is a normal paragraph
 this as well
  just with some spaces before
			`,
			Expect: `
This is a normal paragraph
this as well
just with some spaces before
			`,
		},
		{
			Name: "dont trim nested lists",
			Text: `
- Home
- About
	- People
	- History
		- 2019
		- 2020		
			`,
			Expect: `
- Home
- About
	- People
	- History
		- 2019
		- 2020		
			`,
		},
		{
			Name: "dont trim list with multiple paragraphs",
			Text: `
1.  This is a list item with two paragraphs. Lorem ipsum dolor
	sit amet, consectetuer adipiscing elit. Aliquam hendrerit
	mi posuere lectus.

	Vestibulum enim wisi, viverra nec, fringilla in, laoreet
	vitae, risus. Donec sit amet nisl. Aliquam semper ipsum
	sit amet velit.

2.  Suspendisse id sem consectetuer libero luctus adipiscing.
			`,
			Expect: `
1.  This is a list item with two paragraphs. Lorem ipsum dolor
	sit amet, consectetuer adipiscing elit. Aliquam hendrerit
	mi posuere lectus.

	Vestibulum enim wisi, viverra nec, fringilla in, laoreet
	vitae, risus. Donec sit amet nisl. Aliquam semper ipsum
	sit amet velit.

2.  Suspendisse id sem consectetuer libero luctus adipiscing.
			`,
		},
		{
			Name: "dont trim code blocks",
			Text: `
This is a normal paragraph:

    This is a code block.
			`,
			Expect: `
This is a normal paragraph:

    This is a code block.
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			output := TrimpLeadingSpaces(test.Text)

			if output != test.Expect {
				t.Errorf("expected '%s' but got '%s'", test.Expect, output)
			}
		})
	}

}

func TestTrimTrailingSpaces(t *testing.T) {
	var tests = []struct {
		Name   string
		Text   string
		Expect string
	}{
		{
			Name: "trim after normal text",
			Text: `
1\. xxx 

2\. xxxx	
			`,
			Expect: `
1\. xxx

2\. xxxx
`,
		},
		{
			Name:   "dont trim inside normal text",
			Text:   "When `x = 3`, that means `x + 2 = 5`",
			Expect: "When `x = 3`, that means `x + 2 = 5`",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			output := TrimTrailingSpaces(test.Text)

			if output != test.Expect {
				t.Errorf("expected '%s' but got '%s'", test.Expect, output)
			}
		})
	}
}

func TestEscapeMultiLine(t *testing.T) {
	var tests = []struct {
		Name   string
		Text   string
		Expect string
	}{
		{
			Name: "escape new lines",
			Text: `line1
line2

line3









line4`,
			Expect: `line1\
line2\
\
line3\
\
line4`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			output := EscapeMultiLine(test.Text)

			if output != test.Expect {
				t.Errorf("expected '%s' but got '%s'", test.Expect, output)
			}
		})
	}
}

func TestCalculateCodeFence(t *testing.T) {
	var tests = []struct {
		Name      string
		FenceChar rune

		Text   string
		Expect string
	}{
		{
			Name:      "no occurrences with backtick",
			FenceChar: '`',
			Text:      `normal ~~~ code block`,
			Expect:    "```",
		},
		{
			Name:      "no occurrences with tilde",
			FenceChar: '~',
			Text:      "normal ``` code block",
			Expect:    "~~~",
		},
		{
			Name:      "one exact occurrence",
			FenceChar: '`',
			Text:      "```",
			Expect:    "````",
		},
		{
			Name:      "one occurrences with backtick",
			FenceChar: '`',
			Text:      "normal ``` code block",
			Expect:    "````",
		},
		{
			Name:      "one bigger occurrences with backtick",
			FenceChar: '`',
			Text:      "normal ````` code block",
			Expect:    "``````",
		},
		{
			Name:      "multiple occurrences with backtick",
			FenceChar: '`',
			Text:      "normal ``` code `````` block",
			Expect:    "```````",
		},
		{
			Name:      "multiple occurrences with tilde",
			FenceChar: '~',
			Text:      "normal ~~~ code ~~~~~~~~~~~~ block",
			Expect:    "~~~~~~~~~~~~~",
		},
		{
			Name:      "multiple occurrences on different lines with tilde",
			FenceChar: '~',
			Text: `
normal
	~~~
code ~~~~~~~~~~~~ block
				`,
			Expect: "~~~~~~~~~~~~~",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			output := CalculateCodeFence(test.FenceChar, test.Text)

			if output != test.Expect {
				t.Errorf("expected '%s' (x%d) but got '%s' (x%d)", test.Expect, strings.Count(test.Expect, string(test.FenceChar)), output, strings.Count(output, string(test.FenceChar)))
			}
		})
	}
}

func TestIsListItem(t *testing.T) {
	var tests = []struct {
		name string

		opt  *Options
		line string

		expected bool
	}{
		{
			name: "nothing",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "",
			expected: false,
		},
		{
			name: "just spaces",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "   ",
			expected: false,
		},
		{
			name: "just text",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  text",
			expected: false,
		},
		{
			name: "just numbers",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  123",
			expected: false,
		},
		{
			name: "unordered: with -",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  - item",
			expected: true,
		},
		{
			name: "unordered: with *",

			opt: &Options{
				BulletListMarker: "*",
			},
			line:     "  * item",
			expected: true,
		},
		{
			name: "unordered: with * false positive",

			opt: &Options{
				BulletListMarker: "*",
			},
			line:     "  - item",
			expected: false,
		},
		{
			name: "unordered: with multiple spaces",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  -  item",
			expected: true,
		},
		{
			name: "unordered: without space",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  -item",
			expected: false,
		},
		{
			name: "ordered: without space",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  1.item",
			expected: false,
		},
		{
			name: "ordered: with space",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  1. item",
			expected: true,
		},
		{
			name: "ordered: without dot",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  1 item",
			expected: false,
		},
		{
			name: "ordered: with dot before",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  .1 item",
			expected: false,
		},
		{
			name: "ordered: with big number",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  1001. item",
			expected: true,
		},
		{
			name: "ordered: with date",

			opt: &Options{
				BulletListMarker: "-",
			},
			line:     "  01.01 January",
			expected: false,
		},
		{
			name: "with divider",

			opt: &Options{
				BulletListMarker: "*",
			},
			line:     "***",
			expected: false,
		},
		{
			name: "with divider and spaces",

			opt: &Options{
				BulletListMarker: "*",
			},
			line:     "* * *",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isListItem(test.opt, test.line)
			if result != test.expected {
				t.Errorf("expected '%+v' but got '%+v' for '%s'", test.expected, result, test.line)
			}
		})
	}

}
