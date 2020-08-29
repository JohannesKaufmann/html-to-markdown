package plugin

import (
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

const moveFrontmatterAttr = "movefrontmatter"

// EXPERIMENTALMoveFrontMatter moves a frontmatter block at the beginning
// of the document to the top of the generated markdown block, without touching (and escaping) it.
func EXPERIMENTALMoveFrontMatter(delimiters ...rune) md.Plugin {
	return func(c *md.Converter) []md.Rule {
		if len(delimiters) == 0 {
			delimiters = []rune{'+', '$', '-', '%'}
		}

		var delimitersList []string
		for _, c := range delimiters {
			delimitersList = append(delimitersList, strings.Repeat(string(c), 3))
		}

		isDelimiter := func(line string) bool {
			for _, delimiter := range delimitersList {
				if strings.HasPrefix(line, delimiter) {
					return true
				}
			}
			return false
		}

		c.Before(func(selec *goquery.Selection) {
			selec.Find("body").Contents().EachWithBreak(func(i int, s *goquery.Selection) bool {
				text := s.Text()

				// skip empty strings
				if strings.TrimSpace(text) == "" {
					return true
				}

				var frontmatter string
				var html string = text // if there is no frontmatter, keep the text

				lines := strings.Split(text, "\n")
				for i := 0; i < len(lines); i++ {
					if isDelimiter(lines[i]) {
						if i == 0 {
							continue
						}

						// split the frontmatter
						f := lines[:i+1]
						frontmatter = strings.Join(f, "\n")

						// and the html content AFTER the frontmatter
						h := lines[i+1:]
						html = strings.Join(h, "\n")
						break
					}
				}

				s.SetAttr(moveFrontmatterAttr, frontmatter)
				s.SetText(html)

				// the front matter must be the first thing in the file. So we break out of the loop
				return false
			})
		})

		return []md.Rule{
			{
				Filter: []string{"#text"},
				AdvancedReplacement: func(content string, selec *goquery.Selection, opt *md.Options) (md.AdvancedResult, bool) {
					frontmatter, exists := selec.Attr(moveFrontmatterAttr)

					if !exists {
						return md.AdvancedResult{}, true
					}

					return md.AdvancedResult{
						Header:   frontmatter,
						Markdown: content,
					}, false
				},
			},
		}
	}
}
