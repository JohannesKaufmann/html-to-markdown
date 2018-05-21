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
	converter.AddRules(plugin.Table...)
	converter.AddRules(strongRule)

	convert := func(html string) {
		markdown, err := converter.ConvertString(html)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("md ->", markdown)
	}
	go convert("<p>Hi</p>")
	go convert("<strong>Important</strong>")

	time.Sleep(time.Second * 10)
}
