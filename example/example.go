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

	convert := func(html string) {
		markdown, err := converter.ConvertString(html)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\tmd ->", markdown)
	}
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
