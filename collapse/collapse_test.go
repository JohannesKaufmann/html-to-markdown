package collapse

import (
	"strings"
	"testing"

	"github.com/JohannesKaufmann/html-to-markdown/v2/internal/tester"
	"golang.org/x/net/html"
)

func TestCollapse_DocType(t *testing.T) {
	// The DOCTYPE gets removed
	input := `<!DOCTYPE html><html><head></head><body></body></html>`

	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		t.Error(err)
	}

	Collapse(doc, nil)

	tester.ExpectRepresentation(t, doc, "after", `
#document
├─html
│ ├─head
│ ├─body
	`)
}

func TestCollapse_NoFirstChild(t *testing.T) {
	boldNode := &html.Node{
		Type: html.ElementNode,
		Data: "strong",
	}

	Collapse(boldNode, nil)

	tester.ExpectRepresentation(t, boldNode, "after", `strong`)
}

func TestCollapse_StartWithCode(t *testing.T) {
	textNode := &html.Node{
		Type: html.TextNode,
		Data: "  text  ",
	}
	codeNode := &html.Node{
		Type: html.ElementNode,
		Data: "code",
	}
	codeNode.AppendChild(textNode)

	Collapse(codeNode, nil)

	tester.ExpectRepresentation(t, codeNode, "after", `
code
├─#text "  text  "
	`)
}

func TestCollapse_TwoTextNodes(t *testing.T) {
	node1 := &html.Node{
		Type: html.ElementNode,
		Data: "span",
	}

	node2 := &html.Node{
		Type: html.TextNode,
		Data: "  a  ",
	}
	node3 := &html.Node{
		Type: html.TextNode,
		Data: "  b  ",
	}
	node1.AppendChild(node2)
	node1.AppendChild(node3)

	Collapse(node1, nil)

	tester.ExpectRepresentation(t, node1, "after", `
span
├─#text "a "
├─#text "b"
	`)
}

func TestCollapse_LastTextIsEmpty(t *testing.T) {
	node1 := &html.Node{
		Type: html.ElementNode,
		Data: "span",
	}

	node2 := &html.Node{
		Type: html.TextNode,
		Data: "text",
	}
	node3 := &html.Node{
		Type: html.TextNode,
		Data: " ",
	}
	node1.AppendChild(node2)
	node1.AppendChild(node3)

	Collapse(node1, nil)

	tester.ExpectRepresentation(t, node1, "after", `
span
├─#text "text"
	`)
}

func TestCollapse_Table(t *testing.T) {
	runs := []struct {
		desc  string
		input string

		expectedBefore string // optional
		expectedAfter  string
	}{
		{
			desc:  "basic example",
			input: "   <p>Foo   bar</p>  <p>Words</p> ",
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "Foo   bar"
│ │ ├─#text "  "
│ │ ├─p
│ │ │ ├─#text "Words"
│ │ ├─#text " "
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "Foo bar"
│ │ ├─p
│ │ │ ├─#text "Words"
			`,
		},
		{
			desc:  "without whitespace",
			input: "<p>Some<strong>Text</strong></p>",
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "Some"
│ │ │ ├─strong
│ │ │ │ ├─#text "Text"
			`,
		},
		{
			desc:  "with one space & space in paragraph",
			input: "<p>Some <strong> text. </strong></p>",
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "Some "
│ │ │ ├─strong
│ │ │ │ ├─#text "text."
			`,
		},
		{
			desc:  "with one space",
			input: "<p>Some<strong> text. </strong></p>",
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "Some"
│ │ │ ├─strong
│ │ │ │ ├─#text " text. "
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "Some"
│ │ │ ├─strong
│ │ │ │ ├─#text " text."
			`,
		},
		{
			desc:  "with three space",
			input: "<p>Some<strong>   text.   </strong></p>",
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "Some"
│ │ │ ├─strong
│ │ │ │ ├─#text "   text.   "
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "Some"
│ │ │ ├─strong
│ │ │ │ ├─#text " text."
			`,
		},
		{
			desc:  "with three space (at beginning of paragraph)",
			input: "<p><strong>   text.   </strong></p>",
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─strong
│ │ │ │ ├─#text "text."
			`,
		},
		{
			desc:  "with image between",
			input: `<p><strong>  a  </strong><img src="/img.png" /><strong>  b  </strong></p>`,
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─strong
│ │ │ │ ├─#text "  a  "
│ │ │ ├─img (src="/img.png")
│ │ │ ├─strong
│ │ │ │ ├─#text "  b  "
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─strong
│ │ │ │ ├─#text "a "
│ │ │ ├─img (src="/img.png")
│ │ │ ├─strong
│ │ │ │ ├─#text " b"
			`,
		},
		{
			desc:  "spans directly next to each other",
			input: "<p><span>(Text A)</span><span>(Text B)</span></p>",
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─span
│ │ │ │ ├─#text "(Text A)"
│ │ │ ├─span
│ │ │ │ ├─#text "(Text B)"
			`,
		},
		{
			desc:  "spans with newline between each other",
			input: "<p>\n<span>(Text A)</span>\n<span>(Text B)</span>\n</p>",
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "\n"
│ │ │ ├─span
│ │ │ │ ├─#text "(Text A)"
│ │ │ ├─#text "\n"
│ │ │ ├─span
│ │ │ │ ├─#text "(Text B)"
│ │ │ ├─#text "\n"
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─span
│ │ │ │ ├─#text "(Text A)"
│ │ │ ├─#text " "
│ │ │ ├─span
│ │ │ │ ├─#text "(Text B)"
│ │ │ ├─#text ""
			`,
		},
		{
			desc: "spans with indentation",
			input: `
			<div>
				<span>A</span>
				<span>B</span>
			</div>
			<div>
				<span>C</span>
			</div>
			`,
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─div
│ │ │ ├─#text "\n\t\t\t\t"
│ │ │ ├─span
│ │ │ │ ├─#text "A"
│ │ │ ├─#text "\n\t\t\t\t"
│ │ │ ├─span
│ │ │ │ ├─#text "B"
│ │ │ ├─#text "\n\t\t\t"
│ │ ├─#text "\n\t\t\t"
│ │ ├─div
│ │ │ ├─#text "\n\t\t\t\t"
│ │ │ ├─span
│ │ │ │ ├─#text "C"
│ │ │ ├─#text "\n\t\t\t"
│ │ ├─#text "\n\t\t\t"
			`,

			// TODO: are we expecting empty #text nodes??!
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─div
│ │ │ ├─span
│ │ │ │ ├─#text "A"
│ │ │ ├─#text " "
│ │ │ ├─span
│ │ │ │ ├─#text "B"
│ │ │ ├─#text ""
│ │ ├─div
│ │ │ ├─span
│ │ │ │ ├─#text "C"
│ │ │ ├─#text ""
			`,
		},
		{
			desc:  "code with space",
			input: "<p><code> </code>aaa</p>",
			// Note: This is different then the javascript implementation.
			// We want the space to be preserved.
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─code
│ │ │ │ ├─#text " "
│ │ │ ├─#text "aaa"
			`,
		},
		{
			desc: "#text in sample",
			input: `
			<h2>
			  <div>
				Browse
				<ul>
				  <li><a href="/go">go</a></li>
				</ul>
				or <a href="/ask">ask</a>.
			  </div>
			</h2>
			`,
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─h2
│ │ │ ├─#text "\n\t\t\t  "
│ │ │ ├─div
│ │ │ │ ├─#text "\n\t\t\t\tBrowse\n\t\t\t\t"
│ │ │ │ ├─ul
│ │ │ │ │ ├─#text "\n\t\t\t\t  "
│ │ │ │ │ ├─li
│ │ │ │ │ │ ├─a (href="/go")
│ │ │ │ │ │ │ ├─#text "go"
│ │ │ │ │ ├─#text "\n\t\t\t\t"
│ │ │ │ ├─#text "\n\t\t\t\tor "
│ │ │ │ ├─a (href="/ask")
│ │ │ │ │ ├─#text "ask"
│ │ │ │ ├─#text ".\n\t\t\t  "
│ │ │ ├─#text "\n\t\t\t"
│ │ ├─#text "\n\t\t\t"
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─h2
│ │ │ ├─div
│ │ │ │ ├─#text "Browse"
│ │ │ │ ├─ul
│ │ │ │ │ ├─li
│ │ │ │ │ │ ├─a (href="/go")
│ │ │ │ │ │ │ ├─#text "go"
│ │ │ │ ├─#text "or "
│ │ │ │ ├─a (href="/ask")
│ │ │ │ │ ├─#text "ask"
│ │ │ │ ├─#text "."
			`,
		},

		// - - - - - - //
		{
			desc:  "mdn example: inline formatting context",
			input: "<h1>   Hello \n\t\t\t\t<span> World!</span>\t  </h1>",
			// -> https://developer.mozilla.org/en-US/docs/Web/API/Document_Object_Model/Whitespace

			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─h1
│ │ │ ├─#text "   Hello \n\t\t\t\t"
│ │ │ ├─span
│ │ │ │ ├─#text " World!"
│ │ │ ├─#text "\t  "
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─h1
│ │ │ ├─#text "Hello "
│ │ │ ├─span
│ │ │ │ ├─#text "World!"
│ │ │ ├─#text ""
			`,
		},
		{
			desc:  "mdn example: block formatting contexts",
			input: "<body>\n\t<div>  Hello  </div>\n\n   <div>  World!  </div>  \n</body>",
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─#text "\n\t"
│ │ ├─div
│ │ │ ├─#text "  Hello  "
│ │ ├─#text "\n\n   "
│ │ ├─div
│ │ │ ├─#text "  World!  "
│ │ ├─#text "  \n"
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─div
│ │ │ ├─#text "Hello"
│ │ ├─div
│ │ │ ├─#text "World!"
			`,
		},

		// - - - - - - Comments - - - - - - //
		{
			desc:  "#comment inside paragraph",
			input: `<p>before<!-- my comment -->after</p>`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "before"
│ │ │ ├─#comment
│ │ │ ├─#text "after"
			`,
		},
		{
			desc:  "#comment inside paragraph (with spaces)",
			input: `<p>before  <!-- my comment -->  after</p>`,
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "before  "
│ │ │ ├─#comment
│ │ │ ├─#text "  after"
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─p
│ │ │ ├─#text "before "
│ │ │ ├─#comment
│ │ │ ├─#text "after"
			`,
		},
		{
			desc:  "#comment inside div",
			input: `<div>before<!-- my comment -->after</div>`,
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─div
│ │ │ ├─#text "before"
│ │ │ ├─#comment
│ │ │ ├─#text "after"
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─div
│ │ │ ├─#text "before"
│ │ │ ├─#comment
│ │ │ ├─#text "after"
			`,
		},
		{
			desc:  "#comment inside div (with spaces)",
			input: `<div>before  <!-- my comment -->  after</div>`,
			expectedBefore: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─div
│ │ │ ├─#text "before  "
│ │ │ ├─#comment
│ │ │ ├─#text "  after"
			`,
			expectedAfter: `
#document
├─html
│ ├─head
│ ├─body
│ │ ├─div
│ │ │ ├─#text "before "
│ │ │ ├─#comment
│ │ │ ├─#text "after"
			`,
		},
	}

	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(run.input))
			if err != nil {
				t.Fatal(err)
			}

			if run.expectedBefore != "" {
				tester.ExpectRepresentation(t, doc, "before", run.expectedBefore)
			}

			Collapse(doc, nil)

			tester.ExpectRepresentation(t, doc, "after", run.expectedAfter)
		})
	}
}
