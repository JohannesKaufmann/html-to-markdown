package plugin

import (
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

// Strikethrough converts `<strike>`, `<s>`, and `<del>` elements
func Strikethrough(character string) md.Plugin {
	return func(c *md.Converter) []md.Rule {
		if character == "" {
			character = "~"
		}

		return []md.Rule{
			md.Rule{
				Filter: []string{"del", "s", "strike"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					// trim spaces so that the following does NOT happen: `~ and cake~`
					content = strings.TrimSpace(content)
					return md.String(character + content + character)
				},
			},
		}
	}
}
