package main

import (
	"fmt"
	"log"
	"time"

	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	converter := md.NewConverter("www.google.com", true, nil)

	strongRule := md.Rule{
		Filter: []string{"strong"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			fmt.Println("STRONG")
			return nil
		},
	}
	converter.AddRules(plugin.Strikethrough...)
	converter.AddRules(plugin.TaskListItems...)
	converter.AddRules(plugin.Table...)
	converter.AddRules(strongRule)
	converter.AddRules(plugin.Youtube...)

	// converter.Use(plugin.VimeoEmbed(plugin.VimeoWithTitle))
	converter.Use(plugin.VimeoEmbed(plugin.VimeoWithDescription))

	convert := func(html string) {
		markdown, err := converter.ConvertString(html)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\tmd ->", markdown)
	}
	go convert(`<iframe id="youtube-9742" frameborder="0" allowfullscreen="1" allow="autoplay; encrypted-media" title="Player for Code+Design Camp Berlin 04/2017" width="640" height="360" src="https://www.youtube.com/embed/xGk1PpIbisU?autoplay=0&amp;controls=0&amp;rel=0&amp;showinfo=0&amp;iv_load_policy=3&amp;cc_load_policy=0&amp;cc_lang_pref=en&amp;wmode=transparent&amp;modestbranding=1&amp;disablekb=1&amp;origin=https%3A%2F%2Fcode.design&amp;enablejsapi=1&amp;widgetid=1" tabindex="-1"></iframe>`)
	go convert(`<iframe src="http://player.vimeo.com/video/47387431?title=0&amp;byline=0&amp;portrait=0&amp;autoplay=0" width="1600" height="900" frameborder="0" webkitAllowFullScreen mozallowfullscreen allowFullScreen></iframe>`)
	go convert("<p>Hi</p>")
	go convert("<strong>Important</strong>")
	go convert("<del>Not Important</del>")
	go convert(`<ul><li><input type=checkbox checked>Checked!</li></ul>`)
	go convert(`<ul><li><input type=checkbox>Check Me!</li></ul>`)
	/*
					go convert(`
				<table>
					<thead>
						<tr>
							<th>Column 1</th>
							<th>Column 2</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>Row 1, Column 1</td>
							<td>Row 1, Column 2</td>
						</tr>
						<tr>
							<td>Row 2, Column 1</td>
							<td>Row 2, Column 2</td>
						</tr>
					</tbody>
				</table>
					`)
					go convert(`
						<table>
				      <thead>
				        <tr>
				          <th align="left">Column 1</th>
				          <th align="center">Column 2</th>
				          <th align="right">Column 3</th>
				          <th align="foo">Column 4</th>
				        </tr>
				      </thead>
				      <tbody>
				        <tr>
				          <td>Row 1, Column 1</td>
				          <td>Row 1, Column 2</td>
				          <td>Row 1, Column 3</td>
				          <td>Row 1, Column 4</td>
				        </tr>
				        <tr>
				          <td>Row 2, Column 1</td>
				          <td>Row 2, Column 2</td>
				          <td>Row 2, Column 3</td>
				          <td>Row 2, Column 4</td>
				        </tr>
				      </tbody>
				    </table>
					`)

			go convert(`
			<table>
		      <thead>
		        <tr>
		          <th align="left">Column 1</th>
		          <th align="center">Column 2</th>
		          <th align="right">Column 3</th>
		          <th align="foo">Column 4</th>
		        </tr>
		      </thead>
		      <tbody>
		        <tr>
		          <td></td>
		          <td>Row 1, Column 2</td>
		          <td>Row 1, Column 3</td>
		          <td>Row 1, Column 4</td>
		        </tr>
		        <tr>
		          <td>Row 2, Column 1</td>
		          <td></td>
		          <td>Row 2, Column 3</td>
		          <td>Row 2, Column 4</td>
		        </tr>
		        <tr>
		          <td>Row 3, Column 1</td>
		          <td>Row 3, Column 2</td>
		          <td></td>
		          <td>Row 3, Column 4</td>
		        </tr>
		        <tr>
		          <td>Row 4, Column 1</td>
		          <td>Row 4, Column 2</td>
		          <td>Row 4, Column 3</td>
		          <td></td>
		        </tr>
		        <tr>
		          <td></td>
		          <td></td>
		          <td></td>
		          <td>Row 5, Column 4</td>
		        </tr>
		      </tbody>
		    </table>
			`)
	*/
	go convert(`
	<table>
      <thead>
        <td>Heading 1</td>
        <td>Heading 2</td>
      </thead>
      <tbody>
        <tr>
          <td>Row 1</td>
          <td>Row 1</td>
        </tr>
        <tr>
          <td></td>
          <td></td>
        </tr>
        <tr>
          <td>Row 3</td>
          <td>Row 3</td>
        </tr>
      </tbody>
    </table>
	`)

	time.Sleep(time.Second * 10)
}
