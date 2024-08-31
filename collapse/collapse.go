/*

The function to collapse whitespace was adapted from the "turndown" library by Dom Christie,
which was adapted from the "collapse-whitespace" library by Luc Thevenard.

It was ported from Javascript to Golang by Johannes Kaufmann for the use in the "html-to-markdown" library.
To increase performance the use of regex was replaced by custom code.

https://github.com/wooorm/collapse-white-space
https://github.com/mixmark-io/turndown
https://github.com/JohannesKaufmann/html-to-markdown

-----------

MIT License

Copyright (c) 2017 Dom Christie
Copyright (c) 2014 Luc Thevenard <lucthevenard@gmail.com>
Copyright (c) 2018 Johannes Kaufmann

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

// collapse can collapse whitespace in html elements.
//
// It is a port from the Javascript library "turndown" to Golang.
package collapse

import (
	"strings"

	"github.com/JohannesKaufmann/dom"
	"golang.org/x/net/html"
)

// Note: Originally in the javascript version, this just checked for "pre".
// I changed it, to also return true for "code"
func isPreOrCode(node *html.Node) bool {
	name := dom.NodeName(node)

	return name == "pre" || name == "code"
}

func next(prev *html.Node, current *html.Node) *html.Node {
	if (prev != nil && prev.Parent == current) || isPreOrCode(current) {
		if current.NextSibling != nil {
			return current.NextSibling
		}

		return current.Parent
	}

	if current.FirstChild != nil {
		return current.FirstChild
	}
	if current.NextSibling != nil {
		return current.NextSibling
	}

	return current.Parent
}

var blockElements = []string{
	"address", "article", "aside", "audio", "blockquote", "body", "canvas", "center", "dd", "dir", "div", "dl", "dt", "fieldset", "figcaption", "figure", "footer", "form", "frameset", "h1", "h2", "h3", "h4", "h5", "h6", "header", "hgroup", "hr", "html", "isindex", "li", "main", "menu", "nav", "noframes", "noscript", "ol", "output", "p", "pre", "section", "table", "tbody", "td", "tfoot", "th", "thead", "tr", "ul",
}

var voidElements = []string{
	// Note: Compared to the javascript implementation, I removed "source"
	"area", "base", "br", "col", "command", "embed", "hr", "img", "input", "keygen", "link", "meta", "param" /* "source, "*/, "track", "wbr",
}

func isBlock(node *html.Node) bool {
	name := dom.NodeName(node)

	for _, elem := range blockElements {
		if elem == name {
			return true
		}
	}
	return false
}
func isVoid(node *html.Node) bool {
	name := dom.NodeName(node)

	for _, elem := range voidElements {
		if elem == name {
			return true
		}
	}
	return false
}

func remove(node *html.Node) *html.Node {
	next := node.NextSibling
	if next == nil {
		next = node.Parent
	}

	node.Parent.RemoveChild(node)

	return next

}

func Collapse(element *html.Node) {
	if element.FirstChild == nil || isPreOrCode(element) {
		return
	}

	var prevText *html.Node = nil
	var keepLeadingWs = false

	var prev *html.Node = nil
	var node = next(prev, element)

	for node != element {
		if node.Type == html.TextNode /* node.nodeType == 4 */ { // Node.TEXT_NODE or Node.CDATA_SECTION_NODE
			var text = replaceAnyWhitespaceWithSpace(node.Data)

			if (prevText == nil || strings.HasSuffix(prevText.Data, " ")) &&
				!keepLeadingWs && text[0] == ' ' {
				text = text[1:]
			}

			// `text` might be empty at this point.
			if text == "" {
				node = remove(node)
				continue
			}

			node.Data = text

			prevText = node
		} else if node.Type == html.ElementNode { // Node.ELEMENT_NODE
			if isBlock(node) || dom.NodeName(node) == "br" {
				if prevText != nil {
					prevText.Data = strings.TrimSuffix(prevText.Data, " ")
				}

				prevText = nil
				keepLeadingWs = false
			} else if isVoid(node) || isPreOrCode(node) || node.Data == "code" {
				// Avoid trimming space around non-block, non-BR void elements and inline PRE.
				prevText = nil
				keepLeadingWs = true
			} else if prevText != nil {
				// Drop protection if set previously.
				keepLeadingWs = false
			}
		} else if node.Type == html.CommentNode {
			// TODO: Is this enough to keep the comments? Does this cause other problems?
		} else {
			// E.g. DoctypeNode

			node = remove(node)
			continue
		}

		var nextNode = next(prev, node)
		prev = node
		node = nextNode
	}

	if prevText != nil {
		prevText.Data = strings.TrimSuffix(prevText.Data, " ")
		if prevText.Data == "" {
			remove(prevText)
		}
	}
}
