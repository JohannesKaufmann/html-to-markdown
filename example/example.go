package main

import (
	"fmt"

	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

func init() {
	md.Blacklist = append(md.Blacklist, "input")

	// var Youtube md.Tag = "YOUTUBE"
	md.AddRule(
		md.Rule{
			HTMLNodes: []string{"span"},
			ToElement: func(domain string, isChildren bool, s *goquery.Selection, children md.ChildrenToElement) *md.Element {
				if !s.HasClass("bold") {
					return nil
				}
				return &md.Element{
					Tag:        md.Bold,
					ChildNodes: children(domain, true, s),
				}
			},
		},
	)

	/*
		htmlRules := make(map[string]func(*goquery.Selection) md.Element)
		htmlRules["iframe"] = func(s *goquery.Selection) md.Element {
			// if false {
			// 	return rules.Iframe(s)
			// }
			return md.Element{
				Tag:  Youtube,
				Text: "youtube_id",
			}
		}

		// TODO: parent element instead of tag?
		mdRules := make(map[md.Tag]func(element *md.Element, before *md.Element, after *md.Element, parent md.Tag) string)
		mdRules[Youtube] = func(element *md.Element, before *md.Element, after *md.Element, parent md.Tag) string {
			return "YOUTUBE_VIDEO"
		}

		_, _ = htmlRules, mdRules
	*/
}
func main() {
	domain := md.DomainFromURL("http://www.bonnerruderverein.de/ueber-uns/informationen/")
	fmt.Println("domain:", domain)
}
