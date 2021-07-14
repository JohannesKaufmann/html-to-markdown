package plugin

import (
	"regexp"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

// TableCompat is a compatibility plugon for environments where
// only commonmark markdown (without Tables) is supported.
//
// Note: In an environment that supports "real" Tables, like GitHub's Flavored Markdown
// use `plugin.Table()` instead.
func TableCompat() md.Plugin {
	return func(c *md.Converter) []md.Rule {
		return []md.Rule{
			{
				Filter: []string{"td", "th"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					content = strings.TrimSpace(content)

					if content == "" {
						return &content
					}

					next := selec.Next()
					nextIsEmpty := strings.TrimSpace(next.Text()) == ""
					if (next.Is("td") || next.Is("th")) && !nextIsEmpty {
						content = content + " · "
					}

					return &content
				},
			},
			{
				Filter: []string{"tr"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					content = content + "\n\n"

					return &content
				},
			},
		}
	}
}

// Table converts a html table (using hyphens and pipe characters) to a
// visuall representation in markdown.
//
// Note: This Plugin overrides the default compatibility rules from `commonmark.go`.
// Only use this Plugin in an environment that has extendeded the normal syntax,
// like GitHub's Flavored Markdown.
func Table() md.Plugin {
	return func(c *md.Converter) []md.Rule {
		c.Before(func(selec *goquery.Selection) {
			selec.Find("caption").Each(func(i int, s *goquery.Selection) {
				parent := s.Parent()
				if !parent.Is("table") {
					return
				}

				// move the caption from inside the table to after the table
				parent.AfterSelection(s)
			})
		})

		return []md.Rule{
			{
				Filter: []string{"table"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					noHeader := selec.Find("thead").Length() == 0 && selec.Find("th").Length() == 0
					if noHeader {
						var maxCount int
						selec.Find("tr").Each(func(i int, s *goquery.Selection) {
							count := s.Children().Length()
							if count > maxCount {
								maxCount = count
							}
						})

						// add an empty header, so that the table is recognized.
						header := "|" + strings.Repeat("     |", maxCount)
						divider := "|" + strings.Repeat(" --- |", maxCount)

						content = header + "\n" + divider + content
					}

					content = "\n\n" + content + "\n\n"
					return &content
				},
			},
			{ // TableCell
				Filter: []string{"th", "td"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					return md.String(getCellContent(content, selec))
				},
			},
			{ // TableRow
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

							borderCells += getCellContent(border, s)
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
	}
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

var newLineRe = regexp.MustCompile(`(\r?\n)+`)

func getCellContent(content string, s *goquery.Selection) string {
	content = strings.TrimSpace(content)
	if s.Find("table").Length() == 0 {
		// nested tables not found
		content = newLineRe.ReplaceAllString(content, "<br>")
	}
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
