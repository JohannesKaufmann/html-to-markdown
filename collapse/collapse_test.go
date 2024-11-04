package collapse

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func getBody(doc *html.Node) *html.Node {
	var body *html.Node

	var finder func(*html.Node)
	finder = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			finder(child)
		}
	}
	finder(doc)

	return body
}

func TestCollapse_DocType(t *testing.T) {
	// The DOCTYPE gets removed
	input := `<!DOCTYPE html><html><head></head><body></body></html>`

	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		t.Error(err)
	}

	Collapse(doc, nil)

	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	if err != nil {
		t.Error(err)
	}

	expected := `<html><head></head><body></body></html>`
	if buf.String() != expected {
		t.Errorf("expected %q but got %q", expected, buf.String())
	}
}

func TestCollapse_NoFirstChild(t *testing.T) {
	boldNode := &html.Node{
		Type: html.ElementNode,
		Data: "strong",
	}

	Collapse(boldNode, nil)

	var buf bytes.Buffer
	err := html.Render(&buf, boldNode)
	if err != nil {
		t.Error(err)
	}

	expected := `<strong></strong>`
	if buf.String() != expected {
		t.Errorf("expected %q but got %q", expected, buf.String())
	}
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

	var buf bytes.Buffer
	err := html.Render(&buf, codeNode)
	if err != nil {
		t.Error(err)
	}

	expected := `<code>  text  </code>`
	if buf.String() != expected {
		t.Errorf("expected %q but got %q", expected, buf.String())
	}
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

	var buf bytes.Buffer
	err := html.Render(&buf, node1)
	if err != nil {
		t.Error(err)
	}

	expected := `<span>a b</span>`
	if buf.String() != expected {
		t.Errorf("expected %q but got %q", expected, buf.String())
	}
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

	var buf bytes.Buffer
	err := html.Render(&buf, node1)
	if err != nil {
		t.Error(err)
	}

	expected := `<span>text</span>`
	if buf.String() != expected {
		t.Errorf("expected %q but got %q", expected, buf.String())
	}
}

func TestCollapse_Table(t *testing.T) {
	runs := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "basic example",
			input:    "   <p>Foo   bar</p>  <p>Words</p> ",
			expected: "<body><p>Foo bar</p><p>Words</p></body>",
		},
		{
			desc:     "without whitespace",
			input:    "<p>Some<strong>Text</strong></p>",
			expected: "<body><p>Some<strong>Text</strong></p></body>",
		},
		{
			desc:     "with one space & space in paragraph",
			input:    "<p>Some <strong> text. </strong></p>",
			expected: "<body><p>Some <strong>text.</strong></p></body>",
		},
		{
			desc:     "with one space",
			input:    "<p>Some<strong> text. </strong></p>",
			expected: "<body><p>Some<strong> text.</strong></p></body>",
		},
		{
			desc:     "with three space",
			input:    "<p>Some<strong>   text.   </strong></p>",
			expected: "<body><p>Some<strong> text.</strong></p></body>",
		},
		{
			desc:     "with three space (at beginning of paragraph)",
			input:    "<p><strong>   text.   </strong></p>",
			expected: "<body><p><strong>text.</strong></p></body>",
		},
		{
			desc:     "with image between",
			input:    `<p><strong>  a  </strong><img src="/img.png" /><strong>  b  </strong></p>`,
			expected: `<body><p><strong>a </strong><img src="/img.png"/><strong> b</strong></p></body>`,
		},
		{
			desc:     "spans directly next to each other",
			input:    "<p><span>(Text A)</span><span>(Text B)</span></p>",
			expected: "<body><p><span>(Text A)</span><span>(Text B)</span></p></body>",
		},
		{
			desc:     "spans with newline between each other",
			input:    "<p>\n<span>(Text A)</span>\n<span>(Text B)</span>\n</p>",
			expected: "<body><p><span>(Text A)</span> <span>(Text B)</span></p></body>",
		},
		{
			desc:  "code with space",
			input: "<p><code> </code>aaa</p>",
			// Note: This is different thant the javascript implementation.
			// We want the space to be preserved.
			expected: "<body><p><code> </code>aaa</p></body>",
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
			expected: `<body><h2><div>Browse<ul><li><a href="/go">go</a></li></ul>or <a href="/ask">ask</a>.</div></h2></body>`,
		},

		// - - - - - - //
		{
			desc:     "mdn example: inline formatting context",
			input:    "<h1>   Hello \n\t\t\t\t<span> World!</span>\t  </h1>",
			expected: "<body><h1>Hello <span>World!</span></h1></body>",
			// -> https://developer.mozilla.org/en-US/docs/Web/API/Document_Object_Model/Whitespace
		},
		{
			desc:     "mdn example: block formatting contexts",
			input:    "<body>\n\t<div>  Hello  </div>\n\n   <div>  World!  </div>  \n</body>",
			expected: "<body><div>Hello</div><div>World!</div></body>",
		},

		// - - - - - - Comments - - - - - - //
		{
			desc:     "#comment inside paragraph",
			input:    `<p>before<!-- my comment -->after</p>`,
			expected: `<body><p>before<!-- my comment -->after</p></body>`,
		},
		{
			desc:     "#comment inside paragraph (with spaces)",
			input:    `<p>before  <!-- my comment -->  after</p>`,
			expected: `<body><p>before <!-- my comment -->after</p></body>`,
		},
		{
			desc:     "#comment inside div",
			input:    `<div>before<!-- my comment -->after</div>`,
			expected: `<body><div>before<!-- my comment -->after</div></body>`,
		},
		{
			desc:     "#comment inside div (with spaces)",
			input:    `<div>before  <!-- my comment -->  after</div>`,
			expected: `<body><div>before <!-- my comment -->after</div></body>`,
		},
	}

	for _, run := range runs {
		t.Run(run.desc, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(run.input))
			if err != nil {
				t.Error(err)
			}

			Collapse(doc, nil)

			var buf bytes.Buffer
			err = html.Render(&buf, getBody(doc))
			if err != nil {
				t.Error(err)
			}

			if buf.String() != run.expected {
				t.Errorf("expected %q but got %q", run.expected, buf.String())
			}
		})
	}
}
