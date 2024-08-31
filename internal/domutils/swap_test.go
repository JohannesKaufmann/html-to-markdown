package domutils

import (
	"context"
	"testing"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func generateANodes() *html.Node {
	div := &html.Node{
		Namespace: "",
		Type:      html.ElementNode,
		DataAtom:  atom.Div,

		Attr: []html.Attribute{},
		Data: "div",
	}
	link := &html.Node{
		Namespace: "",
		Type:      html.ElementNode,
		DataAtom:  atom.A,

		Attr: []html.Attribute{
			{
				Namespace: "",
				Key:       "KeyA",
				Val:       "ValA",
			},
		},

		Data: "a",
	}
	text := &html.Node{
		Namespace: "",
		Type:      html.TextNode,
		DataAtom:  0,

		Data: "ContentA",
	}
	link.AppendChild(text)
	div.AppendChild(link)

	return div
}
func generateBNodes() *html.Node {
	div := &html.Node{
		Namespace: "",
		Type:      html.ElementNode,
		DataAtom:  atom.Main,

		Attr: []html.Attribute{},
		Data: "main",
	}
	link := &html.Node{
		Namespace: "",
		Type:      html.ElementNode,
		DataAtom:  atom.B,

		Attr: []html.Attribute{
			{
				Namespace: "",
				Key:       "KeyB",
				Val:       "ValB",
			},
		},

		Data: "b",
	}
	text := &html.Node{
		Namespace: "",
		Type:      html.TextNode,
		DataAtom:  0,

		Data: "ContentB",
	}
	link.AppendChild(text)
	div.AppendChild(link)

	return div
}
func TestSwapTagsOfNodes_Basics(t *testing.T) {
	a := generateANodes()
	b := generateBNodes()

	swapTagsOfNodes(a.FirstChild, b.FirstChild)

	// These fields should have changed:
	if a.FirstChild.DataAtom != atom.B {
		t.Error("expected different a atom")
	}
	if a.FirstChild.Data != "b" {
		t.Error("expected different a data")
	}
	if len(a.FirstChild.Attr) != 1 {
		t.Error("expected different a attributes length")
	}
	if a.FirstChild.Attr[0].Key != "KeyB" {
		t.Error("expected different a attribute key")
	}
	if a.FirstChild.Attr[0].Val != "ValB" {
		t.Error("expected different a attribute key")
	}

	// The pointers should NOT have changed:
	if a.FirstChild.Parent.Data != "div" {
		t.Error("expected the same parent for a")
	}
	if a.FirstChild.FirstChild.Data != "ContentA" {
		t.Error("expected the same text for a")
	}

	// - - - //

	if b.FirstChild.DataAtom != atom.A {
		t.Error("expected different b atom")
	}
	if b.FirstChild.Data != "a" {
		t.Error("expected different b data")
	}
	if len(b.FirstChild.Attr) != 1 {
		t.Error("expected different b attributes length")
	}
	if b.FirstChild.Attr[0].Key != "KeyA" {
		t.Error("expected different b attribute key")
	}
	if b.FirstChild.Attr[0].Val != "ValA" {
		t.Error("expected different b attribute key")
	}

	// The pointers should NOT have changed:
	if b.FirstChild.Parent.Data != "main" {
		t.Error("expected the same parent for b")
	}
	if b.FirstChild.FirstChild.Data != "ContentB" {
		t.Error("expected the same text for b")
	}
}

func TestSwapTags_HeadingLink(t *testing.T) {
	runs := []struct {
		desc string

		input     string
		startFrom string

		expected string
	}{
		{
			desc: "simple",

			input:     `<a href="/page.html"><h3>Heading</h3></a>`,
			startFrom: "body",

			expected: `
├─body
│ ├─h3
│ │ ├─a (href="/page.html")
│ │ │ ├─#text "Heading"
			`,
		},
		{
			desc: "simple with whitespace",

			input: `
<a href="/page.html">
	<h3>Heading</h3>
</a>
			`,
			startFrom: "body",

			expected: `
├─body
│ ├─h3
│ │ ├─#text "\n\t"
│ │ ├─a (href="/page.html")
│ │ │ ├─#text "Heading"
│ │ ├─#text "\n"
			`,
		},
		{
			desc: "more content",

			input: `
<a href="/reisen">
	<h3><span>Reiseinspiration</span><span>Beste Orte in Berlin</span></h3>
</a>
			`,
			startFrom: "body",

			expected: `
├─body
│ ├─h3
│ │ ├─#text "\n\t"
│ │ ├─a (href="/reisen")
│ │ │ ├─span
│ │ │ │ ├─#text "Reiseinspiration"
│ │ │ ├─span
│ │ │ │ ├─#text "Beste Orte in Berlin"
│ │ ├─#text "\n"
			`,
		},
		{
			desc: "not possible",

			input: `
<a href="/page.html">
	<h3>Heading</h3>
	<p>Some other content</p>
</a>
			`,
			startFrom: "body",

			expected: `
├─body
│ ├─a (href="/page.html")
│ │ ├─#text "\n\t"
│ │ ├─h3
│ │ │ ├─#text "Heading"
│ │ ├─#text "\n\t"
│ │ ├─p
│ │ │ ├─#text "Some other content"
│ │ ├─#text "\n"
			`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc := tester.Parse(t, run.input, run.startFrom)

			isLink := func(n *html.Node) bool {
				return dom.NodeName(n) == "a"
			}
			isHeading := func(n *html.Node) bool {
				name := dom.NodeName(n)

				if name == "h1" || name == "h2" || name == "h3" || name == "h4" || name == "h5" || name == "h6" {
					return true
				}
				return false
			}
			SwapTags(context.TODO(), doc, isLink, isHeading)

			tester.ExpectRepresentation(t, doc, "output", run.expected)

		})
	}
}

func TestSwapTags_PreCode(t *testing.T) {
	runs := []struct {
		desc  string
		input string

		expectedBefore string
		expectedAfter  string
	}{

		// - - - - - Pre - - - - - //
		{
			desc:  "div with pre: keep",
			input: "<div><pre>content</pre></div>",

			expectedBefore: `
├─body
│ ├─div
│ │ ├─pre
│ │ │ ├─#text "content"
			`,
			expectedAfter: `
├─body
│ ├─div
│ │ ├─pre
│ │ │ ├─#text "content"
			`,
		},
		{
			desc:  "p with pre: keep",
			input: "<p><pre>content</pre></p>",

			// The <pre> is a block node, so cannot be in a paragraph.
			expectedBefore: `
├─body
│ ├─p
│ ├─pre
│ │ ├─#text "content"
│ ├─p
			`,
			expectedAfter: `
├─body
│ ├─p
│ ├─pre
│ │ ├─#text "content"
│ ├─p
			`,
		},
		// - - - - - Code - - - - - //
		{
			desc:  "div with code: keep",
			input: "<div><code>content</code></div>",

			expectedBefore: `
├─body
│ ├─div
│ │ ├─code
│ │ │ ├─#text "content"
			`,
			expectedAfter: `
├─body
│ ├─div
│ │ ├─code
│ │ │ ├─#text "content"
			`,
		},
		{
			desc:  "p with code: keep",
			input: "<p><code>content</code></p>",

			expectedBefore: `
├─body
│ ├─p
│ │ ├─code
│ │ │ ├─#text "content"
			`,
			expectedAfter: `
├─body
│ ├─p
│ │ ├─code
│ │ │ ├─#text "content"
			`,
		},

		// - - - - - Nested in correct order - - - - - //
		{
			desc:  "keep correct code block",
			input: `<div><pre><code>content</code></pre></div>`,

			expectedBefore: `
├─body
│ ├─div
│ │ ├─pre
│ │ │ ├─code
│ │ │ │ ├─#text "content"
			`,
			expectedAfter: `
├─body
│ ├─div
│ │ ├─pre
│ │ │ ├─code
│ │ │ │ ├─#text "content"
			`,
		},
		// - - - - - Nested in wrong order - - - - - //
		{
			desc:  "swap wrong code block",
			input: `<div><code><pre>content</pre></code></div>`,

			expectedBefore: `
├─body
│ ├─div
│ │ ├─code
│ │ │ ├─pre
│ │ │ │ ├─#text "content"
			`,
			expectedAfter: `
├─body
│ ├─div
│ │ ├─pre
│ │ │ ├─code
│ │ │ │ ├─#text "content"
			`,
		},
		{
			desc:  "html parsing already causes swap",
			input: `<p><code><pre>content</pre></code></p>`,

			// Notice how the html parsing already looks different...
			expectedBefore: `
├─body
│ ├─p
│ │ ├─code
│ ├─pre
│ │ ├─code
│ │ │ ├─#text "content"
│ ├─p
			`,
			expectedAfter: `
├─body
│ ├─p
│ │ ├─code
│ ├─pre
│ │ ├─code
│ │ │ ├─#text "content"
│ ├─p
			`,
		},

		{
			desc:  "different ast then expected",
			input: `<p>before<code>a<pre>b</pre>c</code>after</p>`,

			expectedBefore: `
├─body
│ ├─p
│ │ ├─#text "before"
│ │ ├─code
│ │ │ ├─#text "a"
│ ├─pre
│ │ ├─code
│ │ │ ├─#text "b"
│ ├─code
│ │ ├─#text "c"
│ ├─#text "after"
│ ├─p
			`,
			expectedAfter: `
├─body
│ ├─p
│ │ ├─#text "before"
│ │ ├─code
│ │ │ ├─#text "a"
│ ├─pre
│ │ ├─code
│ │ │ ├─#text "b"
│ ├─code
│ │ ├─#text "c"
│ ├─#text "after"
│ ├─p
			`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc := tester.Parse(t, run.input, "")

			tester.ExpectRepresentation(t, doc, "before", run.expectedBefore)

			isCode := func(n *html.Node) bool {
				return dom.NodeName(n) == "code"
			}
			isPre := func(n *html.Node) bool {
				return dom.NodeName(n) == "pre"
			}
			SwapTags(context.TODO(), doc, isCode, isPre)

			tester.ExpectRepresentation(t, doc, "output", run.expectedAfter)
		})
	}
}

func TestSwapTags_StrongLinks(t *testing.T) {
	runs := []struct {
		desc  string
		input string

		expectedAfter string
	}{
		{
			desc:  "swap strong and link",
			input: `<p>before<strong><a href="/">middle</a></strong>after</p>`,

			expectedAfter: `
├─body
│ ├─p
│ │ ├─#text "before"
│ │ ├─a (href="/")
│ │ │ ├─strong
│ │ │ │ ├─#text "middle"
│ │ ├─#text "after"
			`,
		},
		{
			desc:  "empty span",
			input: `<p>before<strong><span></span><a href="/">with empty span</a><span></span></strong>after</p>`,

			expectedAfter: `
├─body
│ ├─p
│ │ ├─#text "before"
│ │ ├─strong
│ │ │ ├─span
│ │ │ ├─a (href="/")
│ │ │ │ ├─#text "with empty span"
│ │ │ ├─span
│ │ ├─#text "after"
			`,
		},
		{
			desc:  "span with spaces",
			input: `<p>before<strong><span>  </span><a href="/">with empty span</a><span>  </span></strong>after</p>`,

			expectedAfter: `
├─body
│ ├─p
│ │ ├─#text "before"
│ │ ├─strong
│ │ │ ├─span
│ │ │ │ ├─#text "  "
│ │ │ ├─a (href="/")
│ │ │ │ ├─#text "with empty span"
│ │ │ ├─span
│ │ │ │ ├─#text "  "
│ │ ├─#text "after"
			`,
		},
		{
			desc:  "spans nested",
			input: `<p>before<strong><span><span>  </span> </span><a href="/">with empty span</a><span><span>  </span> </span></strong>after</p>`,

			expectedAfter: `
├─body
│ ├─p
│ │ ├─#text "before"
│ │ ├─strong
│ │ │ ├─span
│ │ │ │ ├─span
│ │ │ │ │ ├─#text "  "
│ │ │ │ ├─#text " "
│ │ │ ├─a (href="/")
│ │ │ │ ├─#text "with empty span"
│ │ │ ├─span
│ │ │ │ ├─span
│ │ │ │ │ ├─#text "  "
│ │ │ │ ├─#text " "
│ │ ├─#text "after"
			`,
		},
	}
	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc := tester.Parse(t, run.input, "")

			isBoldOrItalic := func(node *html.Node) bool {
				name := dom.NodeName(node)
				if name == "strong" || name == "b" {
					return true
				}
				if name == "em" || name == "i" {
					return true
				}

				return false
			}

			isLink := func(node *html.Node) bool {
				return dom.NodeName(node) == "a"
			}

			// Remove all unnessesary span tags
			// for _, node := range dom.GetAllNodes(doc) {
			// 	name := dom.NodeName(node)
			// 	if name == "span" {
			// 		dom.UnwrapNode(node)
			// 	}
			// }

			// collapse.Collapse(doc)

			SwapTags(context.TODO(), doc, isBoldOrItalic, isLink)

			tester.ExpectRepresentation(t, doc, "output", run.expectedAfter)
		})
	}
}
