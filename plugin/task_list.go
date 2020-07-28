package plugin

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

// TaskListItems converts checkboxes into task list items.
func TaskListItems() md.Plugin {
	return func(c *md.Converter) []md.Rule {
		return []md.Rule{
			{
				Filter: []string{"input"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					if !selec.Parent().Is("li") {
						return nil
					}
					if selec.AttrOr("type", "") != "checkbox" {
						return nil
					}

					_, ok := selec.Attr("checked")
					if ok {
						return md.String("[x] ")
					}
					return md.String("[ ] ")
				},
			},
		}
	}
}
