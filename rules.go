package md

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func isInlineElement(e string) bool {
	for _, element := range inlineElements {
		if element == e {
			return true
		}
	}
	return false
}
func isBlockElement(e string) bool {
	for _, element := range blockElements {
		if element == e {
			return true
		}
	}
	return false
}

type rule struct {
	Filter []string

	Replacement func(content string, node *goquery.Selection, options map[string]string) string
}


type MarkdownConverter struct {
	Emphasis rule

	Other rule
}
func (conv *MarkdownConverter) X() {
}

var Options = map[string]string{ // TODO: maybe struct
	// headingStyle
	"hr":               "* * *",
	"bulletListMarker": "*",
	// codeBlockStyle
	"fence":           "```",
	"emDelimiter":     "_",
	"strongDelimiter": "**",
	// linkStyle
	// linkReferenceStyle
}
var Remove = []string{"script", "style", "#comment", "head", "svg"}
var Keep = []string{} // for example: iframe

// - - - - - - - - - - - - - //

// TODO: remove
var Blacklist = []string{"script", "style", "#comment", "head", "svg"}

type ToElements func(domain string, isChildren bool, s *goquery.Selection, children ToElements) []*Element

// type ChildrenToElement func(domain string, isChildren bool, s *goquery.Selection) []*Element
type ToMD func(element Element, before *Element, after *Element, parent *Element) *string

var toElementRules map[string][]ToElements
var toMDRules map[Tag][]ToMD

func init() {
	toElementRules = make(map[string][]ToElements)
	toMDRules = make(map[Tag][]ToMD)
	initDefaultRules()
}

type Rule struct {
	HTMLNodes  []string
	ToElements ToElements

	Tag  Tag
	ToMD ToMD
}

func AddRule(rules ...Rule) {
	for _, rule := range rules {
		for _, node := range rule.HTMLNodes {
			val, _ := toElementRules[node]
			val = append(val, rule.ToElements)
			toElementRules[node] = val
		}

		val, _ := toMDRules[rule.Tag]
		val = append(val, rule.ToMD)
		toMDRules[rule.Tag] = val
	}
}

func SelecToElem(domain string, isChildren bool, selec *goquery.Selection, _ ToElements) []*Element {
	var elements []*Element

	selec.Contents().Each(func(i int, s *goquery.Selection) {
		node := goquery.NodeName(s)
		for _, item := range Blacklist {
			if node == item {
				return
			}
		}

		rules, ok := toElementRules[node]
		if !ok {
			fmt.Println("no rule found for", node)
			return
		}
		for i := len(rules) - 1; i >= 0; i-- {
			toElements := rules[i]
			elem := toElements(domain, isChildren, s, SelecToElem)
			if elem != nil {
				elements = append(elements, elem...)
				return
			}
		}
	})
	return elements
}

func initDefaultRules() {
	/*
		e.Tag = Text
			e.ChildNodes = SelecToElem(domain, true, s)
			if len(e.ChildNodes) != 0 {
				elements = append(elements, e)
			}
	*/
	AddRule(
		Rule{
			HTMLNodes: []string{"p"},
			ToElements: func(domain string, isChildren bool, s *goquery.Selection, toElements ToElements) []*Element {
				children := toElements(domain, isChildren, s, toElements)
				if len(children) == 0 {
					return nil
				}

				return []*Element{
					{
						Tag:        Text,
						ChildNodes: children,
					},
				}
			},
		},
		Rule{
			HTMLNodes: []string{"#text"},
			ToElements: func(domain string, isChildren bool, s *goquery.Selection, toElements ToElements) []*Element {
				text := s.Text()
				text = removeTabs(text)

				trimed := strings.TrimSpace(text)
				if trimed == "" {
					return nil
				}

				return []*Element{
					{
						Tag:  TextNode,
						Text: text,
					},
				}
			},
		},
		Rule{ // container
			HTMLNodes: []string{
				"div",
				"ul",
				"ol",
				"section",
				"article",
				"aside",
				"footer",
				"nav",
				"header",
				"body",
				"html",
			},
			ToElements: func(domain string, isChildren bool, s *goquery.Selection, children ToElements) []*Element {
				elems := children(domain, isChildren, s, children)
				if len(elems) != 0 {
					var onlyTextNodes = true
					for _, elem := range elems {
						if !isInlineNode(elem.Tag) {
							onlyTextNodes = false
						}
					}
					if onlyTextNodes {
						return []*Element{
							{
								Tag:        Text,
								ChildNodes: elems,
							},
						}
						// e.Tag = Text
						// e.ChildNodes = elems
						// elements = append(elements, e)
					} else {
						return elems
						// elements = append(elements, elems...)
					}
				}
				return nil
			},
		},
	)
}
