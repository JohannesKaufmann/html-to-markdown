package md

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var blockElements = []string{
	"address",
	"article",
	"aside",
	"audio",
	"video",
	"blockquote",
	"canvas",
	"dd",
	"div",
	"dl",
	"fieldset",
	"figcaption",
	"figure",
	"footer",
	"form",
	"h1", "h2", "h3", "h4", "h5", "h6",
	"header",
	"hgroup",
	"hr",
	"noscript",
	"ol", "ul",
	"output",
	"p",
	"pre",
	"section",
	"table", "tfoot",
}
var inlineElements = []string{ // -> https://developer.mozilla.org/de/docs/Web/HTML/Inline_elemente
	"b", "big", "i", "small", "tt",
	"abbr", "acronym", "cite", "code", "dfn", "em", "kbd", "strong", "samp", "var",
	"a", "bdo", "br", "img", "map", "object", "q", "script", "span", "sub", "sup",
	"button", "input", "label", "select", "textarea",
}

func IsInlineElement(e string) bool {
	for _, element := range inlineElements {
		if element == e {
			return true
		}
	}
	return false
}
func IsBlockElement(e string) bool {
	for _, element := range blockElements {
		if element == e {
			return true
		}
	}
	return false
}
func String(text string) *string {
	return &text
}

type Options struct {
	StrongDelimiter string
	Fence           string
	HR              string
}

type Rule struct {
	Filter      []string
	Replacement func(content string, selec *goquery.Selection, options *Options) *string
}

var leadingNewlinesR = regexp.MustCompile(`^\n+`)
var trailingNewlinesR = regexp.MustCompile(`\n+$`)

var newlinesR = regexp.MustCompile(`\n+`)
var tabR = regexp.MustCompile(`\t+`)
var indentR = regexp.MustCompile(`(?m)\n`)

func DefaultRule(content string, selec *goquery.Selection, opt *Options) *string {
	return &content
}
func KeepRule(content string, selec *goquery.Selection, opt *Options) *string {
	element := selec.Get(0)

	var buf bytes.Buffer
	err := html.Render(&buf, element)
	if err != nil {
		panic(err)
	}

	return String(buf.String())
}

func (c *Converter) selecToMD(domain string, selec *goquery.Selection, opt *Options) string {
	var builder strings.Builder
	// TODO: selec.Contents() Children
	// TODO: Text() or DirectText()
	selec.Contents().Each(func(i int, s *goquery.Selection) {
		name := goquery.NodeName(s)
		rules := c.getRuleFuncs(name)
		// r, ok := rules[name]
		// if !ok {
		// 	content := c.selecToMD(domain, s, opt)
		// 	res := DefaultRule(content, s, opt)
		// 	if name != "html" && name != "body" && name != "head" {
		// 		fmt.Println(name, "\t-> default rule")
		// 	}
		// 	builder.WriteString(*res)
		// 	return
		// }

		for i := len(rules) - 1; i >= 0; i-- {
			rule := rules[i]
			content := c.selecToMD(domain, s, opt)
			res := rule(content, s, opt)
			if res != nil {
				builder.WriteString(*res)
				return
			}
		}
	})
	return builder.String()
}
