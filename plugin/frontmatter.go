package plugin

import (
	"fmt"

	"github.com/JohannesKaufmann/html-to-markdown"
	yaml "gopkg.in/yaml.v2"
)

// type frontMatterCallback func(selec *goquery.Selection) map[string]interface{}

// TODO: automatically convert to formats (look at hugo)

func EXPERIMENTAL_FrontMatter(format string) md.Plugin {
	return func(c *md.Converter) []md.Rule {
		data := make(map[string]interface{})

		d, err := yaml.Marshal(data)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(d))
		/*
			add rule for `head`
				- get title
				- return AdvancedResult{ Header: formated_yaml }, skip
						-> added to head
						-> others rules can apply

		*/

		title := "" // c.Find("head title").Text()

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

		_ = text
		// c.AddLeading(text)
		return nil
	}
}
