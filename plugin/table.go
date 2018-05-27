package plugin

import (
	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

// TODO: maybe something like TableCompat for environments
// where only commonmark markdown is supported.

// EXPERIMENTAL_Table converts a html table to markdown.
var EXPERIMENTAL_Table = []md.Rule{
	md.Rule{ // TableCell
		Filter: []string{"th", "td"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			return md.String(cell(content, selec))
		},
	},
	md.Rule{ // TableRow
		Filter: []string{"tr"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			borderCells := ""

			if isHeadingRow(selec) {
				selec.Children().Each(func(i int, s *goquery.Selection) {
					border := "---"
					if align, ok := s.Attr("align"); ok {
						switch align {
						case "left":
							border = ":--"
						case "right":
							border = "--:"
						case "center":
							border = ":-:"
						}
					}

					borderCells += cell(border, s)
				})
			}

			text := "\n" + content
			if borderCells != "" {
				text += "\n" + borderCells
			}
			return &text
		},
	},
}

// A tr is a heading row if:
// - the parent is a THEAD
// - or if its the first child of the TABLE or the first TBODY (possibly
//   following a blank THEAD)
// - and every cell is a TH
func isHeadingRow(s *goquery.Selection) bool {
	parent := s.Parent()

	if goquery.NodeName(parent) == "thead" {
		return true
	}

	isTableOrBody := parent.Is("table") || isFirstTbody(parent)

	everyTH := true
	s.Children().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) != "th" {
			everyTH = false
		}
	})

	if parent.Children().First().IsSelection(s) && isTableOrBody && everyTH {
		return true
	}

	return false
}
func isFirstTbody(s *goquery.Selection) bool {
	firstSibling := s.Siblings().Eq(0) // TODO: previousSibling
	if s.Is("tbody") && firstSibling.Length() == 0 {
		return true
	}

	return false
}

// function isHeadingRow (tr) {
//   var parentNode = tr.parentNode
//   return (
//     parentNode.nodeName === 'THEAD' ||
//     (
//       parentNode.firstChild === tr &&
//       (parentNode.nodeName === 'TABLE' || isFirstTbody(parentNode)) &&
//       every.call(tr.childNodes, function (n) { return n.nodeName === 'TH' })
//     )
//   )
// }
func cell(content string, s *goquery.Selection) string {
	index := -1
	for i, node := range s.Parent().Children().Nodes {
		if s.IsNodes(node) {
			index = i
			break
		}
	}
	prefix := " "
	if index == 0 {
		prefix = "| "
	}
	return prefix + content + " |"
}
