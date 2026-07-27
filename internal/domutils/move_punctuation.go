package domutils

import (
	"context"
	"unicode"
	"unicode/utf8"

	"github.com/JohannesKaufmann/dom"
	"golang.org/x/net/html"
)

// The commonmark spec (section 6.2) only recognizes a delimiter run as emphasis
// if it "flanks" the content. A "Unicode punctuation character" is a character
// in the P or S general category.
func isPunctuation(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

// The start and the end of a line count as whitespace.
func isWhitespace(r rune) bool {
	return r == 0 || unicode.IsSpace(r)
}

// CanOpen reports whether a delimiter run between prev and next is
// left-flanking — only then can it open emphasis.
func CanOpen(prev, next rune) bool {
	if isWhitespace(next) {
		return false
	}
	if !isPunctuation(next) {
		return true
	}
	return isWhitespace(prev) || isPunctuation(prev)
}

// CanClose reports whether a delimiter run between prev and next is
// right-flanking — only then can it close emphasis.
func CanClose(prev, next rune) bool {
	if isWhitespace(prev) {
		return false
	}
	if !isPunctuation(prev) {
		return true
	}
	return isWhitespace(next) || isPunctuation(next)
}

// A span carries no markdown meaning, so we look through it.
func textNodeAtEdge(node *html.Node, fromStart bool) *html.Node {
	for node != nil {
		if node.Type == html.TextNode {
			return node
		}
		if dom.NodeName(node) != "span" {
			return nil
		}
		if fromStart {
			node = node.FirstChild
		} else {
			node = node.LastChild
		}
	}
	return nil
}

func prevTextNode(node *html.Node) *html.Node {
	return textNodeAtEdge(node.PrevSibling, false)
}
func nextTextNode(node *html.Node) *html.Node {
	return textNodeAtEdge(node.NextSibling, true)
}

func lastRune(node *html.Node) rune {
	if node == nil {
		return 0
	}
	r, size := utf8.DecodeLastRuneInString(node.Data)
	if size == 0 {
		return 0
	}
	return r
}
func firstRune(node *html.Node) rune {
	if node == nil {
		return 0
	}
	r, size := utf8.DecodeRuneInString(node.Data)
	if size == 0 {
		return 0
	}
	return r
}

// MovePunctuation shifts punctuation from the edges of a bold/italic/strikethrough
// node to the outside, but only where the delimiters would otherwise not flank the
// content — "before<em>.middle</em>after" cannot be expressed as "before*.middle*after",
// because that opening "*" neither opens nor closes emphasis, so it renders literally.
// Moving a single character is enough: a punctuation neighbour satisfies the rule.
//
// This mirrors what DelimiterForEveryLine already does with surrounding whitespace.
func MovePunctuation(ctx context.Context, doc *html.Node, matchFn func(*html.Node) bool) {
	node := doc
	for node != nil {
		if matchFn(node) {
			movePunctuationOut(node)
		}

		node = dom.GetNextNeighborElement(node)
	}
}

func movePunctuationOut(node *html.Node) {
	prev := prevTextNode(node)
	next := nextTextNode(node)

	if lead := textNodeAtEdge(node.FirstChild, true); lead != nil {
		r, size := utf8.DecodeRuneInString(lead.Data)
		if size > 0 && isPunctuation(r) && !CanOpen(lastRune(prev), r) {
			// Only reachable with a neighbour that is neither whitespace nor punctuation.
			lead.Data = lead.Data[size:]
			prev.Data += string(r)
		}
	}

	if trail := textNodeAtEdge(node.LastChild, false); trail != nil {
		r, size := utf8.DecodeLastRuneInString(trail.Data)
		if size > 0 && isPunctuation(r) && !CanClose(r, firstRune(next)) {
			trail.Data = trail.Data[:len(trail.Data)-size]
			next.Data = string(r) + next.Data
		}
	}
}
