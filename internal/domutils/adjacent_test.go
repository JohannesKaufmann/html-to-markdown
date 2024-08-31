package domutils

import (
	"testing"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
	"golang.org/x/net/html"
)

func TestMergeAdjacent(t *testing.T) {
	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:  "don't change other tags",
			input: `<span>a</span> <span>b</span>`,
			expected: `
├─body
│ ├─span
│ │ ├─#text "a"
│ ├─#text " "
│ ├─span
│ │ ├─#text "b"
			`,
		},
		{
			desc:  "don't change simple strong",
			input: `<strong>a</strong>`,
			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
			`,
		},
		{
			desc:  "dont merge two adjacent strong tags with space between",
			input: `<strong>a</strong> <strong>b</strong>`,
			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
│ ├─#text " "
│ ├─strong
│ │ ├─#text "b"
			`,
		},
		{
			desc:  "merge two adjacent strong tags without space between",
			input: `<strong>a</strong><strong>b</strong>`,
			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
│ │ ├─#text "b"
			`,
		},
		{
			desc:  "merge three adjacent strong tags without space between",
			input: `<strong>a</strong><strong>b</strong><strong>c</strong>`,
			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
│ │ ├─#text "b"
│ │ ├─#text "c"
			`,
		},
		{
			desc:  "merge four adjacent strong tags without space between",
			input: `<strong>a</strong><strong>b</strong><strong>c</strong><strong>d</strong>`,
			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
│ │ ├─#text "b"
│ │ ├─#text "c"
│ │ ├─#text "d"
			`,
		},
		{
			desc:  "dont merge if there is tag content between",
			input: `<strong>a</strong><p>between</p><strong>b</strong>`,
			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
│ ├─p
│ │ ├─#text "between"
│ ├─strong
│ │ ├─#text "b"
			`,
		},
		{
			desc:  "dont merge if there is #text content between",
			input: `<strong>a</strong> between <strong>b</strong>`,
			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
│ ├─#text " between "
│ ├─strong
│ │ ├─#text "b"
			`,
		},
		{
			desc:  "dont merge if there is break between",
			input: `<strong>a</strong><br/><strong>b</strong>`,
			expected: `
├─body
│ ├─strong
│ │ ├─#text "a"
│ ├─br
│ ├─strong
│ │ ├─#text "b"
			`,
		},
		{
			desc:  "merge three adjacent italic tags without space between",
			input: `<em>a</em><em>b</em><em>c</em>`,
			expected: `
├─body
│ ├─em
│ │ ├─#text "a"
│ │ ├─#text "b"
│ │ ├─#text "c"
			`,
		},

		{
			desc:  "dont merge two nested strong tags with space between",
			input: `<div><strong>A</strong></div> <strong>B</strong>`,
			expected: `
├─body
│ ├─div
│ │ ├─strong
│ │ │ ├─#text "A"
│ ├─#text " "
│ ├─strong
│ │ ├─#text "B"

			`,
		},

		{
			desc:  "(for now) dont merge nested strongs inside div",
			input: `<div><strong>A</strong></div><strong>B</strong>`,
			expected: `
├─body
│ ├─div
│ │ ├─strong
│ │ │ ├─#text "A"
│ ├─strong
│ │ ├─#text "B"
			`,
		},
		{
			desc:  "(for now) dont merge deeply nested strongs inside div",
			input: `<div><div><div><strong>A</strong></div></div><div><strong>b</strong></div></div>`,
			expected: `
├─body
│ ├─div
│ │ ├─div
│ │ │ ├─div
│ │ │ │ ├─strong
│ │ │ │ │ ├─#text "A"
│ │ ├─div
│ │ │ ├─strong
│ │ │ │ ├─#text "b"
			`,
		},

		{
			desc:  "dont merge two nested strong tags enclosed in a",
			input: `<a href="/"><strong>A</strong></a><strong>B</strong>`,
			expected: `
├─body
│ ├─a (href="/")
│ │ ├─strong
│ │ │ ├─#text "A"
│ ├─strong
│ │ ├─#text "B"
			`,
		},

		// - - - - - - - - - - - Span - - - - - - - - - - - //
		{
			desc:  "merge next strong nested in span #1",
			input: `<p><strong>a</strong><span><strong>b</strong></span>other text</p>`,
			expected: `
├─body
│ ├─p
│ │ ├─strong
│ │ │ ├─#text "a"
│ │ │ ├─#text "b"
│ │ ├─span
│ │ ├─#text "other text"
			`,
		},
		{
			desc:  "merge next strong nested in span #2",
			input: `<p><strong>a</strong><span><span><strong>b</strong></span></span>other text</p>`,
			expected: `
├─body
│ ├─p
│ │ ├─strong
│ │ │ ├─#text "a"
│ │ │ ├─#text "b"
│ │ ├─span
│ │ │ ├─span
│ │ ├─#text "other text"
			`,
		},
		{
			desc:  "merge next strong nested in span #3",
			input: `<p><strong>a</strong><span><strong>b</strong></span><span><strong>c</strong>other text</span></p>`,
			expected: `
├─body
│ ├─p
│ │ ├─strong
│ │ │ ├─#text "a"
│ │ │ ├─#text "b"
│ │ │ ├─#text "c"
│ │ ├─span
│ │ ├─span
│ │ │ ├─#text "other text"
			`,
		},
		{
			desc:  "dont merge other span tags",
			input: `<p><strong>a</strong><span>other text</span></p>`,
			expected: `
├─body
│ ├─p
│ │ ├─strong
│ │ │ ├─#text "a"
│ │ ├─span
│ │ │ ├─#text "other text"
			`,
		},
		{
			desc:  "dont merge span content if space between",
			input: `<p><strong>a</strong><span> <strong>b</strong></span></p>`,
			expected: `
├─body
│ ├─p
│ │ ├─strong
│ │ │ ├─#text "a"
│ │ ├─span
│ │ │ ├─#text " "
│ │ │ ├─strong
│ │ │ │ ├─#text "b"
			`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc := tester.Parse(t, run.input, "")

			MergeAdjacent(doc, func(n *html.Node) bool {
				name := dom.NodeName(n)
				return name == "strong" || name == "em"
			})

			tester.ExpectRepresentation(t, doc, "output", run.expected)
		})
	}
}

func TestMergeAdjacentTextNodes(t *testing.T) {
	div := &html.Node{
		Type: html.ElementNode,
		Data: "div",
	}
	textOne := &html.Node{
		Type: html.TextNode,
		Data: "one",
	}
	textTwo := &html.Node{
		Type: html.TextNode,
		Data: "two",
	}
	textThree := &html.Node{
		Type: html.TextNode,
		Data: "three",
	}
	div.AppendChild(textOne)
	div.AppendChild(textTwo)
	div.AppendChild(textThree)

	MergeAdjacentTextNodes(div)

	expected := `
div
├─#text "onetwothree"
	`
	tester.ExpectRepresentation(t, div, "output", expected)
}
