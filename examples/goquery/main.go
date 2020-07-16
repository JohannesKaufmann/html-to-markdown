package main

import (
	"fmt"
	"log"

	md "github.com/JohannesKaufmann/html-to-markdown"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://blog.golang.org/godoc-documenting-go-code"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	content := doc.Find("#content")

	conv := md.NewConverter(md.DomainFromURL(url), true, nil)
	markdown := conv.Convert(content)

	fmt.Println(markdown)
}
