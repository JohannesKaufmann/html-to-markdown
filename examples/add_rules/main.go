package main

import (
	"fmt"
	"log"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	html := `Good soundtrack <span class="bb_strike"> and cake</span>.`
	// -> `Good soundtrack ~~and cake~~.`

	/*
		We want to add a rule when a `span` tag has a class of `bb_strike`.
		Have a look at `plugin/strikethrough.go` to see how it is implemented normally.
	*/
	strikethrough := md.Rule{
		Filter: []string{"span"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			// If the span element has not the classname `bb_strike` return nil.
			// That way the next rules will apply. In this case the commonmark rules.
			// -> return nil -> next rule applies
			if !selec.HasClass("bb_strike") {
				return nil
			}

			// Trim spaces so that the following does NOT happen: `~ and cake~`.
			// Because of the space it is not recognized as strikethrough.
			// -> trim spaces at begin&end of string when inside strong/italic/...
			content = strings.TrimSpace(content)
			return md.String("~~" + content + "~~")
		},
	}

	conv := md.NewConverter("", true, nil)
	conv.AddRules(strikethrough)
	// -> add 1+ rules to the converter. the last added will be used first.

	markdown, err := conv.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n\nmarkdown:'%s'\n", markdown)
}
