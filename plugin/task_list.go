package plugin

import (
	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

var TaskListItems = []md.Rule{
	md.Rule{
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
