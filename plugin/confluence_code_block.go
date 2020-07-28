package plugin

import (
	"fmt"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

// ConfluenceCodeBlock converts `<ac:structured-macro>` elements that are used in Atlassian’s Wiki “Confluence”.
// [Contributed by @Skarlso]
func ConfluenceCodeBlock() md.Plugin {
	return func(c *md.Converter) []md.Rule {
		character := "```"
		return []md.Rule{
			{
				Filter: []string{"ac:structured-macro"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					for _, node := range selec.Nodes {
						if node.Data == "ac:structured-macro" {
							// node's last child -> <ac:plain-text-body>. We don't want to filter on that
							// because we would end up with structured-macro around us.
							// ac:plain-text-body's last child is [CDATA which has the actual content we are looking for.
							data := strings.TrimPrefix(node.LastChild.LastChild.Data, "[CDATA[")
							data = strings.TrimSuffix(data, "]]")
							// content, if set, will contain the language that has been set in the field.
							var language string
							if content != "" {
								language = content
							}
							formatted := fmt.Sprintf("%s%s\n%s\n%s", character, language, data, character)
							return md.String(formatted)
						}
					}
					return md.String(character + content + character)
				},
			},
		}
	}
}
