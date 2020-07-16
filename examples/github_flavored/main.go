package main

import (
	"fmt"
	"log"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
)

func main() {
	html := `
	<ul>
		<li><input type=checkbox checked>Checked!</li>
		<li><input type=checkbox>Check Me!</li>
	</ul>
	`
	/*
		- [x] Checked!
		- [ ] Check Me!
	*/

	conv := md.NewConverter("", true, nil)

	// Use the `GitHubFlavored` plugin from the `plugin` package.
	conv.Use(plugin.GitHubFlavored())

	markdown, err := conv.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markdown)
}
