package rules

import (
	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

/*
NameTag / Name
NameHtmlToElement
NameElementToMD
*/

type R struct {
	Tag       md.Tag
	ToElement func(*goquery.Selection) *md.Element
	ToMD      func(md.Element) *string
}

func AddRule(rules ...R) {

}
func z() {
	AddRule(R{}, R{})
}

type Rule interface {
	ToElement(*goquery.Selection) *md.Element
	ToMD(md.Element) string
}

type Iframe md.Tag

func (i Iframe) ToElement(s *goquery.Selection) *md.Element {
	return nil
}
func (i Iframe) ToMD(e md.Element) string {
	return ""
}

var IframeTag Iframe = "IFRAME"

func x(r Rule) {

}
func y() {

	x(IframeTag)
}
