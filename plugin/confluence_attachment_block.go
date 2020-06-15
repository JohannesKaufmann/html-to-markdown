package plugin

import (
	"fmt"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

// ConfluenceAttachments converts `<ri:attachment ri:filename=""/>` elements
// [Contributed by @Skarlso]
func ConfluenceAttachments() md.Plugin {
	return func(c *md.Converter) []md.Rule {
		return []md.Rule{
			{
				Filter: []string{"ri:attachment"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					if v, ok := selec.Attr("ri:filename"); ok {
						formatted := fmt.Sprintf("![][%s]", v)
						return md.String(formatted)
					}
					return md.String("")
				},
			},
		}
	}
}
