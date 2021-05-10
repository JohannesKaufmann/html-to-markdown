package main

import (
	"fmt"
	"log"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	html := `<my_video>https://youtu.be/1SoMeViD</my_video>
	<my_video>https://youtu.be/2SoMeViD</my_video>
	<my_video>https://youtu.be/3SoMeViD</my_video><my_video>https://youtu.be/4SoMeViD</my_video>
	
	<my_video>https://youtu.be/5SoMeViD</my_video>
	`

	videoRule := md.Rule{
		// We want to add a rule for a `my_video` tag.
		Filter: []string{"my_video"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			text := "click to watch video"

			// in this case, the content inside the tag is the url
			href := strings.TrimSpace(content)

			// format it, so that its `[click to watch video](https://youtu.be/1SoMeViD)\n\n`
			md := fmt.Sprintf("[%s](%s)\n\n", text, href)
			return &md
		},
	}

	conv := md.NewConverter("", true, nil)
	conv.AddRules(videoRule)
	// -> add 1+ rules to the converter. the last added will be used first.

	markdown, err := conv.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n\nresult:'%s'\n", markdown)
}
