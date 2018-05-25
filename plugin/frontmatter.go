package plugin

import (
	"fmt"

	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

type FrontMatterCallback func(selec *goquery.Selection) map[string]interface{}

// TODO: automatically convert to formats (look at hugo)

func FrontMatter(format string) md.Plugin {
	return func(c *md.Converter) []md.Rule {

		title := c.Find("head title").Text()

		var text string
		switch format {
		case "toml": // +++
			text = fmt.Sprintf(`
+++
title = "%s"
+++
`, title)
		case "yaml": // ---
			text = fmt.Sprintf(`
---
title: %s
---
`, title)
		case "json": // { }
			text = fmt.Sprintf(`
{
	"title": "%s"
}
`, title)
		default:
			panic("unknown format")
		}

		c.AddLeading(text)
		return nil
	}
}
