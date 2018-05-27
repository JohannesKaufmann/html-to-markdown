package plugin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

var youtubeID = regexp.MustCompile(`youtube\.com\/embed\/([^\&\?\/]+)`)

var EXPERIMENTAL_YoutubeEmbed = []md.Rule{
	md.Rule{
		Filter: []string{"iframe"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			src := selec.AttrOr("src", "")
			if !strings.Contains(src, "www.youtube.com") {
				return nil
			}
			alt := selec.AttrOr("title", "")

			parts := youtubeID.FindStringSubmatch(src)
			if len(parts) != 2 {
				return nil
			}
			id := parts[1]

			text := fmt.Sprintf("[![%s](https://img.youtube.com/vi/%s/0.jpg)](https://www.youtube.com/watch?v=%s)", alt, id, id)
			return &text
		},
	},
}
