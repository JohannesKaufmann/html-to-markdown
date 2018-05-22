package plugin

import (
	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

// Strikethrough converts `<strike>`, `<s>`, and `<del>` elements
var Strikethrough = []md.Rule{
	md.Rule{
		Filter: []string{"del", "s", "strike"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			return md.String("~" + content + "~")
		},
	},
}
