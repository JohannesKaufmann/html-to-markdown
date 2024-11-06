package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"golang.org/x/net/html"
)

func main() {
	input := `
<h2>
	<span>Golang</span>
	<star-rating count="5">five stars</star-rating>
</h2>
<p>Build simple, secure, <i>scalable</i> systems with Go</p>
	`

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)

	// Here we a registering a custom *renderer* for <star-rating> and pass in our function.
	conv.Register.RendererFor("star-rating", converter.TagTypeInline, renderStarRating, converter.PriorityStandard)

	markdown, err := conv.ConvertString(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markdown)
	// ## Golang ⭐️⭐️⭐️⭐️⭐️
	//
	// Build simple, secure, *scalable* systems with Go
}

func renderStarRating(ctx converter.Context, w converter.Writer, node *html.Node) converter.RenderStatus {
	// The "github.com/JohannesKaufmann/dom" package provides helper functions
	// to interact with the html node, like getting the attribute "count".
	rawCount := dom.GetAttributeOr(node, "count", "0")
	count, _ := strconv.Atoi(rawCount)

	rating := strings.Repeat("⭐️", count)

	// Write the content
	w.WriteString(rating)

	// w.WriteString(" (")
	// ctx.RenderChildNodes(ctx, w, node)
	// w.WriteString(")")

	// And then return whether it was a *success*
	// or if the next renderer should be tried.
	return converter.RenderSuccess
}
